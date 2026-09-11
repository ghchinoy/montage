---
title: Reviewing Applications & Repositories
description: Deep dive into the review workflow, understanding scorecards, and interpreting readiness tiers.
---

The `montage review` command executes the full ingestion, static code analysis, and capability evaluation pipeline.

## Review Flags

| Flag | Shorthand | Default | Description |
|---|:---:|:---:|---|
| `--out` | `-o` | `./out` | Output directory for assessment reports and scaffolded artifacts |
| `--format` | | `md,json` | Comma-separated output report formats: `md`, `json` |
| `--scaffold` | | `false` | Automatically generate Agent Plugin, Skill, and MCP scaffold |
| `--llm` | | `true` | Use Gemini / Vertex AI to enrich assessment if ambient credentials exist |
| `--sources` | | `./sources` | Local cache directory for cloned GitHub repositories |

## Ingestion Modes

### 1. Zero-Config Public GitHub Ingestion
Montage accesses public GitHub repositories anonymously without requiring a personal access token (`GITHUB_TOKEN`):
1. Queries the GitHub Public REST API for metadata (stars, languages, description, license).
2. Performs a shallow `git clone --depth 1` into the local `--sources` cache directory.

### 2. Authenticated GitHub Ingestion
If `GITHUB_TOKEN` or `GH_TOKEN` is exported, Montage uses authenticated API calls (raising rate limits to 5,000 req/hr) and can clone private repositories.

### 3. Local Directory Ingestion
When supplied with a local path, Montage verifies that the target exists and proceeds immediately with AST and regex analysis.
