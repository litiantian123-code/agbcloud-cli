// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/agbcloud/agbcloud-cli/cmd"
)

func TestErrorMessageFormatting(t *testing.T) {
	// Get expected newline for current platform
	expectedNewline := "\n"
	if runtime.GOOS == "windows" {
		expectedNewline = "\r\n"
	}

	tests := []struct {
		name     string
		testFunc func() error
		wantText []string // Text that should appear in the error message
	}{
		{
			name: "ValidateCPUMemoryCombo - both required",
			testFunc: func() error {
				return cmd.ValidateCPUMemoryCombo(2, 0) // CPU specified but memory not
			},
			wantText: []string{
				"[ERROR] Both CPU and memory must be specified together",
				"[TOOL] Supported combinations:",
				"• 2c4g: --cpu 2 --memory 4",
				"• 4c8g: --cpu 4 --memory 8",
				"• 8c16g: --cpu 8 --memory 16",
			},
		},
		{
			name: "ValidateCPUMemoryCombo - invalid combination",
			testFunc: func() error {
				return cmd.ValidateCPUMemoryCombo(3, 6) // Invalid combination
			},
			wantText: []string{
				"[ERROR] Invalid CPU/Memory combination: 3c6g",
				"[TOOL] Supported combinations:",
				"• 2c4g: --cpu 2 --memory 4",
				"• 4c8g: --cpu 4 --memory 8",
				"• 8c16g: --cpu 8 --memory 16",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			if err == nil {
				t.Fatal("Expected error but got nil")
			}

			errMsg := err.Error()
			t.Logf("Error message:\n%s", errMsg)

			// Check that all expected text appears in the error message
			for _, want := range tt.wantText {
				if !strings.Contains(errMsg, want) {
					t.Errorf("Error message missing expected text: %q", want)
				}
			}

			// Check that the error message contains actual newlines, not literal \n
			if strings.Contains(errMsg, "\\n") {
				t.Error("Error message contains literal \\n instead of actual newlines")
			}

			// Check that the error message uses the correct platform-specific newlines
			if !strings.Contains(errMsg, expectedNewline) {
				t.Errorf("Error message does not contain expected newline sequence for platform %s", runtime.GOOS)
			}

			// Count the number of lines in the error message
			lines := strings.Split(errMsg, expectedNewline)
			if len(lines) < 3 {
				t.Errorf("Expected multi-line error message, got %d lines", len(lines))
			}
		})
	}
}

func TestValidateCPUMemoryCombo_ValidCombinations(t *testing.T) {
	validCombos := []struct {
		cpu    int
		memory int
	}{
		{0, 0},  // Default (both zero)
		{2, 4},  // 2c4g
		{4, 8},  // 4c8g
		{8, 16}, // 8c16g
	}

	for _, combo := range validCombos {
		t.Run(fmt.Sprintf("%dc%dg", combo.cpu, combo.memory), func(t *testing.T) {
			err := cmd.ValidateCPUMemoryCombo(combo.cpu, combo.memory)
			if err != nil {
				t.Errorf("Expected valid combination %dc%dg to pass, got error: %v", combo.cpu, combo.memory, err)
			}
		})
	}
}

func TestPlatformSpecificNewlines(t *testing.T) {
	// This test verifies that our newline handling works correctly on different platforms
	err := cmd.ValidateCPUMemoryCombo(2, 0)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	errMsg := err.Error()

	// On Windows, we should see \r\n
	// On other platforms, we should see \n
	if runtime.GOOS == "windows" {
		if !strings.Contains(errMsg, "\r\n") {
			t.Error("On Windows, error message should contain \\r\\n sequences")
		}
		// Should not contain standalone \n that aren't part of \r\n
		withoutCRLF := strings.ReplaceAll(errMsg, "\r\n", "")
		if strings.Contains(withoutCRLF, "\n") {
			t.Error("On Windows, error message should not contain standalone \\n characters")
		}
	} else {
		if !strings.Contains(errMsg, "\n") {
			t.Error("On non-Windows platforms, error message should contain \\n sequences")
		}
		if strings.Contains(errMsg, "\r\n") {
			t.Error("On non-Windows platforms, error message should not contain \\r\\n sequences")
		}
	}
}
