---
allowed-tools: Glob, Grep, Read, Write, MultiEdit, Bash(go run:*), Bash(find:*), TodoWrite, mcp__gopls__*
description: Generate or update config.yaml with all default values from codebase scan
argument-hint: [output file path (default: config.yaml)]
---

## Context

- Project: ccstatus-go - Claude Code statusline generator
- Architecture: Provider-Component pattern with self-registering plugins
- Config system: Uses Koanf with YAML format
- Config search order:
  1. `.claude/ccstatus.local.yaml` - Project-specific local (gitignored)
  2. `.claude/ccstatus.yaml` - Project-specific shared
  3. `~/.claude/ccstatus.yaml` - User default

## Task

Scan the codebase to find all default configuration values and generate a comprehensive config.yaml file with:

1. **Component Configuration**
   - Find all registered components in `internal/components/`
   - Extract default component order from `cmd/ccstatus/main.go`
   - Extract settings from each component's `config.go` file
   - Document each component's template variables and options

2. **Provider Configuration**
   - Find all providers in `internal/providers/`
   - Extract default provider settings (cache TTL, etc.)
   - Document provider-specific cache configurations

3. **Core Configuration**
   - Find core settings in `internal/core/` and `cmd/`
   - Extract cache settings from `internal/cache/cache.go`
   - Extract separator settings from `internal/core/statusline.go`
   - Document global options and defaults

4. **Structure the YAML** with:
   - Comprehensive header explaining config file locations
   - Meaningful inline comments for each option
   - Default values clearly marked
   - Example configurations at the bottom
   - Available color options documented
   - Template function reference

## Output

$ARGUMENTS

Generate the config file to the path provided above (default: config.yaml) with:
- All discovered default values from component/provider config.go files
- Comprehensive documentation comments explaining each option
- Proper YAML structure with consistent formatting
- Template variable documentation for each component
- Example configurations for common use cases

## Implementation Steps

1. Scan component registry and implementations for defaults
2. Scan provider implementations for default settings
3. Scan core configuration loading code for default values
4. Scan for any hardcoded values that should be configurable
5. Generate well-documented YAML with all findings
