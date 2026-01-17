# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build
```bash
make build       # Build to build/ccstatus with optimizations (-s -w flags)
make all         # Run lint and build
make clean       # Clean build artifacts
```

### Lint
```bash
make lint        # Run golangci-lint with strict .golangci.yaml config
make lint-fix    # Auto-fix linting issues
```

### Test
```bash
make test        # Run tests
```

## Architecture

This is a Go implementation of a Claude Code statusline generator that follows a provider-component architecture with self-registering plugins.

### Core Pattern
1. **Providers** (`internal/providers/`) - Fetch data from various sources
   - Each provider implements `Provider` interface with `Key()` and `Provide(ctx)` methods
   - Self-register via init() functions when imported
   - Providers run in parallel during data gathering via goroutines
   - Wrapped with `CachingProvider` for TTL-based caching

2. **Components** (`internal/components/`) - Render specific parts of the status line
   - Each component implements `Component` interface with `Render(ctx)` and `RequiredProviders()` methods
   - Self-register via init() functions with name string when imported
   - Component order determined by config or defaults
   - Optional `OptionalComponent` interface for conditional rendering

3. **StatusLine** (`internal/core/statusline.go`) - Main orchestrator
   - Gathers data from all providers in parallel
   - Renders components in configured order
   - Handles separator formatting between components

### Data Flow
1. Main reads Claude session JSON from stdin
2. Load configuration via Koanf
3. Create cache (FileCache or NullCache based on config)
4. Determine active components from config
5. Create only required providers based on component dependencies
6. StatusLine.Render() executes providers in parallel, then renders components
7. Output printed to stdout with ANSI colors

### Key Interfaces
- `core.Provider` - Data fetching interface (`Key()`, `Provide(ctx)`)
- `core.Component` - Rendering interface (`Render(ctx)`, `RequiredProviders()`)
- `core.Cache` - Caching abstraction (FileCache or NullCache)
- `core.RenderContext` - Shared context for passing data from providers to components
- `core.Registry` - Central registry for self-registering components and providers

### Configuration
- Uses Koanf for flexible config loading (YAML format)
- Config file locations (first found wins):
  - `.claude/ccstatus.local.yaml` - Project-specific local config (gitignored)
  - `.claude/ccstatus.yaml` - Project-specific shared config
  - `~/.claude/ccstatus.yaml` - User default config
- Configuration uses Claude session's project directory, not execution directory
- Full configuration example in `config.yaml` with extensive comments
- Components can be enabled/disabled and reordered via `active` config key

### Input Format
The tool expects JSON input via stdin with the following structure:
```json
{
  "session_id": "session-uuid",
  "transcript_path": "/path/to/transcript.jsonl",
  "cwd": "/current/working/directory",
  "model": {
    "id": "claude-opus-4-5-20251101",
    "display_name": "Opus 4.5"
  },
  "workspace": {
    "current_dir": "/workspace/current",
    "project_dir": "/workspace/project"
  },
  "version": "2.1.9",
  "output_style": {
    "name": "default"
  },
  "cost": {
    "total_cost_usd": 0.3633969,
    "total_duration_ms": 108056,
    "total_api_duration_ms": 30204,
    "total_lines_added": 67,
    "total_lines_removed": 12
  },
  "context_window": {
    "total_input_tokens": 15234,
    "total_output_tokens": 4521,
    "context_window_size": 200000,
    "current_usage": {
      "input_tokens": 8500,
      "output_tokens": 1200,
      "cache_creation_input_tokens": 5000,
      "cache_read_input_tokens": 2000
    },
    "used_percentage": 81,
    "remaining_percentage": 19
  },
  "exceeds_200k_tokens": false
}
```
