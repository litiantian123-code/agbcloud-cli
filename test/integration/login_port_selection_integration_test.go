// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/auth"
	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginPortSelectionIntegration(t *testing.T) {
	// Skip if running in CI without proper setup
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test")
	}

	tests := []struct {
		name                 string
		occupyDefaultPort    bool
		expectedPortUsed     string
		mockAlternativePorts string
	}{
		{
			name:                 "default port available",
			occupyDefaultPort:    false,
			expectedPortUsed:     "3000",
			mockAlternativePorts: "51152,53152,55152,57152",
		},
		{
			name:                 "default port occupied, use alternative",
			occupyDefaultPort:    true,
			expectedPortUsed:     "51152", // First alternative
			mockAlternativePorts: "51152,53152,55152,57152",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: Create a mock server for OAuth API
			var apiCalls []map[string]string
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Capture query parameters
				queryParams := make(map[string]string)
				for key, values := range r.URL.Query() {
					if len(values) > 0 {
						queryParams[key] = values[0]
					}
				}
				queryParams["path"] = r.URL.Path
				apiCalls = append(apiCalls, queryParams)

				// Mock response based on endpoint
				if strings.Contains(r.URL.Path, "login_provider") {
					response := fmt.Sprintf(`{
						"code": "200",
						"requestId": "test-request-id",
						"success": true,
						"data": {
							"invokeUrl": "https://oauth.example.com/auth?code=test&state=test",
							"alternativePorts": "%s"
						},
						"traceId": "test-trace-id",
						"httpStatusCode": 200
					}`, tt.mockAlternativePorts)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(response))
				} else if strings.Contains(r.URL.Path, "login_translate") {
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
				}
			}))
			defer mockServer.Close()

			// Setup: Optionally occupy the default port
			var defaultPortListener net.Listener
			if tt.occupyDefaultPort {
				var err error
				defaultPortListener, err = net.Listen("tcp", ":3000")
				require.NoError(t, err, "Failed to occupy default port for test")
				defer defaultPortListener.Close()
			}

			// Test the port selection logic directly
			cfg := &client.Configuration{
				Servers: client.ServerConfigurations{
					{
						URL: mockServer.URL,
					},
				},
			}

			apiClient := client.NewAPIClient(cfg)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Simulate the login flow
			defaultPort := "3000"

			// First call - without localhostPort parameter
			response1, _, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, fmt.Sprintf("http://localhost:%s", defaultPort), "CLI", "GOOGLE_LOCALHOST")
			require.NoError(t, err)
			require.True(t, response1.Success)

			// Determine which port to use
			var finalPort string
			var finalResponse client.OAuthLoginProviderResponse

			if !auth.IsPortOccupied(defaultPort) {
				finalPort = defaultPort
				finalResponse = response1
			} else {
				// Select available port from alternatives
				selectedPort, err := auth.SelectAvailablePort(defaultPort, response1.Data.AlternativePorts)
				require.NoError(t, err)

				// Make second API call with selected port
				response2, _, err := apiClient.OAuthAPI.GetLoginProviderURLWithPort(ctx, fmt.Sprintf("http://localhost:%s", selectedPort), "CLI", "GOOGLE_LOCALHOST", selectedPort)
				require.NoError(t, err)
				require.True(t, response2.Success)

				finalPort = selectedPort
				finalResponse = response2
			}

			// Verify the expected port was selected
			assert.Equal(t, tt.expectedPortUsed, finalPort)

			// Verify API calls
			if tt.occupyDefaultPort {
				// Should have made two API calls
				require.Len(t, apiCalls, 2)

				// First call should not have localhostPort
				firstCall := apiCalls[0]
				_, hasLocalhostPort := firstCall["localhostPort"]
				assert.False(t, hasLocalhostPort, "First call should not include localhostPort parameter")

				// Second call should have localhostPort
				secondCall := apiCalls[1]
				assert.Equal(t, tt.expectedPortUsed, secondCall["localhostPort"])
			} else {
				// Should have made only one API call
				require.Len(t, apiCalls, 1)

				// First call should not have localhostPort
				firstCall := apiCalls[0]
				_, hasLocalhostPort := firstCall["localhostPort"]
				assert.False(t, hasLocalhostPort, "First call should not include localhostPort parameter")
			}

			// Test LoginTranslate with the selected port
			translateResponse, _, err := apiClient.OAuthAPI.LoginTranslateWithPort(ctx, "CLI", "GOOGLE_LOCALHOST", "test-auth-code", finalPort)
			require.NoError(t, err)
			require.True(t, translateResponse.Success)

			// Verify the LoginTranslate call was made with correct port
			require.Greater(t, len(apiCalls), 1)
			lastCall := apiCalls[len(apiCalls)-1]
			assert.Equal(t, finalPort, lastCall["localhostPort"])
			assert.Equal(t, "test-auth-code", lastCall["authCode"])

			// Verify response data
			assert.NotEmpty(t, finalResponse.Data.InvokeURL)
			assert.Equal(t, tt.mockAlternativePorts, finalResponse.Data.AlternativePorts)
		})
	}
}

func TestPortSelectionWithRealPorts(t *testing.T) {
	// Skip if running in CI without proper setup
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test")
	}

	// Test with real port occupation
	t.Run("real port occupation test", func(t *testing.T) {
		// Find an available port for testing
		listener, err := net.Listen("tcp", ":0")
		require.NoError(t, err)

		addr := listener.Addr().(*net.TCPAddr)
		testPort := fmt.Sprintf("%d", addr.Port)
		listener.Close()

		// Verify port is available
		assert.False(t, auth.IsPortOccupied(testPort), "Test port should be available")

		// Occupy the port
		listener, err = net.Listen("tcp", ":"+testPort)
		require.NoError(t, err)
		defer listener.Close()

		// Verify port is now occupied
		assert.True(t, auth.IsPortOccupied(testPort), "Test port should be occupied")

		// Test port selection with alternatives
		alternativePorts := "51152,53152,55152,57152"
		selectedPort, err := auth.SelectAvailablePort(testPort, alternativePorts)
		require.NoError(t, err)

		// Should select first alternative since test port is occupied
		assert.Equal(t, "51152", selectedPort)

		// Verify selected port is actually available
		assert.False(t, auth.IsPortOccupied(selectedPort), "Selected port should be available")
	})
}

func TestLoginCommandWithPortSelection(t *testing.T) {
	// Skip if running in CI or if binary doesn't exist
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test")
	}

	// Check if binary exists
	if _, err := os.Stat("../../agbcloud-cli"); os.IsNotExist(err) {
		t.Skip("Binary not found, run 'go build' first")
	}

	t.Run("login command help", func(t *testing.T) {
		// Test that the login command still works with the new port selection logic
		cmd := exec.Command("../../agbcloud-cli", "login", "--help")
		output, err := cmd.CombinedOutput()

		// Command should succeed (exit code 0)
		assert.NoError(t, err, "Login command help should work")

		// Output should contain login-related text
		outputStr := string(output)
		assert.Contains(t, outputStr, "login", "Help output should mention login")
		assert.Contains(t, outputStr, "OAuth", "Help output should mention OAuth")
	})
}
