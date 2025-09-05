// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// shouldSkipSSLVerification determines whether SSL verification should be skipped
// based on the environment variable only
func shouldSkipSSLVerification() bool {
	// Check explicit user setting
	if skipSSL := os.Getenv("AGB_CLI_SKIP_SSL_VERIFY"); skipSSL != "" {
		return skipSSL == "true"
	}

	// Default to SSL verification
	return false
}

// NewFromConfig creates a new API client from the CLI configuration
func NewFromConfig(cfg *config.Config) *APIClient {
	configuration := NewConfiguration()

	// Set the server URL from config
	if cfg.Endpoint != "" {
		configuration.Servers[0].URL = cfg.Endpoint
	}

	// Create HTTP client with optional SSL verification skip
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Check if SSL verification should be skipped
	if shouldSkipSSLVerification() {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	configuration.HTTPClient = httpClient

	return NewAPIClient(configuration)
}

// NewDefault creates a new API client with default configuration
func NewDefault() *APIClient {
	return NewFromConfig(config.DefaultConfig())
}
