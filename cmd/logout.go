// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

var LogoutCmd = &cobra.Command{
	Use:     "logout",
	Short:   "Log out from AgbCloud",
	Long:    "Log out from AgbCloud by invalidating server session and clearing local authentication data",
	Args:    cobra.NoArgs,
	GroupID: "core",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLogout(cmd)
	},
}

func init() {
	// No flags needed for logout command
}

func runLogout(cmd *cobra.Command) error {
	fmt.Println("🔓 Logging out from AgbCloud...")

	// Load configuration
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if we have valid tokens for API logout
	hasValidTokens := cfg.Token != nil &&
		cfg.Token.LoginToken != "" &&
		cfg.Token.SessionId != ""

	if hasValidTokens {
		// Attempt to invalidate server session
		fmt.Println("🌐 Invalidating server session...")

		apiClient := client.NewFromConfig(cfg)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Call logout API
		response, httpResp, err := apiClient.OAuthAPI.Logout(ctx,
			cfg.Token.LoginToken,
			cfg.Token.SessionId)

		if err != nil {
			// Log warning but continue with local cleanup
			fmt.Printf("⚠️  Warning: Could not invalidate server session: %v\n", err)
			if httpResp != nil {
				fmt.Printf("📊 HTTP Status: %d\n", httpResp.StatusCode)
			}
		} else if !response.Success {
			// API call succeeded but logout failed
			fmt.Printf("⚠️  Warning: Server session invalidation failed (Code: %s)\n", response.Code)
		} else {
			// Success
			fmt.Println("✅ Server session invalidated successfully")
		}
	} else {
		fmt.Println("ℹ️  No active session found")
	}

	// Always perform local cleanup
	fmt.Println("🧹 Clearing local authentication data...")

	// Clear tokens from config
	err = cfg.ClearTokens()
	if err != nil {
		return fmt.Errorf("failed to clear local authentication data: %w", err)
	}

	// Success message
	if hasValidTokens {
		fmt.Println("✅ Successfully logged out from AgbCloud")
	} else {
		fmt.Println("✅ Successfully logged out from AgbCloud (local session cleared)")
	}

	return nil
}
