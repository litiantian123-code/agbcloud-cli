// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"net/http"
	"time"

	"github.com/liyuebing/agbcloud-cli/internal/config"
)

// NewFromConfig creates a new API client from the CLI configuration
func NewFromConfig(cfg *config.Config) *APIClient {
	configuration := NewConfiguration()

	// Set the server URL from config
	if cfg.Endpoint != "" {
		configuration.Servers[0].URL = cfg.Endpoint
	}

	// Set API key if available
	if cfg.APIKey != "" {
		configuration.APIKey = cfg.APIKey
		configuration.AddDefaultHeader("Authorization", "Bearer "+cfg.APIKey)
	}

	// Set up HTTP client with reasonable defaults
	configuration.HTTPClient = &http.Client{
		Timeout: 30 * time.Second,
	}

	return NewAPIClient(configuration)
}

// NewDefault creates a new API client with default configuration
func NewDefault() *APIClient {
	return NewFromConfig(config.DefaultConfig())
}
