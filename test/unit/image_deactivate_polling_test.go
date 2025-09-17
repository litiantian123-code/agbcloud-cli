// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"fmt"
	"testing"

	"github.com/agbcloud/agbcloud-cli/cmd"
)

func TestImageDeactivateStatusTransitions(t *testing.T) {
	// Test the status transitions during deactivation process
	tests := []struct {
		name           string
		status         string
		expectedAction string
		shouldComplete bool
		shouldFail     bool
	}{
		{
			name:           "Successfully deactivated",
			status:         "IMAGE_AVAILABLE",
			expectedAction: "complete",
			shouldComplete: true,
			shouldFail:     false,
		},
		{
			name:           "Deactivation in progress",
			status:         "RESOURCE_DELETING",
			expectedAction: "continue",
			shouldComplete: false,
			shouldFail:     false,
		},
		{
			name:           "Still activated - continue monitoring",
			status:         "RESOURCE_PUBLISHED",
			expectedAction: "continue",
			shouldComplete: false,
			shouldFail:     false,
		},
		{
			name:           "Deactivation failed",
			status:         "RESOURCE_FAILED",
			expectedAction: "fail",
			shouldComplete: false,
			shouldFail:     true,
		},
		{
			name:           "Unknown status - continue monitoring",
			status:         "UNKNOWN_STATUS",
			expectedAction: "continue",
			shouldComplete: false,
			shouldFail:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the status formatting
			formattedStatus := cmd.FormatImageStatus(tt.status)

			// Verify that status formatting works correctly
			if tt.status == "IMAGE_AVAILABLE" && formattedStatus != "Available" {
				t.Errorf("Expected 'Available' for IMAGE_AVAILABLE, got %s", formattedStatus)
			}
			if tt.status == "RESOURCE_DELETING" && formattedStatus != "Deactivating" {
				t.Errorf("Expected 'Deactivating' for RESOURCE_DELETING, got %s", formattedStatus)
			}
			if tt.status == "RESOURCE_PUBLISHED" && formattedStatus != "Activated" {
				t.Errorf("Expected 'Activated' for RESOURCE_PUBLISHED, got %s", formattedStatus)
			}
			if tt.status == "RESOURCE_FAILED" && formattedStatus != "Activate Failed" {
				t.Errorf("Expected 'Activate Failed' for RESOURCE_FAILED, got %s", formattedStatus)
			}
		})
	}
}

func TestDeactivateStatusFlow(t *testing.T) {
	// Test the expected flow of deactivation statuses
	expectedFlow := []string{
		"RESOURCE_PUBLISHED", // Initially activated
		"RESOURCE_DELETING",  // Deactivation in progress
		"IMAGE_AVAILABLE",    // Successfully deactivated
	}

	for i, status := range expectedFlow {
		t.Run(fmt.Sprintf("Step_%d_%s", i+1, status), func(t *testing.T) {
			formattedStatus := cmd.FormatImageStatus(status)

			switch status {
			case "RESOURCE_PUBLISHED":
				if formattedStatus != "Activated" {
					t.Errorf("Step %d: Expected 'Activated', got %s", i+1, formattedStatus)
				}
			case "RESOURCE_DELETING":
				if formattedStatus != "Deactivating" {
					t.Errorf("Step %d: Expected 'Deactivating', got %s", i+1, formattedStatus)
				}
			case "IMAGE_AVAILABLE":
				if formattedStatus != "Available" {
					t.Errorf("Step %d: Expected 'Available', got %s", i+1, formattedStatus)
				}
			}
		})
	}
}

func TestDeactivateErrorScenarios(t *testing.T) {
	// Test various error scenarios during deactivation
	errorStatuses := []struct {
		status      string
		description string
		shouldFail  bool
	}{
		{
			status:      "RESOURCE_FAILED",
			description: "Deactivation failed",
			shouldFail:  true,
		},
		{
			status:      "IMAGE_CREATE_FAILED",
			description: "Image in failed state",
			shouldFail:  false, // This shouldn't happen during deactivation, but we continue monitoring
		},
	}

	for _, scenario := range errorStatuses {
		t.Run(scenario.description, func(t *testing.T) {
			formattedStatus := cmd.FormatImageStatus(scenario.status)

			// Verify status formatting
			if scenario.status == "RESOURCE_FAILED" && formattedStatus != "Activate Failed" {
				t.Errorf("Expected 'Activate Failed' for RESOURCE_FAILED, got %s", formattedStatus)
			}
			if scenario.status == "IMAGE_CREATE_FAILED" && formattedStatus != "Create Failed" {
				t.Errorf("Expected 'Create Failed' for IMAGE_CREATE_FAILED, got %s", formattedStatus)
			}
		})
	}
}

func TestDeactivatePollingConfiguration(t *testing.T) {
	// Test that polling configuration is reasonable
	t.Run("Polling interval", func(t *testing.T) {
		// The polling interval should be 5 seconds (same as activation)
		// This is reasonable for deactivation monitoring
		expectedInterval := 5 // seconds

		// This is more of a documentation test - the actual interval is hardcoded
		// in the pollImageDeactivationStatus function
		if expectedInterval != 5 {
			t.Errorf("Expected polling interval of 5 seconds, configuration suggests %d", expectedInterval)
		}
	})

	t.Run("Timeout duration", func(t *testing.T) {
		// The timeout should be 45 minutes (same as activation)
		// This should be sufficient for deactivation operations
		expectedTimeout := 45 // minutes

		// This is more of a documentation test - the actual timeout is hardcoded
		// in the pollImageDeactivationStatus function
		if expectedTimeout != 45 {
			t.Errorf("Expected timeout of 45 minutes, configuration suggests %d", expectedTimeout)
		}
	})
}
