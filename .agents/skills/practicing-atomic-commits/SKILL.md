---
name: practicing-atomic-commits
description: Use when staging and committing code changes to Git, or when git history needs to be clean, readable, and easy to roll back. Do not use for local experiments, prototyping, or non-git version control systems.
version: 1.0.0
---

# Practicing Atomic Commits

## Overview
An **Atomic Commit** is a commit that applies a single, complete, and self-contained change to the repository. It represents a single logical unit of work and leaves the codebase in a compilable, clean state with all tests passing.

## When to Use
* Preparing local changes to be committed.
* Staging selective changes from files containing multiple unrelated edits.
* Cleaning up commit history before pushing to remote or creating a Pull Request.

## Step-by-Step Instructions

1. **Review Changes**: Run `git diff` to inspect all unstaged modifications before staging.
2. **Plan Commits and Ask Permission**: Present a list of planned commits, matching files, and proposed conventional commit messages. Ask the user for explicit permission before executing any staging or committing.
3. **Stage Changes Interactively**: Run `git add -p` to stage specific hunks and patches instead of staging whole files at once.
4. **Run Verification Tests**: Execute unit tests (e.g. `go test -count=1 ./...`) to verify that the code compiles and passes all checks with only the staged changes.
5. **Commit Changes**: Execute `git commit -m "type(scope): description"` to commit changes using conventional commit formatting. Repeat from step 3 for subsequent planned commits.

## Common Mistakes
* **"While-I'm-Here" Commits**: Fixing typos or formatting in unrelated files and committing them with the main change. Use `git add -p` to stage them separately.
* **Committing Broken Code**: Committing incomplete work that breaks compilation or fails tests. Keep unfinished work unstaged.
* **Mixing Refactoring and Features**: Combining refactoring and new feature logic in a single commit. Refactor first in one commit, verify tests, and implement features in a separate commit.

## Red Flags
* Using `git add .` or `git add -A` blindly without reviewing diffs.
* Commit messages containing the word `"and"` (e.g. `fix A and add B`).
* Commit messages like `fixes`, `update`, or `wip` on shared/pushed branches.
* Commits that fail compilation or test suites.
