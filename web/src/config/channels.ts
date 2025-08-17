// 渠道配置接口
export interface ChannelConfig {
  keyInputComponent: string;
  defaultUpstream: string;
  defaultTestModel: string;
  validationEndpoint: string;
  showValidationEndpoint: boolean;
  upstreamPlaceholder: string;
  testModelPlaceholder: string;
  regions?: Array<{ label: string; value: string }>;
}

// 渠道配置映射
export const CHANNEL_CONFIGS: Record<string, ChannelConfig> = {
  openai: {
    keyInputComponent: "StandardKeyInput",
    defaultUpstream: "https://api.openai.com",
    defaultTestModel: "gpt-4o-mini",
    validationEndpoint: "/v1/chat/completions",
    showValidationEndpoint: true,
    upstreamPlaceholder: "https://api.openai.com",
    testModelPlaceholder: "gpt-4o-mini",
  },
  anthropic: {
    keyInputComponent: "StandardKeyInput",
    defaultUpstream: "https://api.anthropic.com",
    defaultTestModel: "claude-3-haiku-20240307",
    validationEndpoint: "/v1/messages",
    showValidationEndpoint: true,
    upstreamPlaceholder: "https://api.anthropic.com",
    testModelPlaceholder: "claude-3-haiku-20240307",
  },
  gemini: {
    keyInputComponent: "StandardKeyInput",
    defaultUpstream: "https://generativelanguage.googleapis.com",
    defaultTestModel: "gemini-1.5-flash",
    validationEndpoint: "",
    showValidationEndpoint: false,
    upstreamPlaceholder: "https://generativelanguage.googleapis.com",
    testModelPlaceholder: "gemini-1.5-flash",
  },
  bedrock: {
    keyInputComponent: "BedrockKeyInput",
    defaultUpstream: "https://bedrock-runtime.us-east-1.amazonaws.com",
    defaultTestModel: "anthropic.claude-3-haiku-20240307-v1:0",
    validationEndpoint: "",
    showValidationEndpoint: false,
    upstreamPlaceholder: "https://bedrock-runtime.{region}.amazonaws.com",
    testModelPlaceholder: "anthropic.claude-3-haiku-20240307-v1:0",
    regions: [
      { label: "US East (N. Virginia)", value: "us-east-1" },
      { label: "US West (Oregon)", value: "us-west-2" },
      { label: "Europe (Ireland)", value: "eu-west-1" },
      { label: "Asia Pacific (Tokyo)", value: "ap-northeast-1" },
      { label: "Asia Pacific (Singapore)", value: "ap-southeast-1" },
      { label: "Asia Pacific (Sydney)", value: "ap-southeast-2" },
      { label: "Canada (Central)", value: "ca-central-1" },
      { label: "Europe (Frankfurt)", value: "eu-central-1" },
      { label: "Europe (London)", value: "eu-west-2" },
      { label: "Asia Pacific (Mumbai)", value: "ap-south-1" },
    ],
  },
};
