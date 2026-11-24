package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeCodebase(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "montage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create sample go.mod
	goMod := `module example.com/sample
go 1.25
`
	_ = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0644)

	// Create sample main.go with Cobra and Env vars
	mainGo := `package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "sample",
	Short: "A sample CLI tool",
}

var listCmd = &cobra.Command{
	Use: "list",
	Short: "List items",
}

func main() {
	apiKey := os.Getenv("SAMPLE_API_KEY")
	_ = apiKey
	rootCmd.AddCommand(listCmd)
	_ = rootCmd.Execute()
}
`
	_ = os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(mainGo), 0644)

	target, findings, err := AnalyzeCodebase(tempDir, nil)
	if err != nil {
		t.Fatalf("AnalyzeCodebase failed: %v", err)
	}

	if !findings.HasCLI {
		t.Errorf("AnalyzeCodebase(%q) HasCLI = false, want true", tempDir)
	}
	if findings.CLIFramework != "cobra" {
		t.Errorf("AnalyzeCodebase(%q) CLIFramework = %q, want %q", tempDir, findings.CLIFramework, "cobra")
	}
	if findings.BuildSystem != "go" {
		t.Errorf("AnalyzeCodebase(%q) BuildSystem = %q, want %q", tempDir, findings.BuildSystem, "go")
	}
	if len(findings.EnvVarsDetected) == 0 || findings.EnvVarsDetected[0] != "SAMPLE_API_KEY" {
		t.Errorf("AnalyzeCodebase(%q) EnvVarsDetected = %v, want SAMPLE_API_KEY first", tempDir, findings.EnvVarsDetected)
	}
	if target.PrimaryLang != "Go" {
		t.Errorf("AnalyzeCodebase(%q) PrimaryLang = %q, want %q", tempDir, target.PrimaryLang, "Go")
	}
}
