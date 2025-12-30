package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"gpt-load/internal/models"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	Register("bedrock", newBedrockChannel)
}

// BedrockChannel implements the ChannelProxy interface for AWS Bedrock
type BedrockChannel struct {
	*BaseChannel
	authType      string // "api_key" or "iam"
	region        string
	apiKey        string // API Key authentication
	accessKeyID   string // IAM authentication
	secretKey     string // IAM authentication
	sessionToken  string // IAM authentication (optional)
	authenticator BedrockAuthenticator
}

// BedrockAuthenticator defines the interface for AWS authentication
type BedrockAuthenticator interface {
	SignRequest(req *http.Request, region string) error
	ValidateCredentials(ctx context.Context, region string) error
}

// newBedrockChannel creates a new Bedrock channel instance
func newBedrockChannel(f *Factory, group *models.Group) (ChannelProxy, error) {
	base, err := f.newBaseChannel("bedrock", group)
	if err != nil {
		return nil, err
	}

	// Parse Bedrock-specific configuration from group config
	config, err := parseBedrockConfig(group)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Bedrock config: %w", err)
	}

	// Create authenticator based on auth type
	var authenticator BedrockAuthenticator
	switch config.AuthType {
	case "api_key":
		authenticator = &APIKeyAuth{
			apiKey: config.APIKey,
		}
	case "iam":
		authenticator = &IAMAuth{
			accessKeyID:  config.AccessKeyID,
			secretKey:    config.SecretAccessKey,
			sessionToken: config.SessionToken,
		}
	default:
		return nil, fmt.Errorf("unsupported auth type: %s", config.AuthType)
	}

	return &BedrockChannel{
		BaseChannel:   base,
		authType:      config.AuthType,
		region:        config.Region,
		apiKey:        config.APIKey,
		accessKeyID:   config.AccessKeyID,
		secretKey:     config.SecretAccessKey,
		sessionToken:  config.SessionToken,
		authenticator: authenticator,
	}, nil
}

// BedrockConfig represents the Bedrock-specific configuration
type BedrockConfig struct {
	AuthType        string `json:"auth_type"`
	Region          string `json:"region"`
	APIKey          string `json:"api_key,omitempty"`
	AccessKeyID     string `json:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
	SessionToken    string `json:"session_token,omitempty"`
}

// parseBedrockConfig extracts Bedrock configuration from group config
func parseBedrockConfig(group *models.Group) (*BedrockConfig, error) {
	if group.Config == nil {
		return nil, fmt.Errorf("group config is nil")
	}

	var config BedrockConfig
	// Convert datatypes.JSONMap to JSON bytes
	configBytes, err := json.Marshal(group.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group config to JSON: %w", err)
	}

	// Unmarshal JSON bytes into BedrockConfig struct
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Bedrock config: %w", err)
	}

	// Validate required fields
	if config.AuthType == "" {
		return nil, fmt.Errorf("auth_type is required")
	}
	if config.Region == "" {
		return nil, fmt.Errorf("region is required")
	}

	switch config.AuthType {
	case "api_key":
		if config.APIKey == "" {
			return nil, fmt.Errorf("api_key is required for API key authentication")
		}
	case "iam":
		if config.AccessKeyID == "" || config.SecretAccessKey == "" {
			return nil, fmt.Errorf("access_key_id and secret_access_key are required for IAM authentication")
		}
	default:
		return nil, fmt.Errorf("unsupported auth_type: %s", config.AuthType)
	}

	return &config, nil
}

// BuildUpstreamURL constructs the Bedrock Converse API URL
func (ch *BedrockChannel) BuildUpstreamURL(originalURL *url.URL, groupName string) (string, error) {
	// Extract model from the request path or use a default approach
	proxyPrefix := "/proxy/" + groupName
	requestPath := originalURL.Path
	requestPath = strings.TrimPrefix(requestPath, proxyPrefix)

	// For Bedrock, we need to construct the Converse API endpoint
	// The URL format is: https://bedrock-runtime.{region}.amazonaws.com/model/{model-id}/converse
	baseURL := fmt.Sprintf("https://bedrock-runtime.%s.amazonaws.com", ch.region)

	// If the request path contains a model, extract it
	if strings.Contains(requestPath, "/v1/chat/completions") {
		// This will be handled in ModifyRequest where we have access to the request body
		return baseURL + "/model/PLACEHOLDER/converse", nil
	}

	return baseURL + requestPath, nil
}

// ModifyRequest adds AWS authentication headers and modifies the request for Bedrock
func (ch *BedrockChannel) ModifyRequest(req *http.Request, apiKey *models.APIKey, group *models.Group) {
	// The actual model will be extracted from the request body and URL will be updated
	// This is a placeholder implementation - the full implementation will be in task 3

	// Add required headers for Bedrock
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Authentication will be handled by the authenticator
	if ch.authenticator != nil {
		// Sign the request using the appropriate authentication method
		if err := ch.authenticator.SignRequest(req, ch.region); err != nil {
			// Log error but don't fail the request here
			// Error handling will be improved in later tasks
		}
	}
}

// IsStreamRequest checks if the request is for streaming response
func (ch *BedrockChannel) IsStreamRequest(c *gin.Context, bodyBytes []byte) bool {
	// Check Accept header for streaming
	if strings.Contains(c.GetHeader("Accept"), "text/event-stream") {
		return true
	}

	// Check query parameter
	if c.Query("stream") == "true" {
		return true
	}

	// Check request body for stream parameter
	type streamPayload struct {
		Stream bool `json:"stream"`
	}
	var p streamPayload
	if err := json.Unmarshal(bodyBytes, &p); err == nil {
		return p.Stream
	}

	return false
}

// ExtractModel extracts the model name from the request
func (ch *BedrockChannel) ExtractModel(c *gin.Context, bodyBytes []byte) string {
	type modelPayload struct {
		Model string `json:"model"`
	}
	var p modelPayload
	if err := json.Unmarshal(bodyBytes, &p); err == nil {
		return p.Model
	}
	return ""
}

// ValidateKey validates the Bedrock API credentials
func (ch *BedrockChannel) ValidateKey(ctx context.Context, apiKey *models.APIKey, group *models.Group) (bool, error) {
	// Use the authenticator to validate credentials
	if ch.authenticator == nil {
		return false, fmt.Errorf("no authenticator configured")
	}

	err := ch.authenticator.ValidateCredentials(ctx, ch.region)
	if err != nil {
		return false, fmt.Errorf("credential validation failed: %w", err)
	}

	return true, nil
}
