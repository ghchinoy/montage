package scaffold

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ghchinoy/montage/internal/model"
)

func TestGenerateArtifacts(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "montage-scaffold-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	rep := &model.AssessmentReport{
		GeneratedAt: time.Now(),
		Target: model.TargetInfo{
			Name:        "testapp",
			PrimaryLang: "Go",
			Description: "A test application for agent scaffolding",
		},
		Tier:         model.Tier4AgentCapable,
		OverallScore: 85,
		Findings: model.AnalysisFindings{
			HasCLI:       true,
			CLIFramework: "cobra",
			Commands: []model.CommandSpec{
				{Name: "status", Description: "Show status", IsMutating: false},
				{Name: "deploy", Description: "Deploy resource", IsMutating: true},
			},
		},
		MCPServer: model.MCPServerRecommendation{
			ServerName: "testapp",
			Transport:  "stdio",
			ProposedTools: []model.MCPToolSpec{
				{Name: "testapp_status", Description: "Show status", IsReadOnly: true},
			},
		},
	}

	files, err := GenerateArtifacts(rep, tempDir)
	if err != nil {
		t.Fatalf("GenerateArtifacts failed: %v", err)
	}

	if len(files) < 4 {
		t.Errorf("expected at least 4 generated files, got %d", len(files))
	}

	// Verify plugin.json exists
	pluginFile := filepath.Join(tempDir, "plugins", "testapp", "plugin.json")
	if _, err := os.Stat(pluginFile); err != nil {
		t.Errorf("plugin.json was not created at %s", pluginFile)
	}

	// Verify SKILL.md exists
	skillFile := filepath.Join(tempDir, "plugins", "testapp", "skills", "testapp", "SKILL.md")
	if _, err := os.Stat(skillFile); err != nil {
		t.Errorf("SKILL.md was not created at %s", skillFile)
	}
}
