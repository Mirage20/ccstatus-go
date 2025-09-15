package tokenusage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
)

func init() {
	// Self-register with type factory
	core.RegisterProvider(string(Key), New, func() interface{} {
		return &TokenUsage{}
	})
}

// Provider provides token usage by reading from transcript file.
type Provider struct {
	transcriptPath string
}

// New creates a new token usage provider with config.
func New(cfgReader *config.Reader, session *core.ClaudeSession) (core.Provider, core.CacheConfig) {
	// Load provider config with defaults
	cfg := config.GetProvider(cfgReader, "tokenusage", defaultConfig())

	return &Provider{
		transcriptPath: session.TranscriptPath,
	}, cfg.Cache
}

// Key returns the unique identifier for this provider.
func (p *Provider) Key() core.ProviderKey {
	return Key
}

// Provide reads and returns token usage from transcript.
func (p *Provider) Provide(_ context.Context) (interface{}, error) {
	if p.transcriptPath == "" {
		return &TokenUsage{}, nil
	}

	usage := p.readTranscript()
	return usage, nil
}

// readTranscript reads the transcript file and extracts token usage.
func (p *Provider) readTranscript() *TokenUsage {
	zeroUsage := &TokenUsage{
		InputTokens:              0,
		OutputTokens:             0,
		CacheCreationInputTokens: 0,
		CacheReadInputTokens:     0,
	}

	// Check if file exists
	if _, err := os.Stat(p.transcriptPath); os.IsNotExist(err) {
		return zeroUsage
	}

	file, err := os.Open(p.transcriptPath)
	if err != nil {
		return zeroUsage
	}
	defer file.Close()

	// Read all lines into memory (we need to process from the end)
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	// Find the last message with usage data (excluding sidechains)
	// Process from the end of the file
	for i := len(lines) - 1; i >= 0; i-- {
		var entry TranscriptEntry
		if err = json.Unmarshal([]byte(lines[i]), &entry); err != nil {
			// Skip invalid JSON lines
			continue
		}

		// Skip sidechain entries (these are explorations that don't affect main context)
		if entry.IsSidechain {
			continue
		}

		// Look for assistant messages with usage data
		if entry.Type == "assistant" && entry.Message.Usage != nil {
			usage := entry.Message.Usage
			return &TokenUsage{
				InputTokens:              usage.InputTokens,
				OutputTokens:             usage.OutputTokens,
				CacheCreationInputTokens: usage.CacheCreationInputTokens,
				CacheReadInputTokens:     usage.CacheReadInputTokens,
			}
		}
	}

	return zeroUsage
}

// TranscriptEntry represents a line in the transcript file.
type TranscriptEntry struct {
	Type        string  `json:"type"`
	IsSidechain bool    `json:"isSidechain"`
	Message     Message `json:"message"`
}

// Message contains the message data including usage.
type Message struct {
	Usage *Usage `json:"usage"`
}

// Usage contains token usage information.
type Usage struct {
	InputTokens              int64 `json:"input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
}

// TokenUsage represents token consumption.
type TokenUsage struct {
	InputTokens              int64 `json:"input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
}

// Total returns total token count.
func (t *TokenUsage) Total() int64 {
	return t.InputTokens + t.OutputTokens + t.CacheCreationInputTokens + t.CacheReadInputTokens
}

// Key is the provider key.
const Key = core.ProviderKey("tokenusage")

// GetTokenUsage is a typed getter for components.
func GetTokenUsage(ctx *core.RenderContext) (*TokenUsage, bool) {
	return core.Get[*TokenUsage](ctx, Key)
}
