// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agbcloud/agbcloud-cli/internal/config"
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

		switch key {
		case "endpoint":
			fmt.Printf("Note: Endpoint configuration is read from AGB_CLI_ENDPOINT environment variable\n")
			fmt.Printf("To set endpoint, use: export AGB_CLI_ENDPOINT=%s\n", value)
			fmt.Printf("Note: You can specify just the domain (e.g., 'agb.cloud'), https:// will be added automatically\n")
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		cfg := config.DefaultConfig()

		switch key {
		case "endpoint":
			fmt.Println(cfg.Endpoint)
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.DefaultConfig()

		fmt.Println("Configuration:")
		fmt.Printf("  endpoint = %s\n", cfg.Endpoint)
		fmt.Println("  callback_port = 3000 (automatic port selection enabled)")

		fmt.Println("\nEnvironment Variables:")
		fmt.Println("  Endpoint:")
		fmt.Println("    AGB_CLI_ENDPOINT (domain only, https:// added automatically)")
		fmt.Println("    Default: agb.cloud")
		fmt.Println("  Callback Port:")
		fmt.Println("    Fixed at 3000 with automatic alternative port selection")

		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configGetCmd)
	ConfigCmd.AddCommand(configListCmd)
}
