// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"bytes"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestVerboseFlagFunctionality(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectDebugLog bool
	}{
		{
			name:           "Without verbose flag",
			args:           []string{"test"},
			expectDebugLog: false,
		},
		{
			name:           "With verbose flag -v",
			args:           []string{"test", "-v"},
			expectDebugLog: true,
		},
		{
			name:           "With verbose flag --verbose",
			args:           []string{"test", "--verbose"},
			expectDebugLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)

			// Reset log level to default
			log.SetLevel(log.InfoLevel)

			// Create a test command with verbose flag
			testCmd := &cobra.Command{
				Use: "test",
				PersistentPreRun: func(cmd *cobra.Command, args []string) {
					// Set up logging based on verbose flag
					verbose, _ := cmd.Flags().GetBool("verbose")
					if verbose {
						log.SetLevel(log.DebugLevel)
					} else {
						log.SetLevel(log.InfoLevel)
					}
				},
				Run: func(cmd *cobra.Command, args []string) {
					// Test both info and debug logs
					log.Info("This is an info message")
					log.Debug("This is a debug message")
				},
			}

			// Add verbose flag
			testCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

			// Set command args
			testCmd.SetArgs(tt.args[1:]) // Skip the command name

			// Execute command
			err := testCmd.Execute()
			assert.NoError(t, err)

			// Check output
			output := buf.String()

			// Info message should always be present
			assert.Contains(t, output, "This is an info message")

			// Debug message should only be present when verbose flag is set
			if tt.expectDebugLog {
				assert.Contains(t, output, "This is a debug message", "Debug message should be present with verbose flag")
			} else {
				assert.NotContains(t, output, "This is a debug message", "Debug message should not be present without verbose flag")
			}

			// Reset log output to stderr
			log.SetOutput(os.Stderr)
		})
	}
}

func TestLogLevelConfiguration(t *testing.T) {
	// Test that log level is correctly set based on verbose flag
	testCmd := &cobra.Command{
		Use: "test",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			if verbose {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing, just test the PreRun
		},
	}

	testCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	// Test without verbose flag
	testCmd.SetArgs([]string{})
	err := testCmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, log.InfoLevel, log.GetLevel(), "Log level should be Info without verbose flag")

	// Test with verbose flag
	testCmd.SetArgs([]string{"-v"})
	err = testCmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, log.DebugLevel, log.GetLevel(), "Log level should be Debug with verbose flag")
}
