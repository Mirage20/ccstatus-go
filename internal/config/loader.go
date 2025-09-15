package config

import (
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Reader provides access to configuration values.
type Reader struct {
	k *koanf.Koanf
}

// NewReader creates a new configuration reader.
func NewReader(projectDir string) *Reader {
	k := koanf.New(".")

	// Load user config file if it exists
	configPath := findConfigFile(projectDir)
	if configPath != "" {
		// Try to load config, but continue with defaults if it fails
		_ = k.Load(file.Provider(configPath), yaml.Parser())
	}

	return &Reader{k: k}
}

// findConfigFile searches for config file in order of preference.
func findConfigFile(projectDir string) string {
	// Project-specific configs (using project dir from Claude session)
	if projectDir != "" {
		// Check for local config first (highest priority)
		localConfig := filepath.Join(projectDir, ".claude", "ccstatus.local.yaml")
		if fileExists(localConfig) {
			return localConfig
		}

		// Check for shared project config
		projectConfig := filepath.Join(projectDir, ".claude", "ccstatus.yaml")
		if fileExists(projectConfig) {
			return projectConfig
		}
	}

	// User default fallback
	if home, err := os.UserHomeDir(); err == nil {
		userConfig := filepath.Join(home, ".claude", "ccstatus.yaml")
		if fileExists(userConfig) {
			return userConfig
		}
	}

	return ""
}

// fileExists checks if a file exists and is readable.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Get retrieves configuration at the specified path, returning defaultValue if not found.
func Get[T any](r *Reader, path string, defaultValue T) T {
	if r.k == nil || path == "" || !r.k.Exists(path) {
		return defaultValue
	}

	result := defaultValue
	_ = r.k.UnmarshalWithConf(path, &result, koanf.UnmarshalConf{Tag: "yaml"})
	return result
}

// GetComponent retrieves configuration for a specific component.
func GetComponent[T any](r *Reader, componentName string, defaultValue T) T {
	path := "components." + componentName
	return Get(r, path, defaultValue)
}

// GetProvider retrieves configuration for a specific provider.
func GetProvider[T any](r *Reader, providerName string, defaultValue T) T {
	path := "providers." + providerName
	return Get(r, path, defaultValue)
}
