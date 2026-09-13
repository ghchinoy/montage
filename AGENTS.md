<!-- headroom:rtk-instructions -->
# RTK (Rust Token Killer) - Token-Optimized Commands

When running shell commands, **always prefix with `rtk`**. This reduces context
usage by 60-90% with zero behavior change. If rtk has no filter for a command,
it passes through unchanged — so it is always safe to use.

## Key Commands
```bash
# Git (59-80% savings)
rtk git status          rtk git diff            rtk git log

# Files & Search (60-75% savings)
rtk ls <path>           rtk read <file>         rtk grep <pattern>
rtk find <pattern>      rtk diff <file>

# Test (90-99% savings) — shows failures only
rtk pytest tests/       rtk cargo test          rtk test <cmd>

# Build & Lint (80-90% savings) — shows errors only
rtk tsc                 rtk lint                rtk cargo build
rtk prettier --check    rtk mypy                rtk ruff check

# Analysis (70-90% savings)
rtk err <cmd>           rtk log <file>          rtk json <file>
rtk summary <cmd>       rtk deps                rtk env

# GitHub (26-87% savings)
rtk gh pr view <n>      rtk gh run list         rtk gh issue list

# Infrastructure (85% savings)
rtk docker ps           rtk kubectl get         rtk docker logs <c>

# Package managers (70-90% savings)
rtk pip list            rtk pnpm install        rtk npm run <script>
```

## Rules
- In command chains, prefix each segment: `rtk git add . && rtk git commit -m "msg"`
- For debugging, use raw command without rtk prefix
- `rtk proxy <cmd>` runs command without filtering but tracks usage
<!-- /headroom:rtk-instructions -->

# Montage Developer & Agent Guide

Montage is the **Agent Ecosystem Readiness Assessor & Tooling Scaffolder** (CLI, Lit WebComponent Dashboard, and MCP Server).

## 🧭 Subsystem Layout

| Subsystem | Location | Stack | Primary Commands |
|---|---|---|---|
| **Core CLI & Engine** | Root (`main.go`, `internal/`) | Go 1.25+ | `make build`, `make test`, `make vet` |
| **Web Dashboard** | `ui/` | Lit, Vite, TypeScript | `make build-ui`, `cd ui && npm run dev` |
| **Documentation Site** | `docs-site/` | Astro 5, Starlight, Catppuccin | `make docs-build`, `make docs-dev` |
| **Embedded Binary** | Root | Go + embedded Lit dist | `make build-embedded` |

---

## 🛠️ Build & Verification Workflows

When working in this repository:
- **Go Core:** Always run `rtk go test ./...` and `rtk go vet ./...` before committing changes.
- **Embedded Builds:** Use `make build-embedded` to compile the Go binary with the embedded WebComponent UI (`-tags embedui`).
- **Docs Site:**
  - Build locally with `rtk make docs-build`.
  - Content collections are defined in `docs-site/src/content.config.ts` (Astro 5 Content Layer).
  - **YAML Frontmatter Rule:** Always double-quote `title:` and `description:` in markdown files if they contain colons (e.g. `description: "The 4-tier spectrum: Pure Skills, MCP Servers..."`).
  - **Theming:** The docs site uses Catppuccin Latte (`light: { flavor: 'latte', accent: 'lavender' }`) and Mocha for dark mode via `@catppuccin/starlight`.

---

## 🤖 MCP Server Testing Guidelines

Montage exposes an MCP server over `stdio` via `montage mcp`. When testing or extending MCP tools:
1. **Discovery:**
   ```bash
   rtk mcptools tools ./bin/montage mcp
   ```
2. **Tool Invocation:**
   ```bash
   rtk mcptools call assess_repository ./bin/montage mcp -p '{"target": "ghchinoy/repotographer", "scaffold": false}'
   ```
3. **Pipe & EOF Handling:** Ensure `server.Run(ctx, &mcp.StdioTransport{})` treats `io.EOF` as a clean exit when the client closes stdin.

---

## 🧩 Agent Ecosystem Artifact Spectrum

When modifying the evaluator (`internal/evaluator/`) or scaffolder (`internal/scaffold/`), preserve the 4-tier artifact spectrum:
1. **Pure Skills:** Zero runtime dependencies (`SKILL.md` with guidelines & heuristics).
2. **Skills with Code:** Bundled deterministic scripts to pre-filter data (60–90% token savings).
3. **Clustered MCP Servers:** Tools organized by operational personas rather than flat lists.
4. **Agent Plugins:** Standard-compliant packages (`plugin.json`, `mcp.json`, `skills/`).
