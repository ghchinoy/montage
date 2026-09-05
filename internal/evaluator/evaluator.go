package evaluator

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ghchinoy/montage/internal/model"
)

// Evaluate produces a full AssessmentReport from the analyzed target and findings.
func Evaluate(ctx context.Context, target model.TargetInfo, findings model.AnalysisFindings, useLLM bool) *model.AssessmentReport {
	d1 := evaluateInterfaceAutomation(findings)
	d2 := evaluateDataInterchange(findings)
	d3 := evaluateExecutionSafety(findings)
	d4 := evaluateAuthConfig(findings)
	d5 := evaluatePackagingRuntime(findings)

	dimensions := []model.ReadinessDimension{d1, d2, d3, d4, d5}

	composite := d1.Weight*float64(d1.Score) +
		d2.Weight*float64(d2.Score) +
		d3.Weight*float64(d3.Score) +
		d4.Weight*float64(d4.Score) +
		d5.Weight*float64(d5.Score)

	overallScore := int(math.Round(composite))
	if overallScore > 100 {
		overallScore = 100
	} else if overallScore < 0 {
		overallScore = 0
	}

	tier, summary := determineTier(overallScore, findings)
	pureSkill := buildPureSkillRec(target, findings, overallScore)
	skillWithCode := buildSkillWithCodeRec(target, findings, overallScore)
	mcpServer := buildMCPServerRec(target, findings, overallScore)
	agentPlugin := buildAgentPluginRec(target, findings, overallScore)
	gaps := identifyGaps(findings)

	capabilities, skillBlueprints, toolGroups := ExtractAndGroupCapabilities(target, findings)

	// If tool groups were extracted, refine proposed MCP tools
	if len(toolGroups) > 0 {
		var allTools []model.MCPToolSpec
		for _, g := range toolGroups {
			allTools = append(allTools, g.Tools...)
		}
		if len(allTools) > 0 {
			mcpServer.ProposedTools = allTools
		}
	}

	report := &model.AssessmentReport{
		GeneratedAt:     time.Now().UTC(),
		Target:          target,
		Findings:        findings,
		Capabilities:    capabilities,
		SkillBlueprints: skillBlueprints,
		ToolGroups:      toolGroups,
		OverallScore:    overallScore,
		Tier:            tier,
		TierSummary:     summary,
		Dimensions:      dimensions,
		PureSkill:       pureSkill,
		SkillWithCode:   skillWithCode,
		MCPServer:       mcpServer,
		AgentPlugin:     agentPlugin,
		ActionableGaps:  gaps,
	}

	// Optionally refine with Gemini if enabled and available
	if useLLM {
		EnrichWithGemini(ctx, report)
	}

	return report
}

func evaluateInterfaceAutomation(f model.AnalysisFindings) model.ReadinessDimension {
	score := 20
	var strengths, gaps []string

	if f.HasAPI {
		score += 35
		strengths = append(strengths, fmt.Sprintf("REST/HTTP API detected (%s framework) with %d endpoints", f.APIFramework, len(f.APIEndpoints)))
	} else if f.HasCLI {
		score += 35
		strengths = append(strengths, fmt.Sprintf("CLI entry point detected (%s framework)", f.CLIFramework))
	} else if len(f.Scripts) > 0 {
		score += 25
		strengths = append(strengths, fmt.Sprintf("%d standalone utility scripts detected", len(f.Scripts)))
	} else {
		gaps = append(gaps, "No dedicated CLI interface or REST API detected; interactions currently rely on internal library or GUI")
	}

	if len(f.Commands) > 0 {
		score += 25
		strengths = append(strengths, fmt.Sprintf("%d distinct subcommands identified for modular tool exposure", len(f.Commands)))
	} else if len(f.APIEndpoints) > 0 {
		score += 25
		strengths = append(strengths, fmt.Sprintf("%d route endpoints ready for direct MCP tool mapping", len(f.APIEndpoints)))
	}

	if f.HasInteractiveUI {
		score -= 20
		gaps = append(gaps, "Interactive prompts detected (e.g. TTY scan/readline) which block headless agent execution without non-interactive bypass")
	} else {
		score += 20
		strengths = append(strengths, "Zero interactive terminal blockers detected; fully headless-compatible")
	}

	if f.ExistingMCP {
		score += 15
		strengths = append(strengths, "Existing Model Context Protocol (MCP) server integration detected")
	}

	return clampDimension("Interface Automation", score, 0.25,
		"Evaluates whether the application can be operated headlessly by LLM agents via CLI or API.",
		strengths, gaps)
}

func evaluateDataInterchange(f model.AnalysisFindings) model.ReadinessDimension {
	score := 25
	var strengths, gaps []string

	if f.HasJSONOutput || contains(f.StructuredFormats, "json") {
		score += 40
		strengths = append(strengths, "Structured JSON output serialization implemented")
	} else {
		gaps = append(gaps, "Application outputs unstructured text; agents will need regex/scraping to parse results")
	}

	if len(f.JSONFlags) > 0 {
		score += 20
		strengths = append(strengths, fmt.Sprintf("Explicit CLI flag (%s) available to request structured machine-readable payload", strings.Join(f.JSONFlags, ", ")))
	} else if f.HasCLI {
		gaps = append(gaps, "CLI lacks an explicit `--json` or `--format` flag for deterministic machine consumption")
	}

	if f.CleanOutputStreams {
		score += 15
		strengths = append(strengths, "Separate output streams: logs and diagnostics routed to stderr, data payload to stdout")
	} else {
		gaps = append(gaps, "Logs and data may be interleaved on stdout, requiring agent-side filtering")
	}

	return clampDimension("Data Interchange", score, 0.20,
		"Evaluates machine readability, structured serialization (JSON/YAML), and stream discipline.",
		strengths, gaps)
}

func evaluateExecutionSafety(f model.AnalysisFindings) model.ReadinessDimension {
	score := 50
	var strengths, gaps []string

	if len(f.MutatingCommands) == 0 && !f.HasDangerousActions {
		score += 50
		strengths = append(strengths, "All operations are read-only / analytical; zero risk of unintended mutations")
	} else {
		if f.HasDryRunFlag {
			score += 30
			strengths = append(strengths, "Dry-run / validation flag detected for safe pre-flight simulations")
		} else {
			score -= 20
			gaps = append(gaps, "Mutating operations present without a `--dry-run` safety flag")
		}

		if len(f.ReadOnlyCommands) > 0 {
			score += 20
			strengths = append(strengths, fmt.Sprintf("Clear division: %d read-only commands vs %d mutating commands", len(f.ReadOnlyCommands), len(f.MutatingCommands)))
		}

		if f.HasDangerousActions {
			score -= 15
			gaps = append(gaps, "Potentially destructive actions (delete/drop/kill) detected; requires strict confirmation gating")
		}
	}

	return clampDimension("Execution Safety", score, 0.20,
		"Evaluates risk profile, mutation boundary, and presence of simulation flags (dry-run).",
		strengths, gaps)
}

func evaluateAuthConfig(f model.AnalysisFindings) model.ReadinessDimension {
	score := 40
	var strengths, gaps []string

	if len(f.EnvVarsDetected) > 0 {
		score += 35
		strengths = append(strengths, fmt.Sprintf("Headless configuration via environment variables (%d detected: %s)",
			len(f.EnvVarsDetected), strings.Join(f.EnvVarsDetected[:min(len(f.EnvVarsDetected), 4)], ", ")))
	}

	if len(f.ConfigFiles) > 0 {
		score += 15
		strengths = append(strengths, fmt.Sprintf("Standard configuration file detected (%s)", strings.Join(f.ConfigFiles[:min(len(f.ConfigFiles), 3)], ", ")))
	}

	if contains(f.AuthMethods, "ambient_cli") {
		score += 15
		strengths = append(strengths, "Leverages ambient CLI authentication (e.g. `gh auth` or cloud SDK), eliminating PAT sprawl")
	}

	if f.HasInteractiveAuth {
		score -= 30
		gaps = append(gaps, "Interactive browser or TTY login flow detected; requires pre-authenticated sessions for agents")
	}

	return clampDimension("Auth & Configuration", score, 0.15,
		"Evaluates whether authentication and credentials can be supplied headlessly without GUI prompts.",
		strengths, gaps)
}

func evaluatePackagingRuntime(f model.AnalysisFindings) model.ReadinessDimension {
	score := 35
	var strengths, gaps []string

	if f.HasSingleBinary {
		score += 45
		strengths = append(strengths, "Compiles to a self-contained single binary (Go/Rust), simplifying container/sandbox distribution")
	} else if f.PackagingType == "script" {
		score += 20
		strengths = append(strengths, fmt.Sprintf("Scriptable packaging with standard runtime (%s)", f.BuildSystem))
	}

	if f.BuildSystem != "" {
		score += 20
		strengths = append(strengths, fmt.Sprintf("Standardized build system present (%s)", f.BuildSystem))
	} else {
		gaps = append(gaps, "No standardized build system (Makefile, go.mod, package.json) detected")
	}

	if f.HasContainer {
		score += 15
		strengths = append(strengths, "Containerfile/Dockerfile available for isolated execution")
	}

	return clampDimension("Runtime & Packaging", score, 0.20,
		"Evaluates ease of execution in automated agent sandboxes and containerized runtimes.",
		strengths, gaps)
}

func clampDimension(name string, score int, weight float64, explanation string, strengths, gaps []string) model.ReadinessDimension {
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}
	return model.ReadinessDimension{
		Name:        name,
		Score:       score,
		Weight:      weight,
		Explanation: explanation,
		Strengths:   strengths,
		Gaps:        gaps,
	}
}

func determineTier(score int, f model.AnalysisFindings) (model.ReadinessTier, string) {
	if f.ExistingMCP && score >= 85 {
		return model.Tier5AgentNative, "Fully agent-native codebase with built-in Model Context Protocol server, clean CLI, and headless execution."
	}
	if score >= 90 {
		return model.Tier5AgentNative, "Ready for immediate Agent Plugin and MCP wrapping. Excellent modularity, structured outputs, and safety controls."
	}
	if score >= 75 {
		return model.Tier4AgentCapable, "Substantially ready for agent integration. Commands can be wrapped in an MCP server or Skill with minimal glue."
	}
	if score >= 55 {
		return model.Tier3Scriptable, "Scriptable interface available. Recommended for a Skill with Code or helper utility scripts before formal MCP exposure."
	}
	if score >= 35 {
		return model.Tier2WorkflowOnly, "Primarily procedural or library code. Recommended for a Pure Skill providing workflow guidance and manual execution checklists."
	}
	return model.Tier1ManualGUI, "GUI-bound or manual application. Requires CLI/headless decoupling before agent tool execution is feasible."
}

func buildPureSkillRec(target model.TargetInfo, f model.AnalysisFindings, score int) model.PureSkillRecommendation {
	appName := target.Name
	if appName == "" {
		appName = "application"
	}

	uses := []string{
		fmt.Sprintf("Guiding coding agents on architectural conventions and best practices for %s", appName),
		"Providing procedural workflows and diagnostic checklists without requiring runtime execution",
		"Reviewing repository health and enforcing project standards",
	}

	instructions := []string{
		fmt.Sprintf("Analyze target project using %s design patterns.", appName),
		"Verify prerequisites and environment setup before suggesting actions.",
		"Format all modifications as reviewable diffs.",
	}

	return model.PureSkillRecommendation{
		Feasibility:     "High",
		Score:           95,
		SuggestedName:   fmt.Sprintf("%s-guide", strings.ToLower(appName)),
		Description:     fmt.Sprintf("Workflow and best-practice guidance for working with %s without external tool dependencies.", appName),
		RecommendedUses: uses,
		KeyInstructions: instructions,
	}
}

func buildSkillWithCodeRec(target model.TargetInfo, f model.AnalysisFindings, score int) model.SkillWithCodeRecommendation {
	appName := strings.ToLower(target.Name)
	feasibility := "Medium"
	subScore := 65
	if f.HasCLI || len(f.Commands) > 0 {
		feasibility = "High"
		subScore = 90
	}

	scripts := []string{}
	if f.HasCLI {
		scripts = append(scripts, fmt.Sprintf("bin/%s or bundled helper CLI", appName))
	} else {
		scripts = append(scripts, fmt.Sprintf("scripts/audit_%s.py (lightweight extractor)", appName))
	}

	return model.SkillWithCodeRecommendation{
		Feasibility:     feasibility,
		Score:           subScore,
		SuggestedName:   fmt.Sprintf("%s-operator", appName),
		Description:     fmt.Sprintf("Skill paired with lightweight deterministic execution scripts to pre-process data before agent reasoning for %s.", appName),
		HelperScripts:   scripts,
		TokenSavingsEst: "60-85% token reduction by offloading raw extraction to script",
		RecommendedFlow: []string{
			"Agent executes bundled script via bash to extract filtered JSON",
			"Script filters out noise and returns compact payload",
			"Agent performs high-level synthesis and reasoning on compact data",
		},
	}
}

func buildMCPServerRec(target model.TargetInfo, f model.AnalysisFindings, score int) model.MCPServerRecommendation {
	appName := strings.ToLower(target.Name)
	feasibility := "Medium"
	subScore := score
	if score >= 75 {
		feasibility = "High"
	} else if score < 50 {
		feasibility = "Low"
	}

	var tools []model.MCPToolSpec
	if len(f.Commands) > 0 {
		for _, cmd := range f.Commands {
			tools = append(tools, model.MCPToolSpec{
				Name:        fmt.Sprintf("%s_%s", appName, strings.ReplaceAll(cmd.Name, "-", "_")),
				Description: fmt.Sprintf("Executes %s %s (%s)", appName, cmd.Name, cmd.Description),
				IsReadOnly:  !cmd.IsMutating,
				SourceCmd:   cmd.Name,
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"options": map[string]interface{}{
							"type":        "string",
							"description": "CLI parameters and target spec",
						},
					},
				},
			})
		}
	} else {
		tools = append(tools, model.MCPToolSpec{
			Name:        fmt.Sprintf("run_%s", appName),
			Description: fmt.Sprintf("Invokes %s operations headlessly", appName),
			IsReadOnly:  true,
		})
	}

	return model.MCPServerRecommendation{
		Feasibility:   feasibility,
		Score:         subScore,
		ServerName:    appName,
		Transport:     "stdio",
		ProposedTools: tools,
		ProposedResources: []string{
			fmt.Sprintf("%s://config/current", appName),
			fmt.Sprintf("%s://status", appName),
		},
		ProposedPrompts: []string{
			fmt.Sprintf("curate_%s_workflow", appName),
		},
	}
}

func buildAgentPluginRec(target model.TargetInfo, f model.AnalysisFindings, score int) model.AgentPluginRecommendation {
	appName := strings.ToLower(target.Name)
	feasibility := "Medium"
	subScore := score
	if score >= 75 {
		feasibility = "High"
	}

	return model.AgentPluginRecommendation{
		Feasibility: feasibility,
		Score:       subScore,
		PluginName:  appName,
		Description: fmt.Sprintf("Comprehensive Agent Plugin bundling the %s MCP server and workflow skills.", appName),
		Components: []string{
			"plugin.json (Agent Plugin manifest adhering to 1.0.0 schema)",
			"mcp.json (MCP stdio server registration)",
			fmt.Sprintf("skills/%s/SKILL.md (Autonomous agent execution instructions)", appName),
		},
	}
}

func identifyGaps(f model.AnalysisFindings) []model.GapItem {
	var gaps []model.GapItem

	if !f.HasJSONOutput {
		gaps = append(gaps, model.GapItem{
			Priority:    "High",
			Dimension:   "Data Interchange",
			Title:       "Add `--json` Structured Output Mode",
			Description: "The application currently returns unstructured or human-formatted text, requiring complex LLM parsing.",
			Remediation: "Implement a `--json` or `--format=json` flag that outputs pure unformatted JSON to stdout.",
			ExampleCode: `// Example Cobra flag integration
cmd.Flags().Bool("json", false, "Output results as machine-readable JSON")
if jsonOut {
    json.NewEncoder(os.Stdout).Encode(result)
    return nil
}`,
		})
	}

	if f.HasInteractiveUI {
		gaps = append(gaps, model.GapItem{
			Priority:    "High",
			Dimension:   "Interface Automation",
			Title:       "Support `--non-interactive` / Bypass TTY Prompts",
			Description: "Interactive terminal prompts (e.g. readline, input()) halt autonomous LLM agent execution.",
			Remediation: "Add a `--yes` / `--non-interactive` flag or detect non-TTY stdout/stdin (`term.IsTerminal`).",
			ExampleCode: `if !isatty.IsTerminal(os.Stdin.Fd()) || nonInteractiveFlag {
    // Proceed with default answers or fail deterministically
}`,
		})
	}

	if len(f.MutatingCommands) > 0 && !f.HasDryRunFlag {
		gaps = append(gaps, model.GapItem{
			Priority:    "Medium",
			Dimension:   "Execution Safety",
			Title:       "Add `--dry-run` Pre-flight Simulation",
			Description: "Mutating commands execute directly without an agent-safe validation mode.",
			Remediation: "Add a `--dry-run` flag that validates schemas and returns proposed changes without applying writes.",
			ExampleCode: `cmd.Flags().Bool("dry-run", false, "Simulate operation without writing changes")`,
		})
	}

	if !f.CleanOutputStreams {
		gaps = append(gaps, model.GapItem{
			Priority:    "Low",
			Dimension:   "Data Interchange",
			Title:       "Separate Diagnostic Logs from Data Payload",
			Description: "Logging statements written to stdout corrupt JSON parsers when captured by agents.",
			Remediation: "Direct all informational logs, progress bars, and warnings to `os.Stderr` or `stderr`.",
		})
	}

	return gaps
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, val) {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
