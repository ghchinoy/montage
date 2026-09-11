---
title: Agent Plugins Specification
description: Reference guide for the Agent Plugins 1.0.0 standard used by Montage.
---

Montage adheres to the [Agent Plugins Specification](https://agent-plugins.org/schemas/1.0.0/plugin.schema.json).

### `plugin.json` Example

```json
{
  "$schema": "https://agent-plugins.org/schemas/1.0.0/plugin.schema.json",
  "name": "montage",
  "version": "0.2.0",
  "description": "Agent Ecosystem Readiness Assessor & Tooling Scaffolder",
  "author": {
    "name": "ghchinoy",
    "url": "https://github.com/ghchinoy"
  },
  "homepage": "https://github.com/ghchinoy/montage",
  "repository": "https://github.com/ghchinoy/montage",
  "license": "Apache-2.0",
  "keywords": ["mcp", "agent-plugins", "readiness-assessment", "scaffolding"]
}
```

### `mcp.json` Example

```json
{
  "$schema": "https://agent-plugins.org/schemas/1.0.0/mcp.schema.json",
  "mcpServers": {
    "montage": {
      "type": "stdio",
      "command": "montage",
      "args": ["mcp"]
    }
  }
}
```
