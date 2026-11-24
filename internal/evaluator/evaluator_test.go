package evaluator

import (
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

	rep := Evaluate(t.Context(), target, findings, false)

	if rep.OverallScore < 85 {
		t.Errorf("Evaluate(%q) OverallScore = %d, want >= 85", target.Name, rep.OverallScore)
	}

	if rep.Tier != model.Tier5AgentNative {
		t.Errorf("Evaluate(%q) Tier = %q, want %q", target.Name, rep.Tier, model.Tier5AgentNative)
	}

	if len(rep.Dimensions) != 5 {
		t.Errorf("Evaluate(%q) len(Dimensions) = %d, want 5", target.Name, len(rep.Dimensions))
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

	rep := Evaluate(t.Context(), target, findings, false)

	if rep.OverallScore > 50 {
		t.Errorf("Evaluate(%q) OverallScore = %d, want <= 50", target.Name, rep.OverallScore)
	}
}
