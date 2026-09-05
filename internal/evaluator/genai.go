package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ghchinoy/montage/internal/model"
	"google.golang.org/genai"
)

// EnrichWithGemini optionally calls Gemini to enrich the report with domain synthesis.
func EnrichWithGemini(ctx context.Context, report *model.AssessmentReport) {
	apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	project := strings.TrimSpace(os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if project == "" {
		project = strings.TrimSpace(os.Getenv("PROJECT_ID"))
	}
	if project == "" {
		project = strings.TrimSpace(os.Getenv("GCP_PROJECT"))
	}

	if apiKey == "" && project == "" {
		return // offline mode; retain deterministic evaluation
	}

	var client *genai.Client
	var err error

	if project != "" {
		loc := strings.TrimSpace(os.Getenv("GOOGLE_CLOUD_LOCATION"))
		if loc == "" {
			loc = "global"
		}
		client, err = genai.NewClient(ctx, &genai.ClientConfig{
			Project:  project,
			Location: loc,
			Backend:  genai.BackendVertexAI,
		})
	} else {
		client, err = genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		})
	}

	if err != nil {
		return
	}

	prompt := fmt.Sprintf(`You are an expert AI Agent Architect.
Analyze this application and provide concise strategic recommendations for turning it into an Agent Ecosystem tool (MCP Server, Agent Plugin, Pure Skill, or Skill with Code).

Application Name: %s
Description: %s
Primary Language: %s
CLI Detected: %v (Framework: %s)
Current Readiness Score: %d/100 (%s)
Subcommands: %s

Return a JSON object matching this schema:
{
  "strategic_advice": "A 2-3 sentence strategic roadmap for agent integration",
  "recommended_mcp_tools": [
    {
      "name": "tool_name",
      "description": "tool purpose",
      "is_read_only": true
    }
  ]
}
`, report.Target.Name, report.Target.Description, report.Target.PrimaryLang,
		report.Findings.HasCLI, report.Findings.CLIFramework, report.OverallScore, report.Tier,
		strings.Join(report.Findings.ReadOnlyCommands, ", "))

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	})
	if err != nil {
		return
	}

	partText := strings.TrimSpace(resp.Text())
	if partText == "" {
		return
	}

	type llmEnrichment struct {
		StrategicAdvice     string `json:"strategic_advice"`
		RecommendedMCPTools []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			IsReadOnly  bool   `json:"is_read_only"`
		} `json:"recommended_mcp_tools"`
	}

	var parsed llmEnrichment
	if err := json.Unmarshal([]byte(partText), &parsed); err == nil {
		if parsed.StrategicAdvice != "" {
			report.TierSummary = report.TierSummary + "\n\n**AI Agent Strategic Guidance:** " + parsed.StrategicAdvice
		}
		if len(parsed.RecommendedMCPTools) > 0 {
			var newTools []model.MCPToolSpec
			for _, t := range parsed.RecommendedMCPTools {
				newTools = append(newTools, model.MCPToolSpec{
					Name:        t.Name,
					Description: t.Description,
					IsReadOnly:  t.IsReadOnly,
				})
			}
			report.MCPServer.ProposedTools = newTools
		}
	}
}
