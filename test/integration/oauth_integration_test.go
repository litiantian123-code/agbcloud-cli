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

	"github.com/liyuebing/agbcloud-cli/internal/client"
	"github.com/liyuebing/agbcloud-cli/internal/config"
)

// TestOAuthAuthSourceURL tests the OAuth auth source URL endpoint
func TestOAuthAuthSourceURL(t *testing.T) {
	// Skip if running in CI without proper configuration
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	// Create client WITHOUT API key for OAuth endpoint (OAuth doesn't need auth)
	cfg := &config.Config{
		Endpoint: "https://agb.cloud",
		// Note: OAuth endpoint doesn't require APIKey
	}

	apiClient := client.NewFromConfig(cfg)

	// Test with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("TestOAuthGoogleLogin", func(t *testing.T) {
		t.Logf("Testing OAuth Google login URL with fromUrlPath=https://agb.cloud")

		response, httpResp, err := apiClient.OAuthAPI.GetGoogleLoginURL(ctx, "https://agb.cloud")

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

				// Verify we got a proper HTTP response
				if httpResp == nil {
					t.Fatal("Expected HTTP response, got nil")
				}

				// Log detailed information about the request/response
				if httpResp.Request != nil {
					t.Logf("Request was made to: %s", httpResp.Request.URL.String())
					t.Logf("Request method: %s", httpResp.Request.Method)
					t.Logf("Request headers: %v", httpResp.Request.Header)
				}

				// Check different status codes
				switch httpResp.StatusCode {
				case 200:
					t.Logf("✅ Success: API endpoint exists and responded")
				case 400:
					t.Logf("⚠️  Bad Request: Check parameters or request format")
				case 401:
					t.Logf("⚠️  Unauthorized: Authentication required")
				case 403:
					t.Logf("⚠️  Forbidden: Access denied")
				case 404:
					t.Logf("❌ Not Found: Endpoint does not exist")
				case 500:
					t.Logf("⚠️  Server Error: Internal server error")
				default:
					t.Logf("❓ Unexpected status: %d", httpResp.StatusCode)
				}
			} else {
				t.Fatalf("❌ Network error prevented API communication: %v", err)
			}
		} else {
			// If successful, verify the response structure
			t.Logf("✅ Success! Google OAuth URL: %s", response.Data.InvokeURL)
			t.Logf("Success: %t", response.Success)
			t.Logf("Code: %s", response.Code)
			t.Logf("RequestID: %s", response.RequestID)
			t.Logf("TraceID: %s", response.TraceID)
			t.Logf("HTTPStatusCode: %d", response.HTTPStatusCode)

			if httpResp.StatusCode != 200 {
				t.Errorf("Expected status code 200 for successful response, got %d", httpResp.StatusCode)
			}

			if response.Data.InvokeURL == "" {
				t.Logf("⚠️  Warning: Empty InvokeURL in response (might be expected)")
			}
		}
	})

	t.Run("TestOAuthWithoutFromUrlPath", func(t *testing.T) {
		t.Logf("Testing OAuth without fromUrlPath parameter")

		response, httpResp, err := apiClient.OAuthAPI.GetGoogleLoginURL(ctx, "")

		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
			if httpResp.Request != nil {
				t.Logf("Request URL: %s", httpResp.Request.URL.String())
			}
		} else {
			t.Logf("HTTP Response is nil")
		}

		if err != nil {
			if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
				t.Logf("API Error (without fromUrlPath): %s", apiErr.Error())
				t.Logf("Response Body: %s", string(apiErr.Body()))
			} else {
				t.Fatalf("❌ Network error prevented API communication (without fromUrlPath): %v", err)
			}
		} else {
			t.Logf("Success (without fromUrlPath): %s", response.Data.InvokeURL)
		}
	})
}
