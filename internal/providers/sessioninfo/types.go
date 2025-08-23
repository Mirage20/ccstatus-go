package sessioninfo

import "github.com/mirage20/ccstatus-go/internal/core"

// SessionInfo represents the session data provided by this provider
type SessionInfo struct {
	Model     core.ModelInfo
	SessionID string
	Version   string
	Workspace core.Workspace
}