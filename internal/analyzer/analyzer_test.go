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
		t.Errorf("expected HasCLI=true")
	}
	if findings.CLIFramework != "cobra" {
		t.Errorf("expected CLIFramework=cobra, got %s", findings.CLIFramework)
	}
	if findings.BuildSystem != "go" {
		t.Errorf("expected BuildSystem=go, got %s", findings.BuildSystem)
	}
	if len(findings.EnvVarsDetected) == 0 || findings.EnvVarsDetected[0] != "SAMPLE_API_KEY" {
		t.Errorf("expected SAMPLE_API_KEY detected in env vars, got %v", findings.EnvVarsDetected)
	}
	if target.PrimaryLang != "Go" {
		t.Errorf("expected PrimaryLang=Go, got %s", target.PrimaryLang)
	}
}
