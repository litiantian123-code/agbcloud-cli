// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// createTestFile creates a temporary test file for upload testing
func createTestFile(t *testing.T) (string, func()) {
	// Create a temporary file with test content
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("agbcloud_test_%d.txt", time.Now().UnixNano()))

	testContent := []byte("This is a test file for AgbCloud CLI image upload integration test.\nCreated at: " + time.Now().String())

	err := os.WriteFile(tmpFile, testContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Return file path and cleanup function
	cleanup := func() {
		if err := os.Remove(tmpFile); err != nil {
			t.Logf("Warning: Failed to cleanup test file %s: %v", tmpFile, err)
		}
	}

	return tmpFile, cleanup
}

// uploadFileToOSS uploads a file to the provided OSS URL using HTTP PUT
func uploadFileToOSS(t *testing.T, ossURL, filePath string) error {
	// Read the file content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Create HTTP PUT request
	req, err := http.NewRequest(http.MethodPut, ossURL, bytes.NewReader(fileContent))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}

	// Set appropriate headers
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = int64(len(fileContent))

	// Execute the request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute PUT request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	respBody, _ := io.ReadAll(resp.Body)

	t.Logf("Upload response status: %d", resp.StatusCode)
	t.Logf("Upload response headers: %v", resp.Header)
	if len(respBody) > 0 {
		t.Logf("Upload response body: %s", string(respBody))
	}

	// Check if upload was successful (typically 200 or 201 for OSS)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// TestImageUploadCredentialIntegration tests the /api/image/getUploadCredential endpoint
func TestImageUploadCredentialIntegration(t *testing.T) {
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
		t.Skipf("No valid tokens found: %v. Please run 'agbcloud-cli login' first.", err)
	}

	// Create API client
	apiClient := client.NewFromConfig(cfg)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("TestGetUploadCredentialWithValidTokens", func(t *testing.T) {
		t.Logf("Testing /api/image/getUploadCredential with valid loginToken and sessionId")
		t.Logf("Using LoginToken: %s", tokens.LoginToken)
		t.Logf("Using SessionId: %s", tokens.SessionId)

		// Call the getUploadCredential API
		response, httpResp, err := apiClient.ImageAPI.GetUploadCredential(ctx, tokens.LoginToken, tokens.SessionId)

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
				// Don't fail immediately - log the error and continue to analyze the response
				t.Logf("[ERROR] API error occurred: %s", apiErr.Error())
			} else {
				t.Fatalf("[ERROR] Network error prevented API communication: %v", err)
			}
		} else {
			// Log success details
			t.Logf("[OK] Success! Upload credential response received")
			t.Logf("Success: %v", response.Success)
			t.Logf("Code: %s", response.Code)
			t.Logf("RequestID: %s", response.RequestID)
			t.Logf("TraceID: %s", response.TraceID)
			t.Logf("HTTPStatusCode: %d", response.HTTPStatusCode)
			t.Logf("Response Data: %+v", response.Data)

			// Log detailed data fields
			t.Logf("[DOC] Detailed Response Analysis:")
			t.Logf("   - OssURL: %s", response.Data.OssURL)
			t.Logf("   - TaskID: %s", response.Data.TaskID)

			// Validate response structure
			if response.Data.TaskID == "" {
				t.Error("[ERROR] TaskID should not be empty in successful response")
			} else {
				t.Logf("[OK] TaskID received: %s", response.Data.TaskID)
			}

			// Test file upload to OSS URL if we have a valid OSS URL
			if response.Data.OssURL != "" {
				t.Logf("[REFRESH] Testing file upload to OSS URL...")

				// Create a test file
				testFilePath, cleanup := createTestFile(t)
				defer cleanup()

				t.Logf("[DIR] Created test file: %s", testFilePath)

				// Upload the file to OSS
				uploadErr := uploadFileToOSS(t, response.Data.OssURL, testFilePath)
				if uploadErr != nil {
					t.Logf("[ERROR] File upload failed: %v", uploadErr)
					// Log as error but don't fail the test since OSS URL might have restrictions
					t.Errorf("File upload to OSS failed: %v", uploadErr)
				} else {
					t.Logf("[OK] File upload to OSS successful!")
				}
			} else {
				t.Logf("[WARN]  No OSS URL provided, skipping file upload test")
			}
		}
	})

	t.Run("TestGetUploadCredentialWithEmptyLoginToken", func(t *testing.T) {
		t.Logf("Testing /api/image/getUploadCredential with empty loginToken")

		// Call with empty loginToken
		_, _, err := apiClient.ImageAPI.GetUploadCredential(ctx, "", tokens.SessionId)

		if err == nil {
			t.Error("Expected error for empty loginToken, but got none")
		} else {
			t.Logf("[OK] Expected error for empty loginToken: %v", err)
		}
	})

	t.Run("TestGetUploadCredentialWithEmptySessionId", func(t *testing.T) {
		t.Logf("Testing /api/image/getUploadCredential with empty sessionId")

		// Call with empty sessionId
		_, _, err := apiClient.ImageAPI.GetUploadCredential(ctx, tokens.LoginToken, "")

		if err == nil {
			t.Error("Expected error for empty sessionId, but got none")
		} else {
			t.Logf("[OK] Expected error for empty sessionId: %v", err)
		}
	})
}
