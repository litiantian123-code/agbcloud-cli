// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// TestLogoutIntegration tests the logout API endpoint with real server
func TestLogoutIntegration(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=true to run.")
	}

	// Get configuration
	cfg := config.DefaultConfig()

	// Create API client
	apiClient := client.NewFromConfig(cfg)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test cases for integration testing
	tests := []struct {
		name         string
		sessionToken string
		sessionId    string
		expectError  bool
	}{
		{
			name:         "logout with valid parameters",
			sessionToken: "test-session-token",
			sessionId:    "test-session-id",
			expectError:  false, // May succeed or fail depending on server state
		},
		{
			name:         "logout with empty sessionToken",
			sessionToken: "",
			sessionId:    "test-session-id",
			expectError:  true,
		},
		{
			name:         "logout with empty sessionId",
			sessionToken: "test-session-token",
			sessionId:    "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, httpResp, err := apiClient.OAuthAPI.Logout(ctx, tt.sessionToken, tt.sessionId)

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			// For valid parameters, we may get either success or failure
			// depending on whether the session is valid on the server
			if err != nil {
				t.Logf("Logout failed (expected for invalid session): %v", err)
				if httpResp != nil {
					t.Logf("HTTP Status Code: %d", httpResp.StatusCode)
				}
				return
			}

			// If no error, verify response structure
			if httpResp == nil {
				t.Error("Expected HTTP response but got nil")
				return
			}

			t.Logf("Logout response - Success: %v, Code: %s, Message: %s",
				response.Success, response.Code, response.Data.Message)

			// Basic response structure validation
			if response.Code == "" {
				t.Error("Expected non-empty response code")
			}

			if response.RequestID == "" {
				t.Error("Expected non-empty request ID")
			}
		})
	}
}

// TestLogoutWithRealSession tests logout with a real session (if available)
func TestLogoutWithRealSession(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=true to run.")
	}

	// Get configuration
	cfg, err := config.GetConfig()
	if err != nil {
		t.Skipf("Could not load config: %v", err)
	}

	// Check if we have valid tokens
	tokens, err := cfg.GetTokens()
	if err != nil {
		t.Skipf("No valid tokens found: %v", err)
	}

	// Create API client
	apiClient := client.NewFromConfig(cfg)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Logf("Testing logout with real session - SessionId: %s", tokens.SessionId)

	// Call logout with real session data
	response, httpResp, err := apiClient.OAuthAPI.Logout(ctx, tokens.LoginToken, tokens.SessionId)

	if err != nil {
		t.Logf("Logout failed: %v", err)
		if httpResp != nil {
			t.Logf("HTTP Status Code: %d", httpResp.StatusCode)
		}
		// This is not necessarily a test failure as the session might be invalid
		return
	}

	// Log the response
	t.Logf("Logout successful - Success: %v, Code: %s, Message: %s",
		response.Success, response.Code, response.Data.Message)

	// Verify response structure
	if response.Code == "" {
		t.Error("Expected non-empty response code")
	}

	if response.RequestID == "" {
		t.Error("Expected non-empty request ID")
	}

	// If logout was successful, the session should be invalidated
	if response.Success {
		t.Log("Session successfully logged out")
	}
}
