// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
	log "github.com/sirupsen/logrus"
)

// RefreshTokenIfNeeded checks and refreshes token if it's about to expire (within 5 minutes)
// This provides automatic token management for seamless API access
func RefreshTokenIfNeeded(ctx context.Context) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	// Check if we have tokens
	if cfg.Token == nil {
		return fmt.Errorf("no valid token found, use 'agbcloud-cli login' to reauthenticate")
	}

	// Check if token is about to expire (within 5 minutes)
	if time.Until(cfg.Token.ExpiresAt) > 5*time.Minute {
		log.Debug("Token is still valid, no refresh needed")
		return nil
	}

	log.Info("Token is approaching expiry, refreshing...")

	// Create API client for refresh call
	apiClient := client.NewFromConfig(cfg)

	// Perform token refresh
	response, _, err := apiClient.OAuthAPI.RefreshToken(ctx,
		cfg.Token.KeepAliveToken,
		cfg.Token.SessionId)
	if err != nil {
		// If refresh fails, clear the tokens
		cfg.ClearTokens()
		return fmt.Errorf("use 'agbcloud-cli login' to reauthenticate: %w", err)
	}

	if !response.Success {
		// If refresh fails, clear the tokens
		cfg.ClearTokens()
		return fmt.Errorf("use 'agbcloud-cli login' to reauthenticate: refresh failed with code %s", response.Code)
	}

	// Save new tokens
	err = cfg.SaveTokens(
		response.Data.LoginToken,
		response.Data.SessionId,
		response.Data.KeepAliveToken,
		response.Data.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save refreshed tokens: %w", err)
	}

	log.Info("Token refreshed successfully")
	return nil
}

// HandleAPIError handles API errors and attempts token refresh if needed
func HandleAPIError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// Check if it's an authentication error that might require token refresh
	if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
		errorBody := string(apiErr.Body())

		// Check for token expiry indicators
		if isTokenExpiredError(errorBody) {
			log.Warn("Received token expiry error, attempting refresh...")

			// Try to refresh token
			refreshErr := RefreshTokenIfNeeded(ctx)
			if refreshErr != nil {
				log.Error("Token refresh failed:", refreshErr)
				return fmt.Errorf("authentication failed and token refresh unsuccessful: %w. Run 'agbcloud-cli login' to reauthenticate", refreshErr)
			}

			log.Info("Token refreshed successfully, please retry your request")
			return fmt.Errorf("token was refreshed, please retry your request")
		}
	}

	return err
}

// isTokenExpiredError checks if the error indicates token expiry
func isTokenExpiredError(errorBody string) bool {
	// Check for common token expiry indicators in AgbCloud API responses
	expiredIndicators := []string{
		"UserLogin.Expired",
		"Token is expired",
		"INVALID_TOKEN",
		"Invalid or expired tokens",
		"token expired",
		"authentication failed",
	}

	errorBodyLower := strings.ToLower(errorBody)
	for _, indicator := range expiredIndicators {
		if strings.Contains(errorBodyLower, strings.ToLower(indicator)) {
			return true
		}
	}

	return false
}

// ClearAuthTokens clears stored tokens (for logout or when tokens become invalid)
func ClearAuthTokens() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	// Clear tokens
	return cfg.ClearTokens()
}
