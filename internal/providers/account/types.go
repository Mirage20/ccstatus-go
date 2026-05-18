package account

// Account represents which Claude Code login is currently active.
type Account struct {
	// Profile is the matched profile name (e.g. "personal", "work"),
	// or "unknown" if the active token doesn't match any configured slot,
	// or "" if no active token was found.
	Profile string `json:"profile"`
}
