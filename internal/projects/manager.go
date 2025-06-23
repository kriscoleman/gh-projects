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

package projects

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/kriscoleman/gh-projects/internal/github"
)

type Manager struct {
	client    *github.Client
	projectID string
}

func NewManager(client *github.Client, projectID string) *Manager {
	return &Manager{
		client:    client,
		projectID: projectID,
	}
}

type IterationInfo struct {
	Current  *github.Iteration
	Previous *github.Iteration
	FieldID  string
}

func (m *Manager) GetIterations() (*IterationInfo, error) {
	result, err := m.client.GraphQL(github.GetProjectFieldsQuery, map[string]interface{}{
		"projectId": m.projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project fields: %w", err)
	}

	node, ok := result["data"].(map[string]interface{})["node"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid project response structure")
	}

	fields, ok := node["fields"].(map[string]interface{})["nodes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no fields found in project")
	}

	var iterationField map[string]interface{}
	var fieldID string

	for _, field := range fields {
		f := field.(map[string]interface{})
		if dataType, ok := f["dataType"].(string); ok && dataType == "ITERATION" {
			iterationField = f
			fieldID = f["id"].(string)
			break
		}
	}

	if iterationField == nil {
		return nil, fmt.Errorf("no iteration field found in project")
	}

	config, ok := iterationField["configuration"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("iteration field has no configuration")
	}

	iterations, ok := config["iterations"].([]interface{})
	if !ok {
		iterations = []interface{}{}
	}
	
	completedIterations, ok := config["completedIterations"].([]interface{})
	if !ok {
		completedIterations = []interface{}{}
	}
	
	// Combine all iterations
	allIterations := append(completedIterations, iterations...)
	
	if len(allIterations) == 0 {
		return nil, fmt.Errorf("no iterations found in project")
	}

	parsedIterations := make([]*github.Iteration, 0, len(allIterations))
	now := time.Now()
	
	log.Printf("Current time: %v", now)
	log.Printf("Found %d active iterations from API", len(iterations))
	log.Printf("Found %d completed iterations from API", len(completedIterations))
	log.Printf("Total iterations: %d", len(allIterations))

	for _, iter := range allIterations {
		i := iter.(map[string]interface{})
		startDateStr := i["startDate"].(string)
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			continue
		}

		iteration := &github.Iteration{
			ID:        i["id"].(string),
			Title:     i["title"].(string),
			StartDate: startDate,
			Duration:  int(i["duration"].(float64)),
		}
		iteration.Field.ID = fieldID
		iteration.Field.Name = iterationField["name"].(string)

		parsedIterations = append(parsedIterations, iteration)
		
		endDate := iteration.StartDate.AddDate(0, 0, iteration.Duration)
		log.Printf("Iteration %s: %s to %s", iteration.Title, iteration.StartDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}

	sort.Slice(parsedIterations, func(i, j int) bool {
		return parsedIterations[i].StartDate.Before(parsedIterations[j].StartDate)
	})

	var current, previous *github.Iteration

	// Find the most recent past iteration and the current/future iteration
	for i, iter := range parsedIterations {
		endDate := iter.StartDate.AddDate(0, 0, iter.Duration)
		
		log.Printf("Checking iteration %s: start=%v, end=%v, now=%v", 
			iter.Title, iter.StartDate, endDate, now)
		log.Printf("Start <= now: %v, End > now: %v", 
			!iter.StartDate.After(now), endDate.After(now))
		
		// Check if this is the current iteration (we're within its date range)
		if !iter.StartDate.After(now) && endDate.After(now) {
			// This is the current iteration
			current = iter
			if i > 0 {
				previous = parsedIterations[i-1]
			}
			log.Printf("Found current iteration: %s", current.Title)
			break
		} else if endDate.Before(now) || endDate.Equal(now) {
			// This iteration has ended, keep it as the most recent previous
			previous = iter
			log.Printf("Setting previous iteration: %s", previous.Title)
		} else if iter.StartDate.After(now) && current == nil {
			// This is a future iteration, use it as current if we haven't found one yet
			current = iter
			log.Printf("Found future iteration as current: %s", current.Title)
			// The previous iteration is already set from the loop
			break
		}
	}

	// If we don't have a current iteration, use the next upcoming one
	if current == nil && len(parsedIterations) > 0 {
		// Find the first future iteration
		for _, iter := range parsedIterations {
			if iter.StartDate.After(now) {
				current = iter
				break
			}
		}
	}

	if current == nil {
		return nil, fmt.Errorf("no current or future iteration found")
	}


	if previous == nil {
		return nil, fmt.Errorf("no previous iteration found - need at least 2 iterations to perform rollover")
	}
	
	log.Printf("Selected current iteration: %s", current.Title)
	log.Printf("Selected previous iteration: %s", previous.Title)

	return &IterationInfo{
		Current:  current,
		Previous: previous,
		FieldID:  fieldID,
	}, nil
}

func (m *Manager) GetIterationItems(iterationID string) ([]*github.Issue, error) {
	var allItems []*github.Issue
	var cursor string
	hasNextPage := true

	for hasNextPage {
		variables := map[string]interface{}{
			"projectId": m.projectID,
		}
		if cursor != "" {
			variables["after"] = cursor
		}

		result, err := m.client.GraphQL(github.GetIterationItemsQuery, variables)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch iteration items: %w", err)
		}

		data := result["data"].(map[string]interface{})
		node := data["node"].(map[string]interface{})
		items := node["items"].(map[string]interface{})
		
		pageInfo := items["pageInfo"].(map[string]interface{})
		hasNextPage = pageInfo["hasNextPage"].(bool)
		if endCursor, ok := pageInfo["endCursor"].(string); ok {
			cursor = endCursor
		}

		nodes := items["nodes"].([]interface{})
		
		for _, item := range nodes {
			itemData := item.(map[string]interface{})
			
			
			content, ok := itemData["content"].(map[string]interface{})
			if !ok || content == nil {
				log.Printf("Item has no content or content is nil")
				continue
			}
			
			// Check if this is an issue by looking for required fields
			if _, hasID := content["id"].(string); !hasID {
				continue
			}
			if _, hasNumber := content["number"].(float64); !hasNumber {
				continue
			}
			

			issue := &github.Issue{
				ID:     content["id"].(string),
				Number: int(content["number"].(float64)),
				Title:  content["title"].(string),
				State:  content["state"].(string),
			}

			fieldValues := itemData["fieldValues"].(map[string]interface{})["nodes"].([]interface{})
			itemID := itemData["id"].(string)
			
			hasIterationMatch := false
			var fieldValueNodes []github.FieldValue
			
			for _, fv := range fieldValues {
				fieldValue := fv.(map[string]interface{})
				
				if fieldValue["__typename"] == "ProjectV2ItemFieldIterationValue" {
					if iterationId, ok := fieldValue["iterationId"].(string); ok {
						if iterationId == iterationID {
							hasIterationMatch = true
						}
					} else {
						log.Printf("No iterationId found in field value")
					}
				}
				
				if fieldValue["__typename"] == "ProjectV2ItemFieldSingleSelectValue" {
					field := fieldValue["field"].(map[string]interface{})
					fieldValueNodes = append(fieldValueNodes, github.FieldValue{
						TypeName: "ProjectV2ItemFieldSingleSelectValue",
						Field: struct {
							ID   string
							Name string
						}{
							ID:   field["id"].(string),
							Name: field["name"].(string),
						},
						Title: fieldValue["name"].(string),
					})
				}
			}
			
			if hasIterationMatch {
				issue.ProjectItems.Nodes = append(issue.ProjectItems.Nodes, struct {
					ID         string
					FieldValues struct {
						Nodes []github.FieldValue
					}
				}{
					ID: itemID,
					FieldValues: struct {
						Nodes []github.FieldValue
					}{
						Nodes: fieldValueNodes,
					},
				})
				
				allItems = append(allItems, issue)
			}
		}
	}

	return allItems, nil
}

func (m *Manager) UpdateItemIteration(itemID, fieldID, iterationID string) error {
	_, err := m.client.GraphQL(github.UpdateItemIterationMutation, map[string]interface{}{
		"projectId":   m.projectID,
		"itemId":      itemID,
		"fieldId":     fieldID,
		"iterationId": iterationID,
	})
	
	if err != nil {
		return fmt.Errorf("failed to update item iteration: %w", err)
	}
	
	return nil
}