// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// CallbackServerConfig holds configuration for the callback server
type CallbackServerConfig struct {
	Port string
}

// GetCallbackPort returns the callback port from environment variable, config, or default
func GetCallbackPort(configPort string) string {
	// Check environment variable first
	if envPort := os.Getenv("AGB_CLI_CALLBACK_PORT"); envPort != "" {
		return envPort
	}

	// Check config port
	if configPort != "" {
		return configPort
	}

	// Default port
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
		w.Write([]byte(GetSuccessHTML()))

		// Delay server close to ensure browser receives the success page
		go func() {
			time.Sleep(500 * time.Millisecond)
			server.Close()
		}()
	})

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			// Only set error if it's not the expected server closed error
			if err != nil {
				// Use a channel or other mechanism to communicate startup errors
				// For now, we'll let the timeout handle startup failures
			}
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
        
        .countdown {
            color: #95a5a6;
            font-size: 0.85rem;
            margin-top: 1.5rem;
            font-style: italic;
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
        
        <div class="countdown" id="countdown">This window will close automatically in <span id="timer">10</span> seconds</div>
        <div class="brand">AgbCloud CLI</div>
    </div>

    <script>
        // Countdown timer
        let timeLeft = 10;
        const timerElement = document.getElementById('timer');
        
        const countdown = setInterval(function() {
            timeLeft--;
            timerElement.textContent = timeLeft;
            
            if (timeLeft <= 0) {
                clearInterval(countdown);
                window.close();
            }
        }, 1000);
        
        // Also close after 10 seconds as fallback
        setTimeout(function() {
            window.close();
        }, 10 * 1000);
    </script>
</body>
</html>`
}
