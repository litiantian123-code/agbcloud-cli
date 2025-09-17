// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

// printErrorMessage prints multi-line error messages by printing each line separately
// This avoids Windows line ending issues
func printErrorMessage(lines ...string) error {
	// Print each line to stderr for immediate display
	for _, line := range lines {
		fmt.Fprintln(os.Stderr, line)
	}
	// Return an error containing the full message for testing purposes
	fullMessage := strings.Join(lines, getNewline())
	return fmt.Errorf("%s", fullMessage)
}

// getNewline returns the appropriate newline character(s) for the current platform
func getNewline() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// nl is a convenience variable for newline
var nl = getNewline()

var ImageCmd = &cobra.Command{
	Use:     "image",
	Short:   "Manage images",
	Long:    "Create and manage custom images for AgbCloud",
	GroupID: "management",
}

var imageCreateCmd = &cobra.Command{
	Use:   "create <image-name>",
	Short: "Create a custom image",
	Long:  "Create a custom image using a Dockerfile",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return printErrorMessage(
				"[ERROR] Missing required argument: <image-name>",
				"",
				"[TIP] Usage: agbcloud image create <image-name> --dockerfile <path> --imageId <id>",
				"[NOTE] Example: agbcloud image create myImage --dockerfile ./Dockerfile --imageId agb-code-space-1",
				"[NOTE] Short form: agbcloud image create myImage -f ./Dockerfile -i agb-code-space-1",
			)
		}
		if len(args) > 1 {
			return printErrorMessage(
				fmt.Sprintf("[ERROR] Too many arguments provided. Expected 1 argument (image name), got %d", len(args)),
				"",
				"[TIP] Usage: agbcloud image create <image-name> --dockerfile <path> --imageId <id>",
				"[NOTE] Example: agbcloud image create myImage --dockerfile ./Dockerfile --imageId agb-code-space-1",
				"[NOTE] Short form: agbcloud image create myImage -f ./Dockerfile -i agb-code-space-1",
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runImageCreate(cmd, args)
	},
}

var imageActivateCmd = &cobra.Command{
	Use:   "activate <image-id>",
	Short: "Activate an image",
	Long: `Activate an image with specified resources.

Supported CPU and Memory combinations:
  2c4g  - 2 CPU cores with 4 GB memory
  4c8g  - 4 CPU cores with 8 GB memory  
  8c16g - 8 CPU cores with 16 GB memory

If no CPU/memory is specified, default resources will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return printErrorMessage(
				"[ERROR] Missing required argument: <image-id>",
				"",
				"[TIP] Usage: agbcloud image activate <image-id> [--cpu <cores> --memory <gb>]",
				"[NOTE] Example: agbcloud image activate img-7a8b9c1d0e --cpu 2 --memory 4",
			)
		}
		if len(args) > 1 {
			return printErrorMessage(
				fmt.Sprintf("[ERROR] Too many arguments provided. Expected 1 argument (image ID), got %d", len(args)),
				"",
				"[TIP] Usage: agbcloud image activate <image-id> [--cpu <cores> --memory <gb>]",
				"[NOTE] Example: agbcloud image activate img-7a8b9c1d0e --cpu 2 --memory 4",
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runImageActivate(cmd, args)
	},
}

var imageDeactivateCmd = &cobra.Command{
	Use:   "deactivate <image-id>",
	Short: "Deactivate an image",
	Long:  "Deactivate a running image instance",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return printErrorMessage(
				"[ERROR] Missing required argument: <image-id>",
				"",
				"[TIP] Usage: agbcloud image deactivate <image-id>",
				"[NOTE] Example: agbcloud image deactivate img-7a8b9c1d0e",
			)
		}
		if len(args) > 1 {
			return printErrorMessage(
				fmt.Sprintf("[ERROR] Too many arguments provided. Expected 1 argument (image ID), got %d", len(args)),
				"",
				"[TIP] Usage: agbcloud image deactivate <image-id>",
				"[NOTE] Example: agbcloud image deactivate img-7a8b9c1d0e",
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runImageDeactivate(cmd, args)
	},
}

var imageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List images",
	Long: `List images with pagination support.

Image types:
  User   - Custom images created by users
  System - System-provided base images`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runImageList(cmd, args)
	},
}

func init() {
	// Add flags for create command
	imageCreateCmd.Flags().StringP("dockerfile", "f", "", "Path to Dockerfile (required)")
	imageCreateCmd.Flags().StringP("imageId", "i", "", "Source image ID (required)")
	// Note: We handle required flag validation manually for better error messages

	// Add flags for activate command
	imageActivateCmd.Flags().IntP("cpu", "c", 0, "CPU cores")
	imageActivateCmd.Flags().IntP("memory", "m", 0, "Memory in GB")

	// Add flags for list command
	imageListCmd.Flags().StringP("type", "t", "User", "Image type: User (custom images) or System (base images)")
	imageListCmd.Flags().IntP("page", "p", 1, "Page number (default: 1)")
	imageListCmd.Flags().IntP("size", "s", 10, "Page size (default: 10)")

	// Add subcommands to image command
	ImageCmd.AddCommand(imageCreateCmd)
	ImageCmd.AddCommand(imageActivateCmd)
	ImageCmd.AddCommand(imageDeactivateCmd)
	ImageCmd.AddCommand(imageListCmd)
}

// ValidateCPUMemoryCombo validates that CPU and memory combination is supported
func ValidateCPUMemoryCombo(cpu, memory int) error {
	// If both are 0, use default (no validation needed)
	if cpu == 0 && memory == 0 {
		return nil
	}

	// If only one is specified, both must be specified
	if (cpu == 0 && memory > 0) || (cpu > 0 && memory == 0) {
		return printErrorMessage(
			"[ERROR] Both CPU and memory must be specified together",
			"",
			"[TOOL] Supported combinations:",
			"  • 2c4g: --cpu 2 --memory 4",
			"  • 4c8g: --cpu 4 --memory 8",
			"  • 8c16g: --cpu 8 --memory 16",
		)
	}

	// Check supported combinations
	validCombos := map[int]int{
		2: 4,  // 2c4g
		4: 8,  // 4c8g
		8: 16, // 8c16g
	}

	expectedMemory, exists := validCombos[cpu]
	if !exists || expectedMemory != memory {
		return printErrorMessage(
			fmt.Sprintf("[ERROR] Invalid CPU/Memory combination: %dc%dg", cpu, memory),
			"",
			"[TOOL] Supported combinations:",
			"  • 2c4g: --cpu 2 --memory 4",
			"  • 4c8g: --cpu 4 --memory 8",
			"  • 8c16g: --cpu 8 --memory 16",
		)
	}

	return nil
}

func runImageCreate(cmd *cobra.Command, args []string) error {
	imageName := args[0]
	dockerfilePath, _ := cmd.Flags().GetString("dockerfile")
	sourceImageId, _ := cmd.Flags().GetString("imageId")

	// Validate required flags with friendly messages
	if dockerfilePath == "" {
		return printErrorMessage(
			fmt.Sprintf("[ERROR] Missing required flag: --dockerfile for %s", imageName),
			"",
			"[TIP] Usage: agbcloud image create %s --dockerfile <path> --imageId <id>",
			fmt.Sprintf("[NOTE] Example: agbcloud image create %s --dockerfile ./Dockerfile --imageId agb-code-space-1", imageName),
			fmt.Sprintf("[NOTE] Short form: agbcloud image create %s -f ./Dockerfile -i agb-code-space-1", imageName),
		)
	}
	if sourceImageId == "" {
		return printErrorMessage(
			fmt.Sprintf("[ERROR] Missing required flag: --imageId for %s", imageName),
			"",
			"[TIP] Usage: agbcloud image create %s --dockerfile <path> --imageId <id>",
			fmt.Sprintf("[NOTE] Example: agbcloud image create %s --dockerfile ./Dockerfile --imageId agb-code-space-1", imageName),
			fmt.Sprintf("[NOTE] Short form: agbcloud image create %s -f ./Dockerfile -i agb-code-space-1", imageName),
		)
	}

	fmt.Printf("[BUILD]  Creating image '%s'...\n", imageName)

	// Load configuration and check authentication
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if cfg.Token == nil || cfg.Token.LoginToken == "" || cfg.Token.SessionId == "" {
		return fmt.Errorf("not authenticated. Please run 'agbcloud login' first")
	}

	// Validate dockerfile path
	if !filepath.IsAbs(dockerfilePath) {
		dockerfilePath, err = filepath.Abs(dockerfilePath)
		if err != nil {
			return fmt.Errorf("failed to resolve dockerfile path: %w", err)
		}
	}

	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return fmt.Errorf("dockerfile not found: %s", dockerfilePath)
	}

	// Create API client
	apiClient := client.NewFromConfig(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
	defer cancel()

	// Step 1: Get upload credential
	fmt.Println("[SIGNAL] Getting upload credentials...")
	uploadResp, httpResp, err := apiClient.ImageAPI.GetUploadCredential(ctx, cfg.Token.LoginToken, cfg.Token.SessionId)
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("[ERROR] API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("[DATA] Status Code: %d\n", httpResp.StatusCode)
			}
			return fmt.Errorf("failed to get upload credentials: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	if !uploadResp.Success {
		fmt.Printf("[SEARCH] Request ID: %s\n", uploadResp.RequestID)
		return fmt.Errorf("failed to get upload credentials: %s", uploadResp.Code)
	}

	fmt.Printf("[OK] Upload credentials obtained (Task ID: %s)\n", uploadResp.Data.TaskID)

	// Step 2: Upload dockerfile
	fmt.Println("[UPLOAD] Uploading Dockerfile...")
	err = uploadDockerfile(dockerfilePath, uploadResp.Data.OssURL)
	if err != nil {
		fmt.Printf("[DOC] Task ID: %s\n", uploadResp.Data.TaskID)
		return fmt.Errorf("failed to upload dockerfile: %w", err)
	}

	fmt.Println("[OK] Dockerfile uploaded successfully")

	// Step 3: Create image
	fmt.Println("[WORK] Creating image...")
	createResp, httpResp, err := apiClient.ImageAPI.CreateImage(ctx, cfg.Token.LoginToken, cfg.Token.SessionId, imageName, uploadResp.Data.TaskID, sourceImageId)
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("[ERROR] API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("[DATA] Status Code: %d\n", httpResp.StatusCode)
			}
			fmt.Printf("[DOC] Task ID: %s\n", uploadResp.Data.TaskID)
			return fmt.Errorf("failed to create image: %s", apiErr.Error())
		}
		fmt.Printf("[DOC] Task ID: %s\n", uploadResp.Data.TaskID)
		return fmt.Errorf("network error: %v", err)
	}

	if !createResp.Success {
		fmt.Printf("[DOC] Task ID: %s\n", uploadResp.Data.TaskID)
		fmt.Printf("[SEARCH] Request ID: %s\n", createResp.RequestID)
		return fmt.Errorf("failed to create image: %s", createResp.Code)
	}

	fmt.Println("[OK] Image creation initiated")

	// Step 4: Poll for task status
	fmt.Println("[MONITOR] Monitoring image creation progress...")
	return pollImageTask(ctx, apiClient, cfg.Token.LoginToken, cfg.Token.SessionId, uploadResp.Data.TaskID)
}

func runImageActivate(cmd *cobra.Command, args []string) error {
	imageId := args[0]
	cpu, _ := cmd.Flags().GetInt("cpu")
	memory, _ := cmd.Flags().GetInt("memory")

	// Validate CPU and memory combination
	if err := ValidateCPUMemoryCombo(cpu, memory); err != nil {
		return err
	}

	fmt.Printf("[>>] Activating image '%s'...\n", imageId)
	if cpu > 0 || memory > 0 {
		fmt.Printf("[SAVE] CPU: %d cores, Memory: %d GB\n", cpu, memory)
	}

	// Load configuration and check authentication
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if cfg.Token == nil || cfg.Token.LoginToken == "" || cfg.Token.SessionId == "" {
		return fmt.Errorf("not authenticated. Please run 'agbcloud login' first")
	}

	// Create API client
	apiClient := client.NewFromConfig(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check current image status first
	fmt.Println("[SEARCH] Checking current image status...")
	listResp, httpResp, err := apiClient.ImageAPI.ListImages(ctx, cfg.Token.LoginToken, cfg.Token.SessionId, "User", 1, 1, []string{imageId})
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("[ERROR] API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("[DATA] Status Code: %d\n", httpResp.StatusCode)
			}
			return fmt.Errorf("failed to check image status: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	if !listResp.Success {
		fmt.Printf("[SEARCH] Request ID: %s\n", listResp.RequestID)
		return fmt.Errorf("failed to check image status: %s", listResp.Code)
	}

	// Check if image exists
	if len(listResp.Data.Images) == 0 {
		return fmt.Errorf("image not found: %s", imageId)
	}

	image := listResp.Data.Images[0]
	currentStatus := image.Status
	formattedStatus := FormatImageStatus(currentStatus)

	fmt.Printf("[DATA] Current Status: %s\n", formattedStatus)

	// Handle different current statuses
	switch currentStatus {
	case "RESOURCE_PUBLISHED":
		fmt.Printf("[OK] Image is already activated! Image ID: %s\n", imageId)
		fmt.Printf("[DATA] Status: %s\n", formattedStatus)
		return nil
	case "RESOURCE_DEPLOYING":
		fmt.Printf("[REFRESH] Image is already activating, joining the activation process...\n")
		fmt.Println("[MONITOR] Monitoring image activation status...")
		return pollImageActivationStatus(ctx, apiClient, cfg.Token.LoginToken, cfg.Token.SessionId, imageId)
	case "RESOURCE_FAILED", "RESOURCE_CEASED":
		fmt.Printf("[WARN]  Image is in failed state (%s), attempting to restart activation...\n", formattedStatus)
	case "IMAGE_AVAILABLE":
		fmt.Printf("[OK] Image is available, proceeding with activation...\n")
	default:
		fmt.Printf("[DATA] Image status: %s, proceeding with activation...\n", formattedStatus)
	}

	// Call StartImage API
	fmt.Println("[REFRESH] Starting image activation...")
	startResp, httpResp, err := apiClient.ImageAPI.StartImage(ctx, cfg.Token.LoginToken, cfg.Token.SessionId, imageId, cpu, memory)
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("[ERROR] API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("[DATA] Status Code: %d\n", httpResp.StatusCode)
			}
			return fmt.Errorf("failed to start image: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	if !startResp.Success {
		fmt.Printf("[SEARCH] Request ID: %s\n", startResp.RequestID)
		return fmt.Errorf("failed to start image: %s", startResp.Code)
	}

	// Display success information
	fmt.Printf("[OK] Image activation initiated successfully!\n")
	fmt.Printf("[DATA] Operation Status: %v\n", startResp.Data)
	fmt.Printf("[SEARCH] Request ID: %s\n", startResp.RequestID)

	// Start status polling
	fmt.Println("[MONITOR] Monitoring image activation status...")
	return pollImageActivationStatus(ctx, apiClient, cfg.Token.LoginToken, cfg.Token.SessionId, imageId)
}

func runImageDeactivate(cmd *cobra.Command, args []string) error {
	imageId := args[0]

	fmt.Printf("[STOP] Deactivating image '%s'...\n", imageId)

	// Load configuration and check authentication
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if cfg.Token == nil || cfg.Token.LoginToken == "" || cfg.Token.SessionId == "" {
		return fmt.Errorf("not authenticated. Please run 'agbcloud login' first")
	}

	// Create API client
	apiClient := client.NewFromConfig(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Call StopImage API
	fmt.Println("[REFRESH] Deactivating image instance...")
	stopResp, httpResp, err := apiClient.ImageAPI.StopImage(ctx, cfg.Token.LoginToken, cfg.Token.SessionId, imageId)
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("[ERROR] API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("[DATA] Status Code: %d\n", httpResp.StatusCode)
			}
			return fmt.Errorf("failed to deactivate image: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	if !stopResp.Success {
		fmt.Printf("[SEARCH] Request ID: %s\n", stopResp.RequestID)
		return fmt.Errorf("failed to deactivate image: %s", stopResp.Code)
	}

	// Display success information
	fmt.Printf("[OK] Image deactivation initiated successfully!\n")
	fmt.Printf("[DATA] Operation Status: %v\n", stopResp.Data)
	fmt.Printf("[SEARCH] Request ID: %s\n", stopResp.RequestID)

	// Start status polling
	fmt.Println("[MONITOR] Monitoring image deactivation status...")
	return pollImageDeactivationStatus(ctx, apiClient, cfg.Token.LoginToken, cfg.Token.SessionId, imageId)
}

func runImageList(cmd *cobra.Command, args []string) error {
	imageType, _ := cmd.Flags().GetString("type")
	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("size")

	fmt.Printf("[DOC] Listing %s images (Page %d, Size %d)...\n", imageType, page, pageSize)

	// Load configuration and check authentication
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if cfg.Token == nil || cfg.Token.LoginToken == "" || cfg.Token.SessionId == "" {
		return fmt.Errorf("not authenticated. Please run 'agbcloud login' first")
	}

	// Create API client
	apiClient := client.NewFromConfig(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Call ListImages API
	fmt.Println("[SEARCH] Fetching image list...")
	listResp, httpResp, err := apiClient.ImageAPI.ListImages(ctx, cfg.Token.LoginToken, cfg.Token.SessionId, imageType, page, pageSize, nil)
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("[ERROR] API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("[DATA] Status Code: %d\n", httpResp.StatusCode)
			}
			return fmt.Errorf("failed to list images: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	if !listResp.Success {
		fmt.Printf("[SEARCH] Request ID: %s\n", listResp.RequestID)
		return fmt.Errorf("failed to list images: %s", listResp.Code)
	}

	// Display results
	fmt.Printf("[OK] Found %d images (Total: %d)\n", len(listResp.Data.Images), listResp.Data.Total)
	fmt.Printf("[PAGE] Page %d of %d (Page Size: %d)\n\n", listResp.Data.Page, (listResp.Data.Total+listResp.Data.PageSize-1)/listResp.Data.PageSize, listResp.Data.PageSize)

	if len(listResp.Data.Images) == 0 {
		fmt.Println("[EMPTY] No images found.")
		return nil
	}

	// Display image table with CPU/Memory information
	fmt.Printf("%-25s %-25s %-20s %-15s %-12s %-20s\n", "IMAGE ID", "IMAGE NAME", "STATUS", "TYPE", "CPU/MEMORY", "UPDATED AT")
	fmt.Printf("%-25s %-25s %-20s %-15s %-12s %-20s\n", "--------", "----------", "------", "----", "----------", "----------")

	for _, image := range listResp.Data.Images {
		fmt.Printf("%-25s %-25s %-20s %-15s %-12s %-20s\n",
			truncateString(image.ImageID, 25),
			truncateString(image.ImageName, 25),
			FormatImageStatus(image.Status),
			truncateString(image.Type, 15),
			FormatResources(image.CPU, image.Memory),
			formatTimestamp(image.UpdateTime))
	}

	return nil
}

// uploadDockerfile uploads the dockerfile content to the provided OSS URL with retry mechanism
func uploadDockerfile(dockerfilePath, ossURL string) error {
	// Read dockerfile content
	content, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to read dockerfile: %w", err)
	}

	// Create retry configuration for upload
	retryConfig := &client.RetryConfig{
		MaxRetries:    3,
		InitialDelay:  1 * time.Second,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
	}

	var lastErr error
	delay := retryConfig.InitialDelay

	for attempt := 0; attempt <= retryConfig.MaxRetries; attempt++ {
		fmt.Printf("[UPLOAD] Dockerfile upload attempt %d/%d...\n", attempt+1, retryConfig.MaxRetries+1)

		// Create HTTP PUT request for each attempt
		req, err := http.NewRequest(http.MethodPut, ossURL, strings.NewReader(string(content)))
		if err != nil {
			return fmt.Errorf("failed to create upload request: %w", err)
		}

		// Set appropriate headers
		req.Header.Set("Content-Type", "application/octet-stream")
		req.ContentLength = int64(len(content))

		// Execute the upload
		httpClient := &http.Client{Timeout: 60 * time.Second}
		resp, err := httpClient.Do(req)

		// Success case
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			if attempt > 0 {
				fmt.Printf("[OK] Dockerfile upload succeeded on attempt %d\n", attempt+1)
			}
			return nil
		}

		// Handle error cases
		if err != nil {
			lastErr = fmt.Errorf("failed to upload dockerfile: %w", err)
		} else {
			// Read response body for error details
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
		}

		// Don't retry if this is the last attempt
		if attempt == retryConfig.MaxRetries {
			break
		}

		// Check if the error is retryable
		shouldRetry := false
		if err != nil {
			// Use the same retry logic as the API client
			shouldRetry = client.IsRetryableError(err)
		} else if resp != nil {
			// Check if HTTP status is retryable
			shouldRetry = client.IsRetryableHTTPStatus(resp.StatusCode)
		}

		if !shouldRetry {
			fmt.Printf("[WARN]  Upload error is not retryable, stopping attempts\n")
			break
		}

		// Wait before retrying
		fmt.Printf("[RETRY] Upload failed (attempt %d/%d), retrying in %v...\n",
			attempt+1, retryConfig.MaxRetries+1, delay)

		time.Sleep(delay)

		// Calculate next delay with exponential backoff
		delay = time.Duration(float64(delay) * retryConfig.BackoffFactor)
		if delay > retryConfig.MaxDelay {
			delay = retryConfig.MaxDelay
		}
	}

	fmt.Printf("[ERROR] All %d upload attempts failed\n", retryConfig.MaxRetries+1)
	return fmt.Errorf("dockerfile upload failed after %d attempts, last error: %w",
		retryConfig.MaxRetries+1, lastErr)
}

// pollImageTask polls the image task status until completion or failure
func pollImageTask(ctx context.Context, apiClient *client.APIClient, loginToken, sessionId, taskId string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[DOC] Task ID: %s\n", taskId)
			return fmt.Errorf("timeout waiting for image creation to complete")
		case <-ticker.C:
			taskResp, httpResp, err := apiClient.ImageAPI.GetImageTask(ctx, loginToken, sessionId, taskId)
			if err != nil {
				if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
					fmt.Printf("[WARN]  Warning: Failed to check task status: %s\n", apiErr.Error())
					if httpResp != nil {
						fmt.Printf("[DATA] Status Code: %d\n", httpResp.StatusCode)
					}
					fmt.Printf("[DOC] Task ID: %s\n", taskId)
					continue // Continue polling on API errors
				}
				fmt.Printf("[DOC] Task ID: %s\n", taskId)
				return fmt.Errorf("network error checking task status: %v", err)
			}

			if !taskResp.Success {
				fmt.Printf("[WARN]  Warning: Task status check failed: %s\n", taskResp.Code)
				fmt.Printf("[DOC] Task ID: %s\n", taskId)
				fmt.Printf("[SEARCH] Request ID: %s\n", taskResp.RequestID)
				continue // Continue polling on API errors
			}

			status := taskResp.Data.Status
			message := taskResp.Data.TaskMsg

			fmt.Printf("[DATA] Status: %s", status)
			if message != "" {
				fmt.Printf(" - %s", message)
			}
			fmt.Println()

			switch status {
			case "Finished":
				if taskResp.Data.ImageID != nil {
					fmt.Printf("[SUCCESS] Image created successfully! Image ID: %s\n", *taskResp.Data.ImageID)
				} else {
					fmt.Println("[SUCCESS] Image created successfully!")
				}
				return nil
			case "Failed":
				fmt.Printf("[DOC] Task ID: %s\n", taskId)
				fmt.Printf("[SEARCH] Request ID: %s\n", taskResp.RequestID)
				return fmt.Errorf("image creation failed: %s", message)
			case "Inline":
				// Continue polling - waiting for processing
				continue
			case "Preparing":
				// Continue polling - processing in progress
				continue
			default:
				fmt.Printf("[REFRESH] Unknown status '%s', continuing to monitor...\n", status)
				continue
			}
		}
	}
}

// pollImageDeactivationStatus polls the image deactivation status until completion or failure
func pollImageDeactivationStatus(ctx context.Context, apiClient *client.APIClient, loginToken, sessionId, imageId string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Create a new context with longer timeout for polling
	pollCtx, pollCancel := context.WithTimeout(context.Background(), 45*time.Minute)
	defer pollCancel()

	for {
		select {
		case <-pollCtx.Done():
			fmt.Printf("[DOC] Image ID: %s\n", imageId)
			return fmt.Errorf("timeout waiting for image deactivation to complete")
		case <-ticker.C:
			// Query specific image status using ListImages with imageIds filter
			listResp, httpResp, err := apiClient.ImageAPI.ListImages(pollCtx, loginToken, sessionId, "User", 1, 1, []string{imageId})
			if err != nil {
				if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
					fmt.Printf("[WARN]  Warning: Failed to check image status: %s\n", apiErr.Error())
					if httpResp != nil {
						fmt.Printf("[DATA] Status Code: %d\n", httpResp.StatusCode)
					}
					fmt.Printf("[DOC] Image ID: %s\n", imageId)
					continue // Continue polling on API errors
				}
				fmt.Printf("[DOC] Image ID: %s\n", imageId)
				return fmt.Errorf("network error checking image status: %v", err)
			}

			if !listResp.Success {
				fmt.Printf("[WARN]  Warning: Image status check failed: %s\n", listResp.Code)
				fmt.Printf("[DOC] Image ID: %s\n", imageId)
				fmt.Printf("[SEARCH] Request ID: %s\n", listResp.RequestID)
				continue // Continue polling on API errors
			}

			// Check if we found the image
			if len(listResp.Data.Images) == 0 {
				fmt.Printf("[WARN]  Warning: Image not found: %s\n", imageId)
				continue // Continue polling
			}

			image := listResp.Data.Images[0]
			status := image.Status
			formattedStatus := FormatImageStatus(status)

			fmt.Printf("[DATA] Status: %s", formattedStatus)
			fmt.Println()

			switch status {
			case "IMAGE_AVAILABLE":
				fmt.Printf("[SUCCESS] Image deactivated successfully! Image ID: %s\n", imageId)
				fmt.Printf("[DATA] Final Status: %s\n", formattedStatus)
				return nil
			case "RESOURCE_FAILED":
				fmt.Printf("[DOC] Image ID: %s\n", imageId)
				fmt.Printf("[SEARCH] Request ID: %s\n", listResp.RequestID)
				return fmt.Errorf("image deactivation failed with status: %s", formattedStatus)
			case "RESOURCE_DELETING":
				// Continue polling - deactivation in progress
				continue
			case "RESOURCE_PUBLISHED":
				// Image is still activated, continue polling in case deactivation is delayed
				fmt.Printf("[REFRESH] Image still activated, continuing to monitor deactivation...\n")
				continue
			default:
				fmt.Printf("[REFRESH] Unknown status '%s', continuing to monitor...\n", formattedStatus)
				continue
			}
		}
	}
}

// truncateString truncates a string to the specified length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// formatTimestamp formats a timestamp string for display in local timezone
func formatTimestamp(timestamp string) string {
	if timestamp == "" {
		return "-"
	}
	// Try to parse the UTC timestamp and convert to local timezone
	if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
		// Convert UTC time to local timezone
		localTime := t.Local()
		return localTime.Format("2006-01-02 15:04")
	}
	// If parsing fails, return the original string truncated
	return truncateString(timestamp, 20)
}

// FormatImageStatus formats image status for better readability
func FormatImageStatus(status string) string {
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

// FormatCPU formats CPU value for display, handling null values gracefully
func FormatCPU(cpu *int) string {
	if cpu == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *cpu)
}

// FormatMemory formats Memory value for display, handling null values gracefully
func FormatMemory(memory *int) string {
	if memory == nil {
		return "-"
	}
	return fmt.Sprintf("%dG", *memory)
}

// FormatResources formats CPU and Memory together for compact display
func FormatResources(cpu *int, memory *int) string {
	if cpu == nil && memory == nil {
		return "-"
	}
	if cpu == nil {
		return fmt.Sprintf("-/%dG", *memory)
	}
	if memory == nil {
		return fmt.Sprintf("%d/-", *cpu)
	}
	return fmt.Sprintf("%d/%dG", *cpu, *memory)
}

// pollImageActivationStatus polls the image activation status until completion or failure
func pollImageActivationStatus(ctx context.Context, apiClient *client.APIClient, loginToken, sessionId, imageId string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Create a new context with longer timeout for polling
	pollCtx, pollCancel := context.WithTimeout(context.Background(), 45*time.Minute)
	defer pollCancel()

	for {
		select {
		case <-pollCtx.Done():
			fmt.Printf("[DOC] Image ID: %s\n", imageId)
			return fmt.Errorf("timeout waiting for image activation to complete")
		case <-ticker.C:
			// Query specific image status using ListImages with imageIds filter
			listResp, httpResp, err := apiClient.ImageAPI.ListImages(pollCtx, loginToken, sessionId, "User", 1, 1, []string{imageId})
			if err != nil {
				if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
					fmt.Printf("[WARN]  Warning: Failed to check image status: %s\n", apiErr.Error())
					if httpResp != nil {
						fmt.Printf("[DATA] Status Code: %d\n", httpResp.StatusCode)
					}
					fmt.Printf("[DOC] Image ID: %s\n", imageId)
					continue // Continue polling on API errors
				}
				fmt.Printf("[DOC] Image ID: %s\n", imageId)
				return fmt.Errorf("network error checking image status: %v", err)
			}

			if !listResp.Success {
				fmt.Printf("[WARN]  Warning: Image status check failed: %s\n", listResp.Code)
				fmt.Printf("[DOC] Image ID: %s\n", imageId)
				fmt.Printf("[SEARCH] Request ID: %s\n", listResp.RequestID)
				continue // Continue polling on API errors
			}

			// Check if we found the image
			if len(listResp.Data.Images) == 0 {
				fmt.Printf("[WARN]  Warning: Image not found: %s\n", imageId)
				continue // Continue polling
			}

			image := listResp.Data.Images[0]
			status := image.Status
			formattedStatus := FormatImageStatus(status)

			fmt.Printf("[DATA] Status: %s", formattedStatus)
			fmt.Println()

			switch status {
			case "RESOURCE_PUBLISHED":
				fmt.Printf("[SUCCESS] Image activated successfully! Image ID: %s\n", imageId)
				fmt.Printf("[DATA] Final Status: %s\n", formattedStatus)
				return nil
			case "RESOURCE_FAILED", "RESOURCE_CEASED":
				fmt.Printf("[DOC] Image ID: %s\n", imageId)
				fmt.Printf("[SEARCH] Request ID: %s\n", listResp.RequestID)
				return fmt.Errorf("image activation failed with status: %s", formattedStatus)
			case "RESOURCE_DEPLOYING":
				// Continue polling
				continue
			default:
				fmt.Printf("[REFRESH] Unknown status '%s', continuing to monitor...\n", formattedStatus)
				continue
			}
		}
	}
}
