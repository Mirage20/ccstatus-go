package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mirage20/ccstatus-go/internal/cache"
	"github.com/mirage20/ccstatus-go/internal/components/ccusage/activeblock"
	"github.com/mirage20/ccstatus-go/internal/components/claudecode/changes"
	cccontext "github.com/mirage20/ccstatus-go/internal/components/claudecode/context"
	"github.com/mirage20/ccstatus-go/internal/components/claudecode/duration"
	"github.com/mirage20/ccstatus-go/internal/components/claudecode/model"
	"github.com/mirage20/ccstatus-go/internal/components/claudecode/version"
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	blockusageProvider "github.com/mirage20/ccstatus-go/internal/providers/blockusage"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
	"github.com/mirage20/ccstatus-go/internal/providers/tokenusage"
)

func main() {
	// Handle command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help", "-h", "--help":
			showHelp()
			return
		case "version", "-v", "--version":
			fmt.Println("ccstatus-go v0.1.0")
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

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// Continue with defaults if config loading fails
		cfg = config.Default()
	}

	// Create cache with session isolation
	var c core.Cache
	var fc *cache.FileCache
	if cfg.GetBool("cache.enabled", true) {
		cacheDir := cfg.GetString("cache.dir", os.TempDir())
		fc = cache.NewFileCache(cacheDir, claudeSession.SessionID)
		c = fc

		// Cleanup old cache files occasionally (10% chance)
		if len(os.Args) > 0 && os.Args[0] != "" {
			// Simple random: use last byte of session ID
			if len(claudeSession.SessionID) > 0 && claudeSession.SessionID[len(claudeSession.SessionID)-1]%10 == 0 {
				fc.Cleanup()
			}
		}
	} else {
		c = cache.NewNullCache()
	}

	// Create status line
	statusLine := core.NewStatusLine(cfg, c)

	// Add providers - each provider gets only what it needs from Claude session

	// Session info provider - provides model and session information
	statusLine.AddProvider(sessioninfo.NewProvider(claudeSession))

	// Token usage provider - reads transcript file
	statusLine.AddProvider(tokenusage.NewProvider(claudeSession))

	// Block usage provider - executes ccusage command (with caching)
	blockProvider := blockusageProvider.NewProvider()
	if c != nil {
		// Wrap with caching - 5 seconds TTL for block usage
		cachedBlockProvider := core.NewCachingProvider(blockProvider, c, 10*time.Second)
		statusLine.AddProvider(cachedBlockProvider)
	} else {
		statusLine.AddProvider(blockProvider)
	}

	// TODO: Add git provider when implemented
	// if cfg.GetBool("providers.git.enabled", false) && claudeSession.CWD != "" {
	//     statusLine.AddProvider(git.NewProvider(claudeSession.CWD))
	// }

	// Add components using the component packages
	statusLine.AddComponent(model.New(1))       // Priority 1: Model name
	statusLine.AddComponent(cccontext.New(2))   // Priority 2: Context usage
	statusLine.AddComponent(activeblock.New(3)) // Priority 3: Block usage
	statusLine.AddComponent(changes.New(4))     // Priority 4: Lines changed
	statusLine.AddComponent(duration.New(5))    // Priority 5: Session duration
	statusLine.AddComponent(version.New(6))     // Priority 6: Claude version

	// TODO: Add git component when implemented
	// if cfg.GetBool("components.git.enabled", false) {
	//     statusLine.AddComponent(git.New(4))
	// }

	// Render and output the status line
	output := statusLine.Render(ctx)
	fmt.Println(output)

	// Save cache if using FileCache
	if fc != nil {
		fc.Save()
	}

	return nil
}

// readClaudeSession reads the Claude session information from stdin
func readClaudeSession(reader io.Reader) (*core.ClaudeSession, error) {
	var session core.ClaudeSession
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&session); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("no input provided")
		}
		return nil, fmt.Errorf("failed to parse Claude session: %w", err)
	}
	return &session, nil
}

func showHelp() {
	fmt.Println("ccstatus-go - Status line generator for Claude Code")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ccstatus             Read from stdin and generate status line")
	fmt.Println("  ccstatus help        Show this help message")
	fmt.Println("  ccstatus version     Show version information")
	fmt.Println()
	fmt.Println("Expected JSON input format:")
	example := core.ClaudeSession{
		SessionID:      "session-123",
		TranscriptPath: "/path/to/transcript.json",
		CWD:            "/current/working/directory",
		Model: core.ModelInfo{
			ID:          "claude-3-opus-20240229",
			DisplayName: "Claude 3 Opus",
		},
		Workspace: core.Workspace{
			CurrentDir: "/workspace/current",
			ProjectDir: "/workspace/project",
		},
		Version: "1.0.0",
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("  ", "  ")
	encoder.Encode(example)
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  CCSTATUS_CACHE_DIR   Override cache directory")
	fmt.Println("  CCSTATUS_NO_CACHE    Set to 1 to disable caching")
	fmt.Println("  CLAUDE_SESSION_ID    Session-specific cache")
}
