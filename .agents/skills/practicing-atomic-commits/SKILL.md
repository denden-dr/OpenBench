---
name: practicing-atomic-commits
description: Use when staging and committing code changes to Git, or when git history needs to be clean, readable, and easy to roll back.
---

# Practicing Atomic Commits

## Overview
An **Atomic Commit** is a commit that applies a single, complete, and self-contained change to the repository. It represents a single logical unit of work and leaves the codebase in a compilable, clean state with all tests passing.

---

## When to Use
* Preparing local changes to be committed.
* Staging selective changes from files containing multiple unrelated edits.
* Cleaning up commit history before pushing to remote or creating a Pull Request.

### When NOT to Use
* During rapid prototyping/local experimentation (where you don't commit or only use temporary WIP branches).
* When emergency hotfixing is handled by automated system rollbacks (though even then, atomic commits prevent further issues).

---

## Core Pattern

### Before (Anti-pattern: Bundled/Monolithic Commit)
```
Commit: "feat: add user profile page and fix database query error and update config"
- Modified: apps/frontend/src/routes/profile/+page.svelte (new feature)
- Modified: apps/backend/internal/database/database.go (unrelated bug fix)
- Modified: apps/backend/.env.example (unrelated config update)
```
*Issue:* Hard to revert the bug fix without reverting the feature, hard to review, and hard to write a clear commit message.

### After (Atomic Commits)
```
Commit 1: "refactor(config): clean up deprecated env variables"
Commit 2: "fix(database): resolve query timeout error in NewConnection"
Commit 3: "feat(frontend): implement user profile page UI"
```
*Benefit:* Each commit is isolated, easy to revert, easy to read, and can be cherry-picked.

---

## Quick Reference

| Command | Purpose |
|---------|---------|
| `git diff` | Review unstaged changes before staging. |
| `git add -p <file>` | Stage individual hunks/patches instead of the whole file. |
| `git commit -v` | Open editor for commit message with the diff shown below. |
| `git checkout -p <file>` | Discard specific hunks/patches. |
| `git rebase -i HEAD~<n>` | Interactively squash, reword, edit, or split recent commits. |

---

## Implementation Workflow

### Step 1: Review Changes
Before staging any file, review what has changed:
```bash
git diff
```

### Step 2: Plan and Get User Confirmation (CRITICAL)
Before executing any commits or staging files, present a clear overview of the planned commits to the user:
- List each proposed commit.
- Show which files/changes will go into each commit.
- Show the proposed conventional commit message for each.

**You must ask for explicit confirmation (e.g. "Do you approve these planned commits?") before proceeding.**

### Step 3: Interactive Staging
Once confirmed, stage the files or hunks for the first commit:
```bash
git add -p
```
* During `git add -p`, use:
  - `y` to stage the hunk.
  - `n` to skip the hunk.
  - `s` to split the current hunk into smaller hunks.
  - `e` to manually edit the hunk.

### Step 4: Run Tests
Ensure the code compiles and all tests pass with only the staged changes:
```bash
# Run unit tests
go test -count=1 ./...
```

### Step 5: Commit Staged Changes
Commit the staged changes using the approved conventional commit message:
```bash
git commit -m "type(scope): description"
```
*Common types:* `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `ci`.

Repeat Steps 3 to 5 for each planned commit.

---

## Common Mistakes

### 1. "While-I'm-Here" Commits
* **Problem:** Fixing a typo or cleaning up formatting in an unrelated file while implementing a feature, and committing them together.
* **Fix:** Stash your current work, check out a temporary branch/commit, fix the typo, commit it, then return to your work. Or use `git add -p` to stage them separately.

### 2. Committing Broken Code
* **Problem:** Committing incomplete work that breaks compilation or fails tests.
* **Fix:** Keep unfinished work unstaged or in a local WIP commit that you intend to amend (`git commit --amend`) or squash later.

### 3. Mixing Refactoring and Features
* **Problem:** Changing how a function works (refactoring) and adding new functionality at the same time in one commit.
* **Fix:** First refactor the code in one commit, run tests, commit, and then implement the new feature in a separate commit.

---

## Red Flags
* Using `git add .` or `git add -A` blindly without reviewing diffs.
* Commit messages containing the word `"and"` (e.g. `fix A and add B`).
* Commit messages like `fixes`, `update`, or `wip` on shared/pushed branches.
* Commits that fail compilation or test suites.
