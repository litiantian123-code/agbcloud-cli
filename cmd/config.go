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
		case "callback_port":
			fmt.Printf("Note: Callback port configuration is read from AGB_CLI_CALLBACK_PORT environment variable\n")
			fmt.Printf("To set callback port, use: export AGB_CLI_CALLBACK_PORT=%s\n", value)
			fmt.Printf("Note: Default port is 3000 if not specified\n")
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
		case "callback_port":
			if cfg.CallbackPort != "" {
				fmt.Println(cfg.CallbackPort)
			} else {
				fmt.Println("3000")
			}
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
		if cfg.CallbackPort != "" {
			fmt.Printf("  callback_port = %s\n", cfg.CallbackPort)
		} else {
			fmt.Println("  callback_port = 3000")
		}

		fmt.Println("\nEnvironment Variables:")
		fmt.Println("  Endpoint:")
		fmt.Println("    AGB_CLI_ENDPOINT (domain only, https:// added automatically)")
		fmt.Println("    Default: agb.cloud")
		fmt.Println("  Callback Port:")
		fmt.Println("    AGB_CLI_CALLBACK_PORT")
		fmt.Println("    Default: 3000")

		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configGetCmd)
	ConfigCmd.AddCommand(configListCmd)
}
