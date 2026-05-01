# Planning Rule

Always store implementation plans in the `.agents/plan` directory. 
If the directory does not exist, it must be created before saving the plan.

## File Path Pattern
`.agents/plan/YYYY-MM-DD-<feature-name>.md`

## Constraints
- **Never** include `git commit` as a planning step. Commits are handled by the agent during execution and should not be explicitly part of the task list.
