// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agbcloud/agbcloud-cli/cmd"
)

// captureStderr temporarily redirects stderr to capture output during tests
func captureStderr(f func()) string {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r) // Ignore errors in test helper
	return buf.String()
}

func TestImageCommand(t *testing.T) {
	// Test that image command exists and has correct structure
	assert.Equal(t, "image", cmd.ImageCmd.Use)
	assert.Equal(t, "Manage images", cmd.ImageCmd.Short)
	assert.Equal(t, "management", cmd.ImageCmd.GroupID)

	// Test that subcommands exist
	subcommands := cmd.ImageCmd.Commands()
	assert.Len(t, subcommands, 4)

	var createCmd, activateCmd, deactivateCmd, listCmd *cobra.Command
	for _, subcmd := range subcommands {
		switch {
		case strings.HasPrefix(subcmd.Use, "create"):
			createCmd = subcmd
		case strings.HasPrefix(subcmd.Use, "activate"):
			activateCmd = subcmd
		case strings.HasPrefix(subcmd.Use, "deactivate"):
			deactivateCmd = subcmd
		case subcmd.Use == "list":
			listCmd = subcmd
		}
	}

	require.NotNil(t, createCmd, "create subcommand should exist")
	require.NotNil(t, activateCmd, "activate subcommand should exist")
	require.NotNil(t, deactivateCmd, "deactivate subcommand should exist")
	require.NotNil(t, listCmd, "list subcommand should exist")
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

	// Test flags
	cpuFlag := activateCmd.Flag("cpu")
	require.NotNil(t, cpuFlag, "cpu flag should exist")

	memoryFlag := activateCmd.Flag("memory")
	require.NotNil(t, memoryFlag, "memory flag should exist")
}

func TestImageListCommand(t *testing.T) {
	// Get list subcommand
	var listCmd *cobra.Command
	for _, subcmd := range cmd.ImageCmd.Commands() {
		if subcmd.Use == "list" {
			listCmd = subcmd
			break
		}
	}
	require.NotNil(t, listCmd, "list subcommand should exist")

	// Test command structure
	assert.Equal(t, "list", listCmd.Use)
	assert.Equal(t, "List images", listCmd.Short)
	expectedLong := `List images with pagination support.

Image types:
  User   - Custom images created by users
  System - System-provided base images`
	assert.Equal(t, expectedLong, listCmd.Long)

	// Test flags
	typeFlag := listCmd.Flag("type")
	require.NotNil(t, typeFlag, "type flag should exist")
	assert.Equal(t, "User", typeFlag.DefValue, "type flag default should be 'User'")

	pageFlag := listCmd.Flag("page")
	require.NotNil(t, pageFlag, "page flag should exist")
	assert.Equal(t, "1", pageFlag.DefValue, "page flag default should be '1'")

	sizeFlag := listCmd.Flag("size")
	require.NotNil(t, sizeFlag, "size flag should exist")
	assert.Equal(t, "10", sizeFlag.DefValue, "size flag default should be '10'")
}

func TestImageCommandStructure(t *testing.T) {
	// Test that all expected subcommands exist
	subcommands := cmd.ImageCmd.Commands()
	assert.Len(t, subcommands, 4, "Should have 4 subcommands: create, activate, deactivate, list")

	commandNames := make([]string, len(subcommands))
	for i, subcmd := range subcommands {
		commandNames[i] = strings.Split(subcmd.Use, " ")[0] // Get the first word (command name)
	}

	assert.Contains(t, commandNames, "create", "Should have create subcommand")
	assert.Contains(t, commandNames, "activate", "Should have activate subcommand")
	assert.Contains(t, commandNames, "deactivate", "Should have deactivate subcommand")
	assert.Contains(t, commandNames, "list", "Should have list subcommand")
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
	var err error
	captureStderr(func() {
		err = createCmd.Args(createCmd, []string{})
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Missing required argument: <image-name>")
	assert.Contains(t, err.Error(), "Short form:")

	// Test too many arguments
	captureStderr(func() {
		err = createCmd.Args(createCmd, []string{"image1", "image2"})
	})
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

func TestImageActivateCommandIntegration(t *testing.T) {
	// Test that the activate command properly integrates with the StartImage API
	// This test verifies the command structure and flag handling
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
	expectedLong := `Activate an image with specified resources.

Supported CPU and Memory combinations:
  2c4g  - 2 CPU cores with 4 GB memory
  4c8g  - 4 CPU cores with 8 GB memory  
  8c16g - 8 CPU cores with 16 GB memory

If no CPU/memory is specified, default resources will be used.`
	assert.Equal(t, expectedLong, activateCmd.Long)

	// Test flags exist and have correct properties
	cpuFlag := activateCmd.Flag("cpu")
	require.NotNil(t, cpuFlag)
	assert.Equal(t, "int", cpuFlag.Value.Type())
	assert.Equal(t, "c", cpuFlag.Shorthand)
	assert.Equal(t, "0", cpuFlag.DefValue)

	memoryFlag := activateCmd.Flag("memory")
	require.NotNil(t, memoryFlag)
	assert.Equal(t, "int", memoryFlag.Value.Type())
	assert.Equal(t, "m", memoryFlag.Shorthand)
	assert.Equal(t, "0", memoryFlag.DefValue)

	// Test argument validation
	var err error
	captureStderr(func() {
		err = activateCmd.Args(activateCmd, []string{})
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Missing required argument: <image-id>")

	captureStderr(func() {
		err = activateCmd.Args(activateCmd, []string{"img-123", "extra-arg"})
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Too many arguments provided")

	err = activateCmd.Args(activateCmd, []string{"img-123"})
	assert.NoError(t, err)
}

func TestImageDeactivateCommand(t *testing.T) {
	// Get deactivate subcommand
	var deactivateCmd *cobra.Command
	for _, subcmd := range cmd.ImageCmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "deactivate") {
			deactivateCmd = subcmd
			break
		}
	}
	require.NotNil(t, deactivateCmd)

	// Test command structure
	assert.Equal(t, "deactivate <image-id>", deactivateCmd.Use)
	assert.Equal(t, "Deactivate an image", deactivateCmd.Short)
	assert.Equal(t, "Deactivate a running image instance", deactivateCmd.Long)
}

func TestImageDeactivateCommandArgumentValidation(t *testing.T) {
	// Get the deactivate subcommand specifically
	var deactivateCmd *cobra.Command
	for _, subcmd := range cmd.ImageCmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "deactivate") {
			deactivateCmd = subcmd
			break
		}
	}
	require.NotNil(t, deactivateCmd, "deactivate command should exist")

	// Test missing argument
	var err error
	captureStderr(func() {
		err = deactivateCmd.Args(deactivateCmd, []string{})
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Missing required argument: <image-id>")
	assert.Contains(t, err.Error(), "Usage: agbcloud image deactivate <image-id>")

	// Test too many arguments
	captureStderr(func() {
		err = deactivateCmd.Args(deactivateCmd, []string{"image1", "image2"})
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Too many arguments provided")

	// Test valid argument count
	err = deactivateCmd.Args(deactivateCmd, []string{"test-image-id"})
	assert.NoError(t, err)
}

func TestImageCommandStructureWithDeactivate(t *testing.T) {
	// Test that all expected subcommands exist including deactivate
	subcommands := cmd.ImageCmd.Commands()
	assert.Len(t, subcommands, 4, "Should have 4 subcommands: create, activate, deactivate, list")

	commandNames := make([]string, len(subcommands))
	for i, subcmd := range subcommands {
		commandNames[i] = strings.Split(subcmd.Use, " ")[0] // Get the first word (command name)
	}

	assert.Contains(t, commandNames, "create", "Should have create subcommand")
	assert.Contains(t, commandNames, "activate", "Should have activate subcommand")
	assert.Contains(t, commandNames, "deactivate", "Should have deactivate subcommand")
	assert.Contains(t, commandNames, "list", "Should have list subcommand")
}
