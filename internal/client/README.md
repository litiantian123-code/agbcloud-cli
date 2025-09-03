# AgbCloud HTTP Client

This package provides a structured HTTP client for interacting with the AgbCloud API, inspired by the architecture of the Daytona API client.

## Architecture

The client is built with the following components:

### Core Components

- **APIClient**: Main client struct that manages HTTP communication
- **Configuration**: Client configuration including server URLs, authentication, and HTTP settings
- **Service Pattern**: Each API group is implemented as a separate service (e.g., OAuthAPI)

### Files Structure

```
internal/client/
├── client.go          # Main APIClient implementation
├── configuration.go   # Configuration structures and methods
├── oauth_api.go      # OAuth API service
├── factory.go        # Factory functions for easy client creation
└── README.md         # This documentation
```

## Usage

### Basic Usage

```go
import "github.com/liyuebing/agbcloud-cli/internal/client"

// Create client with default configuration
client := client.NewDefault()

// Make API calls
ctx := context.Background()
response, httpResp, err := client.OAuthAPI.GetGoogleLoginURL(ctx, "https://agb.cloud")
```

### Custom Configuration

```go
import (
    "github.com/liyuebing/agbcloud-cli/internal/client"
    "github.com/liyuebing/agbcloud-cli/internal/config"
)

// Create client from CLI config
cfg := &config.Config{
    APIKey:   "your-api-key",
    Endpoint: "https://agb.cloud",
}
client := client.NewFromConfig(cfg)
```

### Advanced Configuration

```go
// Create custom configuration
cfg := client.NewConfiguration()
cfg.Debug = true
cfg.Servers[0].URL = "https://custom-api.example.com"
cfg.AddDefaultHeader("Custom-Header", "value")

// Create client
apiClient := client.NewAPIClient(cfg)
```

## API Services

### OAuthAPI

Currently implements:
- `GetGoogleLoginURL(ctx context.Context, fromUrlPath string) (OAuthGoogleLoginResponse, *http.Response, error)`

## Authentication

The client supports multiple authentication methods:

1. **API Key**: Set via configuration or context
2. **Bearer Token**: Set via context using `ContextAccessToken`

```go
// Using context for authentication (for APIs that require auth)
ctx := context.WithValue(context.Background(), client.ContextAccessToken, "your-token")
// Note: OAuth API doesn't require authentication
response, _, err := client.OAuthAPI.GetGoogleLoginURL(ctx, "https://agb.cloud")
```

## Error Handling

The client uses `GenericOpenAPIError` for structured error handling:

```go
response, httpResp, err := client.OAuthAPI.GetGoogleLoginURL(ctx, "https://agb.cloud")
if err != nil {
    if apiErr, ok := err.(*client.GenericOpenAPIError); ok {
        fmt.Printf("API Error: %s\n", apiErr.Error())
        fmt.Printf("Response Body: %s\n", string(apiErr.Body()))
        fmt.Printf("HTTP Status: %d\n", httpResp.StatusCode)
    }
}
```

## Testing

Run the tests with:

```bash
go test ./internal/client/... -v
```

The test suite includes:
- Unit tests for client creation
- Mock server tests for API calls
- Error handling tests

## Configuration Options

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| `Host` | string | Override server host | "" |
| `Scheme` | string | Override server scheme | "" |
| `UserAgent` | string | HTTP User-Agent header | "AgbCloud-CLI/1.0.0/go" |
| `Debug` | bool | Enable debug logging | false |
| `HTTPClient` | *http.Client | Custom HTTP client | 30s timeout |
| `APIKey` | string | API key for authentication | "" |

## Adding New API Services

To add a new API service:

1. Create a new interface and service struct:
```go
type NewAPI interface {
    SomeMethod(ctx context.Context) (Response, *http.Response, error)
}

type NewAPIService service
```

2. Add the service to APIClient:
```go
type APIClient struct {
    // ... existing fields
    NewAPI NewAPI
}
```

3. Initialize in NewAPIClient:
```go
c.NewAPI = (*NewAPIService)(&c.common)
```

4. Implement the methods following the existing pattern in oauth_api.go 