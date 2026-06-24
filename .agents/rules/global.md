---
trigger: always_on
---

# Global Agent Rules

## 1. Always Load AGENTS.md
- At the start of any interaction, the agent **MUST** load and read the contents of [AGENTS.md](file:///home/denden/Documents/denden/personal/personal-project/OpenBench/AGENTS.md) at the repository root.
- All guidelines, coding standards, project structures, and testing strategies defined in `AGENTS.md` must be strictly followed.

## 2. Always Use Graphify for Codebase Context
- To understand the codebase, trace connections, locate components, or answer general/specific codebase questions, the agent **MUST** use the `graphify` skill and the `graphify` CLI commands or MCP capabilities.
- **Check Existing Graph**: Before attempting to read or scan files manually, check if `graphify-out/graph.json` exists in the working directory.
- **Query the Graph**: If the graph exists, prioritize using:
  - `graphify query "<question>"` (BFS) or `graphify query "<question>" --dfs` (DFS) to locate context and trace relationships.
  - `graphify explain "<node>"` to analyze specific components.
  - `graphify path "<nodeA>" "<nodeB>"` to find paths between components.
- **Update Graph**: When code is added or modified, update the knowledge graph using:
  - `graphify --update` (incremental update) to ensure context remains fresh and accurate.
