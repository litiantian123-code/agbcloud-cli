// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package apiclient

import (
	"context"

	"github.com/agbcloud/agbcloud-cli/internal/auth"
	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

var apiClient *client.APIClient

// NewClient creates an API client with automatic token refresh for seamless authentication
func NewClient(defaultHeaders map[string]string) (*client.APIClient, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	// If we have an existing client using token auth, refresh if needed
	if apiClient != nil {
		err := auth.RefreshTokenIfNeeded(context.Background())
		if err != nil {
			return nil, err
		}
		return apiClient, nil
	}

	// Create new API client using the factory
	newApiClient := client.NewFromConfig(cfg)

	// Add any additional default headers
	for headerKey, headerValue := range defaultHeaders {
		newApiClient.GetConfig().AddDefaultHeader(headerKey, headerValue)
	}

	// If using token auth, refresh if needed before returning
	if cfg.Token != nil {
		err = auth.RefreshTokenIfNeeded(context.Background())
		if err != nil {
			return nil, err
		}
	}

	apiClient = newApiClient
	return apiClient, nil
}

// NewClientWithDefaults creates an API client with default headers
func NewClientWithDefaults() (*client.APIClient, error) {
	defaultHeaders := map[string]string{
		"X-AgbCloud-Source": "cli",
	}
	return NewClient(defaultHeaders)
}

// HandleAPIError handles API errors and attempts token refresh if needed
func HandleAPIError(ctx context.Context, err error) error {
	return auth.HandleAPIError(ctx, err)
}

// EnsureValidToken ensures the current token is valid and refreshes if needed
func EnsureValidToken(ctx context.Context) error {
	return auth.RefreshTokenIfNeeded(ctx)
}

// ClearCachedClient clears the cached API client (useful for testing or logout)
func ClearCachedClient() {
	apiClient = nil
}
