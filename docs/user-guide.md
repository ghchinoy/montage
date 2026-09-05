# Montage User & Operator Guide

`montage` is a command-line tool and Model Context Protocol (MCP) server that reviews existing applications and determines what **Agent Ecosystem tooling** can be built from them—including **Pure Skills**, **Skills with Code**, **MCP Servers**, and **Agent Plugins**.

---

## 1. Prerequisites & Installation

### A. Prerequisites
1. **Git** (required for shallow cloning repositories)
2. **Go 1.25+** (to compile the binary from source)
3. **GitHub Access (Zero-Config for Public Repositories):**
   - **Public Repositories:** Work out-of-the-box with **zero authentication** via the GitHub Public REST API and anonymous `git clone`. No `gh auth login` required, making it 100% headless and Cloud Run ready.
   - **Optional `GITHUB_TOKEN`:** Set `GITHUB_TOKEN=ghp_...` (or `GH_TOKEN`) to increase rate limits to 5,000 req/hour or access private repositories.
   - **Local GitHub CLI (`gh`):** If authenticated locally via `gh auth login`, Montage automatically uses your ambient credentials as a fallback.
4. **Gemini API Key or Google Cloud Vertex AI [Optional for AI strategic enrichment]**
   - **Gemini Developer API:**
     ```bash
     export GEMINI_API_KEY="your-gemini-api-key"
     ```
   - **Google Cloud Vertex AI:**
     ```bash
     export GOOGLE_CLOUD_PROJECT="your-project-id"
     export GOOGLE_CLOUD_LOCATION="global"
     ```

### B. Building Montage
```bash
git clone https://github.com/ghchinoy/montage.git
cd montage
make build
# Binary is compiled to ./bin/montage
```

Verify the installation:
```bash
./bin/montage version
```

---

## 2. Reviewing Codebases

Montage accepts either a **local filesystem path** or a **GitHub repository reference**.

### A. Review a Local Directory
```bash
montage review /path/to/project --out ./out
```

### B. Review a GitHub Repository Directly
Montage automatically uses `gh` to fetch repository metadata, clone it into `./sources/<owner>/<repo>`, and perform structural analysis:
```bash
# Using full URL:
montage review https://github.com/ghchinoy/repotographer --out ./out

# Using owner/repo shorthand:
montage review ghchinoy/repotographer --out ./out
```

### C. Review and Scaffold Agent Tooling in One Step
Adding the `--scaffold` flag instructs Montage to write out turn-key Agent Plugins, specialized skills, and MCP server starters:
```bash
montage review ~/projects/syntaxis --scaffold --out ./out/syntaxis
```

### CLI Command Options

| Flag | Short | Default | Description |
|---|---|---|---|
| `--out` | `-o` | `./out` | Destination directory for reports and scaffolded tooling |
| `--format` | | `md,json` | Comma-separated output formats: `md` (`MONTAGE_ASSESSMENT.md`), `json` (`montage.json`) |
| `--scaffold` | | `false` | Automatically generates `plugins/`, `skills/`, `scripts/`, and MCP server code |
| `--llm` | | `true` | Enriches assessment with Gemini / Vertex AI strategic guidance if credentials exist |
| `--sources` | | `./sources` | Local cache directory for cloned GitHub repositories |

---

## 3. Deciphering the Assessment Report (`MONTAGE_ASSESSMENT.md`)

Every Montage review generates a comprehensive Markdown report structured into five key sections:

### 1. Executive Summary & Readiness Scorecard
Displays your application's composite score (0–100) and maturity tier:
* **Tier 5: Agent-Native (90–100):** Ready for immediate MCP wrapping and Agent Plugin distribution.
* **Tier 4: Agent-Capable (75–89):** Well-structured interfaces; minor gaps in `--dry-run` or non-interactive flags.
* **Tier 3: Scriptable (55–74):** Functional CLI/scripts; best paired with a *Skill with Code*.
* **Tier 2: Workflow-Only (35–54):** Procedural or library code; recommended for a *Pure Skill*.
* **Tier 1: Manual / GUI-Bound (<35):** Requires refactoring to expose headless interfaces.

The scorecard breaks down the five core evaluation dimensions:
* **Interface Automation (25%):** CLI frameworks, REST routes, subcommands.
* **Data Interchange (20%):** Structured JSON serialization, stdout/stderr purity.
* **Execution Safety (20%):** Mutation boundaries, pre-flight dry-run simulation.
* **Auth & Configuration (15%):** Headless environment variables without TTY blockers.
* **Runtime & Packaging (20%):** Single binaries, containerization, build manifests.

### 2. Capabilities Assessment & Agent Ecosystem Grouping Table
Deconstructs your application into fine-grained functional capabilities and maps each to its ideal agent target:

```markdown
| Capability | Source / Ref | Nature | Target Artifact | Grouping / Domain | Architectural Rationale |
|---|---|---|---|---|---|
| **Analyze Weights** | `POST /analyze_weights` | ⚡ Mutating | **⚡ Skill with Code** | `Clinical Telemetry & Growth` | Performs mathematical analysis over historical records; ideal for a helper script to save LLM context tokens. |
| **Review** | `command: review` | 🔒 Read-Only | **📝 Pure Skill** | `Advisory Review & Feedback` | Advisory review mode preserves original text and surfaces structured improvement recommendations; ideal as a Pure Skill. |
| **Generate Adoption Bio** | `POST /generate_adoption_bio` | ⚡ Mutating | **🔌 MCP Tool** | `Adoption & Copywriting` | Content synthesis service; best paired with an editorial Pure Skill + MCP generation tool. |
```

### 3. Specialized Skill Blueprints
Instead of a single monolithic guide, Montage defines distinct skills tailored to specific operational domains. For each skill, the report outlines:
* Domain classification and operational rules.
* Encompassed capabilities.
* Bundled execution scripts and estimated context token savings (typically 70–90%).
* Pre-flight and safety execution protocols.

### 4. Clustered MCP Toolsets
Groups proposed MCP tools into cohesive operational personas (e.g. *Advisory Review*, *Source Analysis*, *Clinical Telemetry*, *Typesetting*) to avoid bloated, unprincipled tool lists.

### 5. Actionable Gaps & Modernization Roadmap
Highlights specific technical remediations with copy-pasteable code examples (e.g., adding `--json` flags, supporting `--non-interactive`, or introducing `--dry-run` pre-flight simulations).

---

## 4. Using the Scaffolded Artifacts

When `--scaffold` is enabled, Montage creates the following directory structure:

```
out/<app>/
├── MONTAGE_ASSESSMENT.md
├── montage.json
├── cmd/<app>-mcp/main.go        # (If Go target) Go stdio MCP server starter
├── mcp_server.py                # (If Python target) Python FastMCP server starter
├── scripts/
│   └── <domain>_helper.py       # Deterministic data-crunching helper scripts
├── plugins/<app>/
│   ├── plugin.json              # Agent Plugins Specification manifest
│   ├── mcp.json                 # MCP stdio server registration
│   └── skills/                  # Plugin-bundled skills
└── skills/
    ├── <app>-<domain-1>/SKILL.md
    └── <app>-<domain-2>/SKILL.md
```

### Deploying Skills
1. **Local Agent Usage (OpenCode / Claude Code / Antigravity):**
   Copy the generated skills into your project's `.agents/skills` directory:
   ```bash
   cp -r out/<app>/skills/* ~/.agents/skills/
   # Or locally:
   cp -r out/<app>/skills/* .agents/skills/
   ```
2. **Running the MCP Server Starter:**
   - **Python Projects (`mcp_server.py`):**
     ```bash
     pip install "mcp[cli]"
     python out/<app>/mcp_server.py
     ```
   - **Go Projects (`cmd/<app>-mcp/`):**
     ```bash
     go run out/<app>/cmd/<app>-mcp/main.go
     ```

---

## 5. Connecting Montage as an MCP Server

Montage includes a built-in Model Context Protocol server over `stdio`. This allows coding assistants (OpenCode, Claude Desktop, Cursor, Antigravity) to assess codebases dynamically.

### A. Start the Server
```bash
montage mcp
```

### B. Generate Client Configurations
Use `montage mcp config` to print ready-to-paste configuration blocks:

```bash
# For OpenCode:
montage mcp config --client opencode

# For Claude Desktop:
montage mcp config --client claude

# For Cursor:
montage mcp config --client cursor
```

Example configuration for `opencode.json`:
```json
{
  "mcpServers": {
    "montage": {
      "command": "/path/to/montage/bin/montage",
      "args": ["mcp"]
    }
  }
}
```

### C. Available MCP Tools

| Tool | Parameters | Description |
|---|---|---|
| `assess_repository` | `target` (string, required)<br>`scaffold` (bool)<br>`use_llm` (bool)<br>`out_dir` (string) | Ingests a local path or GitHub repo, scores readiness, produces reports, and optionally scaffolds artifacts. |
| `evaluate_readiness` | `target` (string, required) | Fast in-memory evaluation returning 5 dimension scores, tier, and gap analysis without disk side effects. |
| `scaffold_tooling` | `target` (string, required)<br>`out_dir` (string) | Scaffolds the complete `plugins/`, `skills/`, and MCP server starter structure for an assessed target. |

---

## 6. Real-World Case Studies

### Case Study 1: `syntaxis` (AI-Assisted Publication Engine)
* **Target:** Multi-agent Go publication pipeline compiling Typst into PDFs.
* **Readiness Score:** **`93 / 100` (Tier 5: Agent-Native)**
* **Discovered Capabilities:**
  - `pkg/prompts/reviewer.md` + `review` command -> **`syntaxis-advisory-reviewer`** (Pure Skill for academic peer review without text rewrites).
  - `pkg/prompts/*analyst.md` + `pkg/classifier` -> **`syntaxis-source-analyst`** (Skill with Code for token budget management).
  - `pkg/generators/` -> **`syntaxis-multi-modal-visualizer`** (Skill with Code + MCP tools for DOT flowcharts and Gonum plots).
  - `cmd/template.go` + `pkg/registry` -> **`syntaxis-typst-typesetter`** (Pure Skill for Typst syntax + template search tool).
  - `cmd/project.go` + `cmd/execute.go` -> **`syntaxis-publication-orchestrator`** (MCP execution DAG).
  - `pkg/a2api/` -> **`syntaxis-a2a-bridge`** (A2A Protocol v1.0 delegation).

### Case Study 2: `kitten-haven` (Foster Care Application)
* **Target:** FastAPI backend with Firestore database and Lit frontend.
* **Readiness Score:** **`85 / 100` (Tier 4: Agent-Capable)**
* **Discovered Capabilities:**
  - `POST /analyze_weights` -> **`kitten-haven-growth-triage`** (Skill with Code bundling a 30-line Python calculator that compresses 280 raw records into a 3-line summary, saving 99% context tokens).
  - `POST /generate_adoption_bio` -> **`kitten-haven-adoption-copywriter`** (Pure Skill for tone/persona + MCP live API tool).
  - `backend/scripts/set_role.py` -> **`kitten-haven-shelter-administration-security`** (Administrative CLI tool).

### Case Study 3: `repotographer` (Repository Cartographer)
* **Target:** Go Cobra CLI and MCP server mapping GitHub accounts into concept graphs.
* **Readiness Score:** **`80 / 100` (Tier 4: Agent-Capable)**
* **Discovered Capabilities:**
  - Definition B Tri-State Connectivity -> **`repotographer-cartographer`** (Pure Skill for graph architecture).
  - Gemini domain clustering -> **`repotographer-taxonomy-curator`** (Skill with Code + helper script).
  - Graphviz DOT + Cytoscape -> **`repotographer-multi-modal-visualizer`** (MCP rendering tools).

---

## 7. Web Application & Cloud Run Deployment

Montage features a built-in Lit WebComponent web interface where users can type any GitHub URL, click **"Make this part of the Agent Economy"**, inspect interactive scorecards and capability matrices, and download full agent scaffolding zip archives.

### A. Run Locally
```bash
# Start the web service on port 8080
montage serve --port 8080
```
Open `http://localhost:8080` in any browser.

### B. Development with Live HMR (Hot Module Replacement)
If you are developing components in `ui/`:
```bash
# Terminal 1: Run Go backend server
montage serve --port 8080

# Terminal 2: Run Vite development server with proxy
cd ui
npm run dev
# Vite runs on http://localhost:5173 with automatic API proxying to :8080
```

### C. Build Self-Contained Binary
```bash
make build-embedded
# Compiles ui/ into ui/dist and embeds it into bin/montage with -tags embedui
```

### D. Deploy to Google Cloud Run
Montage is packaged with a production-ready, multi-stage `Dockerfile`:
```bash
gcloud run deploy montage \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars GOOGLE_GENAI_USE_VERTEXAI=true,GOOGLE_CLOUD_PROJECT=your-project-id \
  --memory 1Gi \
  --cpu 1
```
Cloud Run automatically sets the `$PORT` environment variable, which Montage listens on by default.

