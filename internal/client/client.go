// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/liyuebing/agbcloud-cli/internal/config"
)

// Client represents the API client
type Client struct {
	httpClient *http.Client
	config     *config.Config
}

// New creates a new API client
func New(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: cfg,
	}
}

// Health checks the API health
func (c *Client) Health() error {
	resp, err := c.httpClient.Get(c.config.Endpoint + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API health check failed with status: %d", resp.StatusCode)
	}

	return nil
}
