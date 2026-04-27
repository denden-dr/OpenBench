---
name: Planning Problem Lister
description: List all planning problems in the project
---

# Planning Problem Lister Skill: Post-Execution Analysis

This skill is designed to be used after a plan (from the Planning Skill) has been executed (via the Planning Execution Skill). Its purpose is to systematically identify, document, and analyze every friction point, error, or unexpected behavior encountered during the implementation phase.

## Core Objective
To create a comprehensive "Post-Implementation Problem Report" that captures the "what" and the "why" of every issue, without focusing on the "how to fix" (which belongs in a new planning phase).

## Implementation Workflow

### 1. Problem Harvesting
- **Review Logs & History**: Analyze the conversation history, terminal outputs, and build/test failures encountered during execution.
- **Trace Mismatches**: Identify where the implementation deviated from the original plan or where the plan itself was insufficient.
- **Identify Friction**: Note any steps that took multiple attempts or required manual intervention that wasn't planned.

### 2. The Problem Report (Output)
Create a new file named `problems.md` (or update it if it exists) with the following structure for each identified problem:

#### A. Problem Definition
- **Description**: A clear, concise summary of the issue encountered.
- **Occurrence**: Where in the process did it happen? (e.g., "During Step 3 of the Execution Plan", "When running `make build`").

#### B. Hypothetical Reasoning
- **Root Cause Analysis**: Provide a hypothetical reasoning for why this happened. 
- **Contextual Factors**: Consider factors like:
    - Plan ambiguity.
    - Unexpected environment state.
    - Dependency version mismatches.
    - Hidden codebase complexity.
    - Model hallucinations or tool limitations.

### 3. Reporting Protocol
- **Be Brutally Honest**: Document even small "gotchas" that were easily fixed.
- **No Solutions**: Do **NOT** provide code fixes or remediation steps in this report. Focus purely on the post-mortem analysis of the problems themselves.
- **Categorization**: Group problems by type (e.g., Tooling, Logic, Environment, Blueprint) for easier review.

## Why This Works
The Planning Problem Lister ensures that:
1. **Lessons are Learned**: Prevents repeating the same mistakes in future planning phases.
2. **Technical Debt is Exposed**: Surfaces underlying issues that might have been "hacked around" during execution.
3. **Better Blueprints**: Provides direct feedback to the Planning Skill on where instructions were unclear or assumptions were wrong.
