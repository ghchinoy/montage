package evaluator

import (
	"context"
	"testing"

	"github.com/ghchinoy/montage/internal/model"
)

func TestEvaluateScoring(t *testing.T) {
	target := model.TargetInfo{
		Name:        "repotographer",
		PrimaryLang: "Go",
	}

	findings := model.AnalysisFindings{
		HasCLI:             true,
		CLIFramework:       "cobra",
		HasJSONOutput:      true,
		JSONFlags:          []string{"--json"},
		StructuredFormats:  []string{"json"},
		CleanOutputStreams: true,
		HasSingleBinary:    true,
		BuildSystem:        "go",
		ExistingMCP:        true,
		EnvVarsDetected:    []string{"GEMINI_API_KEY", "GOOGLE_CLOUD_PROJECT"},
		ReadOnlyCommands:   []string{"map", "suggest", "render"},
	}

	rep := Evaluate(context.Background(), target, findings, false)

	if rep.OverallScore < 85 {
		t.Errorf("expected high agent-readiness score, got %d", rep.OverallScore)
	}

	if rep.Tier != model.Tier5AgentNative {
		t.Errorf("expected Tier 5: Agent-Native, got %s", rep.Tier)
	}

	if len(rep.Dimensions) != 5 {
		t.Errorf("expected 5 dimensions, got %d", len(rep.Dimensions))
	}
}

func TestEvaluateManualApp(t *testing.T) {
	target := model.TargetInfo{
		Name: "legacy-desktop-app",
	}

	findings := model.AnalysisFindings{
		HasCLI:           false,
		HasInteractiveUI: true,
		HasJSONOutput:    false,
	}

	rep := Evaluate(context.Background(), target, findings, false)

	if rep.OverallScore > 50 {
		t.Errorf("expected low readiness score for manual app, got %d", rep.OverallScore)
	}
}
