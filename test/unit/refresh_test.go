// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agbcloud/agbcloud-cli/internal/client"
)

// TestRefreshToken tests the refresh token functionality
func TestRefreshToken(t *testing.T) {
	// Test case 1: Valid refresh request
	t.Run("ValidRefreshRequest", func(t *testing.T) {
		// Mock server that returns successful refresh response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request path
			expectedPath := "/api/biz_login/refresh"
			if r.URL.Path != expectedPath {
				t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
			}

			// Verify HTTP method
			if r.Method != http.MethodGet {
				t.Errorf("Expected GET method, got %s", r.Method)
			}

			// Verify query parameters
			keepAliveToken := r.URL.Query().Get("keepAliveToken")
			sessionId := r.URL.Query().Get("sessionId")

			if keepAliveToken == "" {
				t.Error("keepAliveToken parameter is missing")
			}
			if sessionId == "" {
				t.Error("sessionId parameter is missing")
			}

			// Return mock successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"success": true,
				"code": "SUCCESS",
				"requestId": "test-request-id",
				"traceId": "test-trace-id",
				"httpStatusCode": 200,
				"data": {
					"loginToken": "new_login_token_12345",
					"sessionId": "new_session_id_67890",
					"keepAliveToken": "new_keep_alive_token_abcde",
					"expiresAt": "2025-09-05T20:51:07Z"
				}
			}`)) // Ignore errors in test mock server
		}))
		defer server.Close()

		// Create client with test server
		cfg := client.NewConfiguration()
		cfg.Host = server.URL[7:] // Remove "http://" prefix
		cfg.Scheme = "http"
		apiClient := client.NewAPIClient(cfg)

		// Test refresh with valid parameters
		keepAliveToken := "test_keep_alive_token"
		sessionId := "test_session_id"

		response, httpResp, err := apiClient.OAuthAPI.RefreshToken(context.Background(), keepAliveToken, sessionId)

		// Verify no error occurred
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify HTTP response
		if httpResp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", httpResp.StatusCode)
		}

		// Verify response structure
		if !response.Success {
			t.Error("Expected success to be true")
		}
		if response.Code != "SUCCESS" {
			t.Errorf("Expected code to be SUCCESS, got %s", response.Code)
		}
		if response.Data.LoginToken == "" {
			t.Error("Expected new login token to be present")
		}
		if response.Data.SessionId == "" {
			t.Error("Expected new session ID to be present")
		}
		if response.Data.KeepAliveToken == "" {
			t.Error("Expected new keep alive token to be present")
		}
		if response.Data.ExpiresAt == "" {
			t.Error("Expected expiresAt to be present")
		}
	})

	// Test case 2: Empty keepAliveToken
	t.Run("EmptyKeepAliveToken", func(t *testing.T) {
		cfg := client.NewConfiguration()
		apiClient := client.NewAPIClient(cfg)

		keepAliveToken := ""
		sessionId := "test_session_id"

		_, _, err := apiClient.OAuthAPI.RefreshToken(context.Background(), keepAliveToken, sessionId)

		// Should return error for empty keepAliveToken
		if err == nil {
			t.Error("Expected error for empty keepAliveToken, but got nil")
		}

		if genErr, ok := err.(*client.GenericOpenAPIError); ok {
			expectedMsg := "keepAliveToken parameter is required"
			if genErr.Error() != expectedMsg {
				t.Errorf("Expected error message '%s', got '%s'", expectedMsg, genErr.Error())
			}
		} else {
			t.Errorf("Expected GenericOpenAPIError, got %T", err)
		}
	})

	// Test case 3: Empty sessionId
	t.Run("EmptySessionId", func(t *testing.T) {
		cfg := client.NewConfiguration()
		apiClient := client.NewAPIClient(cfg)

		keepAliveToken := "test_keep_alive_token"
		sessionId := ""

		_, _, err := apiClient.OAuthAPI.RefreshToken(context.Background(), keepAliveToken, sessionId)

		// Should return error for empty sessionId
		if err == nil {
			t.Error("Expected error for empty sessionId, but got nil")
		}

		if genErr, ok := err.(*client.GenericOpenAPIError); ok {
			expectedMsg := "sessionId parameter is required"
			if genErr.Error() != expectedMsg {
				t.Errorf("Expected error message '%s', got '%s'", expectedMsg, genErr.Error())
			}
		} else {
			t.Errorf("Expected GenericOpenAPIError, got %T", err)
		}
	})

	// Test case 4: Server error response
	t.Run("ServerErrorResponse", func(t *testing.T) {
		// Mock server that returns error response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{
				"success": false,
				"code": "INVALID_TOKEN",
				"requestId": "test-request-id",
				"traceId": "test-trace-id",
				"httpStatusCode": 401,
				"message": "Invalid or expired tokens"
			}`)) // Ignore errors in test mock server
		}))
		defer server.Close()

		// Create client with test server
		cfg := client.NewConfiguration()
		cfg.Host = server.URL[7:] // Remove "http://" prefix
		cfg.Scheme = "http"
		apiClient := client.NewAPIClient(cfg)

		keepAliveToken := "invalid_token"
		sessionId := "invalid_session"

		_, httpResp, err := apiClient.OAuthAPI.RefreshToken(context.Background(), keepAliveToken, sessionId)

		// Should return error for server error
		if err == nil {
			t.Error("Expected error for server error response, but got nil")
		}

		// Verify HTTP status code
		if httpResp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got %d", httpResp.StatusCode)
		}

		// Verify error type
		if _, ok := err.(*client.GenericOpenAPIError); !ok {
			t.Errorf("Expected GenericOpenAPIError, got %T", err)
		}
	})
}
