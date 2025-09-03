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

	"github.com/liyuebing/agbcloud-cli/internal/client"
)

// TestOAuthGoogleLoginWithMockServer tests the OAuth Google login API with a mock server
func TestOAuthGoogleLoginWithMockServer(t *testing.T) {
	// Create a mock server that returns the expected response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/api/oauth/google/login" {
			t.Errorf("Expected path /api/oauth/google/login, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Check query parameters
		fromUrlPath := r.URL.Query().Get("fromUrlPath")
		if fromUrlPath != "https://agb.cloud" {
			t.Errorf("Expected fromUrlPath=https://agb.cloud, got %s", fromUrlPath)
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
		response := client.OAuthGoogleLoginResponse{
			Code:      "success",
			RequestID: "test-request-id",
			Success:   true,
			Data: client.OAuthGoogleLoginData{
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

	// Test the OAuth Google login API
	response, httpResp, err := apiClient.OAuthAPI.GetGoogleLoginURL(ctx, "https://agb.cloud")

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

	t.Logf("✅ OAuth Google login test passed!")
	t.Logf("InvokeURL: %s", response.Data.InvokeURL)
}

// TestOAuthGoogleLoginURLConstruction tests URL construction
func TestOAuthGoogleLoginURLConstruction(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Just return a simple response
		response := client.OAuthGoogleLoginResponse{
			Code:    "success",
			Success: true,
			Data: client.OAuthGoogleLoginData{
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
		name        string
		fromUrlPath string
		expectQuery string
	}{
		{
			name:        "Basic URL",
			fromUrlPath: "https://agb.cloud",
			expectQuery: "fromUrlPath=https%3A%2F%2Fagb.cloud",
		},
		{
			name:        "URL with path",
			fromUrlPath: "https://agb.cloud/dashboard",
			expectQuery: "fromUrlPath=https%3A%2F%2Fagb.cloud%2Fdashboard",
		},
		{
			name:        "Empty fromUrlPath",
			fromUrlPath: "",
			expectQuery: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, httpResp, err := apiClient.OAuthAPI.GetGoogleLoginURL(ctx, tc.fromUrlPath)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if httpResp == nil {
				t.Fatal("Expected HTTP response, got nil")
			}

			// Check if the query string matches expectation
			actualQuery := httpResp.Request.URL.RawQuery
			if tc.expectQuery == "" && actualQuery != "" {
				t.Errorf("Expected empty query, got %s", actualQuery)
			} else if tc.expectQuery != "" && actualQuery != tc.expectQuery {
				t.Errorf("Expected query %s, got %s", tc.expectQuery, actualQuery)
			}

			t.Logf("✅ URL construction test passed for %s", tc.name)
		})
	}
}
