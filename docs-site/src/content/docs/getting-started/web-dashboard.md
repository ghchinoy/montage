---
title: Web Dashboard & Cloud Run
description: Operating the interactive Lit WebComponent dashboard locally and deploying to Google Cloud Run.
---

Montage includes a built-in Lit WebComponent dashboard providing an interactive interface for assessing repositories, filtering capability matrices, and downloading turn-key `.zip` agent scaffolding bundles.

## Running Locally

Start the web dashboard on port `8080`:

```bash
montage serve --port 8080
```

Visit `http://localhost:8080` in your browser.

### Development Mode with Live HMR

For frontend development on the Lit WebComponent UI:

```bash
# Terminal 1: Backend API
montage serve --port 8080

# Terminal 2: Vite Dev Server
cd ui && npm run dev
# Open http://localhost:5173
```

## Deploying to Google Cloud Run

Montage includes a multi-stage `Dockerfile` optimized for Cloud Run:

```bash
gcloud run deploy montage \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars GOOGLE_GENAI_USE_VERTEXAI=true,GOOGLE_CLOUD_PROJECT=your-project-id \
  --memory 1Gi \
  --cpu 1
```

Public GitHub repositories work out-of-the-box on Cloud Run with zero authentication via GitHub's public REST API and anonymous shallow git clones.
