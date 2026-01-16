package ratelimit

import (
	"encoding/json"
	"time"
)

// RateLimits represents the rate limit data from Anthropic OAuth API.
type RateLimits struct {
	FiveHour *RateLimit `json:"five_hour"`
	SevenDay *RateLimit `json:"seven_day"`
}

// RateLimit represents a single rate limit window.
type RateLimit struct {
	Utilization float64    `json:"utilization"`
	ResetsAt    *time.Time `json:"-"` // Parsed from string, use custom unmarshal
}

// rateLimitJSON is used for JSON unmarshaling with string time.
type rateLimitJSON struct {
	Utilization float64 `json:"utilization"`
	ResetsAt    *string `json:"resets_at"`
}

// UnmarshalJSON implements custom JSON unmarshaling to parse ISO 8601 timestamps.
func (r *RateLimit) UnmarshalJSON(data []byte) error {
	var aux rateLimitJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.Utilization = aux.Utilization

	if aux.ResetsAt != nil && *aux.ResetsAt != "" {
		// Parse ISO 8601 with fractional seconds and timezone
		t, err := time.Parse(time.RFC3339Nano, *aux.ResetsAt)
		if err != nil {
			// Try without nanoseconds
			t, err = time.Parse(time.RFC3339, *aux.ResetsAt)
			if err != nil {
				return nil //nolint:nilerr // Intentional: don't fail on parse error, just leave ResetsAt nil
			}
		}
		r.ResetsAt = &t
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling to output ISO 8601 timestamps.
func (r *RateLimit) MarshalJSON() ([]byte, error) {
	aux := rateLimitJSON{
		Utilization: r.Utilization,
	}

	if r.ResetsAt != nil {
		s := r.ResetsAt.Format(time.RFC3339Nano)
		aux.ResetsAt = &s
	}

	return json.Marshal(aux)
}

// Credentials represents the Claude Code credentials JSON structure.
type Credentials struct {
	ClaudeAiOauth *OAuthCredentials `json:"claudeAiOauth"`
}

// OAuthCredentials contains the OAuth access token.
type OAuthCredentials struct {
	AccessToken string `json:"accessToken"`
}
