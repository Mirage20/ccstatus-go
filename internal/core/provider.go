package core

import (
	"context"
)

// ProviderKey uniquely identifies a provider
type ProviderKey string

// Provider interface - minimal and focused
type Provider interface {
	// Key returns the unique identifier for this provider
	Key() ProviderKey

	// Provide fetches and returns data
	Provide(ctx context.Context) (interface{}, error)
}
