package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handleHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var data map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &data); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if data["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", data["status"])
	}
	if data["service"] != "montage" {
		t.Errorf("expected service 'montage', got '%s'", data["service"])
	}
}

func TestSanitizeTargetName(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"https://github.com/ghchinoy/repotographer", "ghchinoy_repotographer"},
		{"ghchinoy/syntaxis", "ghchinoy_syntaxis"},
		{"sample-app", "sample-app"},
		{"", "app"},
	}

	for _, c := range cases {
		out := sanitizeTargetName(c.input)
		if out != c.expected {
			t.Errorf("sanitizeTargetName(%q) = %q; want %q", c.input, out, c.expected)
		}
	}
}

func TestZipDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "montage-zip-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_ = os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("hello montage"), 0644)
	subDir := filepath.Join(tempDir, "sub")
	_ = os.MkdirAll(subDir, 0755)
	_ = os.WriteFile(filepath.Join(subDir, "sub.txt"), []byte("nested file"), 0644)

	var buf bytes.Buffer
	if err := ZipDirectory(tempDir, &buf); err != nil {
		t.Fatalf("ZipDirectory failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("failed to read generated zip: %v", err)
	}

	foundFile := false
	for _, f := range zr.File {
		if f.Name == "test.txt" || f.Name == "sub/sub.txt" {
			foundFile = true
		}
	}

	if !foundFile {
		t.Errorf("zip archive did not contain expected files")
	}
}
