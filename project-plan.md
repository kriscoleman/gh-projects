# 📦 Epic: Iteration Rollover Automation Tool

## Overview

Our current GitHub Projects workflow relies on manually reassigning incomplete issues from a past iteration to the latest active iteration. When a new iteration begins and the previous one ends, the iteration view on the board becomes empty—despite carryover work still being relevant. This creates unnecessary friction in our iteration planning and visibility.

To improve developer experience (DevEx) and streamline this process, we propose building a **CLI tool in Golang** that leverages the **GitHub CLI** to automate the reassignment of incomplete issues to the current iteration.

---

## ✨ Goals

* Improve visibility and continuity in GitHub Projects across iteration boundaries.
* Eliminate the need for manual reassignment of carryover issues.
* Build a flexible, developer-friendly CLI tool that supports both interactive and silent modes.
* Lay the groundwork for future GitHub Projects automation features.

---

## 🧩 Scope

### ✅ MVP Features

1. **GitHub CLI Integration**

   * Authenticate and query a given GitHub Project using GitHub CLI.
   * Fetch the current and most recent past iteration cycles.

2. **Issue Evaluation**

   * Identify all issues assigned to the *most recently closed iteration*.
   * Filter out issues that are already marked as “Done”.

3. **Modes of Operation**

   * **Interactive Mode (default):**

     * Prompt the user for each incomplete issue to confirm whether it should be reassigned to the current iteration.
   * **Silent Mode (`--silent`):**

     * Automatically reassign all incomplete issues to the current iteration without prompts.

4. **Update Iteration Assignments**

   * Use GitHub CLI to update each selected issue’s iteration field to the latest cycle.

5. **Dry-Run Support**

   * Add a `--dry-run` flag to preview which issues *would* be reassigned without making changes.

---

## 📐 Technical Details

* **Language:** Go
* **Interface:** CLI
* **Dependency:** GitHub CLI (`gh`), with required permissions to access project fields and edit issues
* **Project Type Support:** GitHub Projects (beta/next-gen)
* **Field Type Handling:** Detect the iteration field from the schema of the project board

---

## 🚀 Future Enhancements (Post-MVP)

* Support for multiple project boards and org-level iteration configurations
* Logging and reporting modes (e.g., summary of changes)
* GitHub Actions integration to run the tool on a schedule or webhook trigger
* Support for other field types (e.g., status, custom labels)

---

## 🧪 Acceptance Criteria

* [ ] The CLI tool can list all issues from the last iteration that are not marked as done.
* [ ] In interactive mode, the tool prompts the user for each issue and reassigns as confirmed.
* [ ] In silent mode, the tool reassigns all incomplete issues without prompts.
* [ ] The reassignment updates the iteration field to the current cycle using GitHub CLI.
* [ ] A dry-run option is available for safe preview of actions.
* [ ] Tool has clear documentation for installation and usage.

---

## 🏁 Milestones

| Milestone              | Description                                                     | Owner | Target Date |
| ---------------------- | --------------------------------------------------------------- | ----- | ----------- |
| Schema Discovery       | Determine how to reliably identify iteration fields in Projects | TBD   | TBD         |
| CLI Skeleton           | Scaffold CLI in Go with basic argument parsing                  | TBD   | TBD         |
| GitHub CLI Integration | Implement logic to list and update issues via `gh`              | TBD   | TBD         |
| Modes & Flags          | Add support for `--silent` and `--dry-run`                      | TBD   | TBD         |
| Testing & Validation   | Functional tests with test projects                             | TBD   | TBD         |
| MVP Release            | Ready for internal use by teams                                 | TBD   | TBD         |

---

## 📚 Resources

* [GitHub CLI Docs](https://cli.github.com/manual/)
* [GitHub Projects API (GraphQL)](https://docs.github.com/en/graphql/overview/explorer)
* [Go CLI Libraries: cobra or urfave/cli](https://github.com/spf13/cobra)
