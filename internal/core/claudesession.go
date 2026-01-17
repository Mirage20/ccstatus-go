package core

// ClaudeSession represents the Claude Code session information passed via stdin
// This is NOT a provider - it's the runtime information from Claude Code.
type ClaudeSession struct {
	HookEventName  string        `json:"hook_event_name,omitempty"`
	SessionID      string        `json:"session_id"`
	TranscriptPath string        `json:"transcript_path"`
	CWD            string        `json:"cwd"`
	Model          ModelInfo     `json:"model"`
	Workspace      Workspace     `json:"workspace"`
	Version        string        `json:"version"`
	OutputStyle    OutputStyle   `json:"output_style"`
	Cost           CostInfo      `json:"cost"`
	ContextWindow  ContextWindow `json:"context_window"`
	Exceeds200K    bool          `json:"exceeds_200k_tokens"`
}

// ContextWindow contains context window information from Claude Code.
type ContextWindow struct {
	// Cumulative totals across the entire session (not current context usage)
	TotalInputTokens  int64 `json:"total_input_tokens"`
	TotalOutputTokens int64 `json:"total_output_tokens"`
	// Maximum context window size for the model
	ContextWindowSize int64 `json:"context_window_size"`
	// Current context window usage from the last API call (maybe nil if no messages yet)
	CurrentUsage *ContextUsage `json:"current_usage,omitempty"`
	// Pre-calculated percentages (maybe nil if no messages yet)
	UsedPercentage      *float64 `json:"used_percentage,omitempty"`
	RemainingPercentage *float64 `json:"remaining_percentage,omitempty"`
}

// ContextUsage contains the current context window usage from the last API call.
type ContextUsage struct {
	InputTokens              int64 `json:"input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
}

// ModelInfo contains information about the Claude model.
type ModelInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// Workspace contains workspace information.
type Workspace struct {
	CurrentDir string `json:"current_dir"`
	ProjectDir string `json:"project_dir"`
}

// OutputStyle contains output style information.
type OutputStyle struct {
	Name string `json:"name"`
}

// CostInfo contains cost and usage metrics.
type CostInfo struct {
	TotalCostUSD       float64 `json:"total_cost_usd"`
	TotalDurationMs    int64   `json:"total_duration_ms"`
	TotalAPIDurationMs int64   `json:"total_api_duration_ms"`
	TotalLinesAdded    int     `json:"total_lines_added"`
	TotalLinesRemoved  int     `json:"total_lines_removed"`
}
