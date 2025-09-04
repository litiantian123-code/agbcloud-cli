// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to AgbCloud",
	Long:  "Authenticate with AgbCloud using OAuth in your browser",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLogin(cmd)
	},
}

func init() {
	// No flags needed for login command
}

func runLogin(cmd *cobra.Command) error {
	fmt.Println("🔐 Starting AgbCloud authentication...")

	// Create client configuration for OAuth (no API key needed for OAuth)
	cfg := config.DefaultConfig()
	// Clear API key for OAuth requests (OAuth doesn't need authentication)
	cfg.APIKey = ""

	apiClient := client.NewFromConfig(cfg)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("🌐 Requesting OAuth login URL...")

	// Get the OAuth URL from the API using the new login provider endpoint
	response, httpResp, err := apiClient.OAuthAPI.GetLoginProviderURL(ctx, "https://agb.cloud", "CLI", "GOOGLE_LOCALHOST")
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("❌ API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("📊 Status Code: %d\n", httpResp.StatusCode)
			}
			if len(apiErr.Body()) > 0 {
				fmt.Printf("📄 Response Body: %s\n", string(apiErr.Body()))
			}
			return fmt.Errorf("failed to get OAuth URL: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	// Verify we got a successful response
	if !response.Success {
		return fmt.Errorf("OAuth request failed: %s", response.Code)
	}

	if response.Data.InvokeURL == "" {
		return fmt.Errorf("received empty OAuth URL from server")
	}

	fmt.Println("✅ Successfully retrieved OAuth URL!")
	fmt.Printf("📋 Request ID: %s\n", response.RequestID)
	fmt.Printf("🔍 Trace ID: %s\n", response.TraceID)
	fmt.Println()

	// Display the URL
	fmt.Println("🔗 OAuth URL:")
	fmt.Printf("  %s\n\n", response.Data.InvokeURL)

	// Open the URL in the browser automatically
	fmt.Println("🌐 Opening the browser for authentication...")
	fmt.Println()
	fmt.Println("If the browser doesn't open automatically, please copy and paste the URL above.")

	err = browser.OpenURL(response.Data.InvokeURL)
	if err != nil {
		fmt.Printf("⚠️  Failed to open browser automatically: %v\n", err)
		fmt.Println("💡 Please copy the URL above and paste it into your browser to complete authentication.")
		return nil
	}

	fmt.Println("✅ Browser opened successfully!")
	fmt.Println("📝 Please complete the authentication process in your browser.")
	fmt.Println("🔄 After authentication, you'll be redirected back to AgbCloud.")

	return nil
}
