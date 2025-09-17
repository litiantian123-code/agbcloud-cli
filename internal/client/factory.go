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

// retryTransport implements http.RoundTripper to integrate retry logic
type retryTransport struct {
	retryClient *RetryableHTTPClient
}

// RoundTrip implements the http.RoundTripper interface
func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.retryClient.Do(req)
}

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

	// Create base HTTP client with optional SSL verification skip
	baseClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Check if SSL verification should be skipped
	if shouldSkipSSLVerification() {
		baseClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	// Wrap with retry functionality
	retryClient := NewRetryableHTTPClient(baseClient, DefaultRetryConfig())

	// Create a wrapper that implements http.Client interface
	configuration.HTTPClient = &http.Client{
		Timeout: baseClient.Timeout,
		Transport: &retryTransport{
			retryClient: retryClient,
		},
	}

	return NewAPIClient(configuration)
}

// NewDefault creates a new API client with default configuration
func NewDefault() *APIClient {
	return NewFromConfig(config.DefaultConfig())
}
