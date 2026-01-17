package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/mirage20/ccstatus-go/internal/cache"
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"

	// Import providers for self-registration.
	_ "github.com/mirage20/ccstatus-go/internal/providers/ratelimit"
	_ "github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"

	// Import components for self-registration.
	_ "github.com/mirage20/ccstatus-go/internal/components/claudecode/changes"
	_ "github.com/mirage20/ccstatus-go/internal/components/claudecode/context"
	_ "github.com/mirage20/ccstatus-go/internal/components/claudecode/cwd"
	_ "github.com/mirage20/ccstatus-go/internal/components/claudecode/duration"
	_ "github.com/mirage20/ccstatus-go/internal/components/claudecode/model"
	_ "github.com/mirage20/ccstatus-go/internal/components/claudecode/version"
	_ "github.com/mirage20/ccstatus-go/internal/components/layout/newline"
	_ "github.com/mirage20/ccstatus-go/internal/components/ratelimit/fivehour"
	_ "github.com/mirage20/ccstatus-go/internal/components/ratelimit/sevenday"
)

// Version information - set via ldflags during build.
var (
	version    = "dev"
	commit     = "unknown"
	commitDate = ""
	buildDate  = "unknown"
)

func main() {
	// Handle command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help", "-h", "--help":
			showHelp()
			return
		case "version", "-v", "--version":
			showVersion()
			return
		}
	}

	// Normal operation - read from stdin and generate status line
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// Read Claude session information from stdin (NOT a provider!)
	claudeSession, err := readClaudeSession(os.Stdin)
	if err != nil {
		// If no valid input, show help
		showHelp()
		return nil
	}

	// Load configuration with project directory from Claude session
	cfgReader := config.NewReader(claudeSession.Workspace.ProjectDir)

	// Create cache with session isolation using the new factory
	c := cache.New(cfgReader, claudeSession.SessionID)
	defer c.Close() // Ignore errors - don't pollute status line output

	// Create status line with configuration
	statusLine := core.NewStatusLine(cfgReader)

	// STEP 1: Get active components from config or use defaults
	componentNames := config.Get(cfgReader, "active", []string{})
	if len(componentNames) == 0 {
		// Default component order if no active list is configured
		componentNames = []string{
			"model",
			"context",
			"ratelimit.fivehour",
			"ratelimit.sevenday",
			"changes",
			"duration",
			"version",
			"newline",
			"cwd",
		}
	}

	// Create components and collect their provider requirements
	var components []core.Component
	providerSet := make(map[string]bool)

	for _, name := range componentNames {
		if comp, exists := core.CreateComponent(name, cfgReader); exists {
			components = append(components, comp)

			// Collect provider dependencies
			for _, providerName := range comp.RequiredProviders() {
				providerSet[providerName] = true
			}
		}
	}

	// STEP 2: Create only the providers that components need
	for providerName := range providerSet {
		// Create provider from registry (registry handles caching)
		if provider, exists := core.CreateProvider(providerName, cfgReader, claudeSession, c); exists {
			statusLine.AddProvider(provider)
		} else {
			// Log warning that a required provider is not registered
			fmt.Fprintf(os.Stderr, "Warning: Component requires provider '%s' but it's not registered\n", providerName)
		}
	}

	// STEP 3: Add components to statusline
	for _, comp := range components {
		statusLine.AddComponent(comp)
	}

	// Render and output the status line
	output := statusLine.Render(ctx)
	fmt.Fprintln(os.Stdout, output)

	return nil
}

// readClaudeSession reads the Claude session information from stdin.
func readClaudeSession(reader io.Reader) (*core.ClaudeSession, error) {
	var session core.ClaudeSession
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&session); err != nil {
		if err == io.EOF {
			return nil, errors.New("no input provided")
		}
		return nil, fmt.Errorf("failed to parse Claude session: %w", err)
	}
	return &session, nil
}

func showVersion() {
	fmt.Fprintf(os.Stdout, "ccstatus %s\n", version)

	// Show commit info if available
	if commit != "unknown" {
		commitInfo := commit
		if commitDate != "" {
			commitInfo = fmt.Sprintf("%s (%s)", commit, commitDate)
		}
		fmt.Fprintf(os.Stdout, "  Commit: %s\n", commitInfo)
	}

	// Show build date
	fmt.Fprintf(os.Stdout, "  Built:  %s\n", buildDate)

	// Show Go version
	fmt.Fprintf(os.Stdout, "  Go:     %s\n", runtime.Version())

	// Show platform
	fmt.Fprintf(os.Stdout, "  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func showHelp() {
	fmt.Fprintln(os.Stdout, "ccstatus - Status line generator for Claude Code")
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Usage:")
	fmt.Fprintln(os.Stdout, "  ccstatus             Read from stdin and generate status line")
	fmt.Fprintln(os.Stdout, "  ccstatus help        Show this help message")
	fmt.Fprintln(os.Stdout, "  ccstatus version     Show version information")
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Expected JSON input format:")
	example := core.ClaudeSession{
		SessionID:      "session-123",
		TranscriptPath: "/path/to/transcript.jsonl",
		CWD:            "/current/working/directory",
		Model: core.ModelInfo{
			ID:          "claude-opus-4-1-20250805",
			DisplayName: "Opus 4.1",
		},
		Workspace: core.Workspace{
			CurrentDir: "/workspace/current",
			ProjectDir: "/workspace/project",
		},
		Version: "1.0.89",
		OutputStyle: core.OutputStyle{
			Name: "default",
		},
		Cost: core.CostInfo{
			TotalCostUSD:       0.12345, //nolint:mnd // Test data value
			TotalDurationMs:    60000,   //nolint:mnd // Test data value
			TotalAPIDurationMs: 15000,   //nolint:mnd // Test data value
			TotalLinesAdded:    50,      //nolint:mnd // Test data value
			TotalLinesRemoved:  10,      //nolint:mnd // Test data value
		},
		Exceeds200K: false,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("  ", "  ")
	_ = encoder.Encode(example)
}
