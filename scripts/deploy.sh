#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Montage Cloud Run Deployment Script
# ==============================================================================

SERVICE_NAME="${SERVICE_NAME:-montage}"
REGION="${REGION:-${GOOGLE_CLOUD_LOCATION:-us-central1}}"

# 1. Resolve GCP Project ID
PROJECT_ID="${1:-${GOOGLE_CLOUD_PROJECT:-${GCP_PROJECT:-}}}"
if [ -z "$PROJECT_ID" ]; then
  PROJECT_ID=$(gcloud config get-value project 2>/dev/null || true)
fi

if [ -z "$PROJECT_ID" ] || [ "$PROJECT_ID" = "(unset)" ]; then
  echo "❌ Error: Google Cloud Project ID is not set."
  echo "Usage: ./scripts/deploy.sh [PROJECT_ID] or set GOOGLE_CLOUD_PROJECT=your-project-id"
  exit 1
fi

echo "===================================================================="
echo "🎞️  Deploying Montage to Google Cloud Run"
echo "===================================================================="
echo "• Project ID:    $PROJECT_ID"
echo "• Service Name:  $SERVICE_NAME"
echo "• Region:        $REGION"
echo "===================================================================="

# 2. Check Prerequisites
if ! command -v gcloud &> /dev/null; then
  echo "❌ Error: Google Cloud SDK ('gcloud') is required but not found in PATH."
  echo "Please install it from https://cloud.google.com/sdk/docs/install"
  exit 1
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$( dirname "$SCRIPT_DIR" )"
cd "$ROOT_DIR"

# 3. Pre-flight Verification (Frontend build & Go test suite)
echo "📦 Running pre-flight checks..."

echo "   1/2 Building Lit frontend..."
(cd ui && npm run build --silent)

echo "   2/2 Running backend test suite..."
go test ./... -v=false

echo "✓ Pre-flight checks passed."

# 4. Deploy to Cloud Run
echo "🚀 Submitting build and deploying to Cloud Run..."
ENV_VARS="GOOGLE_GENAI_USE_VERTEXAI=true,GOOGLE_CLOUD_PROJECT=$PROJECT_ID,GOOGLE_CLOUD_LOCATION=$REGION"
if [ -n "${GITHUB_TOKEN:-}" ]; then
  ENV_VARS="$ENV_VARS,GITHUB_TOKEN=$GITHUB_TOKEN"
elif [ -n "${GH_TOKEN:-}" ]; then
  ENV_VARS="$ENV_VARS,GITHUB_TOKEN=$GH_TOKEN"
fi

gcloud run deploy "$SERVICE_NAME" \
  --source . \
  --project "$PROJECT_ID" \
  --region "$REGION" \
  --allow-unauthenticated \
  --set-env-vars "$ENV_VARS" \
  --memory 1Gi \
  --cpu 1 \
  --format="value(status.url)"

echo "===================================================================="
echo "🎉 Montage is successfully deployed and running on Cloud Run!"
echo "===================================================================="
