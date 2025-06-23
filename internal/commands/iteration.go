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

	"github.com/spf13/cobra"
	"github.com/kriscoleman/gh-projects/internal/github"
	"github.com/kriscoleman/gh-projects/internal/projects"
	"github.com/kriscoleman/gh-projects/internal/ui"
)

func NewIterationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iteration",
		Short: "Manage project iterations",
		Long:  `Commands for managing GitHub project iterations.`,
	}

	cmd.AddCommand(NewIterationRolloverCmd())
	return cmd
}

func NewIterationRolloverCmd() *cobra.Command {
	base := &BaseCommand{}

	cmd := &cobra.Command{
		Use:   "rollover",
		Short: "Roll over incomplete issues from previous iteration",
		Long: `Automatically reassign incomplete issues from the previous iteration 
to the current iteration in GitHub Projects.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIterationRollover(base)
		},
	}

	base.AddCommonFlags(cmd)
	base.RequireProject(cmd)

	return cmd
}

func runIterationRollover(base *BaseCommand) error {
	fmt.Println("🚀 GitHub Projects - Iteration Rollover")
	fmt.Println("======================================")

	client, err := base.GetGitHubClient()
	if err != nil {
		return fmt.Errorf("failed to initialize GitHub client: %w", err)
	}

	owner, number, projectID, err := base.ParseProjectURL(client)
	if err != nil {
		return err
	}

	fmt.Printf("📂 Project: %s/%d\n", owner, number)

	manager := projects.NewManager(client, projectID)

	iterationInfo, err := manager.GetIterations()
	if err != nil {
		return fmt.Errorf("failed to get iterations: %w", err)
	}

	ui.PrintIterationInfo(iterationInfo)

	fmt.Println("\n🔍 Fetching issues from previous iteration...")
	issues, err := manager.GetIterationItems(iterationInfo.Previous.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch issues: %w", err)
	}

	incompleteIssues := projects.FilterIncompleteIssues(issues, projectID)
	
	if len(incompleteIssues) == 0 {
		fmt.Println("\n✅ No incomplete issues found in the previous iteration!")
		return nil
	}

	ui.PrintIssueList(incompleteIssues, "📋 Incomplete issues found")

	var issuesToMove []*github.Issue
	prompter := ui.NewPrompter()

	if base.Silent {
		issuesToMove = incompleteIssues
		fmt.Printf("\n🤖 Silent mode: All %d incomplete issues will be moved\n", len(issuesToMove))
	} else {
		fmt.Println("\n🤔 Please review each issue:")
		for _, issue := range incompleteIssues {
			if prompter.ConfirmIssue(issue) {
				issuesToMove = append(issuesToMove, issue)
			}
		}
	}

	if !base.DryRun && len(issuesToMove) > 0 {
		fmt.Println("\n🔄 Moving issues to current iteration...")
		for i, issue := range issuesToMove {
			for _, item := range issue.ProjectItems.Nodes {
				err := manager.UpdateItemIteration(item.ID, iterationInfo.FieldID, iterationInfo.Current.ID)
				if err != nil {
					fmt.Printf("❌ Failed to move issue #%d: %v\n", issue.Number, err)
				} else {
					fmt.Printf("✅ Moved issue #%d (%d/%d)\n", issue.Number, i+1, len(issuesToMove))
				}
			}
		}
	}

	prompter.ShowSummary(len(incompleteIssues), len(issuesToMove), base.DryRun)

	return nil
}