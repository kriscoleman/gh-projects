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
	"strings"

	"github.com/kriscoleman/gh-projects/internal/github"
)

func FilterIncompleteIssues(issues []*github.Issue, projectID string) []*github.Issue {
	var incomplete []*github.Issue
	
	for _, issue := range issues {
		if issue.State == "CLOSED" {
			continue
		}
		
		isDone := false
		
		for _, projectItem := range issue.ProjectItems.Nodes {
			for _, fieldValue := range projectItem.FieldValues.Nodes {
				if fieldValue.TypeName == "ProjectV2ItemFieldSingleSelectValue" {
					fieldName := strings.ToLower(fieldValue.Field.Name)
					valueName := strings.ToLower(fieldValue.Title)
					
					if (fieldName == "status" || fieldName == "state") && 
					   (valueName == "done" || valueName == "completed" || valueName == "closed") {
						isDone = true
						break
					}
				}
			}
			if isDone {
				break
			}
		}
		
		if !isDone {
			incomplete = append(incomplete, issue)
		}
	}
	
	return incomplete
}

func GetIssueStatus(issue *github.Issue) string {
	for _, projectItem := range issue.ProjectItems.Nodes {
		for _, fieldValue := range projectItem.FieldValues.Nodes {
			if fieldValue.TypeName == "ProjectV2ItemFieldSingleSelectValue" {
				fieldName := strings.ToLower(fieldValue.Field.Name)
				if fieldName == "status" || fieldName == "state" {
					return fieldValue.Title
				}
			}
		}
	}
	return "No Status"
}