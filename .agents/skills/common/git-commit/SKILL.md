---
name: git-commit
description: Use when analyzing current changes with git to create standardized, contextual commits without blindly adding all files.
---

# Git Commit Skill

This skill defines the workflow for safely analyzing workspace changes and creating clean, standard Git commits based on the context of the changes.

## 🚫 CRITICAL RULE: NEVER USE `git add .`
You must **never** use `git add .`, `git add -A`, or `git commit -a`. These commands risk adding unintended changes, debug code, scratch files, or unrelated modifications.

## Workflow

### 1. Check for `.gitignore`
Before proceeding, verify that a `.gitignore` file exists at the root of the repository. If it does not exist, you must recommend creating one to the user to prevent accidental commits of temporary files, secrets, or build artifacts.

### 2. Analyze the Context
Before staging, you must understand exactly what has changed:
- Run `git status` to view modified, untracked, and deleted files.
- Run `git diff <file>` to inspect the line-by-line context of changes.
- Identify the logical groups of changes. Do not bundle unrelated changes.

### 3. Stage Files Selectively
Based on your analysis, explicitly stage **only** the relevant files:
- `git add <file1> <file2> ...`

### 4. Commit with a Standardized Message
Format the commit message using the Conventional Commits specification.

**Format:**
```
<type>[optional scope]: <description>

[optional body]
```

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

**Description Guidelines:**
- Use the imperative, present tense: "change" not "changed" nor "changes".
- Do not capitalize the first letter.
- No period (.) at the end.

### Example
1. `git status`
2. `git diff src/ui.ts`
3. `git add src/ui.ts`
4. `git commit -m "feat(ui): add loading spinner during data fetch"`