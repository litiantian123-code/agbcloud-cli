// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"testing"

	"github.com/agbcloud/agbcloud-cli/cmd"
)

func TestFormatCPU(t *testing.T) {
	tests := []struct {
		name     string
		cpu      *int
		expected string
	}{
		{
			name:     "nil CPU",
			cpu:      nil,
			expected: "-",
		},
		{
			name:     "zero CPU",
			cpu:      intPtr(0),
			expected: "0",
		},
		{
			name:     "normal CPU",
			cpu:      intPtr(4),
			expected: "4",
		},
		{
			name:     "high CPU",
			cpu:      intPtr(16),
			expected: "16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmd.FormatCPU(tt.cpu)
			if result != tt.expected {
				t.Errorf("FormatCPU(%v) = %q, want %q", tt.cpu, result, tt.expected)
			}
		})
	}
}

func TestFormatMemory(t *testing.T) {
	tests := []struct {
		name     string
		memory   *int
		expected string
	}{
		{
			name:     "nil Memory",
			memory:   nil,
			expected: "-",
		},
		{
			name:     "zero Memory",
			memory:   intPtr(0),
			expected: "0G",
		},
		{
			name:     "normal Memory",
			memory:   intPtr(8),
			expected: "8G",
		},
		{
			name:     "high Memory",
			memory:   intPtr(32),
			expected: "32G",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmd.FormatMemory(tt.memory)
			if result != tt.expected {
				t.Errorf("FormatMemory(%v) = %q, want %q", tt.memory, result, tt.expected)
			}
		})
	}
}

func TestFormatResources(t *testing.T) {
	tests := []struct {
		name     string
		cpu      *int
		memory   *int
		expected string
	}{
		{
			name:     "both nil",
			cpu:      nil,
			memory:   nil,
			expected: "-",
		},
		{
			name:     "CPU nil, Memory present",
			cpu:      nil,
			memory:   intPtr(8),
			expected: "-/8G",
		},
		{
			name:     "CPU present, Memory nil",
			cpu:      intPtr(4),
			memory:   nil,
			expected: "4/-",
		},
		{
			name:     "both present",
			cpu:      intPtr(4),
			memory:   intPtr(8),
			expected: "4/8G",
		},
		{
			name:     "both zero",
			cpu:      intPtr(0),
			memory:   intPtr(0),
			expected: "0/0G",
		},
		{
			name:     "high values",
			cpu:      intPtr(16),
			memory:   intPtr(32),
			expected: "16/32G",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmd.FormatResources(tt.cpu, tt.memory)
			if result != tt.expected {
				t.Errorf("FormatResources(%v, %v) = %q, want %q", tt.cpu, tt.memory, result, tt.expected)
			}
		})
	}
}

func TestResourcesDisplayWidth(t *testing.T) {
	// Test that all possible resource formats fit within the allocated column width (12 characters)
	testCases := []struct {
		cpu    *int
		memory *int
	}{
		{nil, nil},
		{nil, intPtr(999)},
		{intPtr(999), nil},
		{intPtr(999), intPtr(999)},
		{intPtr(0), intPtr(0)},
	}

	maxWidth := 12
	for _, tc := range testCases {
		result := cmd.FormatResources(tc.cpu, tc.memory)
		if len(result) > maxWidth {
			t.Errorf("FormatResources(%v, %v) = %q (length %d) exceeds max width %d",
				tc.cpu, tc.memory, result, len(result), maxWidth)
		}
	}
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}
