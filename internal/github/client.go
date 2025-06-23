// Copyright 2025 Kris Coleman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"github.com/google/go-github/v67/github"
)

type Client struct {
	authenticated bool
	ghClient     *github.Client
	useToken     bool
}

func NewClient(token string) (*Client, error) {
	client := &Client{}
	
	if token != "" {
		client.ghClient = github.NewClient(nil).WithAuthToken(token)
		client.useToken = true
		client.authenticated = true
		return client, nil
	}
	
	if tokenFromEnv := os.Getenv("GITHUB_TOKEN"); tokenFromEnv != "" {
		client.ghClient = github.NewClient(nil).WithAuthToken(tokenFromEnv)
		client.useToken = true
		client.authenticated = true
		return client, nil
	}
	
	if err := client.checkAuth(); err != nil {
		return nil, fmt.Errorf("GitHub CLI authentication failed: %w", err)
	}
	client.authenticated = true
	return client, nil
}

func (c *Client) checkAuth() error {
	cmd := exec.Command("gh", "auth", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("not authenticated with GitHub CLI: %s", string(output))
	}
	return nil
}

func (c *Client) GraphQL(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	if c.useToken {
		return c.executeGraphQLWithToken(query, variables)
	}
	return c.executeGraphQLWithCLI(query, variables)
}

func (c *Client) executeGraphQLWithToken(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	req := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}
	
	resp, err := c.ghClient.Client().Post("https://api.github.com/graphql", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("GraphQL request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode GraphQL response: %w", err)
	}
	
	if errors, ok := result["errors"]; ok {
		return nil, fmt.Errorf("GraphQL errors: %v", errors)
	}
	
	return result, nil
}

func (c *Client) executeGraphQLWithCLI(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	args := []string{"api", "graphql", "-f", fmt.Sprintf("query=%s", query)}
	
	for key, value := range variables {
		switch v := value.(type) {
		case int:
			args = append(args, "-F", fmt.Sprintf("%s=%d", key, v))
		case float64:
			args = append(args, "-F", fmt.Sprintf("%s=%d", key, int(v)))
		default:
			args = append(args, "-f", fmt.Sprintf("%s=%v", key, value))
		}
	}
	
	cmd := exec.Command("gh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	
	// Try to parse the output regardless of error status
	// gh returns non-zero exit code even for partial GraphQL errors
	var result map[string]interface{}
	if parseErr := json.Unmarshal(stdout.Bytes(), &result); parseErr != nil {
		// If we can't parse and there was an exec error, return the exec error
		if err != nil {
			return nil, fmt.Errorf("GraphQL query failed: %v\nStderr: %s\nStdout: %s", err, stderr.String(), stdout.String())
		}
		// Otherwise return the parse error
		return nil, fmt.Errorf("failed to parse GraphQL response: %w\nOutput: %s", parseErr, stdout.String())
	}
	
	// Check if we have data - partial errors are OK for queries that try multiple paths
	if data, hasData := result["data"]; hasData && data != nil {
		return result, nil
	}
	
	// No data at all - this is a real error
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %v\nStderr: %s", err, stderr.String())
	}
	
	if errors, ok := result["errors"]; ok {
		return nil, fmt.Errorf("GraphQL errors: %v", errors)
	}
	
	return result, nil
}

func (c *Client) ParseProjectURL(url string) (owner, number string, err error) {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimSuffix(url, "/")
	
	parts := strings.Split(url, "/")
	
	if !strings.Contains(url, "github.com") || !strings.Contains(url, "/projects/") {
		return "", "", fmt.Errorf("invalid GitHub project URL format: must be github.com URL with /projects/")
	}
	
	for i, part := range parts {
		if part == "github.com" && i+2 < len(parts) && parts[i+2] == "projects" && i+3 < len(parts) {
			return parts[i+1], parts[i+3], nil
		}
		if part == "users" && i+1 < len(parts) && i+2 < len(parts) && parts[i+2] == "projects" && i+3 < len(parts) {
			return parts[i+1], parts[i+3], nil
		}
		if part == "orgs" && i+1 < len(parts) && i+2 < len(parts) && parts[i+2] == "projects" && i+3 < len(parts) {
			return parts[i+1], parts[i+3], nil
		}
	}
	
	return "", "", fmt.Errorf("could not parse project URL: expected format https://github.com/{owner}/projects/{number}")
}