---
name: montage
description: Evaluate applications and GitHub repositories for agent readiness, extract fine-grained capabilities (APIs, CLI command trees, scripts, authored sub-agent prompts, package engines), and generate modular Agent Plugins, multi-skill blueprints, and clustered MCP servers.
license: Apache-2.0
metadata:
  version: "0.2.0"
---

# Montage Agent Assessment & Scaffolding Skill

This skill guides coding assistants and developers in evaluating codebases to determine how their capabilities can be integrated into the **Agent Ecosystem** (as **Pure Skills**, **Skills with Code**, **MCP Servers**, or **Agent Plugins**).

Rather than treating target applications as monolithic 1:1 blobs, Montage performs a **Capability-Centric Assessment**, decomposing the codebase into distinct functional domains, clustering them into operational personas, and generating specialized multi-skill blueprints.

---

## The Four Agent Tooling Tiers

When reviewing an application, Montage classifies discovered capabilities across four distinct tiers:

1. **Pure Skill (Guidelines, Decision Rules & Personas):**
   - Natural-language operational manuals, domain rules, validation checklists, and authored sub-agent prompts (`prompts/*.md`).
   - Requires zero code changes or runtime dependencies.
   - Ideal for advisory review modes, clinical/veterinary triage guidelines, editorial style standards, and architectural interpretation.

2. **Skill with Code (Skill + Helper Script):**
   - Skill instructions paired with a lightweight deterministic helper script (Python, Go, Bash).
   - Offloads raw data filtering, log trend crunching, or format conversion to code, reducing LLM context window token usage by 60–90%.
   - Ideal for mathematical analysis (e.g. weight velocity, budget caps, code categorization) before agent reasoning.

3. **MCP Server (Model Context Protocol):**
   - Standalone `stdio` or SSE server exposing typed tools, resources, and prompts with strict schemas.
   - Clustered by operational persona (e.g. `Clinical Telemetry`, `Adoption Services`, `Cartography & Visualization`).
   - Clearly annotates read-only queries versus mutating actions.

4. **Agent Plugin (Skill + MCP Server Bundle):**
   - Packaged bundle adhering to the [Agent Plugins Specification](https://github.com/agentplugins/agent-plugins-spec).
   - Contains `plugin.json` manifest, `mcp.json` server declaration, and domain-specialized `skills/<name>/SKILL.md` instructions.
   - Delivers turn-key agent capability discovery for Claude Desktop, OpenCode, and Antigravity.

---

## Discovered Capability Sources

Montage inspects five distinct architectural layers:
1. **REST & HTTP APIs:** FastAPI, Flask, Express, and Gin routes (`GET`, `POST`), extracting path parameters, docstrings, and request/response models.
2. **CLI Command Hierarchies:** Cobra, Click, Clap, and Commander command trees, honoring declared command groups (e.g. `execution`, `management`, `discovery`).
3. **Standalone Utility Scripts:** Python, Bash, and Node scripts in `scripts/`, `tools/`, and `bin/`.
4. **Authored Sub-Agent Prompts:** Author-designed persona files in `prompts/`, `pkg/prompts/`, and `.prompts/` (e.g. `reviewer.md`, `code_analyst.md`, `writer.md`).
5. **Architectural Package Engines:** Internal engine packages in `pkg/` or `internal/` (e.g. `orchestrator/`, `generators/`, `registry/`, `a2api/`).

---

## Available MCP Tools

When Montage is connected as an MCP server, the following tools are available:

1. `assess_repository`:
   - Inspects a local directory or GitHub repository (`owner/repo` or URL).
   - Evaluates the 5 readiness dimensions (Interface Automation, Data Interchange, Execution Safety, Headless Auth, Packaging).
   - Extracts and groups all capabilities into an itemized inventory table.
   - Generates `MONTAGE_ASSESSMENT.md` and `montage.json`.
   - Optionally scaffolds Agent Plugin, specialized Skills, and MCP server files if `scaffold: true`.

2. `evaluate_readiness`:
   - Fast, in-memory evaluation returning scores, tier classification, strengths, and gap analysis without disk side-effects.

3. `scaffold_tooling`:
   - Generates the complete agent ecosystem directory structure (`plugins/`, `skills/`, `scripts/`, and MCP server starter) for an evaluated target.

---

## Step-by-Step Workflow

### Phase 1: Ingestion & Analysis
- Run Montage against a target:
  ```bash
  montage review /path/to/project --out ./out
  # Or via GitHub:
  montage review https://github.com/owner/repo --out ./out
  ```
- Montage walks the codebase (safely skipping virtual environments like `.venv` and `site-packages`) and discovers all endpoints, command groups, scripts, authored prompts, and engine packages.

### Phase 2: Reviewing the Capabilities Inventory & Scorecard
Inspect the generated `MONTAGE_ASSESSMENT.md`:
- **Readiness Scorecard:** Objective scoring across Interface Automation, Data Interchange, Safety, Auth, and Packaging.
- **Capabilities Assessment Table:** Every capability is tagged with its source, nature (Read-Only vs. Mutating), recommended target artifact, domain grouping, and architectural rationale.
- **Specialized Skill Blueprints:** Multi-skill blueprints mapping to distinct operational domains.
- **Clustered MCP Toolsets:** MCP tools grouped into logical personas.

### Phase 3: Scaffolding Modular Artifacts
Run `montage review <target> --scaffold --out ./out` to generate:
- `plugins/<target>/plugin.json`: Plugin manifest.
- `plugins/<target>/mcp.json`: MCP stdio server definition.
- `skills/<target>-<domain>/SKILL.md`: Multiple specialized skills tailored to each functional domain.
- `scripts/<domain>_helper.py`: Local helper scripts for Skills with Code.
- `cmd/<target>-mcp/main.go` (Go) or `mcp_server.py` (Python): Clustered MCP server starter.

---

## CLI Reference

```bash
# Full review with modular scaffolding
montage review ~/projects/syntaxis --scaffold --out ./out/syntaxis

# Review GitHub repository directly
montage review ghchinoy/repotographer --out ./out/repotographer

# Generate client configurations
montage mcp config --client opencode
```
