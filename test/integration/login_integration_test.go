// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
)

// TestLoginIntegration tests the login integration with OAuth API
func TestLoginIntegration(t *testing.T) {
	// Skip if running in CI without proper configuration
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	apiClient := client.NewDefault()

	// Test with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("TestLoginProviderWithDefaults", func(t *testing.T) {
		t.Logf("Testing OAuth login provider with default parameters")

		response, httpResp, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, "https://agb.cloud", "CLI", "GOOGLE_LOCALHOST")

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
				t.Fatalf("[ERROR] API error occurred: %s", apiErr.Error())
			} else {
				t.Fatalf("[ERROR] Network error prevented API communication: %v", err)
			}
		}

		// Validate successful response
		if !response.Success {
			t.Fatalf("[ERROR] API returned success=false: %+v", response)
		}

		if response.Data.InvokeURL == "" {
			t.Fatalf("[ERROR] Invalid response data - empty InvokeURL: %+v", response)
		}

		// Log success details
		t.Logf("[OK] Success! Login Provider OAuth URL: %s", response.Data.InvokeURL)
		t.Logf("Success: %v", response.Success)
		t.Logf("Code: %s", response.Code)
		t.Logf("RequestID: %s", response.RequestID)
		t.Logf("TraceID: %s", response.TraceID)
		t.Logf("HTTPStatusCode: %d", response.HTTPStatusCode)

		// Validate OAuth URL format
		if !strings.Contains(response.Data.InvokeURL, "accounts.google.com") {
			t.Errorf("[ERROR] OAuth URL doesn't contain Google domain: %s", response.Data.InvokeURL)
		}

		if !strings.Contains(response.Data.InvokeURL, "client_id=") {
			t.Errorf("[ERROR] OAuth URL missing client_id parameter: %s", response.Data.InvokeURL)
		}
	})

	t.Run("TestLoginProviderWithoutFromUrlPath", func(t *testing.T) {
		t.Logf("Testing OAuth login provider without fromUrlPath parameter")

		response, httpResp, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, "", "CLI", "GOOGLE_LOCALHOST")

		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
			if httpResp.Request != nil {
				t.Logf("Request URL: %s", httpResp.Request.URL.String())
			}
		} else {
			t.Logf("HTTP Response is nil")
		}

		if err != nil {
			t.Fatalf("[ERROR] Network error prevented API communication (without fromUrlPath): %v", err)
		}

		if response.Success && response.Data.InvokeURL != "" {
			t.Logf("Success (without fromUrlPath): %s", response.Data.InvokeURL)
		}
	})

	t.Run("TestLoginProviderWithCustomParameters", func(t *testing.T) {
		t.Logf("Testing OAuth login provider with custom parameters")

		response, httpResp, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, "https://agb.cloud/dashboard", "WEB", "GOOGLE_WEB")

		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
			if httpResp.Request != nil {
				t.Logf("Request URL: %s", httpResp.Request.URL.String())
			}
		} else {
			t.Logf("HTTP Response is nil")
		}

		if err != nil {
			t.Fatalf("[ERROR] Network error prevented API communication (custom parameters): %v", err)
		}

		if response.Success && response.Data.InvokeURL != "" {
			t.Logf("Success (custom parameters): %s", response.Data.InvokeURL)
		}
	})
}
