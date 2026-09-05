package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghchinoy/montage/internal/assessor"
	"github.com/ghchinoy/montage/internal/report"
	"github.com/ghchinoy/montage/internal/scaffold"
)

type assessRequest struct {
	Target   string `json:"target"`
	UseLLM   *bool  `json:"use_llm,omitempty"`
	Scaffold *bool  `json:"scaffold,omitempty"`
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "montage",
		"version": assessor.Version,
	})
}

func handleAssess(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req assessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON request: %v", err), http.StatusBadRequest)
		return
	}

	target := strings.TrimSpace(req.Target)
	if target == "" {
		http.Error(w, "target is required (GitHub repo or local path)", http.StatusBadRequest)
		return
	}

	useLLM := true
	if req.UseLLM != nil {
		useLLM = *req.UseLLM
	}

	rep, err := assessor.RunAssessment(r.Context(), target, "./sources", useLLM)
	if err != nil {
		http.Error(w, fmt.Sprintf("Assessment failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Always generate artifacts into cache dir for download
	safeName := sanitizeTargetName(rep.Target.Name)
	cacheDir := filepath.Join("out", safeName)
	_ = os.MkdirAll(cacheDir, 0755)

	// Save reports to cache
	mdPath := filepath.Join(cacheDir, "MONTAGE_ASSESSMENT.md")
	_ = os.WriteFile(mdPath, []byte(report.RenderMarkdown(rep)), 0644)
	if jsonData, err := report.RenderJSON(rep); err == nil {
		_ = os.WriteFile(filepath.Join(cacheDir, "montage.json"), jsonData, 0644)
	}

	// Scaffold artifacts
	_, _ = scaffold.GenerateArtifacts(rep, cacheDir)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rep)
}

func handleExportZip(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	target := strings.TrimSpace(r.URL.Query().Get("target"))
	if target == "" {
		http.Error(w, "target query parameter is required", http.StatusBadRequest)
		return
	}

	// Check if already in cache
	safeName := sanitizeTargetName(target)
	cacheDir := filepath.Join("out", safeName)

	if info, err := os.Stat(cacheDir); err != nil || !info.IsDir() {
		// Run assessment and scaffolding first
		rep, err := assessor.RunAssessment(r.Context(), target, "./sources", false)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate artifacts: %v", err), http.StatusInternalServerError)
			return
		}
		safeName = sanitizeTargetName(rep.Target.Name)
		cacheDir = filepath.Join("out", safeName)
		_ = os.MkdirAll(cacheDir, 0755)
		_ = os.WriteFile(filepath.Join(cacheDir, "MONTAGE_ASSESSMENT.md"), []byte(report.RenderMarkdown(rep)), 0644)
		if data, err := report.RenderJSON(rep); err == nil {
			_ = os.WriteFile(filepath.Join(cacheDir, "montage.json"), data, 0644)
		}
		_, _ = scaffold.GenerateArtifacts(rep, cacheDir)
	}

	zipFilename := fmt.Sprintf("%s-agent-bundle.zip", safeName)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFilename))

	if err := ZipDirectory(cacheDir, w); err != nil {
		http.Error(w, fmt.Sprintf("Failed to stream zip archive: %v", err), http.StatusInternalServerError)
		return
	}
}

func sanitizeTargetName(input string) string {
	clean := strings.TrimSpace(input)
	clean = strings.TrimPrefix(clean, "https://github.com/")
	clean = strings.TrimPrefix(clean, "github.com/")
	clean = strings.ReplaceAll(clean, "/", "_")
	clean = strings.ReplaceAll(clean, "\\", "_")
	clean = strings.ReplaceAll(clean, ":", "_")
	clean = strings.ToLower(clean)
	if clean == "" {
		clean = "app"
	}
	return clean
}
