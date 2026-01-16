package ratelimit

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	keychainService   = "Claude Code-credentials"
	credentialsFile   = ".claude/.credentials.json" //nolint:gosec // Not credentials, just a file path
	credentialTimeout = 2 * time.Second
)

// getOAuthToken retrieves the OAuth access token from platform-specific credential stores.
// Priority: Keychain/secret-tool → file fallback.
func getOAuthToken(ctx context.Context) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		token, err := getTokenFromMacKeychain(ctx)
		if err == nil && token != "" {
			return token, nil
		}
		// Fallback to file
		return getTokenFromFile()

	case "linux":
		token, err := getTokenFromLinuxKeyring(ctx)
		if err == nil && token != "" {
			return token, nil
		}
		// Fallback to file
		return getTokenFromFile()

	default:
		// Windows and others: file only
		return getTokenFromFile()
	}
}

// getTokenFromMacKeychain retrieves credentials from macOS Keychain.
func getTokenFromMacKeychain(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, credentialTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "security", "find-generic-password", "-s", keychainService, "-w")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return parseCredentialsJSON(strings.TrimSpace(string(output)))
}

// getTokenFromLinuxKeyring retrieves credentials from GNOME Keyring via secret-tool.
func getTokenFromLinuxKeyring(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, credentialTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "secret-tool", "lookup", "service", keychainService)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return parseCredentialsJSON(strings.TrimSpace(string(output)))
}

// getTokenFromFile retrieves credentials from ~/.claude/.credentials.json.
func getTokenFromFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	credPath := filepath.Join(homeDir, credentialsFile)
	data, err := os.ReadFile(credPath)
	if err != nil {
		return "", err
	}

	return parseCredentialsJSON(string(data))
}

// parseCredentialsJSON extracts the access token from credentials JSON.
func parseCredentialsJSON(jsonData string) (string, error) {
	var creds Credentials
	if err := json.Unmarshal([]byte(jsonData), &creds); err != nil {
		return "", err
	}

	if creds.ClaudeAiOauth != nil && creds.ClaudeAiOauth.AccessToken != "" {
		return creds.ClaudeAiOauth.AccessToken, nil
	}

	return "", nil
}
