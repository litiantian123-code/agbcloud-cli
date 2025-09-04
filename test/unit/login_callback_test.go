// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/auth"
)

// TestCallbackServerPortConfiguration tests port configuration for callback server
func TestCallbackServerPortConfiguration(t *testing.T) {
	tests := []struct {
		name         string
		envVar       string
		configPort   string
		expectedPort string
	}{
		{
			name:         "DefaultPort",
			envVar:       "",
			configPort:   "",
			expectedPort: "3000",
		},
		{
			name:         "EnvironmentVariablePort",
			envVar:       "8080",
			configPort:   "",
			expectedPort: "8080",
		},
		{
			name:         "ConfigPort",
			envVar:       "",
			configPort:   "9000",
			expectedPort: "9000",
		},
		{
			name:         "EnvironmentVariableOverridesConfig",
			envVar:       "7000",
			configPort:   "9000",
			expectedPort: "7000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env var
			originalEnv := os.Getenv("AGB_CLI_CALLBACK_PORT")
			defer func() {
				if originalEnv != "" {
					os.Setenv("AGB_CLI_CALLBACK_PORT", originalEnv)
				} else {
					os.Unsetenv("AGB_CLI_CALLBACK_PORT")
				}
			}()

			// Set test environment variable
			if tt.envVar != "" {
				os.Setenv("AGB_CLI_CALLBACK_PORT", tt.envVar)
			} else {
				os.Unsetenv("AGB_CLI_CALLBACK_PORT")
			}

			// Test port resolution logic using auth package
			port := auth.GetCallbackPort(tt.configPort)
			if port != tt.expectedPort {
				t.Errorf("Expected port %s, got %s", tt.expectedPort, port)
			}
		})
	}
}

// TestCallbackServerStartStop tests starting and stopping the callback server
func TestCallbackServerStartStop(t *testing.T) {
	port := "3001" // Use a different port to avoid conflicts

	// Test server lifecycle
	t.Run("ServerStartStop", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Start server in background
		codeChan := make(chan string, 1)
		errChan := make(chan error, 1)

		go func() {
			code, err := auth.StartCallbackServer(ctx, port)
			if err != nil {
				errChan <- err
				return
			}
			codeChan <- code
		}()

		// Give server time to start
		time.Sleep(100 * time.Millisecond)

		// Test that server is listening
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/callback?code=test-code", port))
		if err != nil {
			t.Fatalf("Failed to connect to callback server: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Check that we received the code
		select {
		case code := <-codeChan:
			if code != "test-code" {
				t.Errorf("Expected code 'test-code', got '%s'", code)
			}
		case err := <-errChan:
			t.Fatalf("Server returned error: %v", err)
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for callback")
		}
	})
}

// TestCallbackServerWithState tests that callback server ignores state parameter
func TestCallbackServerWithState(t *testing.T) {
	port := "3002"

	t.Run("IgnoreState", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		codeChan := make(chan string, 1)
		errChan := make(chan error, 1)

		go func() {
			code, err := auth.StartCallbackServer(ctx, port)
			if err != nil {
				errChan <- err
				return
			}
			codeChan <- code
		}()

		// Give server time to start
		time.Sleep(100 * time.Millisecond)

		// Send request with any state - should be ignored
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/callback?code=test-code&state=any-state", port))
		if err != nil {
			t.Fatalf("Failed to connect to callback server: %v", err)
		}
		defer resp.Body.Close()

		// Should succeed with 200 status (state is ignored)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Should receive the code
		select {
		case code := <-codeChan:
			if code != "test-code" {
				t.Errorf("Expected code 'test-code', got '%s'", code)
			}
		case err := <-errChan:
			t.Fatalf("Server returned error: %v", err)
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for callback")
		}
	})
}

// TestCallbackServerMissingCode tests handling of missing code parameter
func TestCallbackServerMissingCode(t *testing.T) {
	port := "3003"

	t.Run("MissingCode", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		errChan := make(chan error, 1)

		go func() {
			_, err := auth.StartCallbackServer(ctx, port)
			errChan <- err
		}()

		// Give server time to start
		time.Sleep(100 * time.Millisecond)

		// Send request without code
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/callback", port))
		if err != nil {
			t.Fatalf("Failed to connect to callback server: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		// Check that server returned error
		select {
		case err := <-errChan:
			if err == nil {
				t.Error("Expected error for missing code, got nil")
			}
			if !strings.Contains(err.Error(), "no code") {
				t.Errorf("Expected 'no code' error, got: %v", err)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for error")
		}
	})
}

// TestGenerateRandomState tests the random state generation
func TestGenerateRandomState(t *testing.T) {
	state1, err := auth.GenerateRandomState()
	if err != nil {
		t.Fatalf("Failed to generate random state: %v", err)
	}

	state2, err := auth.GenerateRandomState()
	if err != nil {
		t.Fatalf("Failed to generate second random state: %v", err)
	}

	// States should be different
	if state1 == state2 {
		t.Error("Generated states should be different")
	}

	// States should not be empty
	if state1 == "" || state2 == "" {
		t.Error("Generated states should not be empty")
	}

	// States should be base64 encoded (basic check)
	if len(state1) < 10 || len(state2) < 10 {
		t.Error("Generated states seem too short")
	}
}
