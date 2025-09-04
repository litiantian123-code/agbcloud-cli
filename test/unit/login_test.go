// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/agbcloud/agbcloud-cli/cmd"
	"github.com/agbcloud/agbcloud-cli/internal/client"
)

// TestLoginCommand tests the login command with a mock server
func TestLoginCommand(t *testing.T) {
	// Create a mock server that returns a successful OAuth response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request uses new endpoint
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

		// Return the expected response
		response := client.OAuthLoginProviderResponse{
			Code:      "success",
			RequestID: "test-request-id-123",
			Success:   true,
			Data: client.OAuthLoginProviderData{
				InvokeURL: "https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=test-client-id&redirect_uri=https://agb.cloud/api/oauth/google/login_callback&response_type=code&scope=https://www.googleapis.com/auth/userinfo.profile%20https://www.googleapis.com/auth/userinfo.email&state=https://agb.cloud",
			},
			TraceID:        "test-trace-id-456",
			HTTPStatusCode: 200,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Temporarily override the default endpoint for testing
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Create a root command for testing
	rootCmd := &cobra.Command{Use: "agbcloud"}

	// Create login command and modify it to use mock server
	loginCmd := cmd.LoginCmd

	// Override the RunE function to use mock server
	loginCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// Create client configuration pointing to mock server
		cfg := client.NewConfiguration()
		cfg.Servers[0].URL = mockServer.URL
		cfg.HTTPClient = &http.Client{Timeout: 5 * time.Second}

		apiClient := client.NewAPIClient(cfg)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Use default values since flags are removed
		fromUrlPath := "https://agb.cloud"
		noOpen := true // Always true in tests to avoid opening browser

		// Get the OAuth URL from the API
		response, _, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, fromUrlPath, "CLI", "GOOGLE_LOCALHOST")
		if err != nil {
			return err
		}

		// Verify response
		if !response.Success {
			cmd.Printf("OAuth request failed: %s", response.Code)
			return fmt.Errorf("OAuth request failed: %s", response.Code)
		}

		if response.Data.InvokeURL == "" {
			cmd.Printf("received empty OAuth URL from server")
			return fmt.Errorf("received empty OAuth URL from server")
		}

		// Print success information
		cmd.Printf("✅ Successfully retrieved OAuth URL!\n")
		cmd.Printf("📋 Request ID: %s\n", response.RequestID)
		cmd.Printf("🔍 Trace ID: %s\n", response.TraceID)
		cmd.Printf("🔗 OAuth URL: %s\n", response.Data.InvokeURL)

		if noOpen {
			cmd.Printf("💡 Browser opening disabled by --no-open flag\n")
		} else {
			cmd.Printf("🌐 Would open browser (disabled in test)\n")
		}

		return nil
	}

	rootCmd.AddCommand(loginCmd)

	// Test with --no-open flag
	t.Run("LoginWithNoOpen", func(t *testing.T) {
		// Capture output
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		// Set command arguments
		rootCmd.SetArgs([]string{"login"})

		// Execute command
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("Command execution failed: %v", err)
		}

		output := buf.String()

		// Verify output contains expected elements
		if !strings.Contains(output, "✅ Successfully retrieved OAuth URL!") {
			t.Error("Expected success message not found in output")
		}

		if !strings.Contains(output, "📋 Request ID: test-request-id-123") {
			t.Error("Expected request ID not found in output")
		}

		if !strings.Contains(output, "🔍 Trace ID: test-trace-id-456") {
			t.Error("Expected trace ID not found in output")
		}

		if !strings.Contains(output, "https://accounts.google.com/o/oauth2/auth") {
			t.Error("Expected OAuth URL not found in output")
		}

		if !strings.Contains(output, "💡 Browser opening disabled by --no-open flag") {
			t.Error("Expected no-open message not found in output")
		}

		t.Logf("✅ Login command test passed!")
		t.Logf("Output: %s", output)
	})

	// Test basic functionality (no custom flags needed)
	t.Run("LoginBasicFunctionality", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		rootCmd.SetArgs([]string{"login"})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("Command execution failed: %v", err)
		}

		output := buf.String()

		if !strings.Contains(output, "✅ Successfully retrieved OAuth URL!") {
			t.Error("Expected success message not found in output")
		}

		t.Logf("✅ Basic functionality test passed!")
	})
}

// TestLoginCommandErrorHandling tests error scenarios
func TestLoginCommandErrorHandling(t *testing.T) {
	// Create a mock server that returns an error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer mockServer.Close()

	rootCmd := &cobra.Command{Use: "agbcloud"}
	loginCmd := cmd.LoginCmd

	// Override to use mock server that returns error
	loginCmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg := client.NewConfiguration()
		cfg.Servers[0].URL = mockServer.URL
		cfg.HTTPClient = &http.Client{Timeout: 5 * time.Second}

		apiClient := client.NewAPIClient(cfg)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		fromUrlPath := "https://agb.cloud" // Use default value

		_, _, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, fromUrlPath, "CLI", "GOOGLE_LOCALHOST")
		if err != nil {
			return err
		}

		return nil
	}

	rootCmd.AddCommand(loginCmd)

	t.Run("LoginWithServerError", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		rootCmd.SetArgs([]string{"login"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("Expected command to fail with server error")
		}

		if !strings.Contains(err.Error(), "500") {
			t.Errorf("Expected 500 error, got: %v", err)
		}

		t.Logf("✅ Error handling test passed!")
	})
}
