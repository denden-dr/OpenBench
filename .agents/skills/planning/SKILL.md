---
name: Planning Skill
description: Skill to produce implementation planning
---

# Planning Skill: High-Level Logical Blueprinting

This skill is designed to bridge the gap between human requirements and technical implementation. It focuses on creating a "Logic First, Code Later" plan that ensures architectural integrity and best practices are locked in before a single line of code is written.

## Core Constraint: No Code Generation
**CRITICAL**: When using this skill, do NOT output code blocks or snippets. The output must be purely descriptive, logical, and structural. Use high-level human instructions that define "what" and "how" without providing the "syntax".

## Implementation Workflow

### 1. Context Acquisition
- **Search & Map**: Use research tools to find all relevant files, dependencies, and existing patterns related to the request.
- **Dependency Audit**: Identify if the change impacts existing upstream or downstream services.

### 2. The Implementation Plan (Structure)

#### A. Logical Requirements
- Define the exact problem being solved.
- Identify edge cases that must be handled (e.g., null inputs, network timeouts, race conditions).

#### B. Structural Strategy
- **File System Impact**: List every file to be created, modified, or deleted. 
- **Module Architecture**: Describe the relationship between new and existing components. Do not use code; use descriptions like "Component A will notify Service B via a synchronous call".
- **Interface specs**: Describe types and interfaces purely in natural language (e.g., "The request body should contain an array of session IDs, and the response should be a boolean indicating success").

#### C. Step-by-Step Logic
For each change, provide a numbered list of logical steps:
1. **Validation**: Describe what to check first.
2. **Transformation**: Describe how data changes.
3. **Persistance/Side Effects**: Describe database updates or external API calls.
4. **Resolution**: Describe what is returned or how the flow ends.

#### D. Best Practice & Quality Guardrails
- **Error Handling**: Specify precisely how each failure point should be caught and reported.
- **Security**: Define the validation and sanitization steps required.
- **Performance**: Note any asynchronous operations or caching strategies to implement.
- **Observability**: Specify where logging, metrics, or tracing should be added.

### 3. Verification Plan
- **Test Scenarios**: List specific scenarios (Success, Failure, Edge) that the code must pass.
- **Validation Steps**: Human-readable steps to verify the feature works as intended (e.g., "Check the 'users' table to ensure the password hash is present").

## Why This Works
By separating logic from syntax, the High-Level Planning Skill prevents:
1. **Copy-Paste Errors**: Basic models often replicate incorrect syntax from a plan.
2. **Logical Gaps**: Visualizing the flow in natural language exposes flaws that code hides.
3. **Inconsistent Quality**: Forcing a "Best Practice" section ensures technical debt is addressed upfront.
