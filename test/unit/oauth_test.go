// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuthAPI_GetLoginProviderURL_WithLocalhostPort(t *testing.T) {
	tests := []struct {
		name              string
		localhostPort     string
		expectedInQuery   bool
		expectedPortValue string
	}{
		{
			name:              "with localhost port parameter",
			localhostPort:     "3001",
			expectedInQuery:   true,
			expectedPortValue: "3001",
		},
		{
			name:              "without localhost port parameter",
			localhostPort:     "",
			expectedInQuery:   false,
			expectedPortValue: "",
		},
		{
			name:              "with custom localhost port",
			localhostPort:     "51152",
			expectedInQuery:   true,
			expectedPortValue: "51152",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that captures the request
			var capturedRequest *http.Request
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedRequest = r

				// Mock response with alternativePorts
				response := `{
					"code": "200",
					"requestId": "test-request-id",
					"success": true,
					"data": {
						"invokeUrl": "https://oauth.example.com/auth?code=test",
						"alternativePorts": "51152,53152,55152,57152"
					},
					"traceId": "test-trace-id",
					"httpStatusCode": 200
				}`
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(response))
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

			// Call the API with localhostPort parameter
			ctx := context.Background()
			_, _, err := apiClient.OAuthAPI.GetLoginProviderURLWithPort(ctx, "http://localhost:3000", "CLI", "GOOGLE_LOCALHOST", tt.localhostPort)

			// Verify no error occurred
			require.NoError(t, err)

			// Verify the request was captured
			require.NotNil(t, capturedRequest)

			// Parse query parameters
			queryParams := capturedRequest.URL.Query()

			// Check if localhostPort parameter is present when expected
			if tt.expectedInQuery {
				assert.True(t, queryParams.Has("localhostPort"), "localhostPort parameter should be present in query")
				assert.Equal(t, tt.expectedPortValue, queryParams.Get("localhostPort"), "localhostPort parameter value should match")
			} else {
				assert.False(t, queryParams.Has("localhostPort"), "localhostPort parameter should not be present in query")
			}

			// Verify other required parameters are still present
			assert.Equal(t, "CLI", queryParams.Get("loginClient"))
			assert.Equal(t, "GOOGLE_LOCALHOST", queryParams.Get("oauthProvider"))
		})
	}
}

func TestOAuthAPI_LoginTranslate_WithLocalhostPort(t *testing.T) {
	tests := []struct {
		name              string
		localhostPort     string
		expectedInQuery   bool
		expectedPortValue string
	}{
		{
			name:              "with localhost port parameter",
			localhostPort:     "3001",
			expectedInQuery:   true,
			expectedPortValue: "3001",
		},
		{
			name:              "without localhost port parameter",
			localhostPort:     "",
			expectedInQuery:   false,
			expectedPortValue: "",
		},
		{
			name:              "with alternative port from login provider response",
			localhostPort:     "51152",
			expectedInQuery:   true,
			expectedPortValue: "51152",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that captures the request
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
				w.Write([]byte(response))
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

			// Call the API with localhostPort parameter
			ctx := context.Background()
			_, _, err := apiClient.OAuthAPI.LoginTranslateWithPort(ctx, "CLI", "GOOGLE_LOCALHOST", "test-auth-code", tt.localhostPort)

			// Verify no error occurred
			require.NoError(t, err)

			// Verify the request was captured
			require.NotNil(t, capturedRequest)

			// Parse query parameters
			queryParams := capturedRequest.URL.Query()

			// Check if localhostPort parameter is present when expected
			if tt.expectedInQuery {
				assert.True(t, queryParams.Has("localhostPort"), "localhostPort parameter should be present in query")
				assert.Equal(t, tt.expectedPortValue, queryParams.Get("localhostPort"), "localhostPort parameter value should match")
			} else {
				assert.False(t, queryParams.Has("localhostPort"), "localhostPort parameter should not be present in query")
			}

			// Verify other required parameters are still present
			assert.Equal(t, "CLI", queryParams.Get("loginClient"))
			assert.Equal(t, "GOOGLE_LOCALHOST", queryParams.Get("oauthProvider"))
			assert.Equal(t, "test-auth-code", queryParams.Get("authCode"))
		})
	}
}

func TestOAuthLoginProviderResponse_AlternativePorts(t *testing.T) {
	// Test that the response structure can handle alternativePorts field
	jsonResponse := `{
		"code": "200",
		"requestId": "test-request-id",
		"success": true,
		"data": {
			"invokeUrl": "https://oauth.example.com/auth?code=test",
			"alternativePorts": "51152,53152,55152,57152"
		},
		"traceId": "test-trace-id",
		"httpStatusCode": 200
	}`

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
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

	// Call the API
	ctx := context.Background()
	response, _, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, "http://localhost:3000", "CLI", "GOOGLE_LOCALHOST")

	// Verify no error occurred
	require.NoError(t, err)

	// Verify response structure
	assert.True(t, response.Success)
	assert.Equal(t, "https://oauth.example.com/auth?code=test", response.Data.InvokeURL)
	assert.Equal(t, "51152,53152,55152,57152", response.Data.AlternativePorts)
}

func TestPortSelection_FromAlternativePorts(t *testing.T) {
	tests := []struct {
		name             string
		alternativePorts string
		occupiedPorts    []string
		expectedPort     string
		expectError      bool
	}{
		{
			name:             "first alternative port available",
			alternativePorts: "51152,53152,55152,57152",
			occupiedPorts:    []string{},
			expectedPort:     "51152",
			expectError:      false,
		},
		{
			name:             "second alternative port available",
			alternativePorts: "51152,53152,55152,57152",
			occupiedPorts:    []string{"51152"},
			expectedPort:     "53152",
			expectError:      false,
		},
		{
			name:             "all alternative ports occupied",
			alternativePorts: "51152,53152,55152,57152",
			occupiedPorts:    []string{"51152", "53152", "55152", "57152"},
			expectedPort:     "",
			expectError:      true,
		},
		{
			name:             "empty alternative ports",
			alternativePorts: "",
			occupiedPorts:    []string{},
			expectedPort:     "",
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test will verify the port selection logic
			// The actual implementation will be in auth package
			ports := strings.Split(tt.alternativePorts, ",")
			if tt.alternativePorts == "" {
				ports = []string{}
			}

			selectedPort := ""
			for _, port := range ports {
				port = strings.TrimSpace(port)
				if port == "" {
					continue
				}

				// Check if port is occupied (mock check)
				occupied := false
				for _, occupiedPort := range tt.occupiedPorts {
					if port == occupiedPort {
						occupied = true
						break
					}
				}

				if !occupied {
					selectedPort = port
					break
				}
			}

			if tt.expectError {
				assert.Empty(t, selectedPort, "should not find available port")
			} else {
				assert.Equal(t, tt.expectedPort, selectedPort, "should select correct port")
			}
		})
	}
}
