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

package commands

import (
	"fmt"
	"strconv"
	
	"github.com/spf13/cobra"
	"github.com/kriscoleman/gh-projects/internal/github"
)

// BaseCommand provides common functionality for all commands
type BaseCommand struct {
	ProjectURL string
	DryRun     bool
	Silent     bool
	Token      string
}

// AddCommonFlags adds standard flags that many commands will need
func (b *BaseCommand) AddCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&b.ProjectURL, "project", "p", "", "GitHub project URL")
	cmd.Flags().BoolVar(&b.DryRun, "dry-run", false, "Preview changes without making them")
	cmd.Flags().BoolVarP(&b.Silent, "silent", "s", false, "Run in silent mode (no prompts)")
	cmd.Flags().StringVarP(&b.Token, "token", "t", "", "GitHub token for authentication (can also use GITHUB_TOKEN env var)")
}

// RequireProject marks the project flag as required
func (b *BaseCommand) RequireProject(cmd *cobra.Command) {
	cmd.MarkFlagRequired("project")
}

// GetGitHubClient creates and returns an authenticated GitHub client
func (b *BaseCommand) GetGitHubClient() (*github.Client, error) {
	return github.NewClient(b.Token)
}

// ParseProjectURL parses the project URL and returns owner and project number
func (b *BaseCommand) ParseProjectURL(client *github.Client) (string, int, string, error) {
	owner, numberStr, err := client.ParseProjectURL(b.ProjectURL)
	if err != nil {
		return "", 0, "", err
	}

	projectID, err := getProjectID(client, owner, numberStr)
	if err != nil {
		return "", 0, "", err
	}

	return owner, parseNumber(numberStr), projectID, nil
}

// Helper functions
func getProjectID(client *github.Client, owner, numberStr string) (string, error) {
	number := parseNumber(numberStr)
	projectResult, err := client.GraphQL(github.GetProjectQuery, map[string]interface{}{
		"owner":  owner,
		"number": number,
	})
	if err != nil {
		return "", err
	}

	// Try user projects first
	if userData, ok := projectResult["data"].(map[string]interface{})["user"].(map[string]interface{}); ok {
		if project, ok := userData["projectV2"].(map[string]interface{}); ok {
			return project["id"].(string), nil
		}
	}

	// Try organization projects
	if orgData, ok := projectResult["data"].(map[string]interface{})["organization"].(map[string]interface{}); ok {
		if project, ok := orgData["projectV2"].(map[string]interface{}); ok {
			return project["id"].(string), nil
		}
	}

	return "", fmt.Errorf("project not found")
}

func parseNumber(numberStr string) int {
	number, _ := strconv.Atoi(numberStr)
	return number
}