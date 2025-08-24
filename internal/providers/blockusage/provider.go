package blockusage

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/mirage20/ccstatus-go/internal/core"
)

// Provider provides block usage by executing ccusage command
type Provider struct {
	// This provider doesn't need ClaudeSession!
}

// NewProvider creates a new block usage provider
func NewProvider() *Provider {
	return &Provider{}
}

// Key returns the unique identifier for this provider
func (p *Provider) Key() core.ProviderKey {
	return Key
}

// Provide executes ccusage and returns block usage data
func (p *Provider) Provide(ctx context.Context) (interface{}, error) {
	usage := p.getActiveBlockUsage()
	return usage, nil
}

// getActiveBlockUsage executes ccusage command and parses the result
func (p *Provider) getActiveBlockUsage() *BlockUsage {
	zeroBlockUsage := &BlockUsage{
		InputTokens:              0,
		OutputTokens:             0,
		CacheCreationInputTokens: 0,
		CacheReadInputTokens:     0,
		TotalTokens:              0,
		RemainingMinutes:         0,
		EndTime:                  "",
		UsagePercentage:          0,
	}

	// Execute ccusage command
	cmd := exec.Command("ccusage", "blocks", "-aj")
	output, err := cmd.Output()
	if err != nil {
		// ccusage might not be installed or available
		return zeroBlockUsage
	}

	// Parse the JSON output
	var data CCUsageOutput
	if err := json.Unmarshal(output, &data); err != nil {
		return zeroBlockUsage
	}

	// Find the active block or use the first one
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

	// If no active block, use the first one
	if activeBlock == nil {
		activeBlock = &data.Blocks[0]
	}

	// Calculate usage percentage (assuming 30M tokens per 5-hour block)
	const maxBlockTokens = 30000000
	usagePercentage := float64(activeBlock.TotalTokens) / float64(maxBlockTokens) * 100

	return &BlockUsage{
		InputTokens:              activeBlock.TokenCounts.InputTokens,
		OutputTokens:             activeBlock.TokenCounts.OutputTokens,
		CacheCreationInputTokens: activeBlock.TokenCounts.CacheCreationInputTokens,
		CacheReadInputTokens:     activeBlock.TokenCounts.CacheReadInputTokens,
		TotalTokens:              activeBlock.TotalTokens,
		RemainingMinutes:         activeBlock.Projection.RemainingMinutes,
		EndTime:                  activeBlock.EndTime,
		UsagePercentage:          usagePercentage,
	}
}

// CCUsageOutput represents the JSON output from ccusage command
type CCUsageOutput struct {
	Blocks []Block `json:"blocks"`
}

// Block represents a 5-hour usage block
type Block struct {
	IsActive    bool        `json:"isActive"`
	TotalTokens int64       `json:"totalTokens"`
	TokenCounts TokenCounts `json:"tokenCounts"`
	Projection  Projection  `json:"projection"`
	EndTime     string      `json:"endTime"`
}

// TokenCounts contains detailed token counts
type TokenCounts struct {
	InputTokens              int64 `json:"inputTokens"`
	OutputTokens             int64 `json:"outputTokens"`
	CacheCreationInputTokens int64 `json:"cacheCreationInputTokens"`
	CacheReadInputTokens     int64 `json:"cacheReadInputTokens"`
}

// Projection contains usage projection data
type Projection struct {
	RemainingMinutes int `json:"remainingMinutes"`
}

// BlockUsage represents 5-hour block usage
type BlockUsage struct {
	InputTokens              int64
	OutputTokens             int64
	CacheCreationInputTokens int64
	CacheReadInputTokens     int64
	TotalTokens              int64
	RemainingMinutes         int
	EndTime                  string
	UsagePercentage          float64
}

// Key is the provider key
const Key = core.ProviderKey("blockusage")

// GetBlockUsage is a typed getter for components
func GetBlockUsage(ctx *core.RenderContext) (*BlockUsage, bool) {
	return core.Get[*BlockUsage](ctx, Key)
}
