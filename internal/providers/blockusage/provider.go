package blockusage

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
)

func init() {
	// Self-register with type factory
	core.RegisterProvider(string(Key), New, func() interface{} {
		return &BlockUsage{}
	})
}

// Provider provides block usage by executing ccusage command.
type Provider struct {
	// This provider doesn't need ClaudeSession!
}

// New creates a new block usage provider with config.
func New(cfgReader *config.Reader, _ *core.ClaudeSession) (core.Provider, core.CacheConfig) {
	// Load provider config with defaults
	cfg := config.GetProvider(cfgReader, "blockusage", defaultConfig())

	return &Provider{}, cfg.Cache
}

// Key returns the unique identifier for this provider.
func (p *Provider) Key() core.ProviderKey {
	return Key
}

// Provide executes ccusage and returns block usage data.
func (p *Provider) Provide(ctx context.Context) (interface{}, error) {
	usage := p.getActiveBlockUsage(ctx)
	return usage, nil
}

// getActiveBlockUsage executes ccusage command and parses the result.
func (p *Provider) getActiveBlockUsage(ctx context.Context) *BlockUsage {
	zeroBlockUsage := &BlockUsage{
		InputTokens:              0,
		OutputTokens:             0,
		CacheCreationInputTokens: 0,
		CacheReadInputTokens:     0,
		TotalTokens:              0,
		RemainingMinutes:         0,
		EndTime:                  "",
		MaxBlockTokens:           0,
	}

	// Execute ccusage command (get all blocks in JSON format to calculate max)
	cmd := exec.CommandContext(ctx, "ccusage", "blocks", "-j")
	output, err := cmd.Output()
	if err != nil {
		// ccusage might not be installed or available
		return zeroBlockUsage
	}

	// Parse the JSON output
	var data CCUsageOutput
	if err = json.Unmarshal(output, &data); err != nil {
		return zeroBlockUsage
	}

	// Return zero if no blocks exist
	if len(data.Blocks) == 0 {
		return zeroBlockUsage
	}

	var activeBlock *Block
	for _, block := range data.Blocks {
		if block.IsActive {
			activeBlock = &block
			break
		}
	}

	// If no active block, return zero usage (we're in a gap or no session)
	if activeBlock == nil {
		return zeroBlockUsage
	}

	// Calculate dynamic block limit from historical maximum
	var maxBlockTokens int64
	for _, block := range data.Blocks {
		if !block.IsGap && block.TotalTokens > maxBlockTokens {
			maxBlockTokens = block.TotalTokens
		}
	}

	// If no historical data or all gaps, use a reasonable default
	if maxBlockTokens == 0 {
		maxBlockTokens = 1000000 // 1M tokens as fallback
	}

	return &BlockUsage{
		InputTokens:              activeBlock.TokenCounts.InputTokens,
		OutputTokens:             activeBlock.TokenCounts.OutputTokens,
		CacheCreationInputTokens: activeBlock.TokenCounts.CacheCreationInputTokens,
		CacheReadInputTokens:     activeBlock.TokenCounts.CacheReadInputTokens,
		TotalTokens:              activeBlock.TotalTokens,
		RemainingMinutes:         activeBlock.Projection.RemainingMinutes,
		EndTime:                  activeBlock.EndTime,
		MaxBlockTokens:           maxBlockTokens,
	}
}

// CCUsageOutput represents the JSON output from ccusage command.
type CCUsageOutput struct {
	Blocks []Block `json:"blocks"`
}

// Block represents a 5-hour usage block.
type Block struct {
	IsActive    bool        `json:"isActive"`
	IsGap       bool        `json:"isGap"`
	TotalTokens int64       `json:"totalTokens"`
	TokenCounts TokenCounts `json:"tokenCounts"`
	Projection  Projection  `json:"projection"`
	EndTime     string      `json:"endTime"`
}

// TokenCounts contains detailed token counts.
type TokenCounts struct {
	InputTokens              int64 `json:"inputTokens"`
	OutputTokens             int64 `json:"outputTokens"`
	CacheCreationInputTokens int64 `json:"cacheCreationInputTokens"`
	CacheReadInputTokens     int64 `json:"cacheReadInputTokens"`
}

// Projection contains usage projection data.
type Projection struct {
	RemainingMinutes int `json:"remainingMinutes"`
}

// BlockUsage represents 5-hour block usage.
type BlockUsage struct {
	InputTokens              int64  `json:"input_tokens"`
	OutputTokens             int64  `json:"output_tokens"`
	CacheCreationInputTokens int64  `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64  `json:"cache_read_input_tokens"`
	TotalTokens              int64  `json:"total_tokens"`
	RemainingMinutes         int    `json:"remaining_minutes"`
	EndTime                  string `json:"end_time"`
	MaxBlockTokens           int64  `json:"max_block_tokens"`
}

// Key is the provider key.
const Key = core.ProviderKey("blockusage")

// GetBlockUsage is a typed getter for components.
func GetBlockUsage(ctx *core.RenderContext) (*BlockUsage, bool) {
	return core.Get[*BlockUsage](ctx, Key)
}
