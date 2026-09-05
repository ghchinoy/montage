package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	mcpCmd.AddCommand(mcpConfigCmd)
	mcpConfigCmd.Flags().String("client", "all", "Target MCP client: all, claude, cursor, opencode, antigravity")
}

var mcpConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Print MCP client configuration snippets for Claude Desktop, Cursor, OpenCode, and Antigravity",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := cmd.Flags().GetString("client")

		binPath, err := os.Executable()
		if err != nil {
			binPath = "montage"
		} else {
			if resolved, err := filepath.EvalSymlinks(binPath); err == nil {
				binPath = resolved
			}
		}

		client = strings.ToLower(strings.TrimSpace(client))

		switch client {
		case "claude", "claude-desktop":
			printClaudeConfig(binPath)
		case "cursor":
			printCursorConfig(binPath)
		case "opencode":
			printOpenCodeConfig(binPath)
		case "all":
			printAllConfigs(binPath)
		default:
			return fmt.Errorf("unknown client '%s'. Supported: all, claude, cursor, opencode", client)
		}

		return nil
	},
}

func printClaudeConfig(binPath string) {
	cfg := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"montage": map[string]interface{}{
				"command": binPath,
				"args":    []string{"mcp"},
			},
		},
	}
	bytes, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Printf("### Claude Desktop Configuration (~/Library/Application Support/Claude/claude_desktop_config.json):\n\n```json\n%s\n```\n", string(bytes))
}

func printCursorConfig(binPath string) {
	cfg := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"montage": map[string]interface{}{
				"command": binPath,
				"args":    []string{"mcp"},
			},
		},
	}
	bytes, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Printf("### Cursor MCP Configuration (.cursor/mcp.json):\n\n```json\n%s\n```\n", string(bytes))
}

func printOpenCodeConfig(binPath string) {
	cfg := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"montage": map[string]interface{}{
				"command": binPath,
				"args":    []string{"mcp"},
			},
		},
	}
	bytes, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Printf("### OpenCode Configuration (~/.config/opencode/opencode.json or opencode.json):\n\n```json\n%s\n```\n", string(bytes))
}

func printAllConfigs(binPath string) {
	printClaudeConfig(binPath)
	fmt.Println()
	printCursorConfig(binPath)
	fmt.Println()
	printOpenCodeConfig(binPath)
}
