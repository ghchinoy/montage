package github

import (
	"testing"
)

func TestParseRepoRef(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{"shorthand", "ghchinoy/repotographer", "ghchinoy", "repotographer", false},
		{"https_url", "https://github.com/ghchinoy/repotographer", "ghchinoy", "repotographer", false},
		{"https_url_git_suffix", "https://github.com/ghchinoy/repotographer.git", "ghchinoy", "repotographer", false},
		{"host_shorthand", "github.com/ghchinoy/repotographer", "ghchinoy", "repotographer", false},
		{"missing_slash", "invalid", "", "", true},
		{"empty_url", "http://", "", "", true},
		{"flag_injection_owner", "--upload-pack=evil/repo", "", "", true},
		{"flag_injection_repo", "ghchinoy/--config", "", "", true},
		{"path_traversal_owner", "../etc/passwd", "", "", true},
		{"path_traversal_repo", "ghchinoy/..", "", "", true},
		{"url_path_traversal", "https://github.com/../etc", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := ParseRepoRef(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseRepoRef(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr {
				if owner != tt.wantOwner || repo != tt.wantRepo {
					t.Errorf("ParseRepoRef(%q) = (%q, %q), want (%q, %q)", tt.input, owner, repo, tt.wantOwner, tt.wantRepo)
				}
			}
		})
	}
}

func TestCloneOrPull_RejectsUnsafeIdentifiers(t *testing.T) {
	sourcesDir := t.TempDir()
	tests := []struct {
		name  string
		owner string
		repo  string
	}{
		{"leading_dash_owner", "-evil", "repo"},
		{"leading_dash_repo", "owner", "--upload-pack=sh"},
		{"dotdot_owner", "..", "repo"},
		{"dotdot_repo", "owner", ".."},
		{"slash_in_repo", "owner", "a/b"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if dir, err := CloneOrPull(tc.owner, tc.repo, sourcesDir); err == nil {
				t.Errorf("CloneOrPull(%q, %q, %q) = %q, want error", tc.owner, tc.repo, sourcesDir, dir)
			}
		})
	}
}

func TestIsGitHubTarget(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"ghchinoy/repotographer", true},
		{"https://github.com/ghchinoy/repotographer", true},
		{"./internal", false},
	}
	for _, tc := range tests {
		if got := IsGitHubTarget(tc.input); got != tc.want {
			t.Errorf("IsGitHubTarget(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestCheckPrerequisites(t *testing.T) {
	if err := CheckPrerequisites(); err != nil {
		t.Fatalf("CheckPrerequisites() error = %v, want nil", err)
	}
}

func TestResolveGitHubToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "test-token-123")
	if got := ResolveGitHubToken(); got != "test-token-123" {
		t.Errorf("ResolveGitHubToken() = %q, want %q", got, "test-token-123")
	}
}
