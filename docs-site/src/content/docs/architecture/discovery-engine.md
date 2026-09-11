---
title: Discovery Engine & AST Analysis
description: How Montage deconstructs codebases into granular capabilities.
---

Rather than treating an application as a monolithic tool, Montage decomposes it into discrete capabilities:

1. **CLI Command Trees:** Analyzes Cobra, Click, Clap, and Argparse AST definitions to extract subcommands, descriptions, flags, and command groups.
2. **REST/HTTP Route Extraction:** Inspects route decorators and handler registrations (FastAPI, Express, Gin) to extract HTTP methods, paths, and request/response models.
3. **Standalone Scripts:** Discovers utility scripts in `scripts/`, `tools/`, or `bin/`, parsing docstrings to infer purpose and mutation nature.
4. **Authored Sub-Agent Prompts:** Discovers markdown prompts in `prompts/`, `agents/`, or `roles/` to identify specialized persona capabilities.
5. **Architectural Packages:** Maps packages under `pkg/` or `internal/` to high-level domain responsibilities.
