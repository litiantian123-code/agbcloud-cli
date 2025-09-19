// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"strings"
	"testing"

	"github.com/agbcloud/agbcloud-cli/internal/auth"
)

func TestAuthSuccessPageContainsAutoClose(t *testing.T) {
	// This test verifies that the success HTML contains JavaScript to auto-close the window
	html := auth.GetSuccessHTML()

	// Check that the HTML contains setTimeout function
	if !strings.Contains(html, "setTimeout") {
		t.Error("Success HTML should contain setTimeout function for auto-close")
	}

	// Check that the HTML contains window.close()
	if !strings.Contains(html, "window.close()") {
		t.Error("Success HTML should contain window.close() function")
	}

	// Check that the timeout is set to a reasonable value (3 seconds)
	if !strings.Contains(html, "3000") {
		t.Error("Success HTML should set timeout to 3 seconds")
	}

	// Verify the HTML structure is valid
	if !strings.Contains(html, "<script>") {
		t.Error("Success HTML should contain script tag")
	}

	if !strings.Contains(html, "</script>") {
		t.Error("Success HTML should have closing script tag")
	}

	// Check for visibility change handler
	if !strings.Contains(html, "visibilitychange") {
		t.Error("Success HTML should contain visibilitychange event listener")
	}

	// Verify that countdown and manual close elements are NOT present
	if strings.Contains(html, "countdown") {
		t.Error("Success HTML should not contain countdown elements")
	}

	if strings.Contains(html, "manual-close") {
		t.Error("Success HTML should not contain manual-close elements")
	}

	if strings.Contains(html, "Try to Close") {
		t.Error("Success HTML should not contain manual close button")
	}
}

func TestAuthSuccessPageBasicStructure(t *testing.T) {
	html := auth.GetSuccessHTML()

	// Basic HTML structure checks
	requiredElements := []string{
		"<!DOCTYPE html>",
		"<html lang=\"en\">",
		"</html>",
		"<head>",
		"</head>",
		"<body>",
		"</body>",
		"Authentication Successful",
	}

	for _, element := range requiredElements {
		if !strings.Contains(html, element) {
			t.Errorf("Success HTML should contain: %s", element)
		}
	}
}
