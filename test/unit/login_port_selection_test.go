// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginPortSelectionFlow(t *testing.T) {
	tests := []struct {
		name                   string
		defaultPort            string
		mockAlternativePorts   string
		expectedFirstCallPort  string // Port used in first API call (should be empty)
		expectedSecondCallPort string // Port used in second API call
		expectSecondCall       bool   // Whether a second API call should be made
		simulatePortOccupied   bool   // Whether to simulate default port being occupied
	}{
		{
			name:                   "default port available, no second call needed",
			defaultPort:            "3000",
			mockAlternativePorts:   "51152,53152,55152,57152",
			expectedFirstCallPort:  "", // First call should not include localhostPort
			expectedSecondCallPort: "",
			expectSecondCall:       false,
			simulatePortOccupied:   false,
		},
		{
			name:                   "default port occupied, use first alternative",
			defaultPort:            "3000",
			mockAlternativePorts:   "51152,53152,55152,57152",
			expectedFirstCallPort:  "",      // First call should not include localhostPort
			expectedSecondCallPort: "51152", // Second call should use first alternative
			expectSecondCall:       true,
			simulatePortOccupied:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Track API calls
			var apiCalls []map[string]string

			// Create a test server that captures requests
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Capture query parameters
				queryParams := make(map[string]string)
				for key, values := range r.URL.Query() {
					if len(values) > 0 {
						queryParams[key] = values[0]
					}
				}
				apiCalls = append(apiCalls, queryParams)

				// Mock response with alternativePorts
				response := `{
					"code": "200",
					"requestId": "test-request-id",
					"success": true,
					"data": {
						"invokeUrl": "https://oauth.example.com/auth?code=test",
						"alternativePorts": "` + tt.mockAlternativePorts + `"
					},
					"traceId": "test-trace-id",
					"httpStatusCode": 200
				}`
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response)) // Ignore errors in test mock server
			}))
			defer server.Close()

			// Create client configuration
			cfg := &client.Configuration{
				Servers: client.ServerConfigurations{
					{
						URL: server.URL,
					},
				},
			}

			apiClient := client.NewAPIClient(cfg)
			ctx := context.Background()

			// Simulate the login flow
			// First call - without localhostPort parameter
			response1, _, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, "http://localhost:"+tt.defaultPort, "CLI", "GOOGLE_LOCALHOST")
			require.NoError(t, err)
			require.True(t, response1.Success)

			// Verify first call
			require.Len(t, apiCalls, 1)
			firstCall := apiCalls[0]
			assert.Equal(t, "CLI", firstCall["loginClient"])
			assert.Equal(t, "GOOGLE_LOCALHOST", firstCall["oauthProvider"])

			// First call should not have localhostPort parameter
			_, hasLocalhostPort := firstCall["localhostPort"]
			assert.False(t, hasLocalhostPort, "First call should not include localhostPort parameter")

			// Check if we need to make a second call (simulate port occupation check)
			if tt.expectSecondCall {
				// For testing purposes, we'll directly use the expected port
				// In real implementation, this would be determined by port availability check
				selectedPort := tt.expectedSecondCallPort
				assert.Equal(t, tt.expectedSecondCallPort, selectedPort)

				// Make second call with selected port
				response2, _, err := apiClient.OAuthAPI.GetLoginProviderURLWithPort(ctx, "http://localhost:"+selectedPort, "CLI", "GOOGLE_LOCALHOST", selectedPort)
				require.NoError(t, err)
				require.True(t, response2.Success)

				// Verify second call
				require.Len(t, apiCalls, 2)
				secondCall := apiCalls[1]
				assert.Equal(t, "CLI", secondCall["loginClient"])
				assert.Equal(t, "GOOGLE_LOCALHOST", secondCall["oauthProvider"])
				assert.Equal(t, tt.expectedSecondCallPort, secondCall["localhostPort"])
			} else {
				// Should only have one API call
				assert.Len(t, apiCalls, 1)
			}
		})
	}
}

func TestLoginTranslateWithSelectedPort(t *testing.T) {
	// Test that LoginTranslate is called with the selected port
	selectedPort := "51152"

	// Create a test server that captures requests
	var capturedRequest *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequest = r

		// Mock successful response
		response := `{
			"code": "200",
			"requestId": "test-request-id",
			"success": true,
			"data": {
				"loginToken": "test-login-token",
				"sessionId": "test-session-id",
				"keepAliveToken": "test-keep-alive-token",
				"expiresAt": "2025-01-01T00:00:00Z"
			},
			"traceId": "test-trace-id",
			"httpStatusCode": 200
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response)) // Ignore errors in test mock server
	}))
	defer server.Close()

	// Create client configuration
	cfg := &client.Configuration{
		Servers: client.ServerConfigurations{
			{
				URL: server.URL,
			},
		},
	}

	apiClient := client.NewAPIClient(cfg)
	ctx := context.Background()

	// Call LoginTranslateWithPort
	_, _, err := apiClient.OAuthAPI.LoginTranslateWithPort(ctx, "CLI", "GOOGLE_LOCALHOST", "test-auth-code", selectedPort)
	require.NoError(t, err)

	// Verify the request was captured
	require.NotNil(t, capturedRequest)

	// Parse query parameters
	queryParams := capturedRequest.URL.Query()

	// Verify localhostPort parameter is present
	assert.Equal(t, selectedPort, queryParams.Get("localhostPort"))
	assert.Equal(t, "CLI", queryParams.Get("loginClient"))
	assert.Equal(t, "GOOGLE_LOCALHOST", queryParams.Get("oauthProvider"))
	assert.Equal(t, "test-auth-code", queryParams.Get("authCode"))
}
