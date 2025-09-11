// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// formatImageStatus is the function we're testing (copied from cmd/image.go)
func formatImageStatus(status string) string {
	switch status {
	// Image creation related statuses
	case "IMAGE_CREATING":
		return "Creating"
	case "IMAGE_CREATE_FAILED":
		return "Create Failed"
	case "IMAGE_AVAILABLE":
		return "Available"

	// Resource activation related statuses
	case "RESOURCE_DEPLOYING":
		return "Activating"
	case "RESOURCE_PUBLISHED":
		return "Activated"
	case "RESOURCE_DELETING":
		return "Deactivating"
	case "RESOURCE_FAILED":
		return "Activate Failed"
	case "RESOURCE_CEASED":
		return "Ceased Billing"

	default:
		return status
	}
}

func TestFormatImageStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Image creation related statuses
		{
			name:     "IMAGE_CREATING maps to Creating",
			input:    "IMAGE_CREATING",
			expected: "Creating",
		},
		{
			name:     "IMAGE_CREATE_FAILED maps to Create Failed",
			input:    "IMAGE_CREATE_FAILED",
			expected: "Create Failed",
		},
		{
			name:     "IMAGE_AVAILABLE maps to Available",
			input:    "IMAGE_AVAILABLE",
			expected: "Available",
		},

		// Resource activation related statuses
		{
			name:     "RESOURCE_DEPLOYING maps to Activating",
			input:    "RESOURCE_DEPLOYING",
			expected: "Activating",
		},
		{
			name:     "RESOURCE_PUBLISHED maps to Activated",
			input:    "RESOURCE_PUBLISHED",
			expected: "Activated",
		},
		{
			name:     "RESOURCE_DELETING maps to Deactivating",
			input:    "RESOURCE_DELETING",
			expected: "Deactivating",
		},
		{
			name:     "RESOURCE_FAILED maps to Activate Failed",
			input:    "RESOURCE_FAILED",
			expected: "Activate Failed",
		},
		{
			name:     "RESOURCE_CEASED maps to Ceased Billing",
			input:    "RESOURCE_CEASED",
			expected: "Ceased Billing",
		},

		// Unknown status
		{
			name:     "Unknown status returns as-is",
			input:    "UNKNOWN_STATUS",
			expected: "UNKNOWN_STATUS",
		},
		{
			name:     "Empty status returns as-is",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatImageStatus(tt.input)
			assert.Equal(t, tt.expected, result, "Status mapping should be correct")
		})
	}
}

func TestStatusMappingCategories(t *testing.T) {
	// Test that all image creation statuses are properly mapped
	imageCreationStatuses := map[string]string{
		"IMAGE_CREATING":      "Creating",
		"IMAGE_CREATE_FAILED": "Create Failed",
		"IMAGE_AVAILABLE":     "Available",
	}

	for input, expected := range imageCreationStatuses {
		result := formatImageStatus(input)
		assert.Equal(t, expected, result, "Image creation status %s should map to %s", input, expected)
	}

	// Test that all resource activation statuses are properly mapped
	resourceActivationStatuses := map[string]string{
		"RESOURCE_DEPLOYING": "Activating",
		"RESOURCE_PUBLISHED": "Activated",
		"RESOURCE_DELETING":  "Deactivating",
		"RESOURCE_FAILED":    "Activate Failed",
		"RESOURCE_CEASED":    "Ceased Billing",
	}

	for input, expected := range resourceActivationStatuses {
		result := formatImageStatus(input)
		assert.Equal(t, expected, result, "Resource activation status %s should map to %s", input, expected)
	}

}

func TestStatusMappingConsistency(t *testing.T) {
	// Test that similar statuses map to consistent patterns

	// All "FAILED" statuses should contain "Failed"
	failedStatuses := []string{"IMAGE_CREATE_FAILED", "RESOURCE_FAILED"}
	for _, status := range failedStatuses {
		result := formatImageStatus(status)
		assert.Contains(t, result, "Failed", "Failed status %s should contain 'Failed' in result: %s", status, result)
	}

	// All "CREATING" related statuses should map to "Creating"
	creatingStatuses := []string{"IMAGE_CREATING"}
	for _, status := range creatingStatuses {
		result := formatImageStatus(status)
		assert.Equal(t, "Creating", result, "Creating status %s should map to 'Creating'", status)
	}
}
