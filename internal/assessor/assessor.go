package assessor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghchinoy/montage/internal/analyzer"
	"github.com/ghchinoy/montage/internal/evaluator"
	"github.com/ghchinoy/montage/internal/github"
	"github.com/ghchinoy/montage/internal/model"
)

// Version of Montage.
const Version = "0.2.0"

// RunAssessment conducts the full target ingestion, code analysis, and readiness evaluation.
func RunAssessment(ctx context.Context, targetInput, sourcesDir string, useLLM bool) (*model.AssessmentReport, error) {
	var target *model.TargetInfo
	var localScanPath string

	if sourcesDir == "" {
		sourcesDir = "./sources"
	}

	if github.IsGitHubTarget(targetInput) {
		owner, repo, err := github.ParseRepoRef(targetInput)
		if err != nil {
			return nil, err
		}

		fmt.Printf("🔍 Fetching GitHub metadata for %s/%s...\n", owner, repo)
		meta, err := github.FetchRepoMetadata(owner, repo)
		if err != nil {
			return nil, err
		}
		target = meta

		fmt.Printf("📦 Ingesting repository into %s/%s/%s...\n", sourcesDir, owner, repo)
		clonedPath, err := github.CloneOrPull(owner, repo, sourcesDir)
		if err != nil {
			return nil, err
		}
		localScanPath = clonedPath
		target.LocalPath = clonedPath
	} else {
		absPath, err := filepath.Abs(targetInput)
		if err != nil {
			return nil, fmt.Errorf("invalid local directory: %w", err)
		}
		info, err := os.Stat(absPath)
		if err != nil || !info.IsDir() {
			return nil, fmt.Errorf("path does not exist or is not a directory: %s", absPath)
		}
		localScanPath = absPath
		target = &model.TargetInfo{
			LocalPath:  absPath,
			SourceKind: model.SourceKindLocal,
		}
	}

	fmt.Printf("🔬 Analyzing codebase in %s...\n", localScanPath)
	target, findings, err := analyzer.AnalyzeCodebase(localScanPath, target)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	fmt.Println("⚖️  Evaluating Agent Ecosystem readiness...")
	rep := evaluator.Evaluate(ctx, *target, *findings, useLLM)
	return rep, nil
}
