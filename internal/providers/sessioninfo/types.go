package sessioninfo

import "github.com/mirage20/ccstatus-go/internal/core"

// SessionInfo represents the session data provided by this provider.
type SessionInfo struct {
	Model       core.ModelInfo `json:"model"`
	SessionID   string         `json:"session_id"`
	Version     string         `json:"version"`
	Workspace   core.Workspace `json:"workspace"`
	Cost        core.CostInfo  `json:"cost"`
	Exceeds200K bool           `json:"exceeds_200k"`
}
