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

// Config represents the CLI configuration with profile support
type Config struct {
	ActiveProfileId string    `json:"activeProfile"`
	Profiles        []Profile `json:"profiles"`
	// Legacy fields for backward compatibility
	APIKey       string `json:"api_key,omitempty"`
	Endpoint     string `json:"endpoint,omitempty"`
	CallbackPort string `json:"callback_port,omitempty"`
}

// Profile represents a user profile with authentication information
type Profile struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Api      ServerApi `json:"api"`
	Endpoint string    `json:"endpoint"`
}

// ServerApi holds API configuration and authentication information
type ServerApi struct {
	Url   string `json:"url"`
	Key   string `json:"key,omitempty"`   // API key authentication
	Token *Token `json:"token,omitempty"` // OAuth token authentication
}

// Token represents AgbCloud authentication tokens
type Token struct {
	LoginToken     string    `json:"loginToken"`
	SessionId      string    `json:"sessionId"`
	KeepAliveToken string    `json:"keepAliveToken"`
	CreatedAt      time.Time `json:"createdAt"`
	// Note: AgbCloud tokens don't have explicit expiry, but we track creation time
}

var (
	ErrNoProfilesFound = errors.New("no profiles found. Run 'agbcloud-cli login' to authenticate")
	ErrNoActiveProfile = errors.New("no active profile found. Run 'agbcloud-cli login' to authenticate")
)

// GetConfig loads the configuration from file or creates a new one
func GetConfig() (*Config, error) {
	configFilePath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		config := &Config{}
		return config, config.Save()
	}

	if err != nil {
		return nil, err
	}

	var c Config
	configContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(configContent, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetActiveProfile returns the currently active profile
func (c *Config) GetActiveProfile() (Profile, error) {
	if len(c.Profiles) == 0 {
		return Profile{}, ErrNoProfilesFound
	}

	for _, profile := range c.Profiles {
		if profile.Id == c.ActiveProfileId {
			return profile, nil
		}
	}

	return Profile{}, ErrNoActiveProfile
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

// AddProfile adds a new profile and sets it as active
func (c *Config) AddProfile(profile Profile) error {
	c.Profiles = append(c.Profiles, profile)
	c.ActiveProfileId = profile.Id

	return c.Save()
}

// EditProfile updates an existing profile
func (c *Config) EditProfile(profile Profile) error {
	for i, p := range c.Profiles {
		if p.Id == profile.Id {
			c.Profiles[i] = profile
			return c.Save()
		}
	}

	return fmt.Errorf("profile with id %s not found", profile.Id)
}

// RemoveProfile removes a profile by ID
func (c *Config) RemoveProfile(profileId string) error {
	if c.ActiveProfileId == profileId {
		return errors.New("cannot remove active profile")
	}

	var profiles []Profile
	for _, profile := range c.Profiles {
		if profile.Id != profileId {
			profiles = append(profiles, profile)
		}
	}

	c.Profiles = profiles
	return c.Save()
}

// GetProfile returns a profile by ID
func (c *Config) GetProfile(profileId string) (Profile, error) {
	for _, profile := range c.Profiles {
		if profile.Id == profileId {
			return profile, nil
		}
	}

	return Profile{}, errors.New("profile not found")
}

// SaveTokens saves authentication tokens to the active profile
func (c *Config) SaveTokens(loginToken, sessionId, keepAliveToken string) error {
	activeProfile, err := c.GetActiveProfile()
	if err != nil {
		if err == ErrNoProfilesFound {
			// Create initial profile
			activeProfile, err = c.createInitialProfile()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Update profile with tokens
	activeProfile.Api.Token = &Token{
		LoginToken:     loginToken,
		SessionId:      sessionId,
		KeepAliveToken: keepAliveToken,
		CreatedAt:      time.Now(),
	}
	activeProfile.Api.Key = "" // Clear API key when using tokens

	return c.EditProfile(activeProfile)
}

// GetTokens retrieves authentication tokens from the active profile
func (c *Config) GetTokens() (*Token, error) {
	activeProfile, err := c.GetActiveProfile()
	if err != nil {
		return nil, err
	}

	if activeProfile.Api.Token == nil {
		return nil, errors.New("no authentication tokens found. Run 'agbcloud-cli login' to authenticate")
	}

	return activeProfile.Api.Token, nil
}

// IsAuthenticated checks if the user is authenticated (has tokens or API key)
func (c *Config) IsAuthenticated() bool {
	activeProfile, err := c.GetActiveProfile()
	if err != nil {
		return false
	}

	return activeProfile.Api.Token != nil || activeProfile.Api.Key != ""
}

// createInitialProfile creates the first profile for new users
func (c *Config) createInitialProfile() (Profile, error) {
	endpoint := os.Getenv("AGB_CLI_ENDPOINT")
	if endpoint == "" {
		endpoint = "agb.cloud"
	}

	// Ensure endpoint has https:// prefix
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	profile := Profile{
		Id:       "default",
		Name:     "Default Profile",
		Endpoint: endpoint,
		Api: ServerApi{
			Url: endpoint,
		},
	}

	err := c.AddProfile(profile)
	return profile, err
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

	// Get callback port from environment variable
	callbackPort := os.Getenv("AGB_CLI_CALLBACK_PORT")

	return &Config{
		APIKey:       apiKey,
		Endpoint:     endpoint,
		CallbackPort: callbackPort,
	}
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
