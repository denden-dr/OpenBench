---
description: A standardized workflow for making high-quality, atomic, and conventionally formatted git commits.
---

# Git Commit Workflows: Professional Standard

This workflow ensures that every contribution to the repository is documented with clarity, precision, and adherence to industry best practices. Use this whenever you need to commit changes to the codebase.

## 1. Preparation & Staging

### A. Integrity Check & Hygiene
- **Formatting**: Run `make fmt` (or `go fmt ./...`) to ensure a consistent code style.
- **Dependency Management**: Run `make tidy` (or `go mod tidy`) to clean up `go.mod` and `go.sum`.
- **Verify Build**: Ensure the code compiles and passes basic local checks (e.g., `make build`).
- **Review Diffs**: Run `git diff` to review all changes. Ensure no debug logs, temporary comments, or secrets are leaked.
- **Group Changes**: Identify if the changes are "Atomic". If there are multiple unrelated changes (e.g., a bug fix and a feature), they MUST be committed separately.

### B. Logical Staging
- Stage files that belong to a single logical unit of work.
- Use `git add <file>` or `git add -p` for granular control.

## 2. Commit Message Construction

All commit messages must follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

### A. The Subject Line
**Format**: `<type>(<scope>): <description>`
- **Type**: Must be one of:
    - `feat`: A new feature.
    - `fix`: A bug fix.
    - `docs`: Documentation only changes.
    - `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc).
    - `refactor`: A code change that neither fixes a bug nor adds a feature.
    - `perf`: A code change that improves performance.
    - `test`: Adding missing tests or correcting existing tests.
    - `build`: Changes that affect the build system or external dependencies.
    - `ci`: Changes to CI configuration files and scripts.
    - `chore`: Other changes that don't modify src or test files.
    - `revert`: Reverts a previous commit.
- **Scope**: (Optional) A noun describing a section of the codebase (e.g., `auth`, `api`, `repo`).
- **Description**: 
    - Use the imperative, present tense: "change" not "changed" nor "changes".
    - Don't capitalize the first letter.
    - No dot (.) at the end.
    - Keep it under 50 characters.

### B. The Body (Required for complex changes)
- Separate from subject with a blank line.
- Wrap content at 72 characters.
- Focus on the **WHY** behind the change, rather than the **HOW**.
- Explain the context, motivation, and any trade-offs made.

### C. The Footer (Optional)
- Separate from body with a blank line.
- Use for referencing issue IDs (e.g., `Refs: #123`, `Closes: #456`).
- Use `BREAKING CHANGE: <description>` for API-breaking modifications.

## 3. Execution

### A. Commit Command
Execute the commit using:
```bash
git commit -m "<subject>" -m "<body>"
```
Or simply `git commit` if an editor is available.

### B. Post-Commit Verification
- Run `git log -1` to verify the message looks correct.
- Ensure only the intended files were committed.

## 4. Best Practice Guardrails

- **No "Work in Progress"**: Never commit messages like "WIP", "updates", or "fixing stuff".
- **Atomic Commits**: If you find yourself using "and" in the subject line, you probably should have two commits.
- **Signed-off-by**: If the project requires DCO, append `-s` to the commit command.
- **Keep History Clean**: Prefer `rebase` over `merge` for local feature branch updates to maintain a linear history when requested.

---
*Created by Antigravity - Advanced Agentic Coding Workflow*