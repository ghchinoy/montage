---
title: CLI Command Reference
description: Complete command-line reference for the Montage binary.
---

## Commands

### `montage review <target>`
Evaluates an application or GitHub repository for autonomous agent readiness.

```bash
montage review <path-or-repo> [flags]
```

### `montage scaffold <target>`
Generates Agent Plugin, Skill, and MCP Server artifacts for an application without writing reports.

```bash
montage scaffold <path-or-repo> [flags]
```

### `montage serve`
Starts the Montage interactive web application and REST API server.

```bash
montage serve --port 8080
```

### `montage mcp`
Runs Montage as a Model Context Protocol (MCP) server over `stdio`.

```bash
montage mcp
```

### `montage version`
Prints current Montage version.
