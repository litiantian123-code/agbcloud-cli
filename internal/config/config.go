// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Config represents the CLI configuration
type Config struct {
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	// Get API key from AGB_CLI_API_KEY environment variable only
	apiKey := os.Getenv("AGB_CLI_API_KEY")

	// Get endpoint from environment variable or use default
	endpoint := os.Getenv("AGB_CLI_ENDPOINT")
	if endpoint == "" {
		endpoint = "agb.cloud"
	}

	// Ensure endpoint has https:// prefix
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	return &Config{
		APIKey:   apiKey,
		Endpoint: endpoint,
	}
}

// ConfigDir returns the configuration directory path
func ConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".agbcloud"), nil
}

// ConfigFile returns the configuration file path
func ConfigFile() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}
