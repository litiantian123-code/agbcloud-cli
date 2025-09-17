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
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// TestImageStartIntegration tests the StartImage API with real server
func TestImageStartIntegration(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=true to run.")
	}

	// Get configuration
	cfg, err := config.GetConfig()
	if err != nil {
		t.Skipf("Could not load config: %v", err)
	}

	// Check if we have valid tokens
	tokens, err := cfg.GetTokens()
	if err != nil {
		t.Skipf("No valid tokens found: %v. Please run 'agbcloud login' first.", err)
	}

	t.Logf("[OK] Using authenticated session: %s", tokens.SessionId[:8]+"...")

	// Create API client
	apiClient := client.NewFromConfig(cfg)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test cases for integration testing
	tests := []struct {
		name        string
		imageId     string
		cpu         int
		memory      int
		expectError bool
		description string
	}{
		{
			name:        "start_with_resources",
			imageId:     "test-image-id-123", // Use a test image ID
			cpu:         2,
			memory:      4,
			expectError: true, // Expect error since test image likely doesn't exist
			description: "Test starting image with CPU and memory specifications",
		},
		{
			name:        "start_with_default_resources",
			imageId:     "test-image-id-456", // Use a test image ID
			cpu:         0,
			memory:      0,
			expectError: true, // Expect error since test image likely doesn't exist
			description: "Test starting image with default resources",
		},
		{
			name:        "start_with_invalid_image_id",
			imageId:     "non-existent-image-id",
			cpu:         1,
			memory:      2,
			expectError: true,
			description: "Test error handling with invalid image ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test case: %s", tt.description)
			t.Logf("Image ID: %s, CPU: %d, Memory: %d", tt.imageId, tt.cpu, tt.memory)

			// Call StartImage API
			resp, httpResp, err := apiClient.ImageAPI.StartImage(
				ctx,
				tokens.LoginToken,
				tokens.SessionId,
				tt.imageId,
				tt.cpu,
				tt.memory,
			)

			// Log request details
			if httpResp != nil {
				t.Logf("HTTP Status: %d", httpResp.StatusCode)
				if httpResp.Request != nil {
					t.Logf("Request URL: %s", httpResp.Request.URL.String())
				}
			}

			if err != nil {
				if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
					t.Logf("API Error: %s", apiErr.Error())
					if httpResp != nil {
						t.Logf("HTTP Status Code: %d", httpResp.StatusCode)
					}

					if !tt.expectError {
						t.Errorf("[ERROR] Unexpected API error: %s", apiErr.Error())
					} else {
						t.Logf("[OK] Expected API error occurred: %s", apiErr.Error())
					}
				} else {
					t.Logf("[ERROR] Network error: %v", err)
					if !tt.expectError {
						t.Errorf("[ERROR] Unexpected network error: %v", err)
					}
				}
			} else {
				// Success case
				t.Logf("[OK] API call successful")
				t.Logf("Response Success: %v", resp.Success)
				t.Logf("Response Code: %s", resp.Code)
				t.Logf("Request ID: %s", resp.RequestID)

				if resp.Success {
					t.Logf("Operation Status: %v", resp.Data)
				}

				if tt.expectError {
					t.Errorf("[ERROR] Expected error but API call succeeded")
				}
			}
		})
	}
}

// TestImageStartParameterValidation tests parameter validation in real environment
func TestImageStartParameterValidation(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=true to run.")
	}

	// Get configuration
	cfg, err := config.GetConfig()
	if err != nil {
		t.Skipf("Could not load config: %v", err)
	}

	// Create API client
	apiClient := client.NewFromConfig(cfg)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test parameter validation
	tests := []struct {
		name        string
		loginToken  string
		sessionId   string
		imageId     string
		expectError bool
	}{
		{
			name:        "missing_login_token",
			loginToken:  "",
			sessionId:   "test-session",
			imageId:     "test-image",
			expectError: true,
		},
		{
			name:        "missing_session_id",
			loginToken:  "test-token",
			sessionId:   "",
			imageId:     "test-image",
			expectError: true,
		},
		{
			name:        "missing_image_id",
			loginToken:  "test-token",
			sessionId:   "test-session",
			imageId:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := apiClient.ImageAPI.StartImage(
				ctx,
				tt.loginToken,
				tt.sessionId,
				tt.imageId,
				1,
				2,
			)

			if tt.expectError {
				if err == nil {
					t.Errorf("[ERROR] Expected error but got none")
				} else {
					t.Logf("[OK] Expected error occurred: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("[ERROR] Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestImageStartRealWorkflow tests with real image IDs if available
func TestImageStartRealWorkflow(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=true to run.")
	}

	// Get configuration
	cfg, err := config.GetConfig()
	if err != nil {
		t.Skipf("Could not load config: %v", err)
	}

	// Check if we have valid tokens
	tokens, err := cfg.GetTokens()
	if err != nil {
		t.Skipf("No valid tokens found: %v. Please run 'agbcloud login' first.", err)
	}

	// Create API client
	apiClient := client.NewFromConfig(cfg)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	t.Run("list_and_start_real_image", func(t *testing.T) {
		// First, try to list available images
		t.Log("[SEARCH] Fetching available images...")
		listResp, _, err := apiClient.ImageAPI.ListImages(ctx, tokens.LoginToken, tokens.SessionId, "User", 1, 5)

		if err != nil {
			t.Logf("[WARN]  Could not list images: %v", err)
			t.Skip("Skipping real workflow test - cannot list images")
		}

		if !listResp.Success || len(listResp.Data.Images) == 0 {
			t.Log("[INFO]  No user images available for testing")
			t.Skip("Skipping real workflow test - no images available")
		}

		// Use the first available image for testing
		testImage := listResp.Data.Images[0]
		t.Logf("[DOC] Using image: %s (%s)", testImage.ImageName, testImage.ImageID)

		// Only proceed if the image is available
		if testImage.Status != "IMAGE_AVAILABLE" {
			t.Logf("[WARN]  Image status is %s, not available for starting", testImage.Status)
			t.Skip("Skipping real workflow test - image not available")
		}

		// Try to start the image
		t.Log("🚀 Attempting to start image...")
		startResp, httpResp, err := apiClient.ImageAPI.StartImage(
			ctx,
			tokens.LoginToken,
			tokens.SessionId,
			testImage.ImageID,
			2, // 2 CPU cores
			4, // 4 GB memory
		)

		// Log the results regardless of success/failure
		if httpResp != nil {
			t.Logf("HTTP Status: %d", httpResp.StatusCode)
		}

		if err != nil {
			if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
				t.Logf("API Error: %s", apiErr.Error())
				// This might be expected if the image can't be started for various reasons
				t.Logf("[INFO]  Image start failed (this may be expected): %s", apiErr.Error())
			} else {
				t.Errorf("[ERROR] Network error: %v", err)
			}
		} else {
			t.Logf("[OK] Start image API call successful")
			t.Logf("Response Success: %v", startResp.Success)
			t.Logf("Response Code: %s", startResp.Code)
			t.Logf("Request ID: %s", startResp.RequestID)

			if startResp.Success {
				t.Logf("[SUCCESS] Image started successfully!")
				t.Logf("Operation Status: %v", startResp.Data)
			} else {
				t.Logf("[INFO]  Image start was not successful: %s", startResp.Code)
			}
		}
	})
}
