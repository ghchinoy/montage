---
title: Actionable Gaps & Remediation
description: Code patterns and remediation guides to elevate applications into Tier 5 Agent-Native tools.
---

Montage provides specific code snippets in its assessment reports to remediate common gaps:

### 1. Adding `--json` Output Mode (Cobra Go)

```go
cmd.Flags().Bool("json", false, "Output results as machine-readable JSON")
if jsonOut {
    json.NewEncoder(os.Stdout).Encode(result)
    return nil
}
```

### 2. Supporting Non-Interactive Execution

```go
if !isatty.IsTerminal(os.Stdin.Fd()) || nonInteractiveFlag {
    // Proceed with default answers or fail deterministically
}
```

### 3. Adding `--dry-run` Simulation

```go
cmd.Flags().Bool("dry-run", false, "Simulate operation without writing changes")
if dryRun {
    fmt.Println("Simulation passed: 0 mutations applied.")
    return nil
}
```

### 4. Separating Diagnostic Logs from Data Payloads

Direct all informational messages, spinners, and progress logs to `os.Stderr` so that `os.Stdout` remains pure JSON.
