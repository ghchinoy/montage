---
title: The 5 Readiness Dimensions
description: In-depth breakdown of the 5 architectural dimensions used to evaluate agent ecosystem readiness.
---

Montage evaluates applications across five core architectural dimensions, each weighted according to its impact on autonomous agent reliability:

```
Total Readiness Score (100%) =
    0.25 × Interface Automation +
    0.20 × Data Interchange +
    0.20 × Execution Safety +
    0.15 × Auth & Configuration +
    0.20 × Runtime & Packaging
```

---

## 1. 🤖 Interface Automation (Weight: 25%)
Evaluates whether the application can be operated headlessly by LLM agents via CLI, REST API, or deterministic scripts.

- **Strengths:** REST/HTTP API endpoints, hierarchical CLI command trees (Cobra, Click, Clap), standalone scripts, existing MCP server integration.
- **Blockers:** Interactive TTY prompts (`readline`, `inquirer`, `input()`, `fmt.Scan`) that halt unattended execution.

---

## 2. 📄 Data Interchange (Weight: 20%)
Evaluates machine readability, structured serialization, and output stream separation.

- **Strengths:** Explicit `--json` or `--format=json` flags, structured JSON/YAML serialization, clean stream separation (logs to `stderr`, data to `stdout`).
- **Blockers:** Human-formatted text tables requiring complex LLM regex scraping; interleaving diagnostics and payloads on `stdout`.

---

## 3. 🛡️ Execution Safety & Guardrails (Weight: 20%)
Evaluates risk profile, mutation boundaries, and the presence of pre-flight simulation flags.

- **Strengths:** Purely read-only / analytical operations, explicit `--dry-run` or `--validate-only` simulation flags, clear division between read and write commands.
- **Blockers:** Mutating or destructive operations (`delete`, `drop`, `kill`) without confirmation flags or dry-run simulation.

---

## 4. 🔑 Auth & Configuration (Weight: 15%)
Evaluates whether authentication and credentials can be supplied headlessly without GUI prompts.

- **Strengths:** Headless environment variables, configuration files (`config.yaml`), ambient CLI credentials (`gh auth`, cloud SDK ADC).
- **Blockers:** Interactive browser OAuth flows or interactive password prompts on TTY.

---

## 5. 📦 Runtime & Packaging (Weight: 20%)
Evaluates ease of execution in automated agent sandboxes and containerized runtimes.

- **Strengths:** Self-contained single binaries (Go, Rust), standardized build systems (`go.mod`, `Cargo.toml`, `Makefile`), Dockerfile/Containerfile.
- **Blockers:** Sprawling unpinned dependencies, complex multi-step manual setup procedures.
