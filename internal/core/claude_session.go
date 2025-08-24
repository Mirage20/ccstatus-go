package core

// ClaudeSession represents the Claude Code session information passed via stdin
// This is NOT a provider - it's the runtime information from Claude Code
type ClaudeSession struct {
	SessionID      string      `json:"session_id"`
	TranscriptPath string      `json:"transcript_path"`
	CWD            string      `json:"cwd"`
	Model          ModelInfo   `json:"model"`
	Workspace      Workspace   `json:"workspace"`
	Version        string      `json:"version"`
	OutputStyle    OutputStyle `json:"output_style"`
	Cost           CostInfo    `json:"cost"`
	Exceeds200K    bool        `json:"exceeds_200k_tokens"`
}

// ModelInfo contains information about the Claude model
type ModelInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// Workspace contains workspace information
type Workspace struct {
	CurrentDir string `json:"current_dir"`
	ProjectDir string `json:"project_dir"`
}

// OutputStyle contains output style information
type OutputStyle struct {
	Name string `json:"name"`
}

// CostInfo contains cost and usage metrics
type CostInfo struct {
	TotalCostUSD       float64 `json:"total_cost_usd"`
	TotalDurationMs    int64   `json:"total_duration_ms"`
	TotalAPIDurationMs int64   `json:"total_api_duration_ms"`
	TotalLinesAdded    int     `json:"total_lines_added"`
	TotalLinesRemoved  int     `json:"total_lines_removed"`
}
