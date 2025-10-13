// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// CallbackServerConfig holds configuration for the callback server
type CallbackServerConfig struct {
	Port string
}

// GetCallbackPort returns the default callback port
// Port selection is now handled automatically by the server's alternativePorts mechanism
func GetCallbackPort() string {
	// Always use default port - alternative ports are provided by server
	return "3000"
}

// StartCallbackServer starts a local HTTP server to handle OAuth callbacks
func StartCallbackServer(ctx context.Context, port string) (string, error) {
	var code string
	var err error
	var wg sync.WaitGroup
	wg.Add(1)

	// Create a new ServeMux to avoid conflicts with global handlers
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		// Get authorization code
		code = r.URL.Query().Get("code")
		if code == "" {
			err = fmt.Errorf("no code in callback")
			http.Error(w, "No code", http.StatusBadRequest)
			return
		}

		// Return success page
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		if _, writeErr := w.Write([]byte(GetSuccessHTML())); writeErr != nil {
			// Log the error but don't fail the authentication
			// The code has already been captured successfully
			err = fmt.Errorf("warning: failed to write success page: %w", writeErr)
		}

		// Delay server close to ensure browser receives the success page
		go func() {
			time.Sleep(500 * time.Millisecond)
			server.Close()
		}()
	})

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log startup errors for debugging
			// The callback will timeout if server fails to start
			_ = err // Error is intentionally not propagated
		}
	}()

	// Wait for callback or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Callback received
		if err != nil {
			return "", err
		}
		return code, nil
	case <-ctx.Done():
		server.Close()
		return "", fmt.Errorf("callback timeout: %v", ctx.Err())
	}
}

// GenerateRandomState generates a random state parameter for OAuth
func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// GetSuccessHTML returns the HTML page shown after successful authentication
func GetSuccessHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Authentication Successful - AgbCloud CLI</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            line-height: 1.6;
        }
        
        .container {
            text-align: center;
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            padding: 3rem 2.5rem;
            border-radius: 20px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            max-width: 450px;
            width: 90%;
            position: relative;
            overflow: hidden;
        }
        
        .container::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 4px;
            background: linear-gradient(90deg, #4CAF50, #45a049);
        }
        
        .success-icon {
            width: 80px;
            height: 80px;
            margin: 0 auto 1.5rem;
            background: linear-gradient(135deg, #4CAF50, #45a049);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 2.5rem;
            color: white;
            animation: pulse 2s infinite;
        }
        
        @keyframes pulse {
            0% { transform: scale(1); }
            50% { transform: scale(1.05); }
            100% { transform: scale(1); }
        }
        
        h1 {
            color: #2c3e50;
            font-size: 1.8rem;
            font-weight: 600;
            margin-bottom: 1rem;
            letter-spacing: -0.5px;
        }
        
        .subtitle {
            color: #7f8c8d;
            font-size: 1rem;
            margin-bottom: 2rem;
            font-weight: 400;
        }
        
        .info-box {
            background: #f8f9fa;
            border-left: 4px solid #4CAF50;
            padding: 1rem;
            border-radius: 8px;
            margin-bottom: 1.5rem;
            text-align: left;
        }
        
        .info-title {
            font-weight: 600;
            color: #2c3e50;
            margin-bottom: 0.5rem;
            font-size: 0.9rem;
        }
        
        .info-text {
            color: #5a6c7d;
            font-size: 0.85rem;
        }
        
        .brand {
            position: absolute;
            bottom: 1rem;
            left: 50%;
            transform: translateX(-50%);
            color: #bdc3c7;
            font-size: 0.75rem;
            font-weight: 500;
        }
        
        @media (max-width: 480px) {
            .container {
                padding: 2rem 1.5rem;
                margin: 1rem;
            }
            
            h1 {
                font-size: 1.5rem;
            }
            
            .success-icon {
                width: 60px;
                height: 60px;
                font-size: 2rem;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="success-icon">✓</div>
        <h1>Authentication Successful!</h1>
        <p class="subtitle">You have been successfully authenticated with AgbCloud CLI</p>
        
        <div class="info-box">
            <div class="info-title">Next Steps:</div>
            <div class="info-text">You can now close this window and return to your terminal to continue using the CLI.</div>
        </div>
        
        <div class="brand">AgbCloud CLI</div>
    </div>

    <script>
        // Silently attempt to close the window after 3 seconds
        // No timer display or user notification - just try to close
        setTimeout(function() {
            window.close();
        }, 3000);
        
        // Also try when user switches away from the tab (they're probably done reading)
        document.addEventListener('visibilitychange', function() {
            if (document.hidden) {
                setTimeout(function() {
                    window.close();
                }, 1000);
            }
        });
    </script>
</body>
</html>`
}

// IsPortOccupied checks if a given port is already in use
func IsPortOccupied(port string) bool {
	if !IsValidPort(port) {
		return true // Consider invalid ports as occupied
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true // Port is occupied or invalid
	}

	listener.Close()
	return false // Port is available
}

// IsValidPort checks if a port string is a valid port number
func IsValidPort(port string) bool {
	if port == "" {
		return false
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	return portNum > 0 && portNum <= 65535
}

// ParseAlternativePorts parses a comma-separated string of ports into a slice
func ParseAlternativePorts(alternativePorts string) []string {
	if alternativePorts == "" {
		return []string{}
	}

	ports := strings.Split(alternativePorts, ",")
	var validPorts []string

	for _, port := range ports {
		port = strings.TrimSpace(port)
		if port != "" {
			validPorts = append(validPorts, port)
		}
	}

	return validPorts
}

// SelectAvailablePort selects an available port from default and alternative ports
func SelectAvailablePort(defaultPort, alternativePorts string) (string, error) {
	// First try the default port
	if !IsPortOccupied(defaultPort) {
		return defaultPort, nil
	}

	// Parse alternative ports
	altPorts := ParseAlternativePorts(alternativePorts)

	// Try each alternative port
	for _, port := range altPorts {
		if IsValidPort(port) && !IsPortOccupied(port) {
			return port, nil
		}
	}

	// No available port found - provide detailed error message
	if len(altPorts) == 0 {
		return "", fmt.Errorf("no available port found: default port %s is occupied and no alternative ports provided", defaultPort)
	}

	return "", fmt.Errorf("no available port found: default port %s and all alternative ports [%s] are occupied. Please check if any of these ports can be freed up", defaultPort, alternativePorts)
}
