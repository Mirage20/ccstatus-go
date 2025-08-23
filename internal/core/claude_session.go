package core

// ClaudeSession represents the Claude Code session information passed via stdin
// This is NOT a provider - it's the runtime information from Claude Code
type ClaudeSession struct {
	SessionID      string    `json:"session_id"`
	TranscriptPath string    `json:"transcript_path"`
	CWD            string    `json:"cwd"`
	Model          ModelInfo `json:"model"`
	Workspace      Workspace `json:"workspace"`
	Version        string    `json:"version"`
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