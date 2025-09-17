// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// TestImageListIntegration tests the complete image list workflow using real API calls
func TestImageListIntegration(t *testing.T) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("APIClientListImages", func(t *testing.T) {
		t.Logf("🚀 Testing ListImages API with real backend")
		t.Logf("Using LoginToken: %s", tokens.LoginToken)
		t.Logf("Using SessionId: %s", tokens.SessionId)

		// Test ListImages API call
		t.Logf("\n[DOC] Calling ListImages API...")
		response, httpResp, err := apiClient.ImageAPI.ListImages(ctx, tokens.LoginToken, tokens.SessionId, "User", 1, 10)

		if err != nil {
			if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
				t.Logf("API Error: %s", apiErr.Error())
				t.Logf("Response Body: %s", string(apiErr.Body()))
				if httpResp != nil {
					t.Logf("HTTP Status Code: %d", httpResp.StatusCode)
				}
				t.Fatalf("[ERROR] Failed to list images: %s", apiErr.Error())
			} else {
				t.Fatalf("[ERROR] Network error listing images: %v", err)
			}
		}

		if !response.Success {
			t.Fatalf("[ERROR] ListImages request failed: %+v", response)
		}

		t.Logf("[OK] ListImages API call successful")
		t.Logf("   - Total images: %d", response.Data.Total)
		t.Logf("   - Page: %d", response.Data.Page)
		t.Logf("   - Page size: %d", response.Data.PageSize)
		t.Logf("   - Images returned: %d", len(response.Data.Images))

		// Log image details
		for i, image := range response.Data.Images {
			t.Logf("   - Image %d: %s (%s) - %s", i+1, image.ImageName, image.ImageID, image.Status)
		}

		// Validate response structure
		if response.Data.Page != 1 {
			t.Errorf("Expected page 1, got %d", response.Data.Page)
		}
		if response.Data.PageSize != 10 {
			t.Errorf("Expected page size 10, got %d", response.Data.PageSize)
		}
		if response.Data.Total < 0 {
			t.Errorf("Total should be non-negative, got %d", response.Data.Total)
		}
		if len(response.Data.Images) > response.Data.PageSize {
			t.Errorf("Images count (%d) should not exceed page size (%d)", len(response.Data.Images), response.Data.PageSize)
		}
	})

	t.Run("PaginationTest", func(t *testing.T) {
		t.Logf("[REFRESH] Testing pagination functionality")

		// Test different page sizes
		pageSizes := []int{1, 5, 20}
		for _, pageSize := range pageSizes {
			t.Logf("\n[PAGE] Testing with page size: %d", pageSize)
			response, _, err := apiClient.ImageAPI.ListImages(ctx, tokens.LoginToken, tokens.SessionId, "User", 1, pageSize)

			if err != nil {
				t.Logf("[WARN]  Error with page size %d: %v", pageSize, err)
				continue
			}

			if !response.Success {
				t.Logf("[WARN]  Failed request with page size %d: %+v", pageSize, response)
				continue
			}

			t.Logf("   [OK] Page size %d: got %d images (total: %d)", pageSize, len(response.Data.Images), response.Data.Total)

			// Validate page size
			if response.Data.PageSize != pageSize {
				t.Errorf("Expected page size %d, got %d", pageSize, response.Data.PageSize)
			}

			// Validate images count doesn't exceed page size
			if len(response.Data.Images) > pageSize {
				t.Errorf("Images count (%d) should not exceed requested page size (%d)", len(response.Data.Images), pageSize)
			}
		}
	})

	t.Run("SystemImageTypeTest", func(t *testing.T) {
		t.Logf("🖥️  Testing System image type")

		// Test System image type
		t.Logf("\n[DOC] Calling ListImages API with System type...")
		response, httpResp, err := apiClient.ImageAPI.ListImages(ctx, tokens.LoginToken, tokens.SessionId, "System", 1, 10)

		if err != nil {
			if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
				t.Logf("API Error: %s", apiErr.Error())
				t.Logf("Response Body: %s", string(apiErr.Body()))
				if httpResp != nil {
					t.Logf("HTTP Status Code: %d", httpResp.StatusCode)
				}
				t.Fatalf("[ERROR] Failed to list System images: %s", apiErr.Error())
			} else {
				t.Fatalf("[ERROR] Network error listing System images: %v", err)
			}
		}

		if !response.Success {
			t.Fatalf("[ERROR] ListImages request for System type failed: %+v", response)
		}

		t.Logf("[OK] ListImages API call for System type successful")
		t.Logf("   - Total System images: %d", response.Data.Total)
		t.Logf("   - Page: %d", response.Data.Page)
		t.Logf("   - Page size: %d", response.Data.PageSize)
		t.Logf("   - System images returned: %d", len(response.Data.Images))

		// Log System image details
		for i, image := range response.Data.Images {
			t.Logf("   - System Image %d: %s (%s) - %s", i+1, image.ImageName, image.ImageID, image.Status)
		}

		// Validate response structure
		if response.Data.Page != 1 {
			t.Errorf("Expected page 1, got %d", response.Data.Page)
		}
		if response.Data.PageSize != 10 {
			t.Errorf("Expected page size 10, got %d", response.Data.PageSize)
		}
		if response.Data.Total < 0 {
			t.Errorf("Total should be non-negative, got %d", response.Data.Total)
		}
		if len(response.Data.Images) > response.Data.PageSize {
			t.Errorf("Images count (%d) should not exceed page size (%d)", len(response.Data.Images), response.Data.PageSize)
		}
	})
}

// TestImageListCLIIntegration tests the CLI command end-to-end
func TestImageListCLIIntegration(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=true to run.")
	}

	// Check if we have valid tokens by trying to load config
	cfg, err := config.GetConfig()
	if err != nil {
		t.Skipf("Could not load config: %v", err)
	}

	_, err = cfg.GetTokens()
	if err != nil {
		t.Skipf("No valid tokens found: %v. Please run 'agbcloud login' first.", err)
	}

	// Build the CLI binary for testing
	t.Logf("🔨 Building CLI binary for testing...")
	buildCmd := exec.Command("go", "build", "-o", "agbcloud-test", ".")
	buildCmd.Dir = "../.." // Go up to project root
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI binary: %v", err)
	}
	defer os.Remove("../../agbcloud-test") // Clean up

	t.Run("CLIListImagesDefault", func(t *testing.T) {
		t.Logf("🖥️  Testing CLI: agbcloud image list")

		cmd := exec.Command("./agbcloud-test", "image", "list")
		cmd.Dir = "../.." // Run from project root

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		t.Logf("Command output (stdout):\n%s", stdout.String())
		if stderr.Len() > 0 {
			t.Logf("Command output (stderr):\n%s", stderr.String())
		}

		if err != nil {
			t.Fatalf("CLI command failed: %v", err)
		}

		output := stdout.String()

		// Validate output contains expected elements
		if !strings.Contains(output, "[DOC] Listing User images") {
			t.Errorf("Output should contain listing message")
		}
		if !strings.Contains(output, "IMAGE ID") {
			t.Errorf("Output should contain table header")
		}
		if !strings.Contains(output, "[OK] Found") {
			t.Errorf("Output should contain success message")
		}

		t.Logf("[OK] CLI command executed successfully")
	})

	t.Run("CLIListImagesWithFlags", func(t *testing.T) {
		t.Logf("🖥️  Testing CLI: agbcloud image list --page 1 --size 5 --type User")

		cmd := exec.Command("./agbcloud-test", "image", "list", "--page", "1", "--size", "5", "--type", "User")
		cmd.Dir = "../.." // Run from project root

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		t.Logf("Command output (stdout):\n%s", stdout.String())
		if stderr.Len() > 0 {
			t.Logf("Command output (stderr):\n%s", stderr.String())
		}

		if err != nil {
			t.Fatalf("CLI command with flags failed: %v", err)
		}

		output := stdout.String()

		// Validate output contains expected elements
		if !strings.Contains(output, "[DOC] Listing User images (Page 1, Size 5)") {
			t.Errorf("Output should contain correct listing message with parameters")
		}
		if !strings.Contains(output, "Page Size: 5") {
			t.Errorf("Output should show correct page size")
		}

		t.Logf("[OK] CLI command with flags executed successfully")
	})

	t.Run("CLIListImagesShortFlags", func(t *testing.T) {
		t.Logf("🖥️  Testing CLI: agbcloud image list -p 1 -s 3 -t User")

		cmd := exec.Command("./agbcloud-test", "image", "list", "-p", "1", "-s", "3", "-t", "User")
		cmd.Dir = "../.." // Run from project root

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		t.Logf("Command output (stdout):\n%s", stdout.String())
		if stderr.Len() > 0 {
			t.Logf("Command output (stderr):\n%s", stderr.String())
		}

		if err != nil {
			t.Fatalf("CLI command with short flags failed: %v", err)
		}

		output := stdout.String()

		// Validate output contains expected elements
		if !strings.Contains(output, "[DOC] Listing User images (Page 1, Size 3)") {
			t.Errorf("Output should contain correct listing message with parameters")
		}
		if !strings.Contains(output, "Page Size: 3") {
			t.Errorf("Output should show correct page size")
		}

		t.Logf("[OK] CLI command with short flags executed successfully")
	})

	t.Run("CLIListSystemImages", func(t *testing.T) {
		t.Logf("🖥️  Testing CLI: agbcloud image list --type System")

		cmd := exec.Command("./agbcloud-test", "image", "list", "--type", "System")
		cmd.Dir = "../.." // Run from project root

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		t.Logf("Command output (stdout):\n%s", stdout.String())
		if stderr.Len() > 0 {
			t.Logf("Command output (stderr):\n%s", stderr.String())
		}

		if err != nil {
			t.Fatalf("CLI command for System images failed: %v", err)
		}

		output := stdout.String()

		// Validate output contains expected elements
		if !strings.Contains(output, "[DOC] Listing System images") {
			t.Errorf("Output should contain System images listing message")
		}
		if !strings.Contains(output, "IMAGE ID") {
			t.Errorf("Output should contain table header")
		}
		if !strings.Contains(output, "[OK] Found") {
			t.Errorf("Output should contain success message")
		}

		t.Logf("[OK] CLI command for System images executed successfully")
	})
}

// TestImageListErrorHandling tests error scenarios
func TestImageListErrorHandling(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=true to run.")
	}

	// Create API client with default config (no auth)
	apiClient := client.NewDefault()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("ListImagesWithoutAuth", func(t *testing.T) {
		t.Logf("🔒 Testing ListImages without authentication")

		// This should fail because we don't have valid tokens
		_, _, err := apiClient.ImageAPI.ListImages(ctx, "", "", "User", 1, 10)

		if err == nil {
			t.Errorf("Expected error when calling ListImages without auth, but got none")
		} else {
			t.Logf("[OK] Correctly received error without auth: %v", err)
		}
	})

	t.Run("ListImagesWithInvalidParameters", func(t *testing.T) {
		t.Logf("🚫 Testing ListImages with invalid parameters")

		// Test invalid page
		_, _, err := apiClient.ImageAPI.ListImages(ctx, "fake-token", "fake-session", "User", 0, 10)
		if err == nil {
			t.Errorf("Expected error for invalid page (0), but got none")
		} else {
			t.Logf("[OK] Correctly received error for invalid page: %v", err)
		}

		// Test invalid page size
		_, _, err = apiClient.ImageAPI.ListImages(ctx, "fake-token", "fake-session", "User", 1, 0)
		if err == nil {
			t.Errorf("Expected error for invalid page size (0), but got none")
		} else {
			t.Logf("[OK] Correctly received error for invalid page size: %v", err)
		}

		// Test empty image type
		_, _, err = apiClient.ImageAPI.ListImages(ctx, "fake-token", "fake-session", "", 1, 10)
		if err == nil {
			t.Errorf("Expected error for empty image type, but got none")
		} else {
			t.Logf("[OK] Correctly received error for empty image type: %v", err)
		}
	})
}
