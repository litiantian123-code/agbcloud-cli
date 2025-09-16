// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"strings"
	"testing"

	"github.com/agbcloud/agbcloud-cli/cmd"
)

func TestValidateCPUMemoryCombo(t *testing.T) {
	tests := []struct {
		name        string
		cpu         int
		memory      int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Default configuration (both zero)",
			cpu:         0,
			memory:      0,
			expectError: false,
		},
		{
			name:        "Valid 2c4g combination",
			cpu:         2,
			memory:      4,
			expectError: false,
		},
		{
			name:        "Valid 4c8g combination",
			cpu:         4,
			memory:      8,
			expectError: false,
		},
		{
			name:        "Valid 8c16g combination",
			cpu:         8,
			memory:      16,
			expectError: false,
		},
		{
			name:        "Invalid combination - only CPU specified",
			cpu:         2,
			memory:      0,
			expectError: true,
			errorMsg:    "Both CPU and memory must be specified together",
		},
		{
			name:        "Invalid combination - only memory specified",
			cpu:         0,
			memory:      4,
			expectError: true,
			errorMsg:    "Both CPU and memory must be specified together",
		},
		{
			name:        "Invalid combination - 3c6g",
			cpu:         3,
			memory:      6,
			expectError: true,
			errorMsg:    "Invalid CPU/Memory combination: 3c6g",
		},
		{
			name:        "Invalid combination - 2c8g",
			cpu:         2,
			memory:      8,
			expectError: true,
			errorMsg:    "Invalid CPU/Memory combination: 2c8g",
		},
		{
			name:        "Invalid combination - 4c4g",
			cpu:         4,
			memory:      4,
			expectError: true,
			errorMsg:    "Invalid CPU/Memory combination: 4c4g",
		},
		{
			name:        "Invalid combination - 1c2g",
			cpu:         1,
			memory:      2,
			expectError: true,
			errorMsg:    "Invalid CPU/Memory combination: 1c2g",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.ValidateCPUMemoryCombo(tt.cpu, tt.memory)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for CPU=%d, Memory=%d, but got none", tt.cpu, tt.memory)
					return
				}

				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', but got: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for CPU=%d, Memory=%d, but got: %v", tt.cpu, tt.memory, err)
				}
			}
		})
	}
}
