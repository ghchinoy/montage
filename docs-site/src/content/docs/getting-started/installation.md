---
title: Installation & Build
description: How to install, build, and verify Montage from source.
---

Montage is compiled as a single, self-contained Go binary with zero external runtime dependencies.

## Prerequisites

1. **Go 1.24+** (Recommended: Go 1.25+)
2. **Git** (Required for cloning targets)
3. **GitHub CLI (`gh`)** *(Optional)*: If authenticated locally, Montage automatically uses ambient credentials for private repository analysis.
4. **Gemini API Key or Google Cloud Vertex AI** *(Optional)*: For optional AI-enriched strategic synthesis.

## Building from Source

Clone the repository and compile using `make`:

```bash
git clone https://github.com/ghchinoy/montage.git
cd montage
make build
```

The compiled binary will be placed at `./bin/montage`.

### Verify Installation

```bash
./bin/montage version
# Output: montage v0.2.0
```

### Install to System PATH

To install Montage globally into your `$GOPATH/bin`:

```bash
make install
```

Ensure `$(go env GOPATH)/bin` is included in your `PATH`.
