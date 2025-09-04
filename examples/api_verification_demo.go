// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
)

func main() {
	fmt.Println("🔍 AgbCloud API Verification Demo")
	fmt.Println("==================================")

	// Create API client
	apiClient := client.NewDefault()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test 1: New login provider API with default parameters
	fmt.Println("\n📋 Test 1: Login Provider API with defaults")
	fmt.Println("Parameters: loginClient=CLI, oauthProvider=GOOGLE_LOCALHOST")

	response1, httpResp1, err1 := apiClient.OAuthAPI.GetLoginProviderURL(ctx, "https://agb.cloud", "CLI", "GOOGLE_LOCALHOST")
	if err1 != nil {
		log.Printf("❌ Error: %v", err1)
	} else {
		fmt.Printf("✅ Success! Request ID: %s\n", response1.RequestID)
		fmt.Printf("🔗 OAuth URL: %s\n", response1.Data.InvokeURL)
		if httpResp1 != nil && httpResp1.Request != nil {
			fmt.Printf("📡 Request URL: %s\n", httpResp1.Request.URL.String())
		}
	}

	// Test 2: New login provider API with custom parameters
	fmt.Println("\n📋 Test 2: Login Provider API with custom parameters")
	fmt.Println("Parameters: loginClient=WEB, oauthProvider=GOOGLE_WEB")

	response2, httpResp2, err2 := apiClient.OAuthAPI.GetLoginProviderURL(ctx, "https://agb.cloud/dashboard", "WEB", "GOOGLE_WEB")
	if err2 != nil {
		log.Printf("❌ Error: %v", err2)
	} else {
		fmt.Printf("✅ Success! Request ID: %s\n", response2.RequestID)
		fmt.Printf("🔗 OAuth URL: %s\n", response2.Data.InvokeURL)
		if httpResp2 != nil && httpResp2.Request != nil {
			fmt.Printf("📡 Request URL: %s\n", httpResp2.Request.URL.String())
		}
	}

	// Test 3: Legacy Google login API (backward compatibility)
	fmt.Println("\n📋 Test 3: Legacy Google Login API (backward compatibility)")
	fmt.Println("Should automatically use new endpoint with default parameters")

	response3, httpResp3, err3 := apiClient.OAuthAPI.GetGoogleLoginURL(ctx, "https://agb.cloud")
	if err3 != nil {
		log.Printf("❌ Error: %v", err3)
	} else {
		fmt.Printf("✅ Success! Request ID: %s\n", response3.RequestID)
		fmt.Printf("🔗 OAuth URL: %s\n", response3.Data.InvokeURL)
		if httpResp3 != nil && httpResp3.Request != nil {
			fmt.Printf("📡 Request URL: %s\n", httpResp3.Request.URL.String())
		}
	}

	fmt.Println("\n🎉 API verification completed!")
}
