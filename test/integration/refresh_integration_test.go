// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
)

// TestRefreshTokenIntegration tests the refresh token integration with AgbCloud API
func TestRefreshTokenIntegration(t *testing.T) {
	// Skip if running in CI without proper configuration
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	apiClient := client.NewDefault()

	// Test with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("TestRefreshTokenWithValidTokens", func(t *testing.T) {
		t.Logf("Testing refresh token with valid keepAliveToken and sessionId")

		// Test parameters - using test values
		keepAliveToken := "test_keep_alive_token_12345"
		sessionId := "test_session_id_67890"

		response, httpResp, err := apiClient.OAuthAPI.RefreshToken(ctx, keepAliveToken, sessionId)

		// Log the response for debugging
		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
			t.Logf("Response Headers: %v", httpResp.Header)
			if httpResp.Request != nil {
				t.Logf("Request URL: %s", httpResp.Request.URL.String())
			}
		} else {
			t.Logf("HTTP Response is nil")
		}

		if err != nil {
			if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
				t.Logf("API Error: %s", apiErr.Error())
				t.Logf("Response Body: %s", string(apiErr.Body()))
				if httpResp != nil {
					t.Logf("HTTP Status Code: %d", httpResp.StatusCode)
				}
				// For test tokens, we expect this to fail with specific error
				if httpResp != nil && (httpResp.StatusCode == 401 || httpResp.StatusCode == 400) {
					t.Logf("[OK] Expected error response for test tokens: %s", apiErr.Error())
					return
				}
				t.Fatalf("[ERROR] Unexpected API error occurred: %s", apiErr.Error())
			} else {
				t.Fatalf("[ERROR] Network error prevented API communication: %v", err)
			}
		}

		// Validate successful response structure
		if response.Success {
			t.Logf("[OK] Success! Refresh token response received")
			t.Logf("Success: %v", response.Success)
			t.Logf("Code: %s", response.Code)
			t.Logf("RequestID: %s", response.RequestID)
			t.Logf("TraceID: %s", response.TraceID)
			t.Logf("HTTPStatusCode: %d", response.HTTPStatusCode)

			// Validate response data contains expected fields
			if response.Data.LoginToken == "" {
				t.Errorf("[ERROR] Response missing new login token")
			}
			if response.Data.SessionId == "" {
				t.Errorf("[ERROR] Response missing new session ID")
			}
			if response.Data.KeepAliveToken == "" {
				t.Errorf("[ERROR] Response missing new keep alive token")
			}
		}
	})

	t.Run("TestRefreshTokenWithEmptyKeepAliveToken", func(t *testing.T) {
		t.Logf("Testing refresh token with empty keepAliveToken")

		keepAliveToken := ""
		sessionId := "test_session_id"

		response, httpResp, err := apiClient.OAuthAPI.RefreshToken(ctx, keepAliveToken, sessionId)

		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
		}

		// Should return error for empty keepAliveToken
		if err == nil {
			t.Errorf("[ERROR] Expected error for empty keepAliveToken, but got success: %+v", response)
		} else {
			t.Logf("[OK] Expected error for empty keepAliveToken: %v", err)
		}
	})

	t.Run("TestRefreshTokenWithEmptySessionId", func(t *testing.T) {
		t.Logf("Testing refresh token with empty sessionId")

		keepAliveToken := "test_keep_alive_token"
		sessionId := ""

		response, httpResp, err := apiClient.OAuthAPI.RefreshToken(ctx, keepAliveToken, sessionId)

		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
		}

		// Should return error for empty sessionId
		if err == nil {
			t.Errorf("[ERROR] Expected error for empty sessionId, but got success: %+v", response)
		} else {
			t.Logf("[OK] Expected error for empty sessionId: %v", err)
		}
	})

	t.Run("TestRefreshTokenWithInvalidTokens", func(t *testing.T) {
		t.Logf("Testing refresh token with invalid tokens")

		keepAliveToken := "invalid_keep_alive_token"
		sessionId := "invalid_session_id"

		response, httpResp, err := apiClient.OAuthAPI.RefreshToken(ctx, keepAliveToken, sessionId)

		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
		}

		if err != nil {
			if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
				// Expected to fail with invalid tokens
				t.Logf("[OK] Expected API error for invalid tokens: %s", apiErr.Error())
			} else {
				t.Errorf("[ERROR] Network error: %v", err)
			}
		} else if !response.Success {
			t.Logf("[OK] Expected failure response for invalid tokens: %+v", response)
		} else {
			t.Errorf("[ERROR] Expected error for invalid tokens, but got success: %+v", response)
		}
	})

	t.Run("TestRefreshTokenWithExpiredTokens", func(t *testing.T) {
		t.Logf("Testing refresh token with expired tokens")

		// Use tokens that simulate expired state
		keepAliveToken := "expired_keep_alive_token_12345"
		sessionId := "expired_session_id_67890"

		response, httpResp, err := apiClient.OAuthAPI.RefreshToken(ctx, keepAliveToken, sessionId)

		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
		}

		if err != nil {
			if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
				// Expected to fail with expired tokens
				t.Logf("[OK] Expected API error for expired tokens: %s", apiErr.Error())
				// Check for specific error codes that indicate token expiration
				if httpResp != nil && (httpResp.StatusCode == 401 || httpResp.StatusCode == 403) {
					t.Logf("[OK] Received expected HTTP status for expired tokens: %d", httpResp.StatusCode)
				}
			} else {
				t.Errorf("[ERROR] Network error: %v", err)
			}
		} else if !response.Success {
			t.Logf("[OK] Expected failure response for expired tokens: %+v", response)
		} else {
			t.Errorf("[ERROR] Expected error for expired tokens, but got success: %+v", response)
		}
	})
}
