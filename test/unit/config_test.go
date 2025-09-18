// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"os"
	"testing"

	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// TestGetEndpoint tests the GetEndpoint function
func TestGetEndpoint(t *testing.T) {
	// Save original environment variables
	originalCLIEndpoint := os.Getenv("AGB_CLI_ENDPOINT")

	// Clean up after test
	defer func() {
		os.Setenv("AGB_CLI_ENDPOINT", originalCLIEndpoint)
	}()

	t.Run("DefaultValues", func(t *testing.T) {
		// Clear all environment variables
		os.Unsetenv("AGB_CLI_ENDPOINT")

		endpoint := config.GetEndpoint()

		if endpoint != "https://agb.cloud" {
			t.Errorf("Expected default endpoint https://agb.cloud, got %s", endpoint)
		}

		t.Logf("[OK] Default values test passed")
	})

	t.Run("CLIEnvironmentVariables", func(t *testing.T) {
		// Set CLI-specific environment variables
		os.Setenv("AGB_CLI_ENDPOINT", "cli.agb.cloud")

		endpoint := config.GetEndpoint()

		if endpoint != "https://cli.agb.cloud" {
			t.Errorf("Expected CLI endpoint https://cli.agb.cloud, got %s", endpoint)
		}

		t.Logf("[OK] CLI environment variables test passed")
	})

	t.Run("EndpointAutoHTTPS", func(t *testing.T) {
		// Test various endpoint formats
		testCases := []struct {
			input    string
			expected string
		}{
			{"agb.cloud", "https://agb.cloud"},
			{"custom.domain.com", "https://custom.domain.com"},
			{"https://already.has.protocol", "https://already.has.protocol"},
			{"http://insecure.domain", "http://insecure.domain"},
			{"12.34.56.78", "https://12.34.56.78"},
			{"localhost:8080", "https://localhost:8080"},
		}

		for _, tc := range testCases {
			os.Setenv("AGB_CLI_ENDPOINT", tc.input)
			endpoint := config.GetEndpoint()

			if endpoint != tc.expected {
				t.Errorf("Input %s: expected %s, got %s", tc.input, tc.expected, endpoint)
			} else {
				t.Logf("[OK] %s -> %s", tc.input, endpoint)
			}
		}
	})

	t.Run("EmptyEnvironmentVariables", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("AGB_CLI_ENDPOINT")

		endpoint := config.GetEndpoint()

		if endpoint != "https://agb.cloud" {
			t.Errorf("Expected default endpoint https://agb.cloud, got %s", endpoint)
		}

		t.Logf("[OK] Empty environment variables test passed")
	})
}

// TestConfigTokenOperations tests token-related config operations
func TestConfigTokenOperations(t *testing.T) {
	cfg := &config.Config{}

	// Test initial state
	if cfg.IsAuthenticated() {
		t.Error("Expected config to not be authenticated initially")
	}

	// Test token operations
	err := cfg.SaveTokens("test-login", "test-session", "test-keepalive", "2025-12-31T23:59:59Z")
	if err != nil {
		t.Errorf("Failed to save tokens: %v", err)
	}

	if !cfg.IsAuthenticated() {
		t.Error("Expected config to be authenticated after saving tokens")
	}

	tokens, err := cfg.GetTokens()
	if err != nil {
		t.Errorf("Failed to get tokens: %v", err)
	}

	if tokens.LoginToken != "test-login" {
		t.Errorf("Expected login token 'test-login', got %s", tokens.LoginToken)
	}

	// Test clearing tokens
	err = cfg.ClearTokens()
	if err != nil {
		t.Errorf("Failed to clear tokens: %v", err)
	}

	if cfg.IsAuthenticated() {
		t.Error("Expected config to not be authenticated after clearing tokens")
	}

	t.Logf("[OK] Token operations test passed")
}
