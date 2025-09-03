// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:     "config",
	Short:   "Manage configuration",
	Long:    "Manage AgbCloud CLI configuration settings",
	GroupID: "management",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]
		fmt.Printf("Setting %s = %s\n", key, value)
		// TODO: Implement configuration storage
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		// TODO: Implement configuration retrieval
		if key == "api_key" {
			if apiKey := os.Getenv("AGB_API_KEY"); apiKey != "" {
				fmt.Printf("%s = %s\n", key, apiKey)
			} else {
				fmt.Printf("%s is not set\n", key)
			}
		} else {
			fmt.Printf("%s is not configured\n", key)
		}
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Configuration:")
		if apiKey := os.Getenv("AGB_API_KEY"); apiKey != "" {
			fmt.Printf("  api_key = %s\n", apiKey)
		} else {
			fmt.Println("  api_key = <not set>")
		}
		fmt.Println("  endpoint = https://agb.cloud")
		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configGetCmd)
	ConfigCmd.AddCommand(configListCmd)
}
