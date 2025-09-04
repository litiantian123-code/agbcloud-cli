// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"os"
	"testing"

	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// TestDefaultConfigEnvironmentVariables tests the environment variable configuration
func TestDefaultConfigEnvironmentVariables(t *testing.T) {
	// Save original environment variables
	originalCLIEndpoint := os.Getenv("AGB_CLI_ENDPOINT")

	// Clean up after test
	defer func() {
		os.Setenv("AGB_CLI_ENDPOINT", originalCLIEndpoint)
	}()

	t.Run("DefaultValues", func(t *testing.T) {
		// Clear all environment variables
		os.Unsetenv("AGB_CLI_ENDPOINT")

		cfg := config.DefaultConfig()

		if cfg.Endpoint != "https://agb.cloud" {
			t.Errorf("Expected default endpoint https://agb.cloud, got %s", cfg.Endpoint)
		}

		t.Logf("✅ Default values test passed")
	})

	t.Run("CLIEnvironmentVariables", func(t *testing.T) {
		// Set CLI-specific environment variables
		os.Setenv("AGB_CLI_ENDPOINT", "cli.agb.cloud")

		cfg := config.DefaultConfig()

		if cfg.Endpoint != "https://cli.agb.cloud" {
			t.Errorf("Expected CLI endpoint https://cli.agb.cloud, got %s", cfg.Endpoint)
		}

		t.Logf("✅ CLI environment variables test passed")
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
			cfg := config.DefaultConfig()

			if cfg.Endpoint != tc.expected {
				t.Errorf("Input %s: expected %s, got %s", tc.input, tc.expected, cfg.Endpoint)
			} else {
				t.Logf("✅ %s -> %s", tc.input, cfg.Endpoint)
			}
		}
	})

	t.Run("EmptyEnvironmentVariables", func(t *testing.T) {
		// Set empty environment variables
		os.Setenv("AGB_CLI_ENDPOINT", "")

		cfg := config.DefaultConfig()

		// Should use default when environment variable is empty
		if cfg.Endpoint != "https://agb.cloud" {
			t.Errorf("Expected default endpoint https://agb.cloud, got %s", cfg.Endpoint)
		}

		t.Logf("✅ Empty environment variables test passed")
	})
}

// TestConfigPaths tests configuration file paths
func TestConfigPaths(t *testing.T) {
	t.Run("ConfigDir", func(t *testing.T) {
		configDir, err := config.ConfigDir()
		if err != nil {
			t.Fatalf("Failed to get config directory: %v", err)
		}

		if configDir == "" {
			t.Error("Config directory should not be empty")
		}

		t.Logf("Config directory: %s", configDir)
	})

	t.Run("ConfigFile", func(t *testing.T) {
		configFile, err := config.ConfigFile()
		if err != nil {
			t.Fatalf("Failed to get config file path: %v", err)
		}

		if configFile == "" {
			t.Error("Config file path should not be empty")
		}

		t.Logf("Config file: %s", configFile)
	})
}
