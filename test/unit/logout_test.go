// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agbcloud/agbcloud-cli/internal/client"
)

// TestLogout tests the logout API endpoint
func TestLogout(t *testing.T) {
	tests := []struct {
		name           string
		loginToken     string
		sessionId      string
		mockResponse   client.OAuthLogoutResponse
		mockStatusCode int
		expectError    bool
		errorContains  string
	}{
		{
			name:       "successful logout",
			loginToken: "test-login-token",
			sessionId:  "test-session-id",
			mockResponse: client.OAuthLogoutResponse{
				Code:           "200",
				RequestID:      "test-request-id",
				Success:        true,
				Data:           client.OAuthLogoutData{Message: "Logout successful"},
				TraceID:        "test-trace-id",
				HTTPStatusCode: 200,
			},
			mockStatusCode: 200,
			expectError:    false,
		},
		{
			name:          "logout with empty loginToken",
			loginToken:    "",
			sessionId:     "test-session-id",
			expectError:   true,
			errorContains: "loginToken parameter is required",
		},
		{
			name:          "logout with empty sessionId",
			loginToken:    "test-login-token",
			sessionId:     "",
			expectError:   true,
			errorContains: "sessionId parameter is required",
		},
		{
			name:       "server error response",
			loginToken: "test-login-token",
			sessionId:  "test-session-id",
			mockResponse: client.OAuthLogoutResponse{
				Code:           "500",
				RequestID:      "test-request-id",
				Success:        false,
				Data:           client.OAuthLogoutData{Message: "Internal server error"},
				TraceID:        "test-trace-id",
				HTTPStatusCode: 500,
			},
			mockStatusCode: 500,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the request method and path
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}
				if r.URL.Path != "/api/biz_login/logout" {
					t.Errorf("Expected path /api/biz_login/logout, got %s", r.URL.Path)
				}

				// Verify query parameters
				loginToken := r.URL.Query().Get("loginToken")
				sessionId := r.URL.Query().Get("sessionId")

				if tt.loginToken != "" && loginToken != tt.loginToken {
					t.Errorf("Expected loginToken=%s, got %s", tt.loginToken, loginToken)
				}
				if tt.sessionId != "" && sessionId != tt.sessionId {
					t.Errorf("Expected sessionId=%s, got %s", tt.sessionId, sessionId)
				}

				// Return mock response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer mockServer.Close()

			// Create client configuration
			cfg := client.NewConfiguration()
			cfg.Servers = []client.ServerConfiguration{
				{
					URL: mockServer.URL,
				},
			}

			// Create API client
			apiClient := client.NewAPIClient(cfg)

			// Call the logout API
			ctx := context.Background()
			response, httpResp, err := apiClient.OAuthAPI.Logout(ctx, tt.loginToken, tt.sessionId)

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
				return
			}

			// Check for unexpected errors
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify response
			if httpResp == nil {
				t.Error("Expected HTTP response but got nil")
				return
			}

			if response.Code != tt.mockResponse.Code {
				t.Errorf("Expected code %s, got %s", tt.mockResponse.Code, response.Code)
			}

			if response.Success != tt.mockResponse.Success {
				t.Errorf("Expected success %v, got %v", tt.mockResponse.Success, response.Success)
			}

			if response.Data.Message != tt.mockResponse.Data.Message {
				t.Errorf("Expected message %s, got %s", tt.mockResponse.Data.Message, response.Data.Message)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
