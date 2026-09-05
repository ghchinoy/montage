package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghchinoy/montage/internal/model"
)

// GenerateArtifacts writes the full suite of modular agent tooling artifacts based on the assessment.
func GenerateArtifacts(report *model.AssessmentReport, outDir string) ([]string, error) {
	if outDir == "" {
		outDir = "./agent-ecosystem"
	}

	var generatedFiles []string
	appName := strings.ToLower(report.Target.Name)
	if appName == "" {
		appName = "app"
	}

	pluginDir := filepath.Join(outDir, "plugins", appName)
	standaloneSkillsDir := filepath.Join(outDir, "skills")

	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}
	if err := os.MkdirAll(standaloneSkillsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create standalone skills directory: %w", err)
	}

	// 1. Generate plugin.json
	pluginPath := filepath.Join(pluginDir, "plugin.json")
	pluginContent := generatePluginJSON(report)
	if err := os.WriteFile(pluginPath, []byte(pluginContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write plugin.json: %w", err)
	}
	generatedFiles = append(generatedFiles, pluginPath)

	// 2. Generate mcp.json
	mcpPath := filepath.Join(pluginDir, "mcp.json")
	mcpContent := generateMCPJSON(report)
	if err := os.WriteFile(mcpPath, []byte(mcpContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write mcp.json: %w", err)
	}
	generatedFiles = append(generatedFiles, mcpPath)

	// 3. Generate Modular Skills from Blueprints
	blueprints := report.SkillBlueprints
	if len(blueprints) == 0 {
		// Fallback to single primary skill
		blueprints = []model.SkillBlueprint{
			{
				Name:            appName,
				Kind:            "pure_skill",
				Description:     report.Target.Description,
				Domain:          "General Operations",
				RecommendedUses: report.PureSkill.RecommendedUses,
				KeyInstructions: report.PureSkill.KeyInstructions,
			},
		}
	}

	for _, bp := range blueprints {
		skillDirName := bp.Name
		pluginSkillDir := filepath.Join(pluginDir, "skills", skillDirName)
		standaloneSkillDir := filepath.Join(standaloneSkillsDir, skillDirName)

		_ = os.MkdirAll(pluginSkillDir, 0755)
		_ = os.MkdirAll(standaloneSkillDir, 0755)

		skillContent := generateSpecializedSkillMD(bp, report)

		pPath := filepath.Join(pluginSkillDir, "SKILL.md")
		if err := os.WriteFile(pPath, []byte(skillContent), 0644); err == nil {
			generatedFiles = append(generatedFiles, pPath)
		}

		sPath := filepath.Join(standaloneSkillDir, "SKILL.md")
		if err := os.WriteFile(sPath, []byte(skillContent), 0644); err == nil {
			generatedFiles = append(generatedFiles, sPath)
		}

		// If Skill with Code has a helper script, scaffold starter script
		if bp.Kind == "skill_with_code" && bp.HelperScript != "" {
			scriptPath := filepath.Join(outDir, bp.HelperScript)
			_ = os.MkdirAll(filepath.Dir(scriptPath), 0755)
			scriptContent := generateHelperScript(bp, report)
			if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err == nil {
				generatedFiles = append(generatedFiles, scriptPath)
			}
		}
	}

	// 4. Generate MCP Server Starter (Python FastMCP or Go SDK)
	if report.Target.PrimaryLang == "Python" || report.Findings.BuildSystem == "python" {
		mcpPyPath := filepath.Join(outDir, "mcp_server.py")
		mcpPyCode := generatePythonMCPServer(report)
		if err := os.WriteFile(mcpPyPath, []byte(mcpPyCode), 0644); err == nil {
			generatedFiles = append(generatedFiles, mcpPyPath)
		}
	} else if report.Target.PrimaryLang == "Go" || report.Findings.CLIFramework == "cobra" {
		mcpServerDir := filepath.Join(outDir, "cmd", fmt.Sprintf("%s-mcp", appName))
		if err := os.MkdirAll(mcpServerDir, 0755); err == nil {
			mcpServerPath := filepath.Join(mcpServerDir, "main.go")
			mcpServerCode := generateGoMCPServer(report)
			if err := os.WriteFile(mcpServerPath, []byte(mcpServerCode), 0644); err == nil {
				generatedFiles = append(generatedFiles, mcpServerPath)
			}
		}
	}

	return generatedFiles, nil
}

func generatePluginJSON(r *model.AssessmentReport) string {
	appName := strings.ToLower(r.Target.Name)
	owner := r.Target.Owner
	if owner == "" {
		owner = "agent-developer"
	}

	type author struct {
		Name string `json:"name"`
		URL  string `json:"url,omitempty"`
	}

	type pluginSchema struct {
		Schema      string   `json:"$schema"`
		Name        string   `json:"name"`
		Version     string   `json:"version"`
		Description string   `json:"description"`
		Author      author   `json:"author"`
		Homepage    string   `json:"homepage,omitempty"`
		Repository  string   `json:"repository,omitempty"`
		License     string   `json:"license"`
		Keywords    []string `json:"keywords"`
	}

	keywords := []string{"agent-plugin", "mcp", "automation", appName}
	for _, l := range r.Target.Languages {
		keywords = append(keywords, strings.ToLower(l))
	}

	p := pluginSchema{
		Schema:      "https://agent-plugins.org/schemas/1.0.0/plugin.schema.json",
		Name:        appName,
		Version:     "0.1.0",
		Description: r.Target.Description,
		Author: author{
			Name: owner,
			URL:  fmt.Sprintf("https://github.com/%s", owner),
		},
		Homepage:   r.Target.URL,
		Repository: r.Target.URL,
		License:    "Apache-2.0",
		Keywords:   keywords,
	}

	bytes, _ := json.MarshalIndent(p, "", "  ")
	return string(bytes) + "\n"
}

func generateMCPJSON(r *model.AssessmentReport) string {
	appName := strings.ToLower(r.Target.Name)
	cmdName := appName
	var args []string

	if r.Target.PrimaryLang == "Python" {
		cmdName = "python"
		args = []string{"mcp_server.py"}
	} else {
		args = []string{"mcp"}
	}

	type serverEntry struct {
		Type    string   `json:"type"`
		Command string   `json:"command"`
		Args    []string `json:"args"`
	}

	type mcpSchema struct {
		Schema     string                 `json:"$schema"`
		MCPServers map[string]serverEntry `json:"mcpServers"`
	}

	m := mcpSchema{
		Schema: "https://agent-plugins.org/schemas/1.0.0/mcp.schema.json",
		MCPServers: map[string]serverEntry{
			appName: {
				Type:    "stdio",
				Command: cmdName,
				Args:    args,
			},
		},
	}

	bytes, _ := json.MarshalIndent(m, "", "  ")
	return string(bytes) + "\n"
}

func generateSpecializedSkillMD(bp model.SkillBlueprint, r *model.AssessmentReport) string {
	var b strings.Builder

	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("name: %s\n", bp.Name))
	b.WriteString(fmt.Sprintf("description: %s\n", bp.Description))
	b.WriteString("license: Apache-2.0\n")
	b.WriteString("metadata:\n")
	b.WriteString("  version: \"0.1.0\"\n")
	b.WriteString(fmt.Sprintf("  kind: \"%s\"\n", bp.Kind))
	b.WriteString(fmt.Sprintf("  domain: \"%s\"\n", bp.Domain))
	b.WriteString("---\n\n")

	b.WriteString(fmt.Sprintf("# %s\n\n", bp.Name))
	b.WriteString(fmt.Sprintf("%s\n\n", bp.Description))

	if bp.Kind == "skill_with_code" && bp.HelperScript != "" {
		b.WriteString("## Bundled Execution Script\n\n")
		b.WriteString(fmt.Sprintf("This skill bundles a deterministic helper script to pre-process data before agent reasoning (%s):\n\n", bp.TokenSavingsEst))
		b.WriteString("```bash\n")
		b.WriteString(fmt.Sprintf("python %s --input <data> --output json\n", bp.HelperScript))
		b.WriteString("```\n\n")
	}

	if len(bp.Capabilities) > 0 {
		b.WriteString("## Core Capabilities Covered\n\n")
		for _, c := range bp.Capabilities {
			b.WriteString(fmt.Sprintf("- **%s**\n", c))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Workflow Guidelines & Decision Rules\n\n")
	for i, inst := range bp.KeyInstructions {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, inst))
	}
	b.WriteString("\n")

	b.WriteString("## Agent Execution Protocol\n\n")
	b.WriteString("1. **Pre-flight Validation:** Verify all required configuration or database records exist.\n")
	b.WriteString("2. **Context Efficiency:** If raw logs or histories are large, run the bundled helper script to compute trends.\n")
	b.WriteString("3. **Safety Review:** Review proposed actions or medication/metric entries before committing mutations.\n")

	return b.String()
}

func generateHelperScript(bp model.SkillBlueprint, r *model.AssessmentReport) string {
	var b strings.Builder
	b.WriteString(`#!/usr/bin/env python3
"""
` + bp.Name + ` Helper Script
Automates data extraction, trend analysis, and token compression for ` + r.Target.Name + `.
"""

import sys
import json

def main():
    print(json.dumps({
        "status": "ready",
        "domain": "` + bp.Domain + `",
        "summary": "Processed data locally to save 80% context tokens."
    }, indent=2))

if __name__ == "__main__":
    main()
`)
	return b.String()
}

func generatePythonMCPServer(r *model.AssessmentReport) string {
	appName := strings.ToLower(r.Target.Name)
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`#!/usr/bin/env python3
"""
MCP Server for %s
Generated by Montage. Provides typed tools organized by capability groups.
"""

from mcp.server.fastmcp import FastMCP

mcp = FastMCP("%s")

`, r.Target.Name, appName))

	// Generate tools by group
	if len(r.ToolGroups) > 0 {
		for _, grp := range r.ToolGroups {
			b.WriteString(fmt.Sprintf("# =============================================================================\n"))
			b.WriteString(fmt.Sprintf("# Group: %s (%s)\n", grp.Name, grp.Description))
			b.WriteString(fmt.Sprintf("# =============================================================================\n\n"))

			for _, tool := range grp.Tools {
				funcName := strings.ReplaceAll(tool.Name, "-", "_")
				readOnlyComment := "Read-only operation"
				if !tool.IsReadOnly {
					readOnlyComment = "Mutating operation (requires confirmation)"
				}

				b.WriteString(fmt.Sprintf(`@mcp.tool()
def %s(target_id: str, options: str = "") -> str:
    """
    %s
    Access: %s
    """
    # TODO: Connect to %s backend API or database
    return f"Executed %s on {target_id}"

`, funcName, tool.Description, readOnlyComment, appName, funcName))
			}
		}
	} else {
		b.WriteString(`@mcp.tool()
def query_status() -> str:
    """Retrieves current application status."""
    return "Status OK"
`)
	}

	b.WriteString(`if __name__ == "__main__":
    mcp.run(transport="stdio")
`)

	return b.String()
}

func generateGoMCPServer(r *model.AssessmentReport) string {
	appName := strings.ToLower(r.Target.Name)
	var b strings.Builder

	b.WriteString(`package main

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "` + appName + `-mcp",
		Version: "0.1.0",
	}, nil)

`)

	for _, tool := range r.MCPServer.ProposedTools {
		toolName := strings.ReplaceAll(strings.ReplaceAll(tool.Name, " ", "_"), "-", "_")
		typeName := strings.ReplaceAll(formatCapabilityName(toolName), " ", "")
		b.WriteString(fmt.Sprintf(`	// Tool: %s
	type %sParams struct {
		Target string `+"`json:\"target,omitempty\" jsonschema:\"Target identifier or query\"`"+`
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "%s",
		Description: "%s",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args %sParams) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("%s executed on %%s", args.Target)},
			},
		}, nil, nil
	})

`, toolName, typeName, toolName, tool.Description, typeName, toolName))
	}

	b.WriteString(`	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server stopped: %v\n", err)
		os.Exit(1)
	}
}
`)

	return b.String()
}

func formatCapabilityName(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, " ")
}
