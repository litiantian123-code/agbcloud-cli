// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// runLogoutSimple simulates the logout command logic for testing
func runLogoutSimple(cfg *config.Config) error {
	// Check if we have valid tokens for API logout
	hasValidTokens := cfg.Token != nil &&
		cfg.Token.LoginToken != "" &&
		cfg.Token.SessionId != ""

	if hasValidTokens {
		// Simulate API call (in real implementation, this would call the API)
		fmt.Println("🌐 Invalidating server session...")
		fmt.Println("✅ Server session invalidated successfully")
	} else {
		fmt.Println("ℹ️  No active session found")
	}

	// Always perform local cleanup
	fmt.Println("🧹 Clearing local authentication data...")

	// Clear tokens from config
	err := cfg.ClearTokens()
	if err != nil {
		return fmt.Errorf("failed to clear local authentication data: %w", err)
	}

	fmt.Println("✅ Successfully logged out from AgbCloud")
	return nil
}

func TestLogoutCommandSimple(t *testing.T) {
	tests := []struct {
		name      string
		hasTokens bool
		wantOut   []string
	}{
		{
			name:      "logout with tokens",
			hasTokens: true,
			wantOut: []string{
				"🌐 Invalidating server session...",
				"✅ Server session invalidated successfully",
				"🧹 Clearing local authentication data...",
				"✅ Successfully logged out from AgbCloud",
			},
		},
		{
			name:      "logout without tokens",
			hasTokens: false,
			wantOut: []string{
				"ℹ️  No active session found",
				"🧹 Clearing local authentication data...",
				"✅ Successfully logged out from AgbCloud",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for config
			tempDir, err := os.MkdirTemp("", "agbcloud-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Set config directory
			os.Setenv("AGB_CLI_CONFIG_DIR", tempDir)
			defer os.Unsetenv("AGB_CLI_CONFIG_DIR")

			configPath := filepath.Join(tempDir, "config.json")

			// Create test config
			var cfg *config.Config
			if tt.hasTokens {
				cfg = &config.Config{
					Endpoint: "https://test.agb.cloud",
					Token: &config.Token{
						LoginToken:     "test-login-token",
						SessionId:      "test-session-id",
						KeepAliveToken: "test-keep-alive-token",
						ExpiresAt:      time.Now().Add(24 * time.Hour),
					},
				}
			} else {
				cfg = &config.Config{
					Endpoint: "https://test.agb.cloud",
					Token:    nil, // No tokens
				}
			}

			// Write config to file
			configData, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal config: %v", err)
			}
			err = os.WriteFile(configPath, configData, 0600)
			if err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			// Capture output
			var buf bytes.Buffer
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run logout command
			err = runLogoutSimple(cfg)

			// Restore stdout and get output
			w.Close()
			os.Stdout = oldStdout
			io.Copy(&buf, r)
			output := buf.String()

			if err != nil {
				t.Errorf("runLogoutSimple() error = %v", err)
				return
			}

			// Check output contains expected strings
			for _, expectedOut := range tt.wantOut {
				if !strings.Contains(output, expectedOut) {
					t.Errorf("Expected output to contain %q, got:\n%s", expectedOut, output)
				}
			}

			// Verify tokens are cleared
			updatedCfg, err := config.GetConfig()
			if err != nil {
				t.Errorf("Failed to load updated config: %v", err)
				return
			}

			if updatedCfg.Token != nil {
				t.Errorf("Expected tokens to be cleared, but they still exist")
			}
		})
	}
}

func TestLogoutCommandNoConfig(t *testing.T) {
	// Create temporary directory for config
	tempDir, err := os.MkdirTemp("", "agbcloud-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set config directory
	os.Setenv("AGB_CLI_CONFIG_DIR", tempDir)
	defer os.Unsetenv("AGB_CLI_CONFIG_DIR")

	// Create empty config
	cfg := &config.Config{
		Endpoint: "https://test.agb.cloud",
		Token:    nil,
	}

	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run logout command
	err = runLogoutSimple(cfg)

	// Restore stdout and get output
	w.Close()
	os.Stdout = oldStdout
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Errorf("runLogoutSimple() error = %v", err)
		return
	}

	// Should indicate no active session
	if !strings.Contains(output, "ℹ️  No active session found") {
		t.Errorf("Expected output to contain 'No active session found', got:\n%s", output)
	}

	// Should still show success message
	if !strings.Contains(output, "✅ Successfully logged out from AgbCloud") {
		t.Errorf("Expected output to contain success message, got:\n%s", output)
	}
}
