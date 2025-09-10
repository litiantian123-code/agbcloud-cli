// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agbcloud/agbcloud-cli/cmd"
)

func TestImageCommand(t *testing.T) {
	// Test that image command exists and has correct structure
	assert.Equal(t, "image", cmd.ImageCmd.Use)
	assert.Equal(t, "Manage images", cmd.ImageCmd.Short)
	assert.Equal(t, "management", cmd.ImageCmd.GroupID)

	// Test that subcommands exist
	subcommands := cmd.ImageCmd.Commands()
	assert.Len(t, subcommands, 2)

	var createCmd, activateCmd *cobra.Command
	for _, subcmd := range subcommands {
		switch subcmd.Use {
		case "create <image-name>":
			createCmd = subcmd
		case "activate <image-id>":
			activateCmd = subcmd
		}
	}

	require.NotNil(t, createCmd, "create subcommand should exist")
	require.NotNil(t, activateCmd, "activate subcommand should exist")
}

func TestImageCreateCommand(t *testing.T) {
	// Get create subcommand
	var createCmd *cobra.Command
	for _, subcmd := range cmd.ImageCmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "create") {
			createCmd = subcmd
			break
		}
	}
	require.NotNil(t, createCmd)

	// Test command structure
	assert.Equal(t, "create <image-name>", createCmd.Use)
	assert.Equal(t, "Create a custom image", createCmd.Short)

	// Test required flags
	dockerfileFlag := createCmd.Flag("dockerfile")
	require.NotNil(t, dockerfileFlag, "dockerfile flag should exist")

	imageIdFlag := createCmd.Flag("imageId")
	require.NotNil(t, imageIdFlag, "imageId flag should exist")

	// Test that flags are marked as required by checking if they exist
	// (The actual required validation happens during execution)
	assert.NotNil(t, dockerfileFlag)
	assert.NotNil(t, imageIdFlag)
}

func TestImageActivateCommand(t *testing.T) {
	// Get activate subcommand
	var activateCmd *cobra.Command
	for _, subcmd := range cmd.ImageCmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "activate") {
			activateCmd = subcmd
			break
		}
	}
	require.NotNil(t, activateCmd)

	// Test command structure
	assert.Equal(t, "activate <image-id>", activateCmd.Use)
	assert.Equal(t, "Activate an image", activateCmd.Short)

	// Test optional flags
	cpuFlag := activateCmd.Flag("cpu")
	require.NotNil(t, cpuFlag, "cpu flag should exist")

	memoryFlag := activateCmd.Flag("memory")
	require.NotNil(t, memoryFlag, "memory flag should exist")
}

func TestImageCreateCommandArgumentValidation(t *testing.T) {
	// Get the create subcommand specifically
	var createCmd *cobra.Command
	for _, subcmd := range cmd.ImageCmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "create") {
			createCmd = subcmd
			break
		}
	}
	require.NotNil(t, createCmd, "create command should exist")

	// Test missing argument
	err := createCmd.Args(createCmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Missing required argument: <image-name>")
	assert.Contains(t, err.Error(), "Short form:")

	// Test too many arguments
	err = createCmd.Args(createCmd, []string{"image1", "image2"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Too many arguments provided")

	// Test valid argument count
	err = createCmd.Args(createCmd, []string{"testImage"})
	assert.NoError(t, err)
}

func TestImageActivateCommandFlags(t *testing.T) {
	// Get activate subcommand
	var activateCmd *cobra.Command
	for _, subcmd := range cmd.ImageCmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "activate") {
			activateCmd = subcmd
			break
		}
	}
	require.NotNil(t, activateCmd)

	// Test that flags exist and have correct types
	cpuFlag := activateCmd.Flag("cpu")
	require.NotNil(t, cpuFlag)
	assert.Equal(t, "int", cpuFlag.Value.Type())
	assert.Equal(t, "c", cpuFlag.Shorthand) // Test short form

	memoryFlag := activateCmd.Flag("memory")
	require.NotNil(t, memoryFlag)
	assert.Equal(t, "int", memoryFlag.Value.Type())
	assert.Equal(t, "m", memoryFlag.Shorthand) // Test short form
}

func TestImageCreateCommandFlags(t *testing.T) {
	// Get create subcommand
	var createCmd *cobra.Command
	for _, subcmd := range cmd.ImageCmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "create") {
			createCmd = subcmd
			break
		}
	}
	require.NotNil(t, createCmd)

	// Test that flags exist and have correct types
	dockerfileFlag := createCmd.Flag("dockerfile")
	require.NotNil(t, dockerfileFlag)
	assert.Equal(t, "string", dockerfileFlag.Value.Type())
	assert.Equal(t, "f", dockerfileFlag.Shorthand) // Test short form

	imageIdFlag := createCmd.Flag("imageId")
	require.NotNil(t, imageIdFlag)
	assert.Equal(t, "string", imageIdFlag.Value.Type())
	assert.Equal(t, "i", imageIdFlag.Shorthand) // Test short form
}
