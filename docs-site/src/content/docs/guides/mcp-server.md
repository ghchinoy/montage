---
title: Model Context Protocol (MCP) Server
description: Running Montage as an MCP server to empower coding agents with repository assessment capabilities.
---

Montage can run as an MCP server over `stdio` to empower coding agents (Claude Desktop, OpenCode, Cursor, Antigravity) with repository assessment and tooling scaffolding tools.

## Starting the Server

```bash
montage mcp
```

## Exposed MCP Tools

| Tool | Parameters | Description |
|---|---|---|
| `assess_repository` | `target` (string, required)<br>`scaffold` (bool)<br>`use_llm` (bool)<br>`out_dir` (string) | Ingests a local path or GitHub repo, scores readiness, produces reports, and optionally scaffolds artifacts. |
| `evaluate_readiness` | `target` (string, required) | Fast in-memory evaluation returning 5 dimension scores, tier, and gap analysis without disk side effects. |
| `scaffold_tooling` | `target` (string, required)<br>`out_dir` (string) | Scaffolds the complete `plugins/`, `skills/`, and MCP server starter structure for an assessed target. |

## Client Configuration

### Claude Desktop (`claude_desktop_config.json`)
```json
{
  "mcpServers": {
    "montage": {
      "command": "/path/to/montage",
      "args": ["mcp"]
    }
  }
}
```

### OpenCode (`opencode.json`)
```json
{
  "mcpServers": {
    "montage": {
      "command": "/path/to/montage",
      "args": ["mcp"]
    }
  }
}
```
