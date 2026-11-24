# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.1] - 2026-11-23

### Security

- Upgraded Go module floor from `go 1.25.5` to `go 1.26.10` and refreshed Go dependencies (`golang.org/x/net`, `google.golang.org/grpc`, `golang.org/x/text`, `golang.org/x/crypto`), eliminating all 10 reachable `govulncheck` symbol vulnerabilities.
- Hardened `internal/github` (`ParseRepoRef`, `FetchRepoMetadata`, `CloneOrPull`) against CLI argument injection and path traversal on MCP client inputs, added `--` end-of-options delimiters, and redacted `GITHUB_TOKEN` from `git clone` stderr errors.
- Upgraded `docs-site` `astro` to `^7.4.13` to resolve `GHSA-r3h5-4585-7x3w` and `GHSA-mh32-jxcr-mh9f`.
- Pinned Node.js runtime floors to `22.23.4` in `Dockerfile`, `.github/workflows/deploy-docs.yml`, `ui/package.json`, and `docs-site/package.json`, aligned `Dockerfile` Go builder to `golang:1.26-alpine`, and added `apk upgrade --no-cache` to the Alpine runtime stage.

### Changed

- Upgraded direct Go dependencies (`github.com/modelcontextprotocol/go-sdk` `v1.8.0`, `github.com/spf13/cobra` `v1.10.2`, `google.golang.org/genai` `v1.79.1`) and safe Node patches (`ui` `vite` `^6.4.6`, `docs-site` `@astrojs/starlight` `^0.42.11`).
- Synchronized version markers to `0.2.1` across `internal/assessor`, `internal/github`, `ui`, `docs-site`, and `plugins/montage`.
- Added `.serena/`, `.beads/`, and `.dolt/` to `.gitignore` and added `.gcloudignore`.
- Applied Go readability improvements (`gofmt -s`, replacing legacy `interface{}` with `any`, `var scripts []string`, `%q` formatting, `t.Context()`, and descriptive test failure messages).

## [0.2.0] - 2026-10-21

### Added

- Astro Starlight documentation site (`docs-site/`) with Catppuccin Latte/Mocha theme and GitHub Pages deployment workflow.
- Initial Montage CLI, Lit WebComponent dashboard, and MCP stdio server (`assess_repository`, `evaluate_readiness`, `scaffold_tooling`).
