package evaluator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ghchinoy/montage/internal/model"
)

// ExtractAndGroupCapabilities analyzes discovered endpoints, commands, scripts, prompts, and packages
// to produce an itemized capability inventory, grouped toolsets, and multi-skill blueprints.
func ExtractAndGroupCapabilities(target model.TargetInfo, f model.AnalysisFindings) ([]model.Capability, []model.SkillBlueprint, []model.ToolGroup) {
	var capabilities []model.Capability

	appName := strings.ToLower(target.Name)
	if appName == "" {
		appName = "app"
	}

	// 1. Process Authored Sub-Agent Prompts (High priority: author-designed roles/rules)
	for _, p := range f.Prompts {
		cap := classifyPrompt(appName, p)
		capabilities = append(capabilities, cap)
	}

	// 2. Process Architectural Package Domains (Engine packages)
	for _, pkg := range f.PackageDomains {
		cap := classifyPackageDomain(appName, pkg)
		capabilities = append(capabilities, cap)
	}

	// 3. Process API Endpoints
	for _, ep := range f.APIEndpoints {
		cap := classifyAPIEndpoint(appName, ep)
		capabilities = append(capabilities, cap)
	}

	// 4. Process CLI Commands (Hierarchical with group context)
	for _, cmd := range f.Commands {
		// Skip low-signal subcommands if part of parent
		if isAuxiliaryCommand(cmd.Name) && cmd.Parent != "" {
			continue
		}
		cap := classifyCLICommand(appName, cmd, f.CommandGroups)
		capabilities = append(capabilities, cap)
	}

	// 5. Process Standalone Scripts
	for _, scr := range f.Scripts {
		cap := classifyScript(appName, scr)
		capabilities = append(capabilities, cap)
	}

	// 6. Synthesize Inferred Domain Capabilities (Clinical, cartographic, etc.)
	domainCaps := inferDomainCapabilities(appName, target, f)
	capabilities = append(capabilities, domainCaps...)

	// 7. Deduplicate capabilities by Name
	seenCaps := make(map[string]bool)
	var dedupedCaps []model.Capability
	for _, c := range capabilities {
		if !seenCaps[c.Name] {
			seenCaps[c.Name] = true
			dedupedCaps = append(dedupedCaps, c)
		}
	}
	capabilities = dedupedCaps

	// 8. Cluster into MCP Tool Groups
	toolGroups := clusterMCPToolGroups(capabilities, f)

	// 9. Cluster into Multi-Skill Blueprints
	skillBlueprints := clusterSkillBlueprints(appName, capabilities, f)

	return capabilities, skillBlueprints, toolGroups
}

func isAuxiliaryCommand(name string) bool {
	switch name {
	case "help", "version", "completion":
		return true
	default:
		return false
	}
}

func classifyPrompt(appName string, p model.PromptSpec) model.Capability {
	cleanID := strings.ReplaceAll(p.Name, "-", "_")
	lower := strings.ToLower(p.Name + " " + p.Role + " " + p.Description)

	grouping := "Specialized Sub-Agent Personas"
	switch {
	case strings.Contains(lower, "review") || strings.Contains(lower, "editor") || strings.Contains(lower, "advisory"):
		grouping = "Advisory Review & Feedback"
	case strings.Contains(lower, "analyst") || strings.Contains(lower, "code") || strings.Contains(lower, "data") || strings.Contains(lower, "context"):
		grouping = "Selective Context & Source Analysis"
	case strings.Contains(lower, "write") || strings.Contains(lower, "draft") || strings.Contains(lower, "abstract"):
		grouping = "Document Drafting & Composition"
	case strings.Contains(lower, "image") || strings.Contains(lower, "visual") || strings.Contains(lower, "plot"):
		grouping = "Multi-Modal Visuals & Plots"
	}

	desc := p.Description
	if desc == "" {
		desc = fmt.Sprintf("Authored prompt instructions for %s persona", p.Name)
	}

	return model.Capability{
		ID:                fmt.Sprintf("cap_prompt_%s", cleanID),
		Name:              formatCapabilityName(p.Name + " Persona"),
		Description:       desc,
		Kind:              model.CapabilityKindDomainLogic,
		SourceRef:         p.Path,
		IsMutating:        false,
		RecommendedTarget: model.TargetPureSkill,
		TargetGrouping:    grouping,
		Rationale:         "Authored sub-agent system instructions; provides specialized domain persona, behavioral guardrails, and review criteria as a Pure Skill.",
	}
}

func classifyPackageDomain(appName string, pkg model.PackageSpec) model.Capability {
	cleanID := strings.ReplaceAll(pkg.Name, "-", "_")
	lower := strings.ToLower(pkg.Name + " " + pkg.Domain)

	target := model.TargetSkillWithCode
	rationale := "Internal engine package suitable for wrapping in an agent execution skill or helper utility."

	if strings.Contains(lower, "generator") || strings.Contains(lower, "plot") || strings.Contains(lower, "graphviz") {
		target = model.TargetSkillWithCode
		rationale = "Autonomous visual generators (DOT flowcharts, data plots, illustrations); ideal for a Skill with Code or MCP tool."
	} else if strings.Contains(lower, "orchestrat") || strings.Contains(lower, "pipeline") || strings.Contains(lower, "dag") {
		target = model.TargetMCPTool
		rationale = "Multi-phase agent execution DAG; best triggered and monitored via MCP tools."
	} else if strings.Contains(lower, "a2a") || strings.Contains(lower, "protocol") {
		target = model.TargetMCPTool
		rationale = "Agent-to-Agent protocol server; enables native multi-agent coordination."
	} else if strings.Contains(lower, "classifier") {
		target = model.TargetSkillWithCode
		rationale = "Content classification and context budgeting; token saver before agent prompting."
	} else if strings.Contains(lower, "registry") || strings.Contains(lower, "template") {
		target = model.TargetMCPTool
		rationale = "Template discovery and typesetting package registry."
	}

	return model.Capability{
		ID:                fmt.Sprintf("cap_pkg_%s", cleanID),
		Name:              pkg.Domain,
		Description:       pkg.Description,
		Kind:              model.CapabilityKindDomainLogic,
		SourceRef:         pkg.Path,
		IsMutating:        false,
		RecommendedTarget: target,
		TargetGrouping:    pkg.Domain,
		Rationale:         rationale,
	}
}

func classifyAPIEndpoint(appName string, ep model.APIEndpointSpec) model.Capability {
	cleanPath := strings.TrimPrefix(ep.Path, "/")
	cleanID := strings.ReplaceAll(strings.ReplaceAll(cleanPath, "/", "_"), "-", "_")
	if cleanID == "" {
		cleanID = "root"
	}

	lowerPath := strings.ToLower(ep.Path)
	grouping := "Core Services"
	target := model.TargetMCPTool
	rationale := "Live state query or execution endpoint ideal for an MCP tool."

	if strings.Contains(lowerPath, "weight") || strings.Contains(lowerPath, "metric") || strings.Contains(lowerPath, "health") || strings.Contains(lowerPath, "vitals") {
		grouping = "Clinical Telemetry & Growth"
		if strings.Contains(lowerPath, "analyze") {
			target = model.TargetSkillWithCode
			rationale = "Performs mathematical analysis over historical records; ideal for a Skill with Code helper script to save LLM context tokens."
		}
	} else if strings.Contains(lowerPath, "bio") || strings.Contains(lowerPath, "adoption") || strings.Contains(lowerPath, "content") || strings.Contains(lowerPath, "generate") {
		grouping = "Adoption & Copywriting"
		rationale = "Content synthesis service; best paired with an editorial Pure Skill + MCP generation tool."
	} else if strings.Contains(lowerPath, "user") || strings.Contains(lowerPath, "role") || strings.Contains(lowerPath, "auth") || strings.Contains(lowerPath, "admin") {
		grouping = "Administration & Access"
	} else if strings.Contains(lowerPath, "litter") || strings.Contains(lowerPath, "kitten") || strings.Contains(lowerPath, "item") || strings.Contains(lowerPath, "repo") {
		grouping = "Inventory & Entity Management"
	}

	name := formatCapabilityName(cleanID)
	desc := ep.Description
	if desc == "" {
		desc = fmt.Sprintf("%s %s endpoint", ep.Method, ep.Path)
	}

	return model.Capability{
		ID:                fmt.Sprintf("cap_api_%s", cleanID),
		Name:              name,
		Description:       desc,
		Kind:              model.CapabilityKindAPI,
		SourceRef:         fmt.Sprintf("%s %s", ep.Method, ep.Path),
		HTTPMethod:        ep.Method,
		IsMutating:        ep.IsMutating,
		RecommendedTarget: target,
		TargetGrouping:    grouping,
		Rationale:         rationale,
	}
}

func classifyCLICommand(appName string, cmd model.CommandSpec, commandGroups map[string]string) model.Capability {
	displayName := cmd.FullPath
	if displayName == "" {
		displayName = cmd.Name
	}
	cleanID := strings.ReplaceAll(strings.ReplaceAll(displayName, " ", "_"), "-", "_")
	grouping := ""
	target := model.TargetMCPTool
	rationale := "Exposes an operational command directly over Model Context Protocol."

	lower := strings.ToLower(displayName + " " + cmd.Description)

	// Priority 1: Check declared Cobra GroupID
	if cmd.GroupID != "" && commandGroups != nil {
		if grpTitle, ok := commandGroups[cmd.GroupID]; ok && grpTitle != "" {
			switch strings.ToLower(cmd.GroupID) {
			case "execution":
				if strings.Contains(lower, "review") {
					grouping = "Advisory Review & Feedback"
					target = model.TargetPureSkill
					rationale = "Advisory review mode preserves original text and surfaces structured improvement recommendations; ideal as a Pure Skill."
				} else {
					grouping = "Publication Pipeline Orchestration"
					rationale = "Executes multi-phase document synthesis DAG; ideal as an MCP execution tool."
				}
			case "management":
				grouping = "Project Lifecycle & Source Context"
				rationale = "Project context management; suitable for scriptable MCP tools."
			case "discovery":
				grouping = "Template Catalog & Typesetting"
				target = model.TargetMCPTool
				rationale = "Discovers and queries Typst universe package templates."
			case "system":
				grouping = "Server & Protocol Hosting"
				target = model.TargetCLICommand
				rationale = "Service runner and system configuration."
			default:
				grouping = grpTitle
			}
		}
	}

	// Priority 2: Use parent command and semantic keywords
	if grouping == "" {
		if cmd.Parent == "project" || cmd.Parent == "source" || cmd.Parent == "goal" || cmd.Parent == "asset" {
			grouping = "Project Lifecycle & Source Context"
		} else if cmd.Parent == "template" || strings.Contains(lower, "template") || strings.Contains(lower, "typst") {
			grouping = "Template Catalog & Typesetting"
		} else if strings.Contains(lower, "review") || strings.Contains(lower, "advisory") {
			grouping = "Advisory Review & Feedback"
			target = model.TargetPureSkill
			rationale = "Advisory review mode preserves original text and surfaces structured improvement recommendations; ideal as a Pure Skill."
		} else if strings.Contains(lower, "execute") || strings.Contains(lower, "run") || strings.Contains(lower, "orchestrat") {
			grouping = "Publication Pipeline Orchestration"
		} else if strings.Contains(lower, "map") || strings.Contains(lower, "cartograph") || strings.Contains(lower, "graph") {
			grouping = "Cartography & Concept Visualization"
		} else if strings.Contains(lower, "suggest") || strings.Contains(lower, "taxonomy") {
			grouping = "Thematic Taxonomy Discovery"
		} else if strings.Contains(lower, "render") {
			grouping = "Visual Diagram Rendering"
		} else if strings.Contains(lower, "serve") || strings.Contains(lower, "config") || strings.Contains(lower, "version") {
			grouping = "System & Host Operations"
			target = model.TargetCLICommand
		} else {
			grouping = "General Operations"
		}
	}

	return model.Capability{
		ID:                fmt.Sprintf("cap_cmd_%s", cleanID),
		Name:              formatCapabilityName(displayName),
		Description:       cmd.Description,
		Kind:              model.CapabilityKindCLI,
		SourceRef:         fmt.Sprintf("command: %s", displayName),
		IsMutating:        cmd.IsMutating,
		RecommendedTarget: target,
		TargetGrouping:    grouping,
		Rationale:         rationale,
	}
}

func classifyScript(appName string, scr model.ScriptSpec) model.Capability {
	cleanID := strings.ReplaceAll(scr.Name, "-", "_")
	grouping := "Utility Scripts"
	target := model.TargetSkillWithCode
	rationale := "Standalone utility script suited for bundling into a Skill with Code."

	lower := strings.ToLower(scr.Name + " " + scr.Path)
	if strings.Contains(lower, "role") || strings.Contains(lower, "rule") || strings.Contains(lower, "auth") || strings.Contains(lower, "admin") {
		grouping = "Shelter Administration & Security"
		target = model.TargetCLICommand
		rationale = "Administrative utility suitable for shell execution or operator checklist."
	} else if strings.Contains(lower, "list") || strings.Contains(lower, "dump") || strings.Contains(lower, "export") {
		grouping = "Data Ingestion & Export"
		target = model.TargetSkillWithCode
		rationale = "Data extraction script; saves tokens by filtering raw database outputs."
	}

	return model.Capability{
		ID:                fmt.Sprintf("cap_script_%s", cleanID),
		Name:              formatCapabilityName(scr.Name),
		Description:       scr.Description,
		Kind:              model.CapabilityKindScript,
		SourceRef:         scr.Path,
		IsMutating:        scr.IsMutating,
		RecommendedTarget: target,
		TargetGrouping:    grouping,
		Rationale:         rationale,
	}
}

func inferDomainCapabilities(appName string, target model.TargetInfo, f model.AnalysisFindings) []model.Capability {
	var inferred []model.Capability

	lowerDesc := strings.ToLower(target.Description + " " + target.Name)

	if strings.Contains(lowerDesc, "kitten") || strings.Contains(lowerDesc, "foster") || strings.Contains(lowerDesc, "litter") {
		inferred = append(inferred, model.Capability{
			ID:                "cap_domain_neonate_triage",
			Name:              "Foster Neonatal Growth & Triage Protocol",
			Description:       "Clinical triage rules for evaluating daily weight loss (>10g), feeding frequencies, and health escalation thresholds.",
			Kind:              model.CapabilityKindDomainLogic,
			SourceRef:         "Domain Clinical Standards",
			IsMutating:        false,
			RecommendedTarget: model.TargetPureSkill,
			TargetGrouping:    "Clinical Telemetry & Growth",
			Rationale:         "Clinical judgment and operational decision trees require veterinary guidelines, not code execution.",
		})
		inferred = append(inferred, model.Capability{
			ID:                "cap_domain_adoption_guidelines",
			Name:              "Adoption Persona & Copywriting Guidelines",
			Description:       "Style guide and personality matching rubric for drafting compelling adoption announcements.",
			Kind:              model.CapabilityKindDomainLogic,
			SourceRef:         "Domain Editorial Standards",
			IsMutating:        false,
			RecommendedTarget: model.TargetPureSkill,
			TargetGrouping:    "Adoption & Copywriting",
			Rationale:         "Persona formatting rules and tone instructions are best formulated as a Pure Skill.",
		})
	}

	if strings.Contains(lowerDesc, "cartographer") || strings.Contains(lowerDesc, "taxonomy") || strings.Contains(lowerDesc, "repotographer") {
		inferred = append(inferred, model.Capability{
			ID:                "cap_domain_connectivity_model",
			Name:              "Tri-State Connectivity & Architecture Model",
			Description:       "Architectural rubric for interpreting connected, standalone, and bare repositories in concept graphs.",
			Kind:              model.CapabilityKindDomainLogic,
			SourceRef:         "docs/reference/connectivity-model.mdx",
			IsMutating:        false,
			RecommendedTarget: model.TargetPureSkill,
			TargetGrouping:    "Cartography & Concept Visualization",
			Rationale:         "Domain interpretation rules best given to LLM as architectural guidance.",
		})
	}

	return inferred
}

func clusterMCPToolGroups(caps []model.Capability, f model.AnalysisFindings) []model.ToolGroup {
	groupMap := make(map[string][]model.MCPToolSpec)

	for _, c := range caps {
		if c.RecommendedTarget == model.TargetMCPTool || c.Kind == model.CapabilityKindAPI || (c.Kind == model.CapabilityKindCLI && !c.IsMutating) {
			toolName := strings.ReplaceAll(strings.ToLower(c.Name), " ", "_")
			spec := model.MCPToolSpec{
				Name:        toolName,
				Description: c.Description,
				IsReadOnly:  !c.IsMutating,
				Grouping:    c.TargetGrouping,
			}
			groupMap[c.TargetGrouping] = append(groupMap[c.TargetGrouping], spec)
		}
	}

	var groups []model.ToolGroup
	for grpName, tools := range groupMap {
		groups = append(groups, model.ToolGroup{
			Name:        grpName,
			Description: fmt.Sprintf("Tools supporting %s workflows", strings.ToLower(grpName)),
			Tools:       tools,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups
}

func clusterSkillBlueprints(appName string, caps []model.Capability, f model.AnalysisFindings) []model.SkillBlueprint {
	var blueprints []model.SkillBlueprint
	domainMap := make(map[string][]model.Capability)

	for _, c := range caps {
		domainMap[c.TargetGrouping] = append(domainMap[c.TargetGrouping], c)
	}

	var sortedDomains []string
	for d := range domainMap {
		sortedDomains = append(sortedDomains, d)
	}
	sort.Strings(sortedDomains)

	for _, domain := range sortedDomains {
		domainCaps := domainMap[domain]
		hasPure := false
		hasCode := false
		var capNames []string
		var helperScript string

		for _, c := range domainCaps {
			capNames = append(capNames, c.Name)
			if c.RecommendedTarget == model.TargetPureSkill {
				hasPure = true
			}
			if c.RecommendedTarget == model.TargetSkillWithCode {
				hasCode = true
				if c.Kind == model.CapabilityKindScript {
					helperScript = c.SourceRef
				}
			}
		}
		_ = hasPure

		skillKind := "pure_skill"
		tokenSavings := "n/a (zero execution overhead)"
		if hasCode {
			skillKind = "skill_with_code"
			tokenSavings = "70-90% token reduction via local data aggregation"
			if helperScript == "" {
				cleanD := strings.ToLower(strings.ReplaceAll(domain, " ", "_"))
				cleanD = strings.ReplaceAll(cleanD, "&_", "")
				cleanD = strings.ReplaceAll(cleanD, "&", "")
				helperScript = fmt.Sprintf("scripts/%s_helper.py", cleanD)
			}
		}

		skillSlug := mapDomainToSkillSlug(appName, domain)

		blueprints = append(blueprints, model.SkillBlueprint{
			Name:              skillSlug,
			Kind:              skillKind,
			SuggestedFileName: fmt.Sprintf("%s.skill.md", skillSlug),
			Description:       fmt.Sprintf("Specialized agent skill for %s within %s.", domain, appName),
			Domain:            domain,
			Capabilities:      capNames,
			HelperScript:      helperScript,
			TokenSavingsEst:   tokenSavings,
			RecommendedUses: []string{
				fmt.Sprintf("Autonomous management of %s tasks", domain),
				"Domain verification and structured decision-making",
			},
			KeyInstructions: []string{
				fmt.Sprintf("Evaluate target data according to %s standards.", domain),
				"Verify input preconditions and run read-only validations before modifications.",
			},
		})
	}

	return blueprints
}

func mapDomainToSkillSlug(appName, domain string) string {
	lower := strings.ToLower(domain)
	switch {
	case strings.Contains(lower, "advisory") || strings.Contains(lower, "review"):
		return fmt.Sprintf("%s-advisory-reviewer", appName)
	case strings.Contains(lower, "context") || strings.Contains(lower, "analyst"):
		return fmt.Sprintf("%s-source-analyst", appName)
	case strings.Contains(lower, "visual") || strings.Contains(lower, "plot") || strings.Contains(lower, "generator"):
		return fmt.Sprintf("%s-multi-modal-visualizer", appName)
	case strings.Contains(lower, "typeset") || strings.Contains(lower, "template"):
		return fmt.Sprintf("%s-typst-typesetter", appName)
	case strings.Contains(lower, "pipeline") || strings.Contains(lower, "orchestrat"):
		return fmt.Sprintf("%s-publication-orchestrator", appName)
	case strings.Contains(lower, "a2a") || strings.Contains(lower, "interoperab"):
		return fmt.Sprintf("%s-a2a-bridge", appName)
	case strings.Contains(lower, "clinical") || strings.Contains(lower, "growth") || strings.Contains(lower, "telemetry"):
		return fmt.Sprintf("%s-growth-triage", appName)
	case strings.Contains(lower, "adoption") || strings.Contains(lower, "copywriting"):
		return fmt.Sprintf("%s-adoption-copywriter", appName)
	case strings.Contains(lower, "cartograph") || strings.Contains(lower, "concept"):
		return fmt.Sprintf("%s-cartographer", appName)
	case strings.Contains(lower, "taxonomy") || strings.Contains(lower, "discovery"):
		return fmt.Sprintf("%s-taxonomy-curator", appName)
	case strings.Contains(lower, "lifecycle") || strings.Contains(lower, "project"):
		return fmt.Sprintf("%s-project-manager", appName)
	default:
		clean := strings.ToLower(strings.ReplaceAll(domain, " ", "-"))
		clean = strings.ReplaceAll(clean, "&-", "")
		clean = strings.ReplaceAll(clean, "&", "")
		clean = strings.ReplaceAll(clean, "(", "")
		clean = strings.ReplaceAll(clean, ")", "")
		clean = strings.ReplaceAll(clean, "--", "-")
		clean = strings.Trim(clean, "-")
		return fmt.Sprintf("%s-%s", appName, clean)
	}
}

func formatCapabilityName(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, " ")
}
