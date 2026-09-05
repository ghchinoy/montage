package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghchinoy/montage/internal/cli"
	"github.com/ghchinoy/montage/internal/model"
	"github.com/ghchinoy/montage/internal/report"
	"github.com/ghchinoy/montage/internal/scaffold"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

func init() {
	cli.RootCmd.AddCommand(mcpCmd)
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run Montage as a Model Context Protocol (MCP) stdio server",
	Long:  "Starts Montage as an MCP server over stdio, exposing repository assessment and scaffolding tools to AI agents.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunMCPServer(cmd.Context())
	},
}

// AssessParams defines input parameters for the assess_repository MCP tool.
type AssessParams struct {
	Target   string `json:"target" jsonschema:"Local path or GitHub repository (e.g. ghchinoy/repotographer or https://github.com/owner/repo)"`
	Scaffold bool   `json:"scaffold,omitempty" jsonschema:"Whether to automatically scaffold Agent Plugin, Skill, and MCP files"`
	UseLLM   bool   `json:"use_llm,omitempty" jsonschema:"Whether to use Gemini / Vertex AI to enrich the assessment (default true)"`
	OutDir   string `json:"out_dir,omitempty" jsonschema:"Output directory for report and scaffold (default ./out)"`
}

// AssessResult defines the output of the assess_repository MCP tool.
type AssessResult struct {
	Report         *model.AssessmentReport `json:"report"`
	MarkdownPath   string                  `json:"markdown_path,omitempty"`
	JSONPath       string                  `json:"json_path,omitempty"`
	ScaffoldFiles  []string                `json:"scaffold_files,omitempty"`
}

// EvaluateParams defines input parameters for quick evaluation without file I/O.
type EvaluateParams struct {
	Target string `json:"target" jsonschema:"Local path or GitHub repository"`
}

// ScaffoldParams defines input parameters for the scaffold_tooling MCP tool.
type ScaffoldParams struct {
	Target string `json:"target" jsonschema:"Local path or GitHub repository"`
	OutDir string `json:"out_dir,omitempty" jsonschema:"Output directory for generated agent artifacts"`
}

// RunMCPServer registers all tools and starts the Model Context Protocol stdio server.
func RunMCPServer(ctx context.Context) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "montage",
		Version: cli.Version,
	}, nil)

	// Tool 1: assess_repository
	mcp.AddTool(server, &mcp.Tool{
		Name:        "assess_repository",
		Description: "Reviews an application or GitHub repository to evaluate agent readiness and determine candidate Agent Ecosystem tooling (Plugins, Skills, MCP Servers, CLIs).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AssessParams) (*mcp.CallToolResult, any, error) {
		if strings.TrimSpace(args.Target) == "" {
			return nil, nil, fmt.Errorf("target parameter is required (local path or GitHub repo)")
		}

		outDir := args.OutDir
		if outDir == "" {
			outDir = "./out"
		}

		rep, err := cli.RunAssessment(ctx, args.Target, "./sources", args.UseLLM)
		if err != nil {
			return nil, nil, fmt.Errorf("assessment failed: %w", err)
		}

		_ = os.MkdirAll(outDir, 0755)
		mdPath := filepath.Join(outDir, "MONTAGE_ASSESSMENT.md")
		jsonPath := filepath.Join(outDir, "montage.json")

		_ = os.WriteFile(mdPath, []byte(report.RenderMarkdown(rep)), 0644)
		if data, err := report.RenderJSON(rep); err == nil {
			_ = os.WriteFile(jsonPath, data, 0644)
		}

		var scaffolded []string
		if args.Scaffold {
			scaffolded, _ = scaffold.GenerateArtifacts(rep, outDir)
		}

		res := &AssessResult{
			Report:        rep,
			MarkdownPath:  mdPath,
			JSONPath:      jsonPath,
			ScaffoldFiles: scaffolded,
		}

		summaryText := fmt.Sprintf("Assessed '%s' — Readiness Score: %d/100 (%s). Report saved to %s.",
			rep.Target.Name, rep.OverallScore, rep.Tier, mdPath)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: summaryText},
			},
		}, res, nil
	})

	// Tool 2: evaluate_readiness
	mcp.AddTool(server, &mcp.Tool{
		Name:        "evaluate_readiness",
		Description: "Performs a fast in-memory readiness assessment of a codebase, returning the 5 dimension scores, tier, and gap analysis.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args EvaluateParams) (*mcp.CallToolResult, any, error) {
		if strings.TrimSpace(args.Target) == "" {
			return nil, nil, fmt.Errorf("target parameter is required")
		}

		rep, err := cli.RunAssessment(ctx, args.Target, "./sources", false)
		if err != nil {
			return nil, nil, fmt.Errorf("evaluation failed: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Readiness: %d/100 (%s)\n%s", rep.OverallScore, rep.Tier, rep.TierSummary)},
			},
		}, rep, nil
	})

	// Tool 3: scaffold_tooling
	mcp.AddTool(server, &mcp.Tool{
		Name:        "scaffold_tooling",
		Description: "Generates Agent Plugin manifest (plugin.json), MCP config (mcp.json), workflow skill (SKILL.md), and MCP server boilerplate for an evaluated target.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args ScaffoldParams) (*mcp.CallToolResult, any, error) {
		if strings.TrimSpace(args.Target) == "" {
			return nil, nil, fmt.Errorf("target parameter is required")
		}

		outDir := args.OutDir
		if outDir == "" {
			outDir = "./agent-ecosystem"
		}

		rep, err := cli.RunAssessment(ctx, args.Target, "./sources", false)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to assess target before scaffolding: %w", err)
		}

		files, err := scaffold.GenerateArtifacts(rep, outDir)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate artifacts: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Successfully scaffolded %d agent artifacts in %s", len(files), outDir)},
			},
		}, files, nil
	})

	return server.Run(ctx, &mcp.StdioTransport{})
}
