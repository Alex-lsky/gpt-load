package channel

// Bedrock Converse API request and response types

// ConverseRequest represents the request structure for Bedrock Converse API
type ConverseRequest struct {
	ModelID                      string                         `json:"modelId"`
	Messages                     []ConverseMessage              `json:"messages"`
	System                       []SystemMessage                `json:"system,omitempty"`
	InferenceConfig              *InferenceConfig               `json:"inferenceConfig,omitempty"`
	ToolConfig                   *ToolConfig                    `json:"toolConfig,omitempty"`
	GuardrailConfig              *GuardrailConfig               `json:"guardrailConfig,omitempty"`
	AdditionalModelRequestFields map[string]interface{}         `json:"additionalModelRequestFields,omitempty"`
}

// ConverseMessage represents a message in the conversation
type ConverseMessage struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// ContentBlock represents different types of content in a message
type ContentBlock struct {
	Text     *string                `json:"text,omitempty"`
	Image    *ImageBlock            `json:"image,omitempty"`
	Document *DocumentBlock         `json:"document,omitempty"`
	ToolUse  *ToolUseBlock          `json:"toolUse,omitempty"`
	ToolResult *ToolResultBlock     `json:"toolResult,omitempty"`
	GuardContent *GuardContentBlock `json:"guardContent,omitempty"`
}

// ImageBlock represents image content
type ImageBlock struct {
	Format string `json:"format"`
	Source *ImageSource `json:"source"`
}

// ImageSource represents the source of an image
type ImageSource struct {
	Bytes []byte `json:"bytes,omitempty"`
}

// DocumentBlock represents document content
type DocumentBlock struct {
	Format string `json:"format"`
	Name   string `json:"name"`
	Source *DocumentSource `json:"source"`
}

// DocumentSource represents the source of a document
type DocumentSource struct {
	Bytes []byte `json:"bytes,omitempty"`
}

// ToolUseBlock represents tool usage
type ToolUseBlock struct {
	ToolUseId string                 `json:"toolUseId"`
	Name      string                 `json:"name"`
	Input     map[string]interface{} `json:"input"`
}

// ToolResultBlock represents tool result
type ToolResultBlock struct {
	ToolUseId string         `json:"toolUseId"`
	Content   []ContentBlock `json:"content"`
	Status    string         `json:"status,omitempty"`
}

// GuardContentBlock represents guard content
type GuardContentBlock struct {
	Text *GuardContentText `json:"text,omitempty"`
}

// GuardContentText represents guard content text
type GuardContentText struct {
	Text       string                 `json:"text"`
	Qualifiers []string               `json:"qualifiers,omitempty"`
}

// SystemMessage represents a system message
type SystemMessage struct {
	Text         *string                `json:"text,omitempty"`
	GuardContent *GuardContentBlock     `json:"guardContent,omitempty"`
}

// InferenceConfig represents inference configuration
type InferenceConfig struct {
	MaxTokens     *int32   `json:"maxTokens,omitempty"`
	Temperature   *float32 `json:"temperature,omitempty"`
	TopP          *float32 `json:"topP,omitempty"`
	StopSequences []string `json:"stopSequences,omitempty"`
}

// ToolConfig represents tool configuration
type ToolConfig struct {
	Tools        []Tool        `json:"tools,omitempty"`
	ToolChoice   *ToolChoice   `json:"toolChoice,omitempty"`
}

// Tool represents a tool definition
type Tool struct {
	ToolSpec *ToolSpecification `json:"toolSpec,omitempty"`
}

// ToolSpecification represents tool specification
type ToolSpecification struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	InputSchema  map[string]interface{} `json:"inputSchema"`
}

// ToolChoice represents tool choice configuration
type ToolChoice struct {
	Auto *AutoToolChoice `json:"auto,omitempty"`
	Any  *AnyToolChoice  `json:"any,omitempty"`
	Tool *SpecificToolChoice `json:"tool,omitempty"`
}

// AutoToolChoice represents automatic tool choice
type AutoToolChoice struct{}

// AnyToolChoice represents any tool choice
type AnyToolChoice struct{}

// SpecificToolChoice represents specific tool choice
type SpecificToolChoice struct {
	Name string `json:"name"`
}

// GuardrailConfig represents guardrail configuration
type GuardrailConfig struct {
	GuardrailIdentifier string `json:"guardrailIdentifier"`
	GuardrailVersion    string `json:"guardrailVersion"`
}

// ConverseResponse represents the response from Bedrock Converse API
type ConverseResponse struct {
	Output      ConverseOutput `json:"output"`
	StopReason  string         `json:"stopReason"`
	Usage       Usage          `json:"usage"`
	Metrics     Metrics        `json:"metrics"`
}

// ConverseOutput represents the output of the conversation
type ConverseOutput struct {
	Message *ConverseMessage `json:"message,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	InputTokens  int32 `json:"inputTokens"`
	OutputTokens int32 `json:"outputTokens"`
	TotalTokens  int32 `json:"totalTokens"`
}

// Metrics represents performance metrics
type Metrics struct {
	LatencyMs int64 `json:"latencyMs"`
}

// ConverseStreamResponse represents streaming response from Bedrock
type ConverseStreamResponse struct {
	MessageStart      *MessageStart      `json:"messageStart,omitempty"`
	ContentBlockStart *ContentBlockStart `json:"contentBlockStart,omitempty"`
	ContentBlockDelta *ContentBlockDelta `json:"contentBlockDelta,omitempty"`
	ContentBlockStop  *ContentBlockStop  `json:"contentBlockStop,omitempty"`
	MessageStop       *MessageStop       `json:"messageStop,omitempty"`
	Metadata          *StreamMetadata    `json:"metadata,omitempty"`
}

// MessageStart represents the start of a message in streaming
type MessageStart struct {
	Role string `json:"role"`
}

// ContentBlockStart represents the start of a content block in streaming
type ContentBlockStart struct {
	Start       *ContentBlockStartEvent `json:"start,omitempty"`
	ContentBlockIndex int32             `json:"contentBlockIndex"`
}

// ContentBlockStartEvent represents the start event of a content block
type ContentBlockStartEvent struct {
	ToolUse *ToolUseBlock `json:"toolUse,omitempty"`
}

// ContentBlockDelta represents a delta in content block during streaming
type ContentBlockDelta struct {
	Delta             *ContentBlockDeltaEvent `json:"delta,omitempty"`
	ContentBlockIndex int32                   `json:"contentBlockIndex"`
}

// ContentBlockDeltaEvent represents the delta event of a content block
type ContentBlockDeltaEvent struct {
	Text    *string                `json:"text,omitempty"`
	ToolUse *ToolUseDelta          `json:"toolUse,omitempty"`
}

// ToolUseDelta represents a delta in tool use
type ToolUseDelta struct {
	Input string `json:"input,omitempty"`
}

// ContentBlockStop represents the stop of a content block in streaming
type ContentBlockStop struct {
	ContentBlockIndex int32 `json:"contentBlockIndex"`
}

// MessageStop represents the stop of a message in streaming
type MessageStop struct {
	StopReason        string `json:"stopReason"`
	AdditionalModelResponseFields map[string]interface{} `json:"additionalModelResponseFields,omitempty"`
}

// StreamMetadata represents metadata in streaming response
type StreamMetadata struct {
	Usage   Usage   `json:"usage"`
	Metrics Metrics `json:"metrics"`
}

// OpenAI-compatible request/response types for conversion

// OpenAIRequest represents a standard OpenAI API request
type OpenAIRequest struct {
	Model            string                   `json:"model"`
	Messages         []OpenAIMessage          `json:"messages"`
	MaxTokens        *int32                   `json:"max_tokens,omitempty"`
	Temperature      *float32                 `json:"temperature,omitempty"`
	TopP             *float32                 `json:"top_p,omitempty"`
	Stop             interface{}              `json:"stop,omitempty"`
	Stream           bool                     `json:"stream,omitempty"`
	Tools            []OpenAITool             `json:"tools,omitempty"`
	ToolChoice       interface{}              `json:"tool_choice,omitempty"`
}

// OpenAIMessage represents a message in OpenAI format
type OpenAIMessage struct {
	Role         string                   `json:"role"`
	Content      interface{}              `json:"content"`
	Name         string                   `json:"name,omitempty"`
	ToolCalls    []OpenAIToolCall         `json:"tool_calls,omitempty"`
	ToolCallId   string                   `json:"tool_call_id,omitempty"`
}

// OpenAITool represents a tool in OpenAI format
type OpenAITool struct {
	Type     string                 `json:"type"`
	Function OpenAIFunctionTool     `json:"function"`
}

// OpenAIFunctionTool represents a function tool in OpenAI format
type OpenAIFunctionTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// OpenAIToolCall represents a tool call in OpenAI format
type OpenAIToolCall struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Function OpenAIFunctionCall     `json:"function"`
}

// OpenAIFunctionCall represents a function call in OpenAI format
type OpenAIFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// OpenAIResponse represents a standard OpenAI API response
type OpenAIResponse struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	Choices           []OpenAIChoice         `json:"choices"`
	Usage             OpenAIUsage            `json:"usage"`
	SystemFingerprint string                 `json:"system_fingerprint,omitempty"`
}

// OpenAIChoice represents a choice in OpenAI response
type OpenAIChoice struct {
	Index        int32                  `json:"index"`
	Message      *OpenAIMessage         `json:"message,omitempty"`
	Delta        *OpenAIMessage         `json:"delta,omitempty"`
	FinishReason *string                `json:"finish_reason"`
}

// OpenAIUsage represents usage information in OpenAI format
type OpenAIUsage struct {
	PromptTokens     int32 `json:"prompt_tokens"`
	CompletionTokens int32 `json:"completion_tokens"`
	TotalTokens      int32 `json:"total_tokens"`
}
