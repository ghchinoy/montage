package github

import (
	"testing"
)

func TestParseRepoRef(t *testing.T) {
	tests := []struct {
		input     string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{"ghchinoy/repotographer", "ghchinoy", "repotographer", false},
		{"https://github.com/ghchinoy/repotographer", "ghchinoy", "repotographer", false},
		{"https://github.com/ghchinoy/repotographer.git", "ghchinoy", "repotographer", false},
		{"github.com/ghchinoy/repotographer", "ghchinoy", "repotographer", false},
		{"invalid", "", "", true},
		{"http://", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
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

func TestIsGitHubTarget(t *testing.T) {
	if !IsGitHubTarget("ghchinoy/repotographer") {
		t.Errorf("expected ghchinoy/repotographer to be identified as GitHub target")
	}
	if !IsGitHubTarget("https://github.com/ghchinoy/repotographer") {
		t.Errorf("expected https://github.com/ghchinoy/repotographer to be identified as GitHub target")
	}
	if IsGitHubTarget("./internal") {
		t.Errorf("expected ./internal to not be identified as GitHub target")
	}
}

func TestCheckPrerequisites(t *testing.T) {
	err := CheckPrerequisites()
	if err != nil {
		t.Fatalf("expected CheckPrerequisites to succeed when git is installed, got %v", err)
	}
}

func TestResolveGitHubToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "test-token-123")
	tok := ResolveGitHubToken()
	if tok != "test-token-123" {
		t.Errorf("expected 'test-token-123', got '%s'", tok)
	}
}
