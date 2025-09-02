// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"
)

// Config represents the CLI configuration
type Config struct {
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		APIKey:   os.Getenv("AGB_API_KEY"),
		Endpoint: "https://sdk-api.agb.cloud",
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
