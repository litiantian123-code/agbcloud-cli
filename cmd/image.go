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

func init() {
	// Add flags for create command
	imageCreateCmd.Flags().StringP("dockerfile", "f", "", "Path to Dockerfile (required)")
	imageCreateCmd.Flags().StringP("imageId", "i", "", "Source image ID (required)")
	// Note: We handle required flag validation manually for better error messages

	// Add flags for activate command
	imageActivateCmd.Flags().IntP("cpu", "c", 0, "CPU cores")
	imageActivateCmd.Flags().IntP("memory", "m", 0, "Memory in GB")

	// Add subcommands to image command
	ImageCmd.AddCommand(imageCreateCmd)
	ImageCmd.AddCommand(imageActivateCmd)
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
	fmt.Printf("💾 CPU: %d cores, Memory: %d GB\n", cpu, memory)

	// TODO: Implement activate logic when API is ready
	fmt.Println("⚠️  Image activation API is not yet available")
	fmt.Println("📝 This command will be implemented when the backend API is ready")

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

			switch strings.ToUpper(status) {
			case "SUCCESS", "COMPLETED", "FINISHED":
				if taskResp.Data.ImageID != nil {
					fmt.Printf("🎉 Image created successfully! Image ID: %s\n", *taskResp.Data.ImageID)
				} else {
					fmt.Println("🎉 Image created successfully!")
				}
				return nil
			case "FAILED", "ERROR":
				fmt.Printf("📋 Task ID: %s\n", taskId)
				fmt.Printf("🔍 Request ID: %s\n", taskResp.RequestID)
				return fmt.Errorf("image creation failed: %s", message)
			case "RUNNING", "PENDING", "IN_PROGRESS":
				// Continue polling
				continue
			default:
				fmt.Printf("🔄 Unknown status '%s', continuing to monitor...\n", status)
				continue
			}
		}
	}
}
