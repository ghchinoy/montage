package report

import (
	"strings"
	"testing"
	"time"

	"github.com/ghchinoy/montage/internal/model"
)

func TestRenderMarkdownAndJSON(t *testing.T) {
	rep := &model.AssessmentReport{
		GeneratedAt:  time.Now(),
		Target:       model.TargetInfo{Name: "sample-app", URL: "https://github.com/org/sample-app"},
		OverallScore: 88,
		Tier:         model.Tier4AgentCapable,
		TierSummary:  "Solid CLI foundations.",
		Dimensions: []model.ReadinessDimension{
			{Name: "Interface Automation", Score: 90, Weight: 0.25},
			{Name: "Data Interchange", Score: 85, Weight: 0.20},
		},
		PureSkill: model.PureSkillRecommendation{
			SuggestedName: "sample-guide",
			Feasibility:   "High",
			Score:         90,
		},
		SkillWithCode: model.SkillWithCodeRecommendation{
			SuggestedName: "sample-operator",
			Feasibility:   "High",
			Score:         85,
		},
		MCPServer: model.MCPServerRecommendation{
			ServerName:  "sample-mcp",
			Feasibility: "High",
			Score:       88,
			Transport:   "stdio",
		},
		AgentPlugin: model.AgentPluginRecommendation{
			PluginName:  "sample-app",
			Feasibility: "High",
			Score:       88,
		},
	}

	md := RenderMarkdown(rep)
	if !strings.Contains(md, "sample-app") {
		t.Errorf("markdown output does not contain target name")
	}
	if !strings.Contains(md, "Tier 4: Agent-Capable") {
		t.Errorf("markdown output does not contain tier")
	}

	jsonBytes, err := RenderJSON(rep)
	if err != nil {
		t.Fatalf("RenderJSON failed: %v", err)
	}
	if len(jsonBytes) == 0 {
		t.Fatalf("RenderJSON returned empty byte slice")
	}
}
