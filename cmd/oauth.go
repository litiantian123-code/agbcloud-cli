// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/liyuebing/agbcloud-cli/internal/client"
	"github.com/liyuebing/agbcloud-cli/internal/config"
)

var OAuthCmd = &cobra.Command{
	Use:   "oauth",
	Short: "OAuth authentication commands",
	Long:  "Commands for OAuth authentication with various providers",
}

var googleLoginCmd = &cobra.Command{
	Use:   "google-login [fromUrlPath]",
	Short: "Get Google OAuth login URL",
	Long:  "Retrieve the Google OAuth login URL for authentication",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get fromUrlPath from args or use default
		fromUrlPath := "https://agb.cloud"
		if len(args) > 0 {
			fromUrlPath = args[0]
		}

		// Create client configuration for OAuth (no API key needed)
		cfg := &config.Config{
			Endpoint: "", // Will use agb.cloud fallback for OAuth
		}

		apiClient := client.NewFromConfig(cfg)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Printf("🔐 Requesting Google OAuth login URL...\n")
		fmt.Printf("📍 From URL Path: %s\n\n", fromUrlPath)

		// Call the OAuth API
		response, httpResp, err := apiClient.OAuthAPI.GetGoogleLoginURL(ctx, fromUrlPath)

		if err != nil {
			if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
				fmt.Printf("❌ API Error: %s\n", apiErr.Error())
				if httpResp != nil {
					fmt.Printf("📊 Status Code: %d\n", httpResp.StatusCode)
					fmt.Printf("🌐 Request URL: %s\n", httpResp.Request.URL.String())
				}
				if len(apiErr.Body()) > 0 {
					fmt.Printf("📄 Response Body: %s\n", string(apiErr.Body()))
				}
				return fmt.Errorf("OAuth API error: %s", apiErr.Error())
			}
			return fmt.Errorf("network error: %v", err)
		}

		// Success - display the results
		fmt.Printf("✅ Success!\n\n")
		fmt.Printf("📋 Response Details:\n")
		fmt.Printf("  Code: %s\n", response.Code)
		fmt.Printf("  Success: %t\n", response.Success)
		fmt.Printf("  Request ID: %s\n", response.RequestID)
		fmt.Printf("  Trace ID: %s\n", response.TraceID)
		fmt.Printf("  HTTP Status Code: %d\n", response.HTTPStatusCode)
		fmt.Printf("\n🔗 Google OAuth URL:\n")
		fmt.Printf("  %s\n\n", response.Data.InvokeURL)

		fmt.Printf("💡 You can copy this URL and paste it into your browser to start the OAuth flow.\n")

		return nil
	},
}

func init() {
	OAuthCmd.AddCommand(googleLoginCmd)
}
