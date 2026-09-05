package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghchinoy/montage/internal/assessor"
	"github.com/ghchinoy/montage/internal/model"
	"github.com/ghchinoy/montage/internal/report"
	"github.com/ghchinoy/montage/internal/scaffold"
	"github.com/spf13/cobra"
)

var Version = assessor.Version

// RootCmd is the primary CLI command for Montage.
var RootCmd = &cobra.Command{
	Use:   "montage",
	Short: "Review applications and determine Agent Ecosystem tooling (Plugins, Skills, MCP Servers, CLIs)",
	Long: `Montage assesses codebases and repositories for autonomous agent readiness.
It evaluates interface automation, structured data interchange, execution safety,
headless auth, and packaging to produce actionable readiness reports and
scaffold Agent Plugins, Skills, and Model Context Protocol (MCP) servers.`,
}

func init() {
	RootCmd.AddCommand(reviewCmd)
	RootCmd.AddCommand(scaffoldCmd)
	RootCmd.AddCommand(versionCmd)

	// Flags for review
	reviewCmd.Flags().StringP("out", "o", "./out", "Output directory for assessment reports and artifacts")
	reviewCmd.Flags().String("format", "md,json", "Comma-separated output report formats: md,json")
	reviewCmd.Flags().Bool("scaffold", false, "Automatically generate Agent Plugin, Skill, and MCP scaffold")
	reviewCmd.Flags().Bool("llm", true, "Use Gemini / Vertex AI to enrich assessment if ambient credentials exist")
	reviewCmd.Flags().String("sources", "./sources", "Local cache directory for cloned GitHub repositories")

	// Flags for scaffold
	scaffoldCmd.Flags().StringP("out", "o", "./out", "Output directory for generated agent artifacts")
	scaffoldCmd.Flags().String("sources", "./sources", "Local cache directory for cloned GitHub repositories")
}

var reviewCmd = &cobra.Command{
	Use:   "review <path-or-github-repo>",
	Short: "Review an application or GitHub repo and generate an Agent Readiness report",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetInput := args[0]
		outDir, _ := cmd.Flags().GetString("out")
		formatsStr, _ := cmd.Flags().GetString("format")
		doScaffold, _ := cmd.Flags().GetBool("scaffold")
		useLLM, _ := cmd.Flags().GetBool("llm")
		sourcesDir, _ := cmd.Flags().GetString("sources")

		rep, err := RunAssessment(cmd.Context(), targetInput, sourcesDir, useLLM)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		formats := strings.Split(formatsStr, ",")
		formatSet := make(map[string]bool)
		for _, f := range formats {
			formatSet[strings.TrimSpace(strings.ToLower(f))] = true
		}

		if formatSet["md"] || formatSet["markdown"] {
			mdContent := report.RenderMarkdown(rep)
			mdPath := filepath.Join(outDir, "MONTAGE_ASSESSMENT.md")
			if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
				return fmt.Errorf("failed to write Markdown report: %w", err)
			}
			fmt.Printf("📄 Assessment Report: %s\n", mdPath)
		}

		if formatSet["json"] {
			jsonData, err := report.RenderJSON(rep)
			if err != nil {
				return err
			}
			jsonPath := filepath.Join(outDir, "montage.json")
			if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
				return fmt.Errorf("failed to write JSON report: %w", err)
			}
			fmt.Printf("📊 Assessment Data:   %s\n", jsonPath)
		}

		// Print summary to terminal
		fmt.Printf("\n🏆 Readiness Score: %d/100 (%s)\n", rep.OverallScore, rep.Tier)
		for _, dim := range rep.Dimensions {
			fmt.Printf("   • %-22s %3d/100 (Weight: %.0f%%)\n", dim.Name, dim.Score, dim.Weight*100)
		}

		if doScaffold {
			fmt.Println("\n🏗️  Scaffolding Agent Ecosystem tooling...")
			files, err := scaffold.GenerateArtifacts(rep, outDir)
			if err != nil {
				return fmt.Errorf("scaffolding failed: %w", err)
			}
			for _, f := range files {
				fmt.Printf("   ✓ Created: %s\n", f)
			}
		}

		fmt.Printf("\n🎉 Review complete! Results saved in %s\n", outDir)
		return nil
	},
}

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold <path-or-github-repo>",
	Short: "Generate Agent Plugin, Skill, and MCP Server artifacts for an application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetInput := args[0]
		outDir, _ := cmd.Flags().GetString("out")
		sourcesDir, _ := cmd.Flags().GetString("sources")

		rep, err := RunAssessment(cmd.Context(), targetInput, sourcesDir, false)
		if err != nil {
			return err
		}

		fmt.Printf("🏗️  Scaffolding agent artifacts for '%s'...\n", rep.Target.Name)
		files, err := scaffold.GenerateArtifacts(rep, outDir)
		if err != nil {
			return err
		}
		for _, f := range files {
			fmt.Printf("   ✓ %s\n", f)
		}
		fmt.Printf("\n🎉 Scaffolding completed in %s\n", outDir)
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Montage version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("montage v%s\n", Version)
	},
}

// RunAssessment conducts the full target ingestion, code analysis, and readiness evaluation.
func RunAssessment(ctx context.Context, targetInput, sourcesDir string, useLLM bool) (*model.AssessmentReport, error) {
	return assessor.RunAssessment(ctx, targetInput, sourcesDir, useLLM)
}
