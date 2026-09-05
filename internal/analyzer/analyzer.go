package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ghchinoy/montage/internal/model"
)

// AnalyzeCodebase inspects a local directory and populates TargetInfo and AnalysisFindings.
func AnalyzeCodebase(localPath string, existingTarget *model.TargetInfo) (*model.TargetInfo, *model.AnalysisFindings, error) {
	absPath, err := filepath.Abs(localPath)
	if err != nil {
		return nil, nil, err
	}

	target := existingTarget
	if target == nil {
		target = &model.TargetInfo{
			LocalPath:  absPath,
			SourceKind: model.SourceKindLocal,
		}
	} else {
		target.LocalPath = absPath
	}

	findings := &model.AnalysisFindings{
		StructuredFormats:  make([]string, 0),
		EnvVarsDetected:    make([]string, 0),
		ConfigFiles:        make([]string, 0),
		ExistingAgentFiles: make([]string, 0),
		CommandGroups:      make(map[string]string),
		Prompts:            make([]model.PromptSpec, 0),
		PackageDomains:     make([]model.PackageSpec, 0),
	}

	// Step 1: Walk filesystem to detect languages, build systems, config files, and candidate code files
	fileMap := make(map[string]int) // extension -> count
	var codeFiles []string
	configSet := make(map[string]bool)
	agentFiles := make([]string, 0)

	skipDirs := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, "dist": true,
		"build": true, ".astro": true, ".next": true, "sources": true,
		"bin": true, "obj": true, ".cache": true, ".venv": true,
		"venv": true, "env": true, "__pycache__": true, ".pytest_cache": true,
		".ruff_cache": true, "site-packages": true, ".tox": true,
		".beads": true, ".serena": true, ".firebase": true, "out": true,
	}

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}

		rel, err := filepath.Rel(absPath, path)
		if err != nil {
			return nil
		}

		parts := strings.Split(rel, string(filepath.Separator))
		for _, p := range parts {
			if skipDirs[p] {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if info.IsDir() {
			// Check if architectural domain package under pkg/ or internal/ or services/
			parts := strings.Split(rel, string(filepath.Separator))
			if len(parts) == 2 && (parts[0] == "pkg" || parts[0] == "internal" || parts[0] == "services" || parts[0] == "modules") {
				pkgName := parts[1]
				domainName := inferPackageDomain(pkgName)
				if domainName != "" {
					findings.PackageDomains = append(findings.PackageDomains, model.PackageSpec{
						Path:        rel,
						Name:        pkgName,
						Domain:      domainName,
						Description: fmt.Sprintf("Architectural domain package %s: %s", pkgName, domainName),
					})
				}
			}
			return nil
		}

		baseName := filepath.Base(path)
		lowerBase := strings.ToLower(baseName)
		ext := strings.ToLower(filepath.Ext(path))

		// Detect authored sub-agent prompts
		if ext == ".md" && (strings.Contains(rel, "prompts/") || strings.Contains(rel, "agents/") || strings.Contains(rel, "roles/")) {
			promptName := strings.TrimSuffix(baseName, ext)
			desc, role := extractPromptMetadata(path)
			findings.Prompts = append(findings.Prompts, model.PromptSpec{
				Path:        rel,
				Name:        promptName,
				Description: desc,
				Role:        role,
			})
		}

		// Check for agent ecosystem files
		if lowerBase == "skill.md" || strings.HasSuffix(lowerBase, ".skill.md") {
			agentFiles = append(agentFiles, rel)
			findings.ExistingSkill = true
		}
		if lowerBase == "plugin.json" || lowerBase == "mcp.json" || lowerBase == "mcp_config.json" {
			agentFiles = append(agentFiles, rel)
			if lowerBase == "plugin.json" {
				findings.ExistingPlugin = true
			}
			if lowerBase == "mcp.json" || lowerBase == "mcp_config.json" {
				findings.ExistingMCP = true
			}
		}

		// Check for build systems & configs
		switch lowerBase {
		case "go.mod":
			findings.BuildSystem = "go"
			findings.PackagingType = "binary"
			findings.HasSingleBinary = true
		case "package.json":
			if findings.BuildSystem == "" {
				findings.BuildSystem = "npm"
				findings.PackagingType = "script"
			}
		case "cargo.toml":
			findings.BuildSystem = "cargo"
			findings.PackagingType = "binary"
			findings.HasSingleBinary = true
		case "pyproject.toml", "setup.py", "requirements.txt":
			if findings.BuildSystem == "" {
				findings.BuildSystem = "python"
				findings.PackagingType = "script"
			}
		case "makefile":
			if findings.BuildSystem == "" {
				findings.BuildSystem = "make"
			}
		case "dockerfile", "containerfile":
			findings.HasContainer = true
		case ".goreleaser.yaml", ".goreleaser.yml":
			findings.HasSingleBinary = true
		}

		// Config files detection
		if strings.HasPrefix(lowerBase, ".env") ||
			lowerBase == "config.json" || lowerBase == "config.yaml" ||
			lowerBase == "config.yml" || lowerBase == "config.toml" ||
			strings.HasSuffix(lowerBase, ".config.json") || strings.HasSuffix(lowerBase, ".config.yaml") {
			configSet[rel] = true
		}

		// Tally language extensions
		if ext != "" {
			fileMap[ext]++
		}

		// Detect candidate standalone scripts in scripts/, tools/, or bin/
		if (strings.Contains(rel, "scripts/") || strings.Contains(rel, "tools/") || strings.HasPrefix(rel, "bin/")) &&
			(ext == ".py" || ext == ".sh" || ext == ".js" || ext == ".ts" || ext == ".go") {
			scriptName := strings.TrimSuffix(baseName, ext)
			isMut := isCommandMutating(scriptName, "")
			desc := extractScriptDescription(path)
			if desc == "" {
				desc = fmt.Sprintf("Standalone script: %s", scriptName)
			}
			findings.Scripts = append(findings.Scripts, model.ScriptSpec{
				Path:        rel,
				Name:        scriptName,
				Description: desc,
				IsMutating:  isMut,
			})
		}

		// Select code files for deeper pattern inspection
		switch ext {
		case ".go", ".py", ".ts", ".js", ".rs", ".sh", ".rb":
			if info.Size() < 500*1024 { // skip giant files > 500KB
				codeFiles = append(codeFiles, path)
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	for cfg := range configSet {
		findings.ConfigFiles = append(findings.ConfigFiles, cfg)
	}
	sort.Strings(findings.ConfigFiles)
	findings.ExistingAgentFiles = agentFiles

	// Populate target languages if missing
	detectLanguages(target, fileMap)

	// Step 2: Deep code scanning for CLI, JSON, Env vars, Safety, and Prompts
	scanCodebase(absPath, codeFiles, findings)

	// Step 3: Infer target title & description if not set
	if target.Name == "" {
		target.Name = filepath.Base(absPath)
	}
	if target.Description == "" {
		target.Description = extractReadmeDescription(absPath)
	}

	return target, findings, nil
}

func detectLanguages(target *model.TargetInfo, fileMap map[string]int) {
	langScores := make(map[string]int)
	for ext, count := range fileMap {
		switch ext {
		case ".go":
			langScores["Go"] += count
		case ".py":
			langScores["Python"] += count
		case ".ts", ".tsx":
			langScores["TypeScript"] += count
		case ".js", ".jsx", ".mjs":
			langScores["JavaScript"] += count
		case ".rs":
			langScores["Rust"] += count
		case ".sh", ".bash":
			langScores["Shell"] += count
		case ".rb":
			langScores["Ruby"] += count
		}
	}

	type langCount struct {
		Name  string
		Count int
	}
	var sortedLangs []langCount
	for l, c := range langScores {
		sortedLangs = append(sortedLangs, langCount{l, c})
	}
	sort.Slice(sortedLangs, func(i, j int) bool {
		return sortedLangs[i].Count > sortedLangs[j].Count
	})

	if len(target.Languages) == 0 {
		for _, sl := range sortedLangs {
			target.Languages = append(target.Languages, sl.Name)
		}
	}
	if target.PrimaryLang == "" && len(sortedLangs) > 0 {
		target.PrimaryLang = sortedLangs[0].Name
	}
}

func extractReadmeDescription(root string) string {
	readmeCandidates := []string{"README.md", "readme.md", "README", "readme.markdown"}
	for _, cand := range readmeCandidates {
		p := filepath.Join(root, cand)
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[!") || strings.HasPrefix(line, "<") {
				continue
			}
			// First informative line
			if len(line) > 15 {
				return line
			}
		}
	}
	return ""
}

// Regex patterns for inspection
var (
	// CLI framework detectors
	reCobraCmd = regexp.MustCompile(`(&cobra\.Command|cobra\.Command)\s*\{`)
	reCobraUse = regexp.MustCompile(`Use:\s*"([^"\s]+)`)
	reCobraShort = regexp.MustCompile(`Short:\s*"([^"]+)`)
	reUrfaveCmd = regexp.MustCompile(`(&cli\.Command|cli\.Command)\s*\{`)
	reClickCmd  = regexp.MustCompile(`@(click|cli)\.(command|group)\s*\(`)
	reArgparse  = regexp.MustCompile(`ArgumentParser\s*\(`)
	reCommander = regexp.MustCompile(`\.(command|description|action)\s*\(`)
	reClap      = regexp.MustCompile(`(#\[derive\((.*Parser|.*Subcommand)\)\]|clap::Command::new\("([^"]+)"\))`)

	// Flag detectors
	reFlagJSON   = regexp.MustCompile(`(?i)(--(json|format)|"json"|"format")`)
	reFlagDryRun = regexp.MustCompile(`(?i)(--(dry-run|dryrun|validate-only)|"dry-run"|"dryrun")`)
	reFlagQuiet  = regexp.MustCompile(`(?i)(--(quiet|silent)|"quiet"|"silent"|"-q")`)

	// Serialization / Output detectors
	reJSONMarshal = regexp.MustCompile(`(json\.Marshal|json\.NewEncoder|json\.dumps|JSON\.stringify|serde_json::to_)`)
	reStderrPrint = regexp.MustCompile(`(os\.Stderr|fmt\.Fprintf\(os\.Stderr|console\.error|sys\.stderr)`)

	// Interactive prompts detectors
	reInteractiveGo   = regexp.MustCompile(`(survey\.|promptui\.|huh\.|bubbletea|fmt\.Scan|bufio\.NewReader\(os\.Stdin\))`)
	reInteractivePy   = regexp.MustCompile(`(input\(|inquirer\.|questionary\.|rich\.prompt)`)
	reInteractiveNode = regexp.MustCompile(`(inquirer\.prompt|prompts\(|readline\.createInterface)`)

	// Env var detectors
	reEnvGo   = regexp.MustCompile(`os\.(Getenv|LookupEnv)\("([A-Z0-9_]{3,})"\)`)
	reEnvPy   = regexp.MustCompile(`os\.environ(\.get)?\[?\(?["']([A-Z0-9_]{3,})["']`)
	reEnvNode = regexp.MustCompile(`process\.env\.([A-Z0-9_]{3,})`)
	reEnvRust = regexp.MustCompile(`env::var\("([A-Z0-9_]{3,})"\)`)

	// API Framework & Route detectors
	reFastAPIRoute = regexp.MustCompile(`(?m)@(app|router)\.(get|post|put|delete|patch)\(\s*["']([^"']+)["'](?:[^\)]*response_model\s*=\s*([A-Za-z0-9_]+))?`)
	reExpressRoute = regexp.MustCompile(`(?m)(app|router)\.(get|post|put|delete|patch)\(\s*["']([^"']+)["']`)
	reGinRoute     = regexp.MustCompile(`(?m)\.(GET|POST|PUT|DELETE|PATCH)\(\s*["']([^"']+)["']`)

	// Data model & storage detectors
	reFirestoreCollection = regexp.MustCompile(`\.collection\(["']([a-zA-Z0-9_-]+)["']\)`)
	rePydanticModel       = regexp.MustCompile(`class\s+([A-Za-z0-9_]+)\s*\(\s*BaseModel\s*\):`)

	// Command group detector
	reCobraAddGroup = regexp.MustCompile(`AddGroup\(&cobra\.Group\{\s*ID:\s*"([^"]+)",\s*Title:\s*"([^"]+)"`)
	reCobraGroupID  = regexp.MustCompile(`GroupID:\s*"([^"]+)"`)

	// MCP detection
	reMCPImport = regexp.MustCompile(`("github\.com/modelcontextprotocol/go-sdk|@modelcontextprotocol/sdk|modelcontextprotocol)`)
)

func inferPackageDomain(pkg string) string {
	lower := strings.ToLower(pkg)
	switch {
	case strings.Contains(lower, "orchestrat"):
		return "Publication Pipeline Orchestration"
	case strings.Contains(lower, "generator"):
		return "Multi-Modal Visuals & Plots"
	case strings.Contains(lower, "registry") || strings.Contains(lower, "template"):
		return "Template Catalog & Typesetting"
	case strings.Contains(lower, "classifier"):
		return "Selective Context & Source Analysis"
	case strings.Contains(lower, "a2a") || strings.Contains(lower, "a2api"):
		return "Agent-to-Agent (A2A) Protocol Bridge"
	case strings.Contains(lower, "cartograph"):
		return "Cartography & Concept Visualization"
	case strings.Contains(lower, "taxonom"):
		return "Thematic Taxonomy Discovery"
	default:
		return ""
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

func extractPromptMetadata(path string) (string, string) {
	f, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var role string
	var desc string
	linesRead := 0

	for scanner.Scan() && linesRead < 30 {
		linesRead++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "You are ") {
			role = strings.TrimPrefix(line, "You are ")
			role = strings.TrimSuffix(role, ".")
		} else if desc == "" && !strings.HasPrefix(line, "#") {
			desc = line
		}
	}
	if desc == "" && role != "" {
		desc = role
	}
	return desc, role
}

func extractScriptDescription(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	inDoc := false
	var lines []string
	count := 0

	for scanner.Scan() && count < 25 {
		count++
		l := strings.TrimSpace(scanner.Text())
		if l == "" || strings.HasPrefix(l, "#!") {
			continue
		}
		if strings.HasPrefix(l, `"""`) || strings.HasPrefix(l, `'''`) {
			if strings.Count(l, `"""`) >= 2 || strings.Count(l, `'''`) >= 2 {
				clean := strings.Trim(l, `"' `)
				if len(clean) > 5 {
					return clean
				}
			}
			inDoc = !inDoc
			clean := strings.Trim(l, `"' `)
			if clean != "" {
				lines = append(lines, clean)
			}
			continue
		}
		if inDoc {
			if strings.HasSuffix(l, `"""`) || strings.HasSuffix(l, `'''`) {
				clean := strings.Trim(l, `"' `)
				if clean != "" {
					lines = append(lines, clean)
				}
				break
			}
			lines = append(lines, l)
			continue
		}
		if strings.HasPrefix(l, "#") || strings.HasPrefix(l, "//") {
			comment := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(l, "#"), "//"))
			if len(comment) > 10 && !strings.HasPrefix(comment, "Copyright") {
				return comment
			}
		}
	}
	if len(lines) > 0 {
		return strings.Join(lines, " ")
	}
	return ""
}

func extractRouteDocstring(content, routePath string) string {
	idx := strings.Index(content, routePath)
	if idx == -1 {
		return ""
	}
	sub := content[idx:]
	defIdx := strings.Index(sub, "def ")
	if defIdx == -1 || defIdx > 250 {
		return ""
	}
	subDef := sub[defIdx:]
	colonIdx := strings.Index(subDef, ":")
	if colonIdx == -1 || colonIdx > 350 {
		return ""
	}
	afterColon := strings.TrimSpace(subDef[colonIdx+1:])
	if strings.HasPrefix(afterColon, `"""`) {
		afterDoc := afterColon[3:]
		endDoc := strings.Index(afterDoc, `"""`)
		if endDoc != -1 {
			doc := strings.TrimSpace(afterDoc[:endDoc])
			doc = strings.ReplaceAll(doc, "\n", " ")
			words := strings.Fields(doc)
			return strings.Join(words, " ")
		}
	}
	return ""
}

func extractRequestModel(content, routePath string) string {
	idx := strings.Index(content, routePath)
	if idx == -1 {
		return ""
	}
	sub := content[idx:]
	defIdx := strings.Index(sub, "def ")
	if defIdx == -1 || defIdx > 250 {
		return ""
	}
	colonIdx := strings.Index(sub[defIdx:], ":")
	if colonIdx == -1 || colonIdx > 350 {
		return ""
	}
	sig := sub[defIdx : defIdx+colonIdx]
	reParam := regexp.MustCompile(`request\s*:\s*([A-Za-z0-9_]+)`)
	m := reParam.FindStringSubmatch(sig)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func scanCodebase(root string, files []string, findings *model.AnalysisFindings) {
	envSet := make(map[string]bool)
	commandsMap := make(map[string]model.CommandSpec)
	interactivePromptSet := make(map[string]bool)

	for _, file := range files {
		contentBytes, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(contentBytes)

		// 1. MCP detection
		if reMCPImport.MatchString(content) {
			findings.ExistingMCP = true
		}

		// 2. CLI Frameworks & Commands
		if reCobraCmd.MatchString(content) {
			findings.HasCLI = true
			if findings.CLIFramework == "" {
				findings.CLIFramework = "cobra"
			}
			// Extract command uses
			useMatches := reCobraUse.FindAllStringSubmatch(content, -1)
			shortMatches := reCobraShort.FindAllStringSubmatch(content, -1)
			for i, m := range useMatches {
				cmdName := m[1]
				if cmdName == "help" || strings.HasPrefix(cmdName, "_") {
					continue
				}
				desc := ""
				if i < len(shortMatches) {
					desc = shortMatches[i][1]
				}
				// Determine parent and fullPath
				fileBase := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
				parent := ""
				fullPath := cmdName
				if fileBase != "root" && fileBase != "main" && fileBase != cmdName {
					parent = fileBase
					fullPath = parent + " " + cmdName
				}
				// Look for GroupID
				groupID := ""
				if gm := reCobraGroupID.FindStringSubmatch(content); len(gm) > 1 {
					groupID = gm[1]
				}
				spec := model.CommandSpec{
					Name:        cmdName,
					Parent:      parent,
					GroupID:     groupID,
					FullPath:    fullPath,
					Description: desc,
					IsMutating:  isCommandMutating(fullPath, desc),
					HasDryRun:   reFlagDryRun.MatchString(content),
					HasJSONOut:  reFlagJSON.MatchString(content),
				}
				commandsMap[fullPath] = spec
			}
		}

		// Detect Cobra command groups
		for _, m := range reCobraAddGroup.FindAllStringSubmatch(content, -1) {
			if len(m) > 2 {
				findings.CommandGroups[m[1]] = strings.TrimSpace(strings.TrimSuffix(m[2], ":"))
			}
		}

		// Detect A2A protocol and multi-modal generators
		if strings.Contains(strings.ToLower(content), "a2a") || strings.Contains(content, "AgentCard") {
			findings.HasA2A = true
		}
		if strings.Contains(content, "graphviz") || strings.Contains(content, "gonum") || strings.Contains(content, "nanobanana") || strings.Contains(content, "dotData") {
			findings.HasGenerators = true
		}

		if reClickCmd.MatchString(content) {
			findings.HasCLI = true
			if findings.CLIFramework == "" {
				findings.CLIFramework = "click"
			}
		} else if reArgparse.MatchString(content) {
			findings.HasCLI = true
			if findings.CLIFramework == "" {
				findings.CLIFramework = "argparse"
			}
		} else if reCommander.MatchString(content) {
			findings.HasCLI = true
			if findings.CLIFramework == "" {
				findings.CLIFramework = "commander"
			}
		} else if reClap.MatchString(content) {
			findings.HasCLI = true
			if findings.CLIFramework == "" {
				findings.CLIFramework = "clap"
			}
		}

		// 2b. API Frameworks & Routes
		if reFastAPIRoute.MatchString(content) {
			findings.HasAPI = true
			if findings.APIFramework == "" {
				findings.APIFramework = "fastapi"
			}
			matches := reFastAPIRoute.FindAllStringSubmatch(content, -1)
			for _, m := range matches {
				routePath := m[3]
				method := strings.ToUpper(m[2])
				respModel := ""
				if len(m) > 4 {
					respModel = m[4]
				}
				reqModel := extractRequestModel(content, routePath)
				doc := extractRouteDocstring(content, routePath)
				if doc == "" {
					doc = fmt.Sprintf("%s %s endpoint", method, routePath)
				}
				isMut := method != "GET" && method != "HEAD" && method != "OPTIONS"
				if isCommandMutating(routePath, doc) {
					isMut = true
				}

				exists := false
				for _, ep := range findings.APIEndpoints {
					if ep.Path == routePath && ep.Method == method {
						exists = true
						break
					}
				}
				if !exists {
					findings.APIEndpoints = append(findings.APIEndpoints, model.APIEndpointSpec{
						Path:          routePath,
						Method:        method,
						Description:   doc,
						RequestModel:  reqModel,
						ResponseModel: respModel,
						IsMutating:    isMut,
					})
				}
			}
		}

		if reExpressRoute.MatchString(content) {
			findings.HasAPI = true
			if findings.APIFramework == "" {
				findings.APIFramework = "express"
			}
			matches := reExpressRoute.FindAllStringSubmatch(content, -1)
			for _, m := range matches {
				routePath := m[3]
				method := strings.ToUpper(m[2])
				isMut := method != "GET"
				exists := false
				for _, ep := range findings.APIEndpoints {
					if ep.Path == routePath && ep.Method == method {
						exists = true
						break
					}
				}
				if !exists {
					findings.APIEndpoints = append(findings.APIEndpoints, model.APIEndpointSpec{
						Path:        routePath,
						Method:      method,
						Description: fmt.Sprintf("%s %s endpoint", method, routePath),
						IsMutating:  isMut,
					})
				}
			}
		}

		if reGinRoute.MatchString(content) {
			findings.HasAPI = true
			if findings.APIFramework == "" {
				findings.APIFramework = "gin"
			}
			matches := reGinRoute.FindAllStringSubmatch(content, -1)
			for _, m := range matches {
				method := strings.ToUpper(m[1])
				routePath := m[2]
				isMut := method != "GET"
				exists := false
				for _, ep := range findings.APIEndpoints {
					if ep.Path == routePath && ep.Method == method {
						exists = true
						break
					}
				}
				if !exists {
					findings.APIEndpoints = append(findings.APIEndpoints, model.APIEndpointSpec{
						Path:        routePath,
						Method:      method,
						Description: fmt.Sprintf("%s %s handler", method, routePath),
						IsMutating:  isMut,
					})
				}
			}
		}

		// 2c. Data Models & Storage Collections
		for _, m := range reFirestoreCollection.FindAllStringSubmatch(content, -1) {
			if len(m) > 1 {
				findings.DataModels = appendUnique(findings.DataModels, "collection:"+m[1])
			}
		}
		for _, m := range rePydanticModel.FindAllStringSubmatch(content, -1) {
			if len(m) > 1 {
				findings.DataModels = appendUnique(findings.DataModels, "model:"+m[1])
			}
		}

		// 3. Flags and Data Interchange
		if reFlagJSON.MatchString(content) {
			findings.HasJSONOutput = true
			findings.JSONFlags = appendUnique(findings.JSONFlags, "--json")
		}
		if reFlagDryRun.MatchString(content) {
			findings.HasDryRunFlag = true
		}
		if reJSONMarshal.MatchString(content) {
			findings.StructuredFormats = appendUnique(findings.StructuredFormats, "json")
		}
		if reStderrPrint.MatchString(content) {
			findings.CleanOutputStreams = true
		}

		// 4. Interactive prompts
		if m := reInteractiveGo.FindString(content); m != "" {
			findings.HasInteractiveUI = true
			interactivePromptSet[m] = true
		}
		if m := reInteractivePy.FindString(content); m != "" {
			findings.HasInteractiveUI = true
			interactivePromptSet[m] = true
		}
		if m := reInteractiveNode.FindString(content); m != "" {
			findings.HasInteractiveUI = true
			interactivePromptSet[m] = true
		}

		// 5. Env variables
		for _, match := range reEnvGo.FindAllStringSubmatch(content, -1) {
			if len(match) > 2 && isMeaningfulEnvVar(match[2]) {
				envSet[match[2]] = true
			}
		}
		for _, match := range reEnvPy.FindAllStringSubmatch(content, -1) {
			if len(match) > 2 && isMeaningfulEnvVar(match[2]) {
				envSet[match[2]] = true
			}
		}
		for _, match := range reEnvNode.FindAllStringSubmatch(content, -1) {
			if len(match) > 1 && isMeaningfulEnvVar(match[1]) {
				envSet[match[1]] = true
			}
		}
		for _, match := range reEnvRust.FindAllStringSubmatch(content, -1) {
			if len(match) > 1 && isMeaningfulEnvVar(match[1]) {
				envSet[match[1]] = true
			}
		}
	}

	for env := range envSet {
		findings.EnvVarsDetected = append(findings.EnvVarsDetected, env)
	}
	sort.Strings(findings.EnvVarsDetected)
	if len(findings.EnvVarsDetected) > 0 {
		findings.AuthMethods = appendUnique(findings.AuthMethods, "environment_variables")
	}

	for p := range interactivePromptSet {
		findings.InteractivePrompts = append(findings.InteractivePrompts, p)
	}

	for _, cmd := range commandsMap {
		findings.Commands = append(findings.Commands, cmd)
		cmdPath := cmd.FullPath
		if cmdPath == "" {
			cmdPath = cmd.Name
		}
		if cmd.IsMutating {
			findings.MutatingCommands = append(findings.MutatingCommands, cmdPath)
			if isCommandDangerous(cmd.Name) {
				findings.HasDangerousActions = true
			}
		} else {
			findings.ReadOnlyCommands = append(findings.ReadOnlyCommands, cmdPath)
		}
	}

	sort.Slice(findings.Commands, func(i, j int) bool {
		return findings.Commands[i].FullPath < findings.Commands[j].FullPath
	})
}

func isMeaningfulEnvVar(name string) bool {
	switch name {
	case "PATH", "HOME", "USER", "SHELL", "PWD", "TMPDIR":
		return false
	default:
		return true
	}
}

func isCommandMutating(name, desc string) bool {
	combined := strings.ToLower(name + " " + desc)
	mutatingVerbs := []string{
		"create", "delete", "remove", "drop", "destroy", "update",
		"write", "modify", "apply", "set", "kill", "purge", "clean",
		"deploy", "scaffold", "generate", "install", "publish",
	}
	for _, v := range mutatingVerbs {
		if strings.Contains(combined, v) {
			return true
		}
	}
	return false
}

func isCommandDangerous(name string) bool {
	lower := strings.ToLower(name)
	dangerVerbs := []string{"delete", "remove", "drop", "destroy", "kill", "purge", "rm"}
	for _, v := range dangerVerbs {
		if strings.Contains(lower, v) {
			return true
		}
	}
	return false
}

func appendUnique(slice []string, val string) []string {
	for _, item := range slice {
		if item == val {
			return slice
		}
	}
	return append(slice, val)
}
