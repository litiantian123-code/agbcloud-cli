// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config represents the CLI configuration
// Only stores authentication tokens, endpoint is determined from environment variables
type Config struct {
	Token *Token `json:"token,omitempty"` // OAuth token authentication
}

// Token represents AgbCloud authentication tokens
type Token struct {
	LoginToken     string    `json:"loginToken"`
	SessionId      string    `json:"sessionId"`
	KeepAliveToken string    `json:"keepAliveToken"`
	ExpiresAt      time.Time `json:"expiresAt"`
}

var (
	ErrNoTokenFound = errors.New("no authentication token found. Run 'agbcloud-cli login' to authenticate")
)

// GetConfig loads the configuration from file or creates a new one
func GetConfig() (*Config, error) {
	configFilePath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	var c Config

	// Try to load existing config file
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		// No config file exists, create new config
		c = Config{}
	} else if err != nil {
		return nil, err
	} else {
		// Config file exists, load it
		configContent, err := os.ReadFile(configFilePath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(configContent, &c)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// GetEndpoint returns the endpoint from environment variable or default
func GetEndpoint() string {
	endpoint := os.Getenv("AGB_CLI_ENDPOINT")
	if endpoint == "" {
		endpoint = "agb.cloud"
	}

	// Ensure endpoint has https:// prefix
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	return endpoint
}

// Save writes the configuration to file
func (c *Config) Save() error {
	configFilePath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(configFilePath), 0755)
	if err != nil {
		return err
	}

	configContent, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath, configContent, 0600) // More secure permissions for auth data
}

// GetTokens retrieves authentication tokens
func (c *Config) GetTokens() (*Token, error) {
	if c.Token == nil {
		return nil, ErrNoTokenFound
	}
	return c.Token, nil
}

// SaveTokens saves authentication tokens to the configuration
func (c *Config) SaveTokens(loginToken, sessionId, keepAliveToken, expiresAt string) error {
	// Parse expiresAt time
	var expiresAtTime time.Time
	if expiresAt != "" {
		var err error
		expiresAtTime, err = time.Parse(time.RFC3339, expiresAt)
		if err != nil {
			return fmt.Errorf("failed to parse expiresAt time: %w", err)
		}
	}

	// Update config with tokens
	c.Token = &Token{
		LoginToken:     loginToken,
		SessionId:      sessionId,
		KeepAliveToken: keepAliveToken,
		ExpiresAt:      expiresAtTime,
	}

	return c.Save()
}

// ClearTokens removes authentication tokens from the configuration
func (c *Config) ClearTokens() error {
	c.Token = nil
	return c.Save()
}

// IsAuthenticated checks if the user is authenticated (has tokens)
func (c *Config) IsAuthenticated() bool {
	return c.Token != nil
}

// getConfigPath returns the path to the configuration file
func getConfigPath() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{}
}

// ConfigDir returns the configuration directory path
func ConfigDir() (string, error) {
	// Check for environment variable override first
	agbConfigDir := os.Getenv("AGB_CLI_CONFIG_DIR")
	if agbConfigDir != "" {
		return agbConfigDir, nil
	}

	// Use OS-specific standard config directory
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userConfigDir, "agbcloud"), nil
}

// ConfigFile returns the configuration file path
func ConfigFile() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}
