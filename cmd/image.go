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
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

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
			return fmt.Errorf("❌ Missing required argument: <image-name>\n\n💡 Usage: agbcloud image create <image-name> --dockerfile <path> --imageId <id>\n📝 Example: agbcloud image create myImage --dockerfile ./Dockerfile --imageId agb-code-space-1\n📝 Short form: agbcloud image create myImage -f ./Dockerfile -i agb-code-space-1")
		}
		if len(args) > 1 {
			return fmt.Errorf("❌ Too many arguments provided. Expected 1 argument (image name), got %d\n\n💡 Usage: agbcloud image create <image-name> --dockerfile <path> --imageId <id>\n📝 Example: agbcloud image create myImage --dockerfile ./Dockerfile --imageId agb-code-space-1\n📝 Short form: agbcloud image create myImage -f ./Dockerfile -i agb-code-space-1", len(args))
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
	Long:  "Activate an image with specified resources",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("❌ Missing required argument: <image-id>\n\n💡 Usage: agbcloud image activate <image-id> --cpu <cores> --memory <gb>\n📝 Example: agbcloud image activate img-7a8b9c1d0e --cpu 2 --memory 4\n📝 Short form: agbcloud image activate img-7a8b9c1d0e -c 2 -m 4")
		}
		if len(args) > 1 {
			return fmt.Errorf("❌ Too many arguments provided. Expected 1 argument (image ID), got %d\n\n💡 Usage: agbcloud image activate <image-id> --cpu <cores> --memory <gb>\n📝 Example: agbcloud image activate img-7a8b9c1d0e --cpu 2 --memory 4\n📝 Short form: agbcloud image activate img-7a8b9c1d0e -c 2 -m 4", len(args))
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
			return fmt.Errorf("❌ Missing required argument: <image-id>\n\n💡 Usage: agbcloud image deactivate <image-id>\n📝 Example: agbcloud image deactivate img-7a8b9c1d0e")
		}
		if len(args) > 1 {
			return fmt.Errorf("❌ Too many arguments provided. Expected 1 argument (image ID), got %d\n\n💡 Usage: agbcloud image deactivate <image-id>\n📝 Example: agbcloud image deactivate img-7a8b9c1d0e", len(args))
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

func runImageCreate(cmd *cobra.Command, args []string) error {
	imageName := args[0]
	dockerfilePath, _ := cmd.Flags().GetString("dockerfile")
	sourceImageId, _ := cmd.Flags().GetString("imageId")

	// Validate required flags with friendly messages
	if dockerfilePath == "" {
		return fmt.Errorf("❌ Missing required flag: --dockerfile\n\n💡 Usage: agbcloud image create %s --dockerfile <path> --imageId <id>\n📝 Example: agbcloud image create %s --dockerfile ./Dockerfile --imageId agb-code-space-1\n📝 Short form: agbcloud image create %s -f ./Dockerfile -i agb-code-space-1", imageName, imageName, imageName)
	}
	if sourceImageId == "" {
		return fmt.Errorf("❌ Missing required flag: --imageId\n\n💡 Usage: agbcloud image create %s --dockerfile <path> --imageId <id>\n📝 Example: agbcloud image create %s --dockerfile ./Dockerfile --imageId agb-code-space-1\n📝 Short form: agbcloud image create %s -f ./Dockerfile -i agb-code-space-1", imageName, imageName, imageName)
	}

	fmt.Printf("🏗️  Creating image '%s'...\n", imageName)

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Step 1: Get upload credential
	fmt.Println("📡 Getting upload credentials...")
	uploadResp, httpResp, err := apiClient.ImageAPI.GetUploadCredential(ctx, cfg.Token.LoginToken, cfg.Token.SessionId)
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("❌ API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("📊 Status Code: %d\n", httpResp.StatusCode)
			}
			return fmt.Errorf("failed to get upload credentials: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	if !uploadResp.Success {
		fmt.Printf("🔍 Request ID: %s\n", uploadResp.RequestID)
		return fmt.Errorf("failed to get upload credentials: %s", uploadResp.Code)
	}

	fmt.Printf("✅ Upload credentials obtained (Task ID: %s)\n", uploadResp.Data.TaskID)

	// Step 2: Upload dockerfile
	fmt.Println("📤 Uploading Dockerfile...")
	err = uploadDockerfile(dockerfilePath, uploadResp.Data.OssURL)
	if err != nil {
		fmt.Printf("📋 Task ID: %s\n", uploadResp.Data.TaskID)
		return fmt.Errorf("failed to upload dockerfile: %w", err)
	}

	fmt.Println("✅ Dockerfile uploaded successfully")

	// Step 3: Create image
	fmt.Println("🔨 Creating image...")
	createResp, httpResp, err := apiClient.ImageAPI.CreateImage(ctx, cfg.Token.LoginToken, cfg.Token.SessionId, imageName, uploadResp.Data.TaskID, sourceImageId)
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("❌ API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("📊 Status Code: %d\n", httpResp.StatusCode)
			}
			fmt.Printf("📋 Task ID: %s\n", uploadResp.Data.TaskID)
			return fmt.Errorf("failed to create image: %s", apiErr.Error())
		}
		fmt.Printf("📋 Task ID: %s\n", uploadResp.Data.TaskID)
		return fmt.Errorf("network error: %v", err)
	}

	if !createResp.Success {
		fmt.Printf("📋 Task ID: %s\n", uploadResp.Data.TaskID)
		fmt.Printf("🔍 Request ID: %s\n", createResp.RequestID)
		return fmt.Errorf("failed to create image: %s", createResp.Code)
	}

	fmt.Println("✅ Image creation initiated")

	// Step 4: Poll for task status
	fmt.Println("⏳ Monitoring image creation progress...")
	return pollImageTask(ctx, apiClient, cfg.Token.LoginToken, cfg.Token.SessionId, uploadResp.Data.TaskID)
}

func runImageActivate(cmd *cobra.Command, args []string) error {
	imageId := args[0]
	cpu, _ := cmd.Flags().GetInt("cpu")
	memory, _ := cmd.Flags().GetInt("memory")

	fmt.Printf("🚀 Activating image '%s'...\n", imageId)
	if cpu > 0 || memory > 0 {
		fmt.Printf("💾 CPU: %d cores, Memory: %d GB\n", cpu, memory)
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
	fmt.Println("🔍 Checking current image status...")
	listResp, httpResp, err := apiClient.ImageAPI.ListImages(ctx, cfg.Token.LoginToken, cfg.Token.SessionId, "User", 1, 1, []string{imageId})
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("❌ API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("📊 Status Code: %d\n", httpResp.StatusCode)
			}
			return fmt.Errorf("failed to check image status: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	if !listResp.Success {
		fmt.Printf("🔍 Request ID: %s\n", listResp.RequestID)
		return fmt.Errorf("failed to check image status: %s", listResp.Code)
	}

	// Check if image exists
	if len(listResp.Data.Images) == 0 {
		return fmt.Errorf("image not found: %s", imageId)
	}

	image := listResp.Data.Images[0]
	currentStatus := image.Status
	formattedStatus := formatImageStatus(currentStatus)

	fmt.Printf("📊 Current Status: %s\n", formattedStatus)

	// Handle different current statuses
	switch currentStatus {
	case "RESOURCE_PUBLISHED":
		fmt.Printf("✅ Image is already activated! Image ID: %s\n", imageId)
		fmt.Printf("📊 Status: %s\n", formattedStatus)
		return nil
	case "RESOURCE_DEPLOYING":
		fmt.Printf("🔄 Image is already activating, joining the activation process...\n")
		fmt.Println("⏳ Monitoring image activation status...")
		return pollImageActivationStatus(ctx, apiClient, cfg.Token.LoginToken, cfg.Token.SessionId, imageId)
	case "RESOURCE_FAILED", "RESOURCE_CEASED":
		fmt.Printf("⚠️  Image is in failed state (%s), attempting to restart activation...\n", formattedStatus)
	case "IMAGE_AVAILABLE":
		fmt.Printf("✅ Image is available, proceeding with activation...\n")
	default:
		fmt.Printf("📊 Image status: %s, proceeding with activation...\n", formattedStatus)
	}

	// Call StartImage API
	fmt.Println("🔄 Starting image activation...")
	startResp, httpResp, err := apiClient.ImageAPI.StartImage(ctx, cfg.Token.LoginToken, cfg.Token.SessionId, imageId, cpu, memory)
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("❌ API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("📊 Status Code: %d\n", httpResp.StatusCode)
			}
			return fmt.Errorf("failed to start image: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	if !startResp.Success {
		fmt.Printf("🔍 Request ID: %s\n", startResp.RequestID)
		return fmt.Errorf("failed to start image: %s", startResp.Code)
	}

	// Display success information
	fmt.Printf("✅ Image activation initiated successfully!\n")
	fmt.Printf("📊 Operation Status: %v\n", startResp.Data)
	fmt.Printf("🔍 Request ID: %s\n", startResp.RequestID)

	// Start status polling
	fmt.Println("⏳ Monitoring image activation status...")
	return pollImageActivationStatus(ctx, apiClient, cfg.Token.LoginToken, cfg.Token.SessionId, imageId)
}

func runImageDeactivate(cmd *cobra.Command, args []string) error {
	imageId := args[0]

	fmt.Printf("🛑 Deactivating image '%s'...\n", imageId)

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
	fmt.Println("🔄 Deactivating image instance...")
	stopResp, httpResp, err := apiClient.ImageAPI.StopImage(ctx, cfg.Token.LoginToken, cfg.Token.SessionId, imageId)
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("❌ API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("📊 Status Code: %d\n", httpResp.StatusCode)
			}
			return fmt.Errorf("failed to deactivate image: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	if !stopResp.Success {
		fmt.Printf("🔍 Request ID: %s\n", stopResp.RequestID)
		return fmt.Errorf("failed to deactivate image: %s", stopResp.Code)
	}

	// Display success information
	fmt.Printf("✅ Image deactivation initiated successfully!\n")
	fmt.Printf("📊 Operation Status: %v\n", stopResp.Data)
	fmt.Printf("🔍 Request ID: %s\n", stopResp.RequestID)

	return nil
}

func runImageList(cmd *cobra.Command, args []string) error {
	imageType, _ := cmd.Flags().GetString("type")
	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("size")

	fmt.Printf("📋 Listing %s images (Page %d, Size %d)...\n", imageType, page, pageSize)

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
	fmt.Println("🔍 Fetching image list...")
	listResp, httpResp, err := apiClient.ImageAPI.ListImages(ctx, cfg.Token.LoginToken, cfg.Token.SessionId, imageType, page, pageSize, nil)
	if err != nil {
		if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
			fmt.Printf("❌ API Error: %s\n", apiErr.Error())
			if httpResp != nil {
				fmt.Printf("📊 Status Code: %d\n", httpResp.StatusCode)
			}
			return fmt.Errorf("failed to list images: %s", apiErr.Error())
		}
		return fmt.Errorf("network error: %v", err)
	}

	if !listResp.Success {
		fmt.Printf("🔍 Request ID: %s\n", listResp.RequestID)
		return fmt.Errorf("failed to list images: %s", listResp.Code)
	}

	// Display results
	fmt.Printf("✅ Found %d images (Total: %d)\n", len(listResp.Data.Images), listResp.Data.Total)
	fmt.Printf("📄 Page %d of %d (Page Size: %d)\n\n", listResp.Data.Page, (listResp.Data.Total+listResp.Data.PageSize-1)/listResp.Data.PageSize, listResp.Data.PageSize)

	if len(listResp.Data.Images) == 0 {
		fmt.Println("📭 No images found.")
		return nil
	}

	// Display image table
	fmt.Printf("%-25s %-25s %-20s %-15s %-20s\n", "IMAGE ID", "IMAGE NAME", "STATUS", "TYPE", "UPDATED AT")
	fmt.Printf("%-25s %-25s %-20s %-15s %-20s\n", "--------", "----------", "------", "----", "----------")

	for _, image := range listResp.Data.Images {
		fmt.Printf("%-25s %-25s %-20s %-15s %-20s\n",
			image.ImageID, // Show full IMAGE ID without truncation
			truncateString(image.ImageName, 25),
			formatImageStatus(image.Status),
			truncateString(image.Type, 15),
			formatTimestamp(image.UpdateTime))
	}

	return nil
}

// uploadDockerfile uploads the dockerfile content to the provided OSS URL
func uploadDockerfile(dockerfilePath, ossURL string) error {
	// Read dockerfile content
	content, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to read dockerfile: %w", err)
	}

	// Create HTTP PUT request
	req, err := http.NewRequest(http.MethodPut, ossURL, strings.NewReader(string(content)))
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	// Set appropriate headers
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = int64(len(content))

	// Execute the upload
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload dockerfile: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read response body for error details
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// pollImageTask polls the image task status until completion or failure
func pollImageTask(ctx context.Context, apiClient *client.APIClient, loginToken, sessionId, taskId string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("📋 Task ID: %s\n", taskId)
			return fmt.Errorf("timeout waiting for image creation to complete")
		case <-ticker.C:
			taskResp, httpResp, err := apiClient.ImageAPI.GetImageTask(ctx, loginToken, sessionId, taskId)
			if err != nil {
				if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
					fmt.Printf("⚠️  Warning: Failed to check task status: %s\n", apiErr.Error())
					if httpResp != nil {
						fmt.Printf("📊 Status Code: %d\n", httpResp.StatusCode)
					}
					fmt.Printf("📋 Task ID: %s\n", taskId)
					continue // Continue polling on API errors
				}
				fmt.Printf("📋 Task ID: %s\n", taskId)
				return fmt.Errorf("network error checking task status: %v", err)
			}

			if !taskResp.Success {
				fmt.Printf("⚠️  Warning: Task status check failed: %s\n", taskResp.Code)
				fmt.Printf("📋 Task ID: %s\n", taskId)
				fmt.Printf("🔍 Request ID: %s\n", taskResp.RequestID)
				continue // Continue polling on API errors
			}

			status := taskResp.Data.Status
			message := taskResp.Data.TaskMsg

			fmt.Printf("📊 Status: %s", status)
			if message != "" {
				fmt.Printf(" - %s", message)
			}
			fmt.Println()

			switch status {
			case "Finished":
				if taskResp.Data.ImageID != nil {
					fmt.Printf("🎉 Image created successfully! Image ID: %s\n", *taskResp.Data.ImageID)
				} else {
					fmt.Println("🎉 Image created successfully!")
				}
				return nil
			case "Failed":
				fmt.Printf("📋 Task ID: %s\n", taskId)
				fmt.Printf("🔍 Request ID: %s\n", taskResp.RequestID)
				return fmt.Errorf("image creation failed: %s", message)
			case "Inline":
				// Continue polling - waiting for processing
				continue
			case "Preparing":
				// Continue polling - processing in progress
				continue
			default:
				fmt.Printf("🔄 Unknown status '%s', continuing to monitor...\n", status)
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

// formatImageStatus formats image status for better readability
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

// pollImageActivationStatus polls the image activation status until completion or failure
func pollImageActivationStatus(ctx context.Context, apiClient *client.APIClient, loginToken, sessionId, imageId string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Create a new context with longer timeout for polling
	pollCtx, pollCancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer pollCancel()

	for {
		select {
		case <-pollCtx.Done():
			fmt.Printf("📋 Image ID: %s\n", imageId)
			return fmt.Errorf("timeout waiting for image activation to complete")
		case <-ticker.C:
			// Query specific image status using ListImages with imageIds filter
			listResp, httpResp, err := apiClient.ImageAPI.ListImages(pollCtx, loginToken, sessionId, "User", 1, 1, []string{imageId})
			if err != nil {
				if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
					fmt.Printf("⚠️  Warning: Failed to check image status: %s\n", apiErr.Error())
					if httpResp != nil {
						fmt.Printf("📊 Status Code: %d\n", httpResp.StatusCode)
					}
					fmt.Printf("📋 Image ID: %s\n", imageId)
					continue // Continue polling on API errors
				}
				fmt.Printf("📋 Image ID: %s\n", imageId)
				return fmt.Errorf("network error checking image status: %v", err)
			}

			if !listResp.Success {
				fmt.Printf("⚠️  Warning: Image status check failed: %s\n", listResp.Code)
				fmt.Printf("📋 Image ID: %s\n", imageId)
				fmt.Printf("🔍 Request ID: %s\n", listResp.RequestID)
				continue // Continue polling on API errors
			}

			// Check if we found the image
			if len(listResp.Data.Images) == 0 {
				fmt.Printf("⚠️  Warning: Image not found: %s\n", imageId)
				continue // Continue polling
			}

			image := listResp.Data.Images[0]
			status := image.Status
			formattedStatus := formatImageStatus(status)

			fmt.Printf("📊 Status: %s", formattedStatus)
			fmt.Println()

			switch status {
			case "RESOURCE_PUBLISHED":
				fmt.Printf("🎉 Image activated successfully! Image ID: %s\n", imageId)
				fmt.Printf("📊 Final Status: %s\n", formattedStatus)
				return nil
			case "RESOURCE_FAILED", "RESOURCE_CEASED":
				fmt.Printf("📋 Image ID: %s\n", imageId)
				fmt.Printf("🔍 Request ID: %s\n", listResp.RequestID)
				return fmt.Errorf("image activation failed with status: %s", formattedStatus)
			case "RESOURCE_DEPLOYING":
				// Continue polling
				continue
			default:
				fmt.Printf("🔄 Unknown status '%s', continuing to monitor...\n", formattedStatus)
				continue
			}
		}
	}
}
