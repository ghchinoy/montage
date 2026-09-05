package model

import "time"

// SourceKind indicates if the target was loaded from local filesystem or GitHub.
type SourceKind string

const (
	SourceKindLocal  SourceKind = "local"
	SourceKindGitHub SourceKind = "github"
)

// TargetInfo contains metadata about the inspected application/repository.
type TargetInfo struct {
	Name          string     `json:"name"`
	Owner         string     `json:"owner,omitempty"`
	URL           string     `json:"url,omitempty"`
	LocalPath     string     `json:"local_path"`
	SourceKind    SourceKind `json:"source_kind"`
	Description   string     `json:"description,omitempty"`
	DefaultBranch string     `json:"default_branch,omitempty"`
	PrimaryLang   string     `json:"primary_language,omitempty"`
	Languages     []string   `json:"languages,omitempty"`
	Topics        []string   `json:"topics,omitempty"`
	Stars         int        `json:"stars,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	PushedAt      *time.Time `json:"pushed_at,omitempty"`
}

// CommandSpec represents a discovered CLI command or subcommand.
type CommandSpec struct {
	Name        string   `json:"name"`
	Parent      string   `json:"parent,omitempty"`
	GroupID     string   `json:"group_id,omitempty"`
	FullPath    string   `json:"full_path,omitempty"`
	Description string   `json:"description,omitempty"`
	Flags       []string `json:"flags,omitempty"`
	IsMutating  bool     `json:"is_mutating"`
	HasDryRun   bool     `json:"has_dry_run"`
	HasJSONOut  bool     `json:"has_json_out"`
}

// APIEndpointSpec represents a discovered HTTP/REST route.
type APIEndpointSpec struct {
	Path          string `json:"path"`
	Method        string `json:"method"`
	HandlerName   string `json:"handler_name,omitempty"`
	Description   string `json:"description,omitempty"`
	RequestModel  string `json:"request_model,omitempty"`
	ResponseModel string `json:"response_model,omitempty"`
	IsMutating    bool   `json:"is_mutating"`
}

// ScriptSpec represents a standalone utility script found in the repository.
type ScriptSpec struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsMutating  bool   `json:"is_mutating"`
}

// PromptSpec represents an authored sub-agent prompt or system instructions file.
type PromptSpec struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Role        string `json:"role,omitempty"`
}

// PackageSpec represents an internal architectural domain or engine package.
type PackageSpec struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Domain      string `json:"domain"`
	Description string `json:"description,omitempty"`
}

// AnalysisFindings contains all discovered technical capabilities and patterns.
type AnalysisFindings struct {
	// CLI & Interface Findings
	HasCLI             bool              `json:"has_cli"`
	CLIFramework       string            `json:"cli_framework,omitempty"`
	Commands           []CommandSpec     `json:"commands,omitempty"`
	CommandGroups      map[string]string `json:"command_groups,omitempty"` // GroupID -> Title
	HasInteractiveUI   bool              `json:"has_interactive_ui"`       // TUI, prompt-based inputs
	InteractivePrompts []string          `json:"interactive_prompts,omitempty"`

	// API Findings
	HasAPI       bool              `json:"has_api"`
	APIFramework string            `json:"api_framework,omitempty"`
	APIEndpoints []APIEndpointSpec `json:"api_endpoints,omitempty"`

	// Script & Utility Findings
	Scripts []ScriptSpec `json:"scripts,omitempty"`

	// Prompt & Sub-Agent Persona Findings
	Prompts []PromptSpec `json:"prompts,omitempty"`

	// Package Domain Findings
	PackageDomains []PackageSpec `json:"package_domains,omitempty"`

	// Protocols & Advanced Capability Flags
	HasA2A        bool `json:"has_a2a"`
	HasGenerators bool `json:"has_generators"`

	// Data Interchange Findings
	HasJSONOutput       bool     `json:"has_json_output"`
	JSONFlags           []string `json:"json_flags,omitempty"`
	StructuredFormats   []string `json:"structured_formats,omitempty"` // json, yaml, csv, etc.
	CleanOutputStreams  bool     `json:"clean_output_streams"`        // logs to stderr, data to stdout
	DataModels          []string `json:"data_models,omitempty"`       // Firestore, SQL, Pydantic entities

	// Execution Safety Findings
	HasDryRunFlag       bool     `json:"has_dry_run_flag"`
	MutatingCommands    []string `json:"mutating_commands,omitempty"`
	ReadOnlyCommands    []string `json:"read_only_commands,omitempty"`
	HasDangerousActions bool     `json:"has_dangerous_actions"` // delete, drop, destroy

	// Auth & Configuration Findings
	AuthMethods        []string `json:"auth_methods,omitempty"` // env, flag, config_file, ambient_cli
	EnvVarsDetected    []string `json:"env_vars_detected,omitempty"`
	ConfigFiles        []string `json:"config_files,omitempty"`
	HasInteractiveAuth bool     `json:"has_interactive_auth"`

	// Packaging & Runtime Findings
	BuildSystem       string `json:"build_system,omitempty"`   // make, goreleaser, cargo, npm, pip
	PackagingType     string `json:"packaging_type,omitempty"` // binary, script, container, library
	HasSingleBinary   bool   `json:"has_single_binary"`
	HasContainer      bool   `json:"has_container"`
	DependenciesCount int    `json:"dependencies_count"`

	// Existing Agent Tooling Findings
	ExistingMCP        bool     `json:"existing_mcp"`
	ExistingSkill      bool     `json:"existing_skill"`
	ExistingPlugin     bool     `json:"existing_plugin"`
	ExistingAgentFiles []string `json:"existing_agent_files,omitempty"`
}

// CapabilityKind indicates the originating mechanism of a capability.
type CapabilityKind string

const (
	CapabilityKindAPI         CapabilityKind = "api_endpoint"
	CapabilityKindCLI         CapabilityKind = "cli_command"
	CapabilityKindScript      CapabilityKind = "script"
	CapabilityKindDomainLogic CapabilityKind = "domain_logic"
	CapabilityKindDataModel   CapabilityKind = "data_model"
)

// AgentTargetCategory indicates where a capability best belongs in the Agent Ecosystem.
type AgentTargetCategory string

const (
	TargetPureSkill     AgentTargetCategory = "pure_skill"
	TargetSkillWithCode AgentTargetCategory = "skill_with_code"
	TargetMCPTool       AgentTargetCategory = "mcp_tool"
	TargetCLICommand    AgentTargetCategory = "cli_command"
)

// Capability represents an evaluated capability of the application mapped to an agent target.
type Capability struct {
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	Kind              CapabilityKind      `json:"kind"`
	SourceRef         string              `json:"source_ref"`
	HTTPMethod        string              `json:"http_method,omitempty"`
	IsMutating        bool                `json:"is_mutating"`
	RecommendedTarget AgentTargetCategory `json:"recommended_target"`
	TargetGrouping    string              `json:"target_grouping"`
	Rationale         string              `json:"rationale"`
}

// ToolGroup organizes MCP tools into a cohesive persona or functional group.
type ToolGroup struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Tools       []MCPToolSpec `json:"tools"`
}

// SkillBlueprint represents an actionable, distinct skill blueprint.
type SkillBlueprint struct {
	Name              string   `json:"name"`
	Kind              string   `json:"kind"` // "pure_skill" or "skill_with_code"
	SuggestedFileName string   `json:"suggested_file_name"`
	Description       string   `json:"description"`
	Domain            string   `json:"domain"`
	Capabilities      []string `json:"capabilities"`
	HelperScript      string   `json:"helper_script,omitempty"`
	TokenSavingsEst   string   `json:"token_savings_est,omitempty"`
	RecommendedUses   []string `json:"recommended_uses"`
	KeyInstructions   []string `json:"key_instructions"`
}

// ReadinessDimension represents one evaluated dimension.
type ReadinessDimension struct {
	Name        string   `json:"name"`
	Score       int      `json:"score"` // 0-100
	Weight      float64  `json:"weight"`
	Explanation string   `json:"explanation"`
	Strengths   []string `json:"strengths,omitempty"`
	Gaps        []string `json:"gaps,omitempty"`
}

// ReadinessTier categorizes the overall maturity level.
type ReadinessTier string

const (
	Tier5AgentNative  ReadinessTier = "Tier 5: Agent-Native"
	Tier4AgentCapable ReadinessTier = "Tier 4: Agent-Capable"
	Tier3Scriptable   ReadinessTier = "Tier 3: Scriptable"
	Tier2WorkflowOnly ReadinessTier = "Tier 2: Workflow-Only"
	Tier1ManualGUI    ReadinessTier = "Tier 1: Manual / GUI-Bound"
)

// MCPToolSpec defines a candidate MCP tool suggested by Montage.
type MCPToolSpec struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	IsReadOnly  bool                   `json:"is_read_only"`
	SourceCmd   string                 `json:"source_cmd,omitempty"`
	Grouping    string                 `json:"grouping,omitempty"`
}

// PureSkillRecommendation details the Pure Skill opportunity.
type PureSkillRecommendation struct {
	Feasibility     string   `json:"feasibility"` // High, Medium, Low
	Score           int      `json:"score"`
	SuggestedName   string   `json:"suggested_name"`
	Description     string   `json:"description"`
	RecommendedUses []string `json:"recommended_uses"`
	KeyInstructions []string `json:"key_instructions"`
}

// SkillWithCodeRecommendation details the Skill with Code opportunity.
type SkillWithCodeRecommendation struct {
	Feasibility     string   `json:"feasibility"` // High, Medium, Low
	Score           int      `json:"score"`
	SuggestedName   string   `json:"suggested_name"`
	Description     string   `json:"description"`
	HelperScripts   []string `json:"helper_scripts"`
	TokenSavingsEst string   `json:"token_savings_est"`
	RecommendedFlow []string `json:"recommended_flow"`
}

// MCPServerRecommendation details the MCP Server opportunity.
type MCPServerRecommendation struct {
	Feasibility       string        `json:"feasibility"` // High, Medium, Low
	Score             int           `json:"score"`
	ServerName        string        `json:"server_name"`
	Transport         string        `json:"transport"` // stdio, sse, stream
	ProposedTools     []MCPToolSpec `json:"proposed_tools"`
	ProposedResources []string      `json:"proposed_resources,omitempty"`
	ProposedPrompts   []string      `json:"proposed_prompts,omitempty"`
}

// AgentPluginRecommendation details the Agent Plugin bundle opportunity.
type AgentPluginRecommendation struct {
	Feasibility string   `json:"feasibility"` // High, Medium, Low
	Score       int      `json:"score"`
	PluginName  string   `json:"plugin_name"`
	Description string   `json:"description"`
	Components  []string `json:"components"`
}

// GapItem is an actionable remediation item to increase agent readiness.
type GapItem struct {
	Priority    string `json:"priority"` // High, Medium, Low
	Dimension   string `json:"dimension"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Remediation string `json:"remediation"`
	ExampleCode string `json:"example_code,omitempty"`
}

// AssessmentReport is the comprehensive assessment output.
type AssessmentReport struct {
	GeneratedAt time.Time        `json:"generated_at"`
	Target      TargetInfo       `json:"target"`
	Findings    AnalysisFindings `json:"findings"`

	// Capabilities Inventory & Grouping
	Capabilities    []Capability     `json:"capabilities"`
	SkillBlueprints []SkillBlueprint `json:"skill_blueprints"`
	ToolGroups      []ToolGroup      `json:"tool_groups"`

	// Scores
	OverallScore int                  `json:"overall_score"` // 0-100
	Tier         ReadinessTier        `json:"tier"`
	TierSummary  string               `json:"tier_summary"`
	Dimensions   []ReadinessDimension `json:"dimensions"`

	// Artifact Spectrum Recommendations
	PureSkill     PureSkillRecommendation     `json:"pure_skill"`
	SkillWithCode SkillWithCodeRecommendation `json:"skill_with_code"`
	MCPServer     MCPServerRecommendation     `json:"mcp_server"`
	AgentPlugin   AgentPluginRecommendation   `json:"agent_plugin"`

	// Actionable Gaps & Modernization Roadmap
	ActionableGaps []GapItem `json:"actionable_gaps"`
}
