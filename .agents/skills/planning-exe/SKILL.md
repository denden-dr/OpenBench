---
name: Planning Execution
description: Standarize planning execution
---

# Planning Execution Skill: From Blueprint to Code

This skill defines the standardized workflow for translating a high-level logical blueprint (created via the Planning Skill) into working, production-grade codebase. It is strictly designed to ensure consistent execution, enabling AI models to implement features step-by-step without losing context, hallucinating code, or deviating from architectural intent.

## Core Constraint: Strict Adherence
**CRITICAL**: You must follow the provided Implementation Plan exactly as written. Do not invent new features, alter the agreed-upon architecture, or skip steps. If the plan contains a logical flaw or relies on deprecated APIs that prevent execution, **STOP** and ask the user for clarification before proceeding.

## Execution Workflow

### 1. Initialization and Task Breakdown
- **Digest the Blueprint**: Carefully read the entire Implementation Plan, paying special attention to the "Structural Strategy" and "Best Practice & Quality Guardrails".
- **Locate Supplementary Context**: Actively check the same directory as the provided plan for additional contextual files. Specifically look for `flow.md` (for sequence/architecture diagrams) and `reference.md` (for best practice code templates). Review these thoroughly before writing any code to ensure precise execution.
- **Checklist Creation**: Create a localized tracking checklist (e.g., via a `task.md` artifact or internal list). Break down the "Step-by-Step Logic" from the plan into actionable task items.

### 2. Iterative Implementation (One Step at a Time)
Do not attempt to write all the code at once. Execute the plan sequentially:
1. **Understand Target Scope**: Identify the specific files and structural interfaces required for the current step.
2. **Translate Logic to Syntax**: Convert the human-readable instructions into the appropriate codebase. Ensure variable names and interfaces match the plan exactly.
3. **Embed Guardrails Immediately**: As you write the code, insert the prescribed error handling, logging (e.g., `logrus`), and validation checks defined in the plan. Do not leave "TODO" comments for guardrails; implement them natively.
4. **Mark Complete**: Check the item off your internal or artifact checklist before moving to the next.

### 3. Safe File Modification Protocols
- **Atomic Edits**: When modifying existing files, use precise edit tools. Do not overwrite entire files unless you are creating them from scratch.
- **Context Preservation**: Never delete or modify existing code, comments, or docstrings that fall outside the scope of the current task.
- **Centralized Dependencies**: If the plan calls for package installations (e.g., `go get`), ensure they are executed promptly so the IDE/linter context is accurate for subsequent steps.

### 4. Continuous Verification
- **Run Checks Promptly**: Execute required build commands (e.g., `make build`), formatters (e.g., `make fmt`), or linters as soon as a logical component is complete, rather than waiting until the very end.
- **Address Failures Fast**: If an intermediate build or test fails, diagnose based on the error output. Fix the syntax while remaining true to the logical blueprint.

### 5. Problem Recording
- **Do Not Solve Plan Flaws**: If you encounter any logical flaws, deprecated APIs, or architectural blockers during execution, **DO NOT** attempt to "solve" the plan or re-engineer the logic.
- **Update problem.md**: Create or update a `problem.md` file in the exact same directory as the plan.
- **Detail the Issue**: Record every single problem encountered. For each issue, write a clear problem description and provide your hypothetical reasoning for why it occurred.

### 6. Final Hand-off
Once the 'Step-by-Step Logic' is fully implemented:
- **Summarize**: Provide a brief walkthrough of what was accomplished mapping back to the Implementation Plan.
- **Trigger Verification**: Guide the user on how to run the manual steps outlined in the plan's "Verification Plan" (e.g., providing the `curl` commands or expected visual outcomes).

## Why This Works
This rigid execution framework ensures that the AI model acts as a highly disciplined "Builder" following an architect's blueprint. It prevents scope creep, maintains coding standards, ensures that the final output perfectly mirrors the agreed-upon logical design, and safely documents architectural issues for human review without hallucinating fixes.