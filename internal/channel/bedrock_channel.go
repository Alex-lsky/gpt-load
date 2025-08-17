package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"gpt-load/internal/models"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	Register("bedrock", newBedrockChannel)
}

// AWSCredentials represents the AWS credentials stored in APIKey.KeyValue
type AWSCredentials struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token,omitempty"`
	Region          string `json:"region"`
}

// FoundationModel represents a Bedrock foundation model
type FoundationModel struct {
	ModelID      string    `json:"model_id"`
	ModelName    string    `json:"model_name"`
	ProviderName string    `json:"provider_name"`
	CachedAt     time.Time `json:"cached_at"`
}

// ModelCache manages the cache of available foundation models
type ModelCache struct {
	models    map[string][]FoundationModel
	cacheTTL  time.Duration
	cacheLock sync.RWMutex
}

// BedrockChannel implements the ChannelProxy interface for AWS Bedrock
type BedrockChannel struct {
	*BaseChannel
	modelCache *ModelCache
}

// newBedrockChannel creates a new Bedrock channel instance
func newBedrockChannel(f *Factory, group *models.Group) (ChannelProxy, error) {
	base, err := f.newBaseChannel("bedrock", group)
	if err != nil {
		return nil, err
	}

	// Initialize model cache with 1 hour TTL
	modelCache := &ModelCache{
		models:   make(map[string][]FoundationModel),
		cacheTTL: time.Hour,
	}

	return &BedrockChannel{
		BaseChannel: base,
		modelCache:  modelCache,
	}, nil
}

// ModifyRequest sets the required headers and authentication for AWS Bedrock API
func (ch *BedrockChannel) ModifyRequest(req *http.Request, apiKey *models.APIKey, group *models.Group) {
	// This will be implemented in task 3.1
	// For now, just set basic headers
	req.Header.Set("Content-Type", "application/json")
}

// IsStreamRequest checks if the request is for a streaming response
func (ch *BedrockChannel) IsStreamRequest(c *gin.Context, bodyBytes []byte) bool {
	// Check for streaming endpoints
	path := c.Request.URL.Path
	if strings.Contains(path, "/converse-stream") || strings.Contains(path, "/invoke-with-response-stream") {
		return true
	}

	// Check Accept header
	if strings.Contains(c.GetHeader("Accept"), "text/event-stream") {
		return true
	}

	// Check for stream parameter in request body
	type streamPayload struct {
		Stream bool `json:"stream"`
	}
	var p streamPayload
	if err := json.Unmarshal(bodyBytes, &p); err == nil {
		return p.Stream
	}

	return false
}

// ValidateKey checks if the given API key is valid by calling ListFoundationModels
func (ch *BedrockChannel) ValidateKey(ctx context.Context, key string) (bool, error) {
	// Parse AWS credentials from the key
	var creds AWSCredentials
	if err := json.Unmarshal([]byte(key), &creds); err != nil {
		return false, fmt.Errorf("failed to parse AWS credentials: %w", err)
	}

	// Validate required fields
	if creds.AccessKeyID == "" || creds.SecretAccessKey == "" || creds.Region == "" {
		return false, fmt.Errorf("missing required AWS credentials fields")
	}

	// This will be fully implemented in task 4.3
	// For now, just return true if credentials are properly formatted
	logrus.Debugf("Validating AWS credentials for region: %s", creds.Region)
	return true, nil
}

// BuildUpstreamURL constructs the target URL for the AWS Bedrock service
func (ch *BedrockChannel) BuildUpstreamURL(originalURL *url.URL, group *models.Group) (string, error) {
	// Extract the region from the first upstream URL or use default
	base := ch.getUpstreamURL()
	if base == nil {
		return "", fmt.Errorf("no upstream URL configured for channel %s", ch.Name)
	}

	// Parse the original request path
	proxyPrefix := "/proxy/" + group.Name
	requestPath := originalURL.Path
	requestPath = strings.TrimPrefix(requestPath, proxyPrefix)

	// Map the request path to Bedrock API endpoints
	var finalPath string
	switch {
	case strings.HasPrefix(requestPath, "/converse"):
		finalPath = requestPath // Keep as is: /converse or /converse-stream
	case strings.Contains(requestPath, "/model/") && strings.Contains(requestPath, "/invoke"):
		finalPath = requestPath // Keep as is: /model/{model-id}/invoke or /model/{model-id}/invoke-with-response-stream
	default:
		return "", fmt.Errorf("unsupported Bedrock API path: %s", requestPath)
	}

	// Construct the final URL
	finalURL := *base
	finalURL.Path = strings.TrimRight(finalURL.Path, "/") + finalPath
	finalURL.RawQuery = originalURL.RawQuery

	return finalURL.String(), nil
}

// parseAWSCredentials parses AWS credentials from the API key value
func parseAWSCredentials(keyValue string) (*AWSCredentials, error) {
	var creds AWSCredentials
	if err := json.Unmarshal([]byte(keyValue), &creds); err != nil {
		return nil, fmt.Errorf("failed to parse AWS credentials: %w", err)
	}

	// Validate required fields
	if creds.AccessKeyID == "" {
		return nil, fmt.Errorf("access_key_id is required")
	}
	if creds.SecretAccessKey == "" {
		return nil, fmt.Errorf("secret_access_key is required")
	}
	if creds.Region == "" {
		return nil, fmt.Errorf("region is required")
	}

	return &creds, nil
}

// createAWSConfig creates an AWS config from credentials
func createAWSConfig(creds *AWSCredentials) aws.Config {
	var credProvider aws.CredentialsProvider
	if creds.SessionToken != "" {
		credProvider = credentials.NewStaticCredentialsProvider(
			creds.AccessKeyID,
			creds.SecretAccessKey,
			creds.SessionToken,
		)
	} else {
		credProvider = credentials.NewStaticCredentialsProvider(
			creds.AccessKeyID,
			creds.SecretAccessKey,
			"",
		)
	}

	return aws.Config{
		Region:      creds.Region,
		Credentials: credProvider,
	}
}
