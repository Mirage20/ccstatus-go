package ratelimit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiURL        = "https://api.anthropic.com/api/oauth/usage"
	apiBetaHeader = "oauth-2025-04-20"
	apiTimeout    = 5 * time.Second
)

// fetchRateLimits fetches rate limits from the Anthropic OAuth API.
func fetchRateLimits(ctx context.Context, token string) (*RateLimits, error) {
	// Create request with timeout
	ctx, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Anthropic-Beta", apiBetaHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rate limits: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var limits RateLimits
	if err = json.Unmarshal(body, &limits); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &limits, nil
}
