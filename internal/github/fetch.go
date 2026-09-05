package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghchinoy/montage/internal/model"
)

// ResolveGitHubToken returns a GitHub personal access token if present in environment variables.
func ResolveGitHubToken() string {
	for _, env := range []string{"GITHUB_TOKEN", "GH_TOKEN"} {
		if tok := strings.TrimSpace(os.Getenv(env)); tok != "" {
			return tok
		}
	}
	return ""
}

// CheckPrerequisites verifies that `git` is available.
func CheckPrerequisites() error {
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git command is required but was not found in PATH")
	}
	return nil
}

// ParseRepoRef extracts owner and repo name from GitHub URL or "owner/repo" shorthand.
func ParseRepoRef(input string) (string, string, error) {
	clean := strings.TrimSpace(input)
	clean = strings.TrimSuffix(clean, ".git")

	// Check if URL
	if strings.HasPrefix(clean, "http://") || strings.HasPrefix(clean, "https://") {
		u, err := url.Parse(clean)
		if err != nil {
			return "", "", fmt.Errorf("invalid URL: %w", err)
		}
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], nil
		}
		return "", "", fmt.Errorf("URL does not contain owner/repo: %s", input)
	}

	if strings.HasPrefix(clean, "github.com/") {
		parts := strings.Split(strings.TrimPrefix(clean, "github.com/"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], nil
		}
		return "", "", fmt.Errorf("invalid github.com path: %s", input)
	}

	parts := strings.Split(clean, "/")
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unrecognized GitHub repo specification: %s (expected owner/repo or https://github.com/owner/repo)", input)
}

// IsGitHubTarget checks if the input is a GitHub URL or owner/repo format rather than a local path.
func IsGitHubTarget(input string) bool {
	// If it exists locally on the filesystem, it's local
	if _, err := os.Stat(input); err == nil {
		return false
	}
	if strings.HasPrefix(input, "./") || strings.HasPrefix(input, "../") || strings.HasPrefix(input, "/") {
		return false
	}
	if strings.Contains(input, "github.com") {
		return true
	}
	parts := strings.Split(strings.TrimSpace(input), "/")
	return len(parts) == 2 && !strings.Contains(parts[0], ".")
}

type ghPublicRepoResponse struct {
	Name            string     `json:"name"`
	FullName        string     `json:"full_name"`
	Description     string     `json:"description"`
	HTMLURL         string     `json:"html_url"`
	DefaultBranch   string     `json:"default_branch"`
	Language        string     `json:"language"`
	Topics          []string   `json:"topics"`
	StargazersCount int        `json:"stargazers_count"`
	CreatedAt       *time.Time `json:"created_at"`
	PushedAt        *time.Time `json:"pushed_at"`
	Message         string     `json:"message,omitempty"` // For API error messages
}

type ghRepoView struct {
	Name            string `json:"name"`
	NameWithOwner   string `json:"nameWithOwner"`
	Description     string `json:"description"`
	URL             string `json:"url"`
	DefaultBranchRef struct {
		Name string `json:"name"`
	} `json:"defaultBranchRef"`
	PrimaryLanguage *struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`
	Languages []struct {
		Node struct {
			Name string `json:"name"`
		} `json:"node"`
	} `json:"languages"`
	RepositoryTopics []struct {
		Name string `json:"name"`
	} `json:"repositoryTopics"`
	StargazerCount int        `json:"stargazerCount"`
	CreatedAt      *time.Time `json:"createdAt"`
	PushedAt       *time.Time `json:"pushedAt"`
}

// FetchRepoMetadata retrieves repository information using the GitHub public REST API (unauthenticated or with GITHUB_TOKEN),
// with a fallback to `gh repo view` if the CLI is authenticated.
func FetchRepoMetadata(owner, repo string) (*model.TargetInfo, error) {
	if err := CheckPrerequisites(); err != nil {
		return nil, err
	}

	// Strategy 1: GitHub Public REST API (works unauthenticated on Cloud Run)
	targetInfo, httpErr := fetchMetadataHTTP(owner, repo)
	if httpErr == nil && targetInfo != nil {
		return targetInfo, nil
	}

	// Strategy 2: Fallback to GitHub CLI (`gh repo view`) if installed and authenticated
	if _, err := exec.LookPath("gh"); err == nil {
		if cliInfo, cliErr := fetchMetadataCLI(owner, repo); cliErr == nil && cliInfo != nil {
			return cliInfo, nil
		}
	}

	if httpErr != nil {
		return nil, fmt.Errorf("failed to fetch metadata for %s/%s: %w", owner, repo, httpErr)
	}
	return nil, fmt.Errorf("failed to fetch metadata for %s/%s", owner, repo)
}

func fetchMetadataHTTP(owner, repo string) (*model.TargetInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Montage/0.2.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if token := ResolveGitHubToken(); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errData)
		msg, _ := errData["message"].(string)
		if resp.StatusCode == http.StatusForbidden && strings.Contains(strings.ToLower(msg), "rate limit") {
			return nil, fmt.Errorf("GitHub API rate limit exceeded. Set GITHUB_TOKEN environment variable to enable 5,000 requests/hour: %s", msg)
		}
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("repository %s/%s not found (or private; set GITHUB_TOKEN for private repos)", owner, repo)
		}
		return nil, fmt.Errorf("GitHub API returned HTTP %d: %s", resp.StatusCode, msg)
	}

	var data ghPublicRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode GitHub API response: %w", err)
	}

	var langs []string
	if data.Language != "" {
		langs = append(langs, data.Language)
	}

	return &model.TargetInfo{
		Name:          data.Name,
		Owner:         owner,
		URL:           data.HTMLURL,
		SourceKind:    model.SourceKindGitHub,
		Description:   data.Description,
		DefaultBranch: data.DefaultBranch,
		PrimaryLang:   data.Language,
		Languages:     langs,
		Topics:        data.Topics,
		Stars:         data.StargazersCount,
		CreatedAt:     data.CreatedAt,
		PushedAt:      data.PushedAt,
	}, nil
}

func fetchMetadataCLI(owner, repo string) (*model.TargetInfo, error) {
	nwo := fmt.Sprintf("%s/%s", owner, repo)
	args := []string{
		"repo", "view", nwo,
		"--json", "name,nameWithOwner,description,url,defaultBranchRef,primaryLanguage,languages,repositoryTopics,stargazerCount,createdAt,pushedAt",
	}

	cmd := exec.Command("gh", args...)
	if token := ResolveGitHubToken(); token != "" {
		cmd.Env = append(os.Environ(), "GH_TOKEN="+token)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gh repo view failed: %s", strings.TrimSpace(stderr.String()))
	}

	var data ghRepoView
	if err := json.Unmarshal(stdout.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("failed to parse repo view json: %w", err)
	}

	primaryLang := ""
	if data.PrimaryLanguage != nil {
		primaryLang = data.PrimaryLanguage.Name
	}

	var langs []string
	for _, l := range data.Languages {
		if l.Node.Name != "" {
			langs = append(langs, l.Node.Name)
		}
	}
	if len(langs) == 0 && primaryLang != "" {
		langs = append(langs, primaryLang)
	}

	var topics []string
	for _, t := range data.RepositoryTopics {
		if t.Name != "" {
			topics = append(topics, t.Name)
		}
	}

	return &model.TargetInfo{
		Name:          data.Name,
		Owner:         owner,
		URL:           data.URL,
		SourceKind:    model.SourceKindGitHub,
		Description:   data.Description,
		DefaultBranch: data.DefaultBranchRef.Name,
		PrimaryLang:   primaryLang,
		Languages:     langs,
		Topics:        topics,
		Stars:         data.StargazerCount,
		CreatedAt:     data.CreatedAt,
		PushedAt:      data.PushedAt,
	}, nil
}

// CloneOrPull clones the target GitHub repository into sourcesDir (default: ./sources/owner/repo).
// It uses standard anonymous `git clone --depth 1` for public repositories, or attaches GITHUB_TOKEN if present.
func CloneOrPull(owner, repo, sourcesDir string) (string, error) {
	if err := CheckPrerequisites(); err != nil {
		return "", err
	}

	if sourcesDir == "" {
		sourcesDir = "sources"
	}

	absSourcesDir, err := filepath.Abs(sourcesDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path of sources dir: %w", err)
	}

	targetDir := filepath.Join(absSourcesDir, owner, repo)

	// Check if already cloned and has files
	if entries, err := os.ReadDir(targetDir); err == nil && len(entries) > 0 {
		return targetDir, nil
	}

	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return "", fmt.Errorf("failed to create parent dir: %w", err)
	}

	// Strategy 1: Standard git clone (Anonymous public clone or token-authenticated)
	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
	if token := ResolveGitHubToken(); token != "" {
		cloneURL = fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s.git", token, owner, repo)
	}

	cmd := exec.Command("git", "clone", "--depth", "1", cloneURL, targetDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err == nil {
		return targetDir, nil
	}

	gitErr := strings.TrimSpace(stderr.String())

	// Strategy 2: Fallback to `gh repo clone` if installed and authenticated
	if _, err := exec.LookPath("gh"); err == nil {
		nwo := fmt.Sprintf("%s/%s", owner, repo)
		ghCmd := exec.Command("gh", "repo", "clone", nwo, targetDir, "--", "--depth", "1")
		if token := ResolveGitHubToken(); token != "" {
			ghCmd.Env = append(os.Environ(), "GH_TOKEN="+token)
		}
		var ghStderr bytes.Buffer
		ghCmd.Stderr = &ghStderr
		if err := ghCmd.Run(); err == nil {
			return targetDir, nil
		}
	}

	return "", fmt.Errorf("failed to clone %s/%s via git: %s", owner, repo, gitErr)
}
