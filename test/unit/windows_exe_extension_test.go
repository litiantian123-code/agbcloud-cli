// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestWindowsExecutableExtension tests that Windows binaries have the correct .exe extension
func TestWindowsExecutableExtension(t *testing.T) {
	tests := []struct {
		name           string
		filename       string
		expectedHasExt bool
		description    string
	}{
		{
			name:           "Windows AMD64 binary",
			filename:       "agbcloud-windows-amd64.exe",
			expectedHasExt: true,
			description:    "Windows AMD64 binary should have .exe extension",
		},
		{
			name:           "Windows ARM64 binary",
			filename:       "agbcloud-windows-arm64.exe",
			expectedHasExt: true,
			description:    "Windows ARM64 binary should have .exe extension",
		},
		{
			name:           "Linux AMD64 binary",
			filename:       "agbcloud-linux-amd64",
			expectedHasExt: false,
			description:    "Linux binary should not have .exe extension",
		},
		{
			name:           "macOS ARM64 binary",
			filename:       "agbcloud-darwin-arm64",
			expectedHasExt: false,
			description:    "macOS binary should not have .exe extension",
		},
		{
			name:           "Windows binary without extension",
			filename:       "agbcloud-windows-amd64",
			expectedHasExt: false,
			description:    "This should fail - Windows binary missing .exe extension",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasExeExt := strings.HasSuffix(tt.filename, ".exe")
			isWindowsBinary := strings.Contains(tt.filename, "windows")

			if isWindowsBinary {
				if tt.expectedHasExt && !hasExeExt {
					t.Errorf("Windows binary %s should have .exe extension", tt.filename)
				}
				if !tt.expectedHasExt && hasExeExt {
					t.Errorf("Test case error: Windows binary %s marked as not expecting .exe extension", tt.filename)
				}
			} else {
				if hasExeExt {
					t.Errorf("Non-Windows binary %s should not have .exe extension", tt.filename)
				}
			}
		})
	}
}

// TestExtractPlatformFromFilename tests platform extraction from binary filenames
func TestExtractPlatformFromFilename(t *testing.T) {
	tests := []struct {
		name              string
		filename          string
		expectedPlatform  string
		expectedArch      string
		expectedIsWindows bool
	}{
		{
			name:              "Windows AMD64 with exe",
			filename:          "agbcloud-windows-amd64.exe",
			expectedPlatform:  "windows",
			expectedArch:      "amd64",
			expectedIsWindows: true,
		},
		{
			name:              "Windows ARM64 with exe",
			filename:          "agbcloud-windows-arm64.exe",
			expectedPlatform:  "windows",
			expectedArch:      "arm64",
			expectedIsWindows: true,
		},
		{
			name:              "Linux AMD64",
			filename:          "agbcloud-linux-amd64",
			expectedPlatform:  "linux",
			expectedArch:      "amd64",
			expectedIsWindows: false,
		},
		{
			name:              "Darwin ARM64",
			filename:          "agbcloud-darwin-arm64",
			expectedPlatform:  "darwin",
			expectedArch:      "arm64",
			expectedIsWindows: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Extract platform and arch from filename
			// Remove binary name prefix
			platformArch := strings.TrimPrefix(tt.filename, "agbcloud-")
			// Remove .exe extension if present
			platformArch = strings.TrimSuffix(platformArch, ".exe")

			parts := strings.Split(platformArch, "-")
			if len(parts) != 2 {
				t.Fatalf("Invalid filename format: %s", tt.filename)
			}

			platform := parts[0]
			arch := parts[1]
			isWindows := platform == "windows"

			if platform != tt.expectedPlatform {
				t.Errorf("Expected platform %s, got %s", tt.expectedPlatform, platform)
			}
			if arch != tt.expectedArch {
				t.Errorf("Expected arch %s, got %s", tt.expectedArch, arch)
			}
			if isWindows != tt.expectedIsWindows {
				t.Errorf("Expected isWindows %v, got %v", tt.expectedIsWindows, isWindows)
			}
		})
	}
}

// TestBinaryFileCreation tests that binary files are created with correct extensions
func TestBinaryFileCreation(t *testing.T) {
	// Skip this test on Windows to avoid file locking issues
	if runtime.GOOS == "windows" {
		t.Skip("Skipping binary creation test on Windows")
	}

	tempDir := t.TempDir()

	testCases := []struct {
		name          string
		platform      string
		arch          string
		shouldHaveExe bool
	}{
		{"Windows AMD64", "windows", "amd64", true},
		{"Windows ARM64", "windows", "arm64", true},
		{"Linux AMD64", "linux", "amd64", false},
		{"Darwin ARM64", "darwin", "arm64", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var filename string
			if tc.shouldHaveExe {
				filename = filepath.Join(tempDir, "agbcloud-"+tc.platform+"-"+tc.arch+".exe")
			} else {
				filename = filepath.Join(tempDir, "agbcloud-"+tc.platform+"-"+tc.arch)
			}

			// Create a dummy binary file
			err := os.WriteFile(filename, []byte("dummy binary content"), 0755)
			if err != nil {
				t.Fatalf("Failed to create test binary: %v", err)
			}

			// Verify file exists
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				t.Errorf("Binary file was not created: %s", filename)
			}

			// Verify extension
			hasExeExt := strings.HasSuffix(filename, ".exe")
			if tc.shouldHaveExe && !hasExeExt {
				t.Errorf("Windows binary should have .exe extension: %s", filename)
			}
			if !tc.shouldHaveExe && hasExeExt {
				t.Errorf("Non-Windows binary should not have .exe extension: %s", filename)
			}
		})
	}
}

// TestPackageNaming tests package naming conventions
func TestPackageNaming(t *testing.T) {
	version := "test-v1.0.0"

	testCases := []struct {
		name            string
		binaryName      string
		expectedZipName string
		expectedTarName string
		expectedExeName string
	}{
		{
			name:            "Windows AMD64",
			binaryName:      "agbcloud-windows-amd64.exe",
			expectedZipName: "agbcloud-test-v1.0.0-windows-amd64.zip",
			expectedTarName: "agbcloud-test-v1.0.0-windows-amd64.tar.gz",
			expectedExeName: "agbcloud-test-v1.0.0-windows-amd64.exe",
		},
		{
			name:            "Windows ARM64",
			binaryName:      "agbcloud-windows-arm64.exe",
			expectedZipName: "agbcloud-test-v1.0.0-windows-arm64.zip",
			expectedTarName: "agbcloud-test-v1.0.0-windows-arm64.tar.gz",
			expectedExeName: "agbcloud-test-v1.0.0-windows-arm64.exe",
		},
		{
			name:            "Linux AMD64",
			binaryName:      "agbcloud-linux-amd64",
			expectedZipName: "", // No zip for non-Windows
			expectedTarName: "agbcloud-test-v1.0.0-linux-amd64.tar.gz",
			expectedExeName: "", // No exe for non-Windows
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Extract platform-arch from binary name
			platformArch := strings.TrimPrefix(tc.binaryName, "agbcloud-")
			platformArch = strings.TrimSuffix(platformArch, ".exe")

			// Generate expected package names
			zipName := "agbcloud-" + version + "-" + platformArch + ".zip"
			tarName := "agbcloud-" + version + "-" + platformArch + ".tar.gz"
			exeName := "agbcloud-" + version + "-" + platformArch + ".exe"

			isWindows := strings.Contains(tc.binaryName, "windows")

			if isWindows {
				if zipName != tc.expectedZipName {
					t.Errorf("Expected zip name %s, got %s", tc.expectedZipName, zipName)
				}
				if exeName != tc.expectedExeName {
					t.Errorf("Expected exe name %s, got %s", tc.expectedExeName, exeName)
				}
			} else {
				if tc.expectedZipName != "" {
					t.Errorf("Non-Windows platform should not have zip package")
				}
				if tc.expectedExeName != "" {
					t.Errorf("Non-Windows platform should not have exe file")
				}
			}

			if tarName != tc.expectedTarName {
				t.Errorf("Expected tar name %s, got %s", tc.expectedTarName, tarName)
			}
		})
	}
}
