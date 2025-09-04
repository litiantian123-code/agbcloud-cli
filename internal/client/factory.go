// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// shouldSkipSSLVerification determines whether SSL verification should be skipped
// based on the endpoint and environment
func shouldSkipSSLVerification(endpoint string) bool {
	// Always respect explicit user setting
	if skipSSL := os.Getenv("AGB_CLI_SKIP_SSL_VERIFY"); skipSSL != "" {
		return skipSSL == "true"
	}

	// Parse the endpoint URL
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		// If we can't parse the URL, be conservative and verify SSL
		return false
	}

	hostname := parsedURL.Hostname()

	// Skip SSL verification for IP addresses (certificates usually don't include IP SANs)
	if net.ParseIP(hostname) != nil {
		return true
	}

	// Skip SSL verification for localhost and local development domains
	if isLocalDevelopmentHost(hostname) {
		return true
	}

	// Skip SSL verification for non-standard ports (likely development environments)
	if parsedURL.Port() != "" && parsedURL.Port() != "443" {
		return true
	}

	// For production domains, always verify SSL
	return false
}

// isLocalDevelopmentHost checks if the hostname is a local development host
func isLocalDevelopmentHost(hostname string) bool {
	localHosts := []string{
		"localhost",
		"127.0.0.1",
		"::1",
	}

	// Check exact matches
	for _, localHost := range localHosts {
		if hostname == localHost {
			return true
		}
	}

	// Check for .local domains (common in development)
	if strings.HasSuffix(hostname, ".local") {
		return true
	}

	// Check for internal/private network domains
	if strings.HasSuffix(hostname, ".internal") ||
		strings.HasSuffix(hostname, ".dev") ||
		strings.HasSuffix(hostname, ".test") {
		return true
	}

	return false
}

// isDebuggerEnvironment detects if we're running in a debugger environment
func isDebuggerEnvironment() bool {
	// Check for Delve debugger environment variable
	if os.Getenv("DELVE_DEBUGGER") == "true" {
		return true
	}

	// Check command line arguments for debugger indicators
	for _, arg := range os.Args {
		if strings.Contains(arg, "dlv") ||
			strings.Contains(arg, "__debug_bin") ||
			strings.Contains(arg, "debug_test") {
			return true
		}
	}

	// Check if we're running under a debugger by looking at the process name
	// This is a heuristic approach for common debugger scenarios
	return false
}

// createDebuggerFriendlyTransport creates a transport optimized for debugger environments
func createDebuggerFriendlyTransport(skipSSL bool) *http.Transport {
	transport := &http.Transport{
		// Completely disable connection pooling to avoid file descriptor issues
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 0,
		MaxConnsPerHost:     1,
		IdleConnTimeout:     1 * time.Second,

		// Use a very simple dialer with minimal features
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second, // Very short timeout
			KeepAlive: 0,               // Disable keep-alive completely
		}).DialContext,

		// Shorter timeouts for faster failure
		TLSHandshakeTimeout:   3 * time.Second,
		ResponseHeaderTimeout: 3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		// Disable all advanced features that might cause issues in debugger
		ForceAttemptHTTP2:  false,
		DisableKeepAlives:  true,
		DisableCompression: true,

		// Disable proxy to avoid additional complexity
		Proxy: nil,
	}

	if skipSSL {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
			// Use minimal TLS configuration for debugger compatibility
			MinVersion: tls.VersionTLS10,
		}
	}

	return transport
}

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

	// Determine SSL verification strategy
	skipSSL := shouldSkipSSLVerification(cfg.Endpoint)

	// Check if we're in a debugger environment and use appropriate configuration
	var httpClient *http.Client

	if isDebuggerEnvironment() {
		// Use debugger-friendly configuration with minimal features
		transport := createDebuggerFriendlyTransport(skipSSL)
		httpClient = &http.Client{
			Timeout:   10 * time.Second, // Shorter timeout for debugger
			Transport: transport,
		}
	} else {
		// Use standard configuration for normal environments
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}

		if skipSSL {
			httpClient.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}
	}

	configuration.HTTPClient = httpClient

	return NewAPIClient(configuration)
}

// NewDefault creates a new API client with default configuration
func NewDefault() *APIClient {
	return NewFromConfig(config.DefaultConfig())
}
