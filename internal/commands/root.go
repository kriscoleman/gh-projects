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
	"github.com/spf13/cobra"
)

var version = "dev"

func SetVersion(v string) {
	version = v
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gh-projects",
		Short:   "A CLI tool for managing GitHub Projects",
		Long:    `A comprehensive CLI tool for automating GitHub Projects management tasks.`,
		Version: version,
	}

	// Add subcommands
	cmd.AddCommand(NewIterationCmd())

	return cmd
}