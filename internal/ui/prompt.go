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

package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kriscoleman/gh-projects/internal/github"
	"github.com/kriscoleman/gh-projects/internal/projects"
)

type Prompter struct {
	scanner *bufio.Scanner
}

func NewPrompter() *Prompter {
	return &Prompter{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (p *Prompter) ConfirmIssue(issue *github.Issue) bool {
	status := projects.GetIssueStatus(issue)
	fmt.Printf("\n📋 Issue #%d: %s\n", issue.Number, issue.Title)
	fmt.Printf("   Status: %s\n", status)
	fmt.Printf("   State: %s\n", issue.State)
	fmt.Print("   Move to current iteration? (y/n/q): ")
	
	p.scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(p.scanner.Text()))
	
	switch response {
	case "y", "yes":
		return true
	case "q", "quit":
		fmt.Println("\n❌ Operation cancelled by user")
		os.Exit(0)
	}
	
	return false
}

func (p *Prompter) ShowSummary(totalIssues, movedIssues int, dryRun bool) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 Summary")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total incomplete issues found: %d\n", totalIssues)
	
	if dryRun {
		fmt.Printf("Issues that would be moved: %d\n", movedIssues)
		fmt.Println("\n🔍 This was a dry run. No changes were made.")
	} else {
		fmt.Printf("Issues moved to current iteration: %d\n", movedIssues)
		fmt.Printf("Issues skipped: %d\n", totalIssues-movedIssues)
	}
}

func PrintIssueList(issues []*github.Issue, title string) {
	fmt.Printf("\n%s (%d issues):\n", title, len(issues))
	fmt.Println(strings.Repeat("-", 50))
	
	for _, issue := range issues {
		status := projects.GetIssueStatus(issue)
		fmt.Printf("• #%d: %s [%s]\n", issue.Number, issue.Title, status)
	}
}

func PrintIterationInfo(info *projects.IterationInfo) {
	fmt.Println("\n🔄 Iteration Information")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Previous iteration: %s\n", info.Previous.Title)
	fmt.Printf("Current iteration: %s\n", info.Current.Title)
	fmt.Printf("Iteration field: %s\n", info.Current.Field.Name)
}