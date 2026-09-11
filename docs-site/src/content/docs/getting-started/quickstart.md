---
title: Quick Start
description: Start assessing codebases and generating agent tooling in minutes.
---

Montage can analyze local directories or fetch public GitHub repositories directly.

## 1. Review a Local Codebase

Run `montage review` against any directory:

```bash
montage review /path/to/project --out ./out
```

### Generated Outputs in `./out`
- `MONTAGE_ASSESSMENT.md`: Detailed markdown report containing the executive scorecard, dimensional analysis, capability matrix, and remediation code snippets.
- `montage.json`: Complete machine-readable assessment dataset.

## 2. Review a GitHub Repository Directly

Montage supports direct GitHub repository references without manual cloning:

```bash
# Full URL
montage review https://github.com/ghchinoy/repotographer --out ./out

# Or shorthand owner/repo:
montage review ghchinoy/repotographer --out ./out
```

## 3. Review and Scaffold Tooling in One Step

Pass the `--scaffold` flag to automatically generate the complete Agent Plugin, MCP server starter, and Skills:

```bash
montage review ghchinoy/repotographer --scaffold --out ./out
```

### Scaffolded Files
- `plugins/<target>/plugin.json`: Agent Plugin manifest adhering to the Agent Plugins Spec.
- `plugins/<target>/mcp.json`: MCP stdio server registration.
- `plugins/<target>/skills/<target>/SKILL.md`: Autonomous agent instructions and workflow.
- `skills/<target>/SKILL.md`: Standalone skill ready for `.agents/skills`.
- `cmd/<target>-mcp/main.go`: Go MCP server starter with candidate tools registered.
