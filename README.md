# gh-projects

A comprehensive CLI tool for managing GitHub Projects with automation capabilities.

## Features

### Iteration Management
- **Rollover automation**: Automatically move incomplete issues from previous iterations to current iteration
- **Interactive mode**: Review each issue before moving
- **Silent mode**: Batch move all incomplete issues
- **Dry-run support**: Preview changes before executing

### Extensible Architecture
- Subcommand structure for future GitHub Projects features
- Consistent flag patterns across commands
- Reusable components for GitHub API interactions

## Requirements

- Go 1.19 or higher
- [GitHub CLI](https://cli.github.com/) (`gh`) installed and authenticated
- Access to the GitHub Project you want to manage

## Installation

### From Source

```bash
git clone https://github.com/kriscoleman/gh-projects.git
cd gh-projects
go build -o gh-projects ./cmd/gh-projects
```

### Using Go Install

```bash
go install github.com/kriscoleman/gh-projects/cmd/gh-projects@latest
```

## Usage

### Iteration Rollover

Move incomplete issues from the previous iteration to the current one:

```bash
gh-projects iteration rollover --project https://github.com/users/USERNAME/projects/NUMBER
```

or for organizations:

```bash
gh-projects iteration rollover --project https://github.com/orgs/ORGNAME/projects/NUMBER
```

### Command Line Options

- `-p, --project` (required): GitHub project URL
- `-s, --silent`: Run in silent mode (automatically move all incomplete issues without prompts)
- `--dry-run`: Preview changes without making them

### Interactive Mode (Default)

In interactive mode, the tool will:
1. Show you each incomplete issue from the previous iteration
2. Display the issue's current status
3. Ask whether you want to move it to the current iteration

Example:
```bash
$ gh-projects iteration rollover -p https://github.com/users/myuser/projects/1

🚀 GitHub Projects - Iteration Rollover
======================================
📂 Project: myuser/1

🔄 Iteration Information
==================================================
Previous iteration: Sprint 23
Current iteration: Sprint 24
Iteration field: Iteration

🔍 Fetching issues from previous iteration...

📋 Incomplete issues found (3 issues):
--------------------------------------------------
• #123: Add user authentication [In Progress]
• #124: Update documentation [Todo]
• #125: Fix navigation bug [In Review]

🤔 Please review each issue:

📋 Issue #123: Add user authentication
   Status: In Progress
   State: OPEN
   Move to current iteration? (y/n/q): y
```

### Silent Mode

Automatically moves all incomplete issues without prompting:

```bash
gh-projects iteration rollover -p https://github.com/users/myuser/projects/1 --silent
```

### Dry Run Mode

Preview what would happen without making any changes:

```bash
gh-projects iteration rollover -p https://github.com/users/myuser/projects/1 --dry-run
```

## How It Works

1. **Authentication**: Uses your existing GitHub CLI authentication
2. **Project Discovery**: Finds the specified project and its iteration field
3. **Iteration Detection**: Identifies the current and most recent past iterations
4. **Issue Filtering**: Fetches all issues from the past iteration and filters out completed ones
5. **User Interaction**: In interactive mode, prompts for each issue; in silent mode, processes all automatically
6. **Updates**: Uses GitHub's GraphQL API to update the iteration field for selected issues

## Error Handling

The tool provides clear error messages for common issues:
- Missing GitHub CLI authentication
- Invalid project URLs
- Projects without iteration fields
- Network connectivity issues

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
