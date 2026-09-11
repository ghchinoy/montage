---
title: Scaffolding Agent Tooling
description: Generating Agent Plugins, Skills, and Model Context Protocol servers with Montage.
---

Montage scaffolds production-ready agent tooling tailored to the architecture of the analyzed target.

## Running the Scaffolder

```bash
montage scaffold <path-or-github-repo> --out ./agent-tooling
```

## Generated Artifacts

### 1. Agent Plugin Manifest (`plugin.json`)
Adheres to the [Agent Plugins Specification 1.0.0](https://agent-plugins.org/schemas/1.0.0/plugin.schema.json). Declares plugin capabilities, author, license, and registered MCP servers and skills.

### 2. MCP Server Configuration (`mcp.json`)
Specifies the stdio transport configuration for the scaffolded MCP server, ready for host agent registration.

### 3. Domain Skills (`skills/<target>/SKILL.md`)
Generated markdown skills formatted for `.agents/skills` or Claude Desktop:
- Frontmatter specifying `name`, `description`, `kind`, and `domain`.
- Core capabilities covered.
- Workflow guidelines, safety rules, and decision protocols.

### 4. Go MCP Server Starter (`cmd/<target>-mcp/main.go`)
A runnable Go MCP server powered by `github.com/modelcontextprotocol/go-sdk`:
- Registers candidate tools extracted from the application's CLI subcommands or REST endpoints.
- Groups tools into operational personas.
- Supports `stdio` transport.
