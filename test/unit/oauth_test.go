// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
)

// TestOAuthLoginProviderWithMockServer tests the OAuth login provider API with a mock server
func TestOAuthLoginProviderWithMockServer(t *testing.T) {
	// Create a mock server that returns the expected response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/api/oauth/login_provider" {
			t.Errorf("Expected path /api/oauth/login_provider, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Check query parameters
		fromUrlPath := r.URL.Query().Get("fromUrlPath")
		if fromUrlPath != "https://agb.cloud" {
			t.Errorf("Expected fromUrlPath=https://agb.cloud, got %s", fromUrlPath)
		}

		// Check new parameters with default values
		loginClient := r.URL.Query().Get("loginClient")
		if loginClient != "CLI" {
			t.Errorf("Expected loginClient=CLI, got %s", loginClient)
		}

		oauthProvider := r.URL.Query().Get("oauthProvider")
		if oauthProvider != "GOOGLE_LOCALHOST" {
			t.Errorf("Expected oauthProvider=GOOGLE_LOCALHOST, got %s", oauthProvider)
		}

		// Check headers
		accept := r.Header.Get("Accept")
		if accept != "application/json" {
			t.Errorf("Expected Accept: application/json, got %s", accept)
		}

		// Verify no authorization header (OAuth endpoint doesn't need auth)
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("Expected no Authorization header, got %s", auth)
		}

		// Return the expected response
		response := client.OAuthLoginProviderResponse{
			Code:      "success",
			RequestID: "test-request-id",
			Success:   true,
			Data: client.OAuthLoginProviderData{
				InvokeURL: "https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=test-client-id",
			},
			TraceID:        "test-trace-id",
			HTTPStatusCode: 200,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create a client configuration pointing to the mock server
	cfg := client.NewConfiguration()
	cfg.Servers[0].URL = mockServer.URL
	cfg.HTTPClient = &http.Client{Timeout: 5 * time.Second}

	apiClient := client.NewAPIClient(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test the OAuth login provider API with new method
	response, httpResp, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, "https://agb.cloud", "CLI", "GOOGLE_LOCALHOST")

	// Verify no error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify HTTP response
	if httpResp == nil {
		t.Fatal("Expected HTTP response, got nil")
	}
	if httpResp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", httpResp.StatusCode)
	}

	// Verify response structure
	if response.Code != "success" {
		t.Errorf("Expected code 'success', got %s", response.Code)
	}
	if !response.Success {
		t.Errorf("Expected success true, got %t", response.Success)
	}
	if response.RequestID != "test-request-id" {
		t.Errorf("Expected requestId 'test-request-id', got %s", response.RequestID)
	}
	if response.Data.InvokeURL == "" {
		t.Error("Expected non-empty InvokeURL")
	}
	if response.TraceID != "test-trace-id" {
		t.Errorf("Expected traceId 'test-trace-id', got %s", response.TraceID)
	}
	if response.HTTPStatusCode != 200 {
		t.Errorf("Expected httpStatusCode 200, got %d", response.HTTPStatusCode)
	}

	t.Logf("✅ OAuth login provider test passed!")
	t.Logf("InvokeURL: %s", response.Data.InvokeURL)
}

// TestOAuthLoginProviderURLConstruction tests URL construction with new parameters
func TestOAuthLoginProviderURLConstruction(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the new endpoint
		if r.URL.Path != "/api/oauth/login_provider" {
			t.Errorf("Expected path /api/oauth/login_provider, got %s", r.URL.Path)
		}

		// Just return a simple response
		response := client.OAuthLoginProviderResponse{
			Code:    "success",
			Success: true,
			Data: client.OAuthLoginProviderData{
				InvokeURL: "https://accounts.google.com/o/oauth2/auth",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	cfg := client.NewConfiguration()
	cfg.Servers[0].URL = mockServer.URL
	apiClient := client.NewAPIClient(cfg)

	ctx := context.Background()

	testCases := []struct {
		name          string
		fromUrlPath   string
		loginClient   string
		oauthProvider string
		expectQuery   map[string]string
	}{
		{
			name:          "Basic URL with defaults",
			fromUrlPath:   "https://agb.cloud",
			loginClient:   "CLI",
			oauthProvider: "GOOGLE_LOCALHOST",
			expectQuery: map[string]string{
				"fromUrlPath":   "https://agb.cloud",
				"loginClient":   "CLI",
				"oauthProvider": "GOOGLE_LOCALHOST",
			},
		},
		{
			name:          "URL with custom parameters",
			fromUrlPath:   "https://agb.cloud/dashboard",
			loginClient:   "WEB",
			oauthProvider: "GOOGLE_WEB",
			expectQuery: map[string]string{
				"fromUrlPath":   "https://agb.cloud/dashboard",
				"loginClient":   "WEB",
				"oauthProvider": "GOOGLE_WEB",
			},
		},
		{
			name:          "Empty fromUrlPath with defaults",
			fromUrlPath:   "",
			loginClient:   "CLI",
			oauthProvider: "GOOGLE_LOCALHOST",
			expectQuery: map[string]string{
				"loginClient":   "CLI",
				"oauthProvider": "GOOGLE_LOCALHOST",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, httpResp, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, tc.fromUrlPath, tc.loginClient, tc.oauthProvider)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if httpResp == nil {
				t.Fatal("Expected HTTP response, got nil")
			}

			// Check if the query parameters match expectation
			actualQuery := httpResp.Request.URL.Query()
			for key, expectedValue := range tc.expectQuery {
				actualValue := actualQuery.Get(key)
				if actualValue != expectedValue {
					t.Errorf("Expected query param %s=%s, got %s", key, expectedValue, actualValue)
				}
			}

			t.Logf("✅ URL construction test passed for %s", tc.name)
		})
	}
}
