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

// TestLoginTranslateIntegration tests the login translate integration with OAuth API
func TestLoginTranslateIntegration(t *testing.T) {
	// Skip if running in CI without proper configuration
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	apiClient := client.NewDefault()

	// Test with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("TestLoginTranslateWithValidAuthCode", func(t *testing.T) {
		t.Logf("Testing OAuth login translate with valid auth code")

		// Test parameters - using correct values
		loginClient := "CLI"
		oauthProvider := "GOOGLE_LOCALHOST"
		authCode := "test_auth_code_12345"

		response, httpResp, err := apiClient.OAuthAPI.LoginTranslate(ctx, loginClient, oauthProvider, authCode)

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
				// For test auth code, we expect this to fail with specific error
				if httpResp != nil && httpResp.StatusCode == 200 {
					t.Logf("✅ Expected error response for test auth code: %s", apiErr.Error())
					return
				}
				t.Fatalf("❌ Unexpected API error occurred: %s", apiErr.Error())
			} else {
				t.Fatalf("❌ Network error prevented API communication: %v", err)
			}
		}

		// Validate successful response structure
		if response.Success {
			t.Logf("✅ Success! Login translate response received")
			t.Logf("Success: %v", response.Success)
			t.Logf("Code: %s", response.Code)
			t.Logf("RequestID: %s", response.RequestID)
			t.Logf("TraceID: %s", response.TraceID)
			t.Logf("HTTPStatusCode: %d", response.HTTPStatusCode)

			// Validate response data contains expected fields
			if response.Data.LoginToken == "" {
				t.Errorf("❌ Response missing login token")
			}
			if response.Data.SessionId == "" {
				t.Errorf("❌ Response missing session ID")
			}
			if response.Data.KeepAliveToken == "" {
				t.Errorf("❌ Response missing keep alive token")
			}
		}
	})

	t.Run("TestLoginTranslateWithEmptyAuthCode", func(t *testing.T) {
		t.Logf("Testing OAuth login translate with empty auth code")

		loginClient := "CLI"
		oauthProvider := "GOOGLE_LOCALHOST"
		authCode := ""

		response, httpResp, err := apiClient.OAuthAPI.LoginTranslate(ctx, loginClient, oauthProvider, authCode)

		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
		}

		// Should return error for empty auth code
		if err == nil {
			t.Errorf("❌ Expected error for empty auth code, but got success: %+v", response)
		} else {
			t.Logf("✅ Expected error for empty auth code: %v", err)
		}
	})

	t.Run("TestLoginTranslateWithInvalidParameters", func(t *testing.T) {
		t.Logf("Testing OAuth login translate with invalid parameters")

		loginClient := ""
		oauthProvider := "INVALID_PROVIDER"
		authCode := "invalid_code"

		response, httpResp, err := apiClient.OAuthAPI.LoginTranslate(ctx, loginClient, oauthProvider, authCode)

		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
		}

		// Should return error for invalid parameters
		if err == nil {
			t.Errorf("❌ Expected error for invalid parameters, but got success: %+v", response)
		} else {
			t.Logf("✅ Expected error for invalid parameters: %v", err)
		}
	})

	t.Run("TestLoginTranslateWithDifferentProviders", func(t *testing.T) {
		t.Logf("Testing OAuth login translate with different OAuth providers")

		testCases := []struct {
			name          string
			loginClient   string
			oauthProvider string
		}{
			{"GOOGLE_Provider", "CLI", "GOOGLE"},
			{"GOOGLE_LOCALHOST_Provider", "CLI", "GOOGLE_LOCALHOST"},
			{"GITHUB_Provider", "CLI", "GITHUB"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				authCode := "test_code_for_" + tc.oauthProvider

				response, httpResp, err := apiClient.OAuthAPI.LoginTranslate(ctx, tc.loginClient, tc.oauthProvider, authCode)

				if httpResp != nil {
					t.Logf("HTTP Status for %s: %d", tc.oauthProvider, httpResp.StatusCode)
				}

				if err != nil {
					if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
						// Expected to fail with test codes
						t.Logf("Expected API error for %s: %s", tc.oauthProvider, apiErr.Error())
					} else {
						t.Errorf("❌ Network error for %s: %v", tc.oauthProvider, err)
					}
				} else {
					t.Logf("✅ Success for %s: %+v", tc.oauthProvider, response)
				}
			})
		}
	})

	t.Run("TestLoginTranslateWithUnsupportedClient", func(t *testing.T) {
		t.Logf("Testing OAuth login translate with unsupported login client")

		// Test with unsupported client type
		loginClient := "WEB"
		oauthProvider := "GOOGLE"
		authCode := "test_code"

		response, httpResp, err := apiClient.OAuthAPI.LoginTranslate(ctx, loginClient, oauthProvider, authCode)

		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
		}

		if err != nil {
			if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
				t.Logf("✅ Expected error for unsupported client: %s", apiErr.Error())
			} else {
				t.Errorf("❌ Network error: %v", err)
			}
		} else if !response.Success {
			t.Logf("✅ Expected failure response for unsupported client: %+v", response)
		} else {
			t.Errorf("❌ Expected error for unsupported client, but got success: %+v", response)
		}
	})
}
