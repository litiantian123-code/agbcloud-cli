// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/config"
)

func TestEndpointPriority(t *testing.T) {
	tests := []struct {
		name           string
		envEndpoint    string
		configEndpoint string
		expectedResult string
	}{
		{
			name:           "environment_variable_priority",
			envEndpoint:    "env.example.com",
			configEndpoint: "config.example.com",
			expectedResult: "https://env.example.com",
		},
		{
			name:           "config_file_fallback",
			envEndpoint:    "",
			configEndpoint: "config.example.com",
			expectedResult: "config.example.com",
		},
		{
			name:           "default_fallback",
			envEndpoint:    "",
			configEndpoint: "",
			expectedResult: "https://agb.cloud",
		},
		{
			name:           "env_with_https_prefix",
			envEndpoint:    "https://secure.example.com",
			configEndpoint: "config.example.com",
			expectedResult: "https://secure.example.com",
		},
		{
			name:           "env_with_http_prefix",
			envEndpoint:    "http://insecure.example.com",
			configEndpoint: "config.example.com",
			expectedResult: "http://insecure.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config directory
			tempDir, err := os.MkdirTemp("", "agbcloud-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Set config directory environment variable
			originalConfigDir := os.Getenv("AGB_CLI_CONFIG_DIR")
			os.Setenv("AGB_CLI_CONFIG_DIR", tempDir)
			defer func() {
				if originalConfigDir == "" {
					os.Unsetenv("AGB_CLI_CONFIG_DIR")
				} else {
					os.Setenv("AGB_CLI_CONFIG_DIR", originalConfigDir)
				}
			}()

			// Set endpoint environment variable
			originalEndpoint := os.Getenv("AGB_CLI_ENDPOINT")
			if tt.envEndpoint != "" {
				os.Setenv("AGB_CLI_ENDPOINT", tt.envEndpoint)
			} else {
				os.Unsetenv("AGB_CLI_ENDPOINT")
			}
			defer func() {
				if originalEndpoint == "" {
					os.Unsetenv("AGB_CLI_ENDPOINT")
				} else {
					os.Setenv("AGB_CLI_ENDPOINT", originalEndpoint)
				}
			}()

			// Create config file with endpoint if specified
			if tt.configEndpoint != "" {
				configPath := filepath.Join(tempDir, "config.json")
				configData := map[string]interface{}{
					"endpoint": tt.configEndpoint,
					"token": map[string]interface{}{
						"loginToken":     "test-login-token",
						"sessionId":      "test-session-id",
						"keepAliveToken": "test-keep-alive-token",
						"expiresAt":      time.Now().Add(time.Hour).Format(time.RFC3339),
					},
				}

				configBytes, err := json.MarshalIndent(configData, "", "  ")
				if err != nil {
					t.Fatalf("Failed to marshal config: %v", err)
				}

				err = os.WriteFile(configPath, configBytes, 0600)
				if err != nil {
					t.Fatalf("Failed to write config file: %v", err)
				}
			}

			// Load config and check endpoint
			cfg, err := config.GetConfig()
			if err != nil {
				t.Fatalf("Failed to get config: %v", err)
			}

			if cfg.Endpoint != tt.expectedResult {
				t.Errorf("Expected endpoint %q, got %q", tt.expectedResult, cfg.Endpoint)
			}

			t.Logf("✅ %s: endpoint correctly set to %q", tt.name, cfg.Endpoint)
		})
	}
}
