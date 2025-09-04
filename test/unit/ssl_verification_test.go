// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"net/http"
	"os"
	"testing"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// TestSSLVerificationStrategy tests the automatic SSL verification strategy
func TestSSLVerificationStrategy(t *testing.T) {
	// Save original environment variable
	originalSkipSSL := os.Getenv("AGB_CLI_SKIP_SSL_VERIFY")
	defer func() {
		os.Setenv("AGB_CLI_SKIP_SSL_VERIFY", originalSkipSSL)
	}()

	testCases := []struct {
		name            string
		endpoint        string
		envVar          string
		expectedSkipSSL bool
		description     string
	}{
		// Production domains - should verify SSL
		{
			name:            "ProductionDomain",
			endpoint:        "https://agb.cloud",
			envVar:          "",
			expectedSkipSSL: false,
			description:     "Production domain should verify SSL",
		},
		{
			name:            "ProductionSubdomain",
			endpoint:        "https://api.agb.cloud",
			envVar:          "",
			expectedSkipSSL: false,
			description:     "Production subdomain should verify SSL",
		},

		// IP addresses - should skip SSL verification
		{
			name:            "IPv4Address",
			endpoint:        "https://12.34.56.78",
			envVar:          "",
			expectedSkipSSL: true,
			description:     "IPv4 address should skip SSL verification",
		},
		{
			name:            "IPv6Address",
			endpoint:        "https://[2001:db8::1]",
			envVar:          "",
			expectedSkipSSL: true,
			description:     "IPv6 address should skip SSL verification",
		},

		// Local development hosts - should skip SSL verification
		{
			name:            "Localhost",
			endpoint:        "https://localhost",
			envVar:          "",
			expectedSkipSSL: true,
			description:     "Localhost should skip SSL verification",
		},
		{
			name:            "LocalhostWithPort",
			endpoint:        "https://localhost:8080",
			envVar:          "",
			expectedSkipSSL: true,
			description:     "Localhost with port should skip SSL verification",
		},
		{
			name:            "LocalDomain",
			endpoint:        "https://api.local",
			envVar:          "",
			expectedSkipSSL: true,
			description:     ".local domain should skip SSL verification",
		},
		{
			name:            "DevDomain",
			endpoint:        "https://api.dev",
			envVar:          "",
			expectedSkipSSL: true,
			description:     ".dev domain should skip SSL verification",
		},
		{
			name:            "TestDomain",
			endpoint:        "https://api.test",
			envVar:          "",
			expectedSkipSSL: true,
			description:     ".test domain should skip SSL verification",
		},
		{
			name:            "InternalDomain",
			endpoint:        "https://api.internal",
			envVar:          "",
			expectedSkipSSL: true,
			description:     ".internal domain should skip SSL verification",
		},

		// Non-standard ports - should skip SSL verification
		{
			name:            "NonStandardPort",
			endpoint:        "https://agb.cloud:8443",
			envVar:          "",
			expectedSkipSSL: true,
			description:     "Non-standard port should skip SSL verification",
		},

		// Explicit environment variable override
		{
			name:            "ExplicitSkipTrue",
			endpoint:        "https://agb.cloud",
			envVar:          "true",
			expectedSkipSSL: true,
			description:     "Explicit AGB_CLI_SKIP_SSL_VERIFY=true should override",
		},
		{
			name:            "ExplicitSkipFalse",
			endpoint:        "https://12.34.56.78",
			envVar:          "false",
			expectedSkipSSL: false,
			description:     "Explicit AGB_CLI_SKIP_SSL_VERIFY=false should override",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable
			if tc.envVar != "" {
				os.Setenv("AGB_CLI_SKIP_SSL_VERIFY", tc.envVar)
			} else {
				os.Unsetenv("AGB_CLI_SKIP_SSL_VERIFY")
			}

			// Create client configuration
			cfg := &config.Config{
				Endpoint: tc.endpoint,
			}

			// Create API client
			apiClient := client.NewFromConfig(cfg)

			// Check if SSL verification is skipped by examining the transport
			transport, ok := apiClient.GetConfig().HTTPClient.Transport.(*http.Transport)
			var actualSkipSSL bool
			if ok && transport.TLSClientConfig != nil {
				actualSkipSSL = transport.TLSClientConfig.InsecureSkipVerify
			} else {
				actualSkipSSL = false
			}

			if actualSkipSSL != tc.expectedSkipSSL {
				t.Errorf("%s: expected skipSSL=%v, got %v",
					tc.description, tc.expectedSkipSSL, actualSkipSSL)
			} else {
				t.Logf("✅ %s: skipSSL=%v", tc.description, actualSkipSSL)
			}
		})
	}
}

// TestSSLVerificationEdgeCases tests edge cases for SSL verification
func TestSSLVerificationEdgeCases(t *testing.T) {
	// Save original environment variable
	originalSkipSSL := os.Getenv("AGB_CLI_SKIP_SSL_VERIFY")
	defer func() {
		os.Setenv("AGB_CLI_SKIP_SSL_VERIFY", originalSkipSSL)
	}()

	t.Run("InvalidURL", func(t *testing.T) {
		os.Unsetenv("AGB_CLI_SKIP_SSL_VERIFY")

		cfg := &config.Config{
			Endpoint: "not-a-valid-url",
		}

		apiClient := client.NewFromConfig(cfg)

		// Should default to SSL verification for invalid URLs
		transport, ok := apiClient.GetConfig().HTTPClient.Transport.(*http.Transport)
		var actualSkipSSL bool
		if ok && transport.TLSClientConfig != nil {
			actualSkipSSL = transport.TLSClientConfig.InsecureSkipVerify
		}

		if actualSkipSSL {
			t.Error("Invalid URL should default to SSL verification")
		} else {
			t.Log("✅ Invalid URL defaults to SSL verification")
		}
	})

	t.Run("EmptyEndpoint", func(t *testing.T) {
		os.Unsetenv("AGB_CLI_SKIP_SSL_VERIFY")

		cfg := &config.Config{
			Endpoint: "",
		}

		apiClient := client.NewFromConfig(cfg)

		// Should default to SSL verification for empty endpoint
		transport, ok := apiClient.GetConfig().HTTPClient.Transport.(*http.Transport)
		var actualSkipSSL bool
		if ok && transport.TLSClientConfig != nil {
			actualSkipSSL = transport.TLSClientConfig.InsecureSkipVerify
		}

		if actualSkipSSL {
			t.Error("Empty endpoint should default to SSL verification")
		} else {
			t.Log("✅ Empty endpoint defaults to SSL verification")
		}
	})
}
