export interface TargetInfo {
  name: string;
  owner?: string;
  url?: string;
  local_path: string;
  source_kind: "local" | "github";
  description?: string;
  default_branch?: string;
  primary_language?: string;
  languages?: string[];
  topics?: string[];
  stars?: number;
}

export interface ReadinessDimension {
  name: string;
  score: number;
  weight: number;
  explanation: string;
  strengths?: string[];
  gaps?: string[];
}

export interface Capability {
  id: string;
  name: string;
  description: string;
  kind: "api_endpoint" | "cli_command" | "script" | "domain_logic" | "data_model";
  source_ref: string;
  http_method?: string;
  is_mutating: boolean;
  recommended_target: "pure_skill" | "skill_with_code" | "mcp_tool" | "cli_command";
  target_grouping: string;
  rationale: string;
}

export interface SkillBlueprint {
  name: string;
  kind: "pure_skill" | "skill_with_code";
  suggested_file_name: string;
  description: string;
  domain: string;
  capabilities: string[];
  helper_script?: string;
  token_savings_est?: string;
  recommended_uses?: string[];
  key_instructions?: string[];
}

export interface MCPToolSpec {
  name: string;
  description: string;
  parameters?: Record<string, any>;
  is_read_only: boolean;
  source_cmd?: string;
  grouping?: string;
}

export interface ToolGroup {
  name: string;
  description: string;
  tools: MCPToolSpec[];
}

export interface GapItem {
  priority: "High" | "Medium" | "Low";
  dimension: string;
  title: string;
  description: string;
  remediation: string;
  example_code?: string;
}

export interface AssessmentReport {
  generated_at: string;
  target: TargetInfo;
  overall_score: number;
  tier: string;
  tier_summary: string;
  dimensions: ReadinessDimension[];
  capabilities: Capability[];
  skill_blueprints: SkillBlueprint[];
  tool_groups: ToolGroup[];
  actionable_gaps: GapItem[];
}
