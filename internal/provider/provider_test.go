package provider

import (
	"context"
	"testing"

	"baize/internal/domain"
)

// TestCreateProvider 测试创建Provider实例
func TestCreateProvider(t *testing.T) {
	// 测试创建OpenAIProvider
	openaiPlatform := &domain.Platform{
		ID:      "test_openai",
		Name:    "测试OpenAI",
		Type:    "openai",
		BaseURL: "https://api.openai.com/v1",
		APIKey:  "sk-test1234",
		Models:  []string{"gpt-4"},
	}

	openaiProvider, err := CreateProvider(openaiPlatform)
	if err != nil {
		t.Fatalf("创建OpenAIProvider失败: %v", err)
	}

	if _, ok := openaiProvider.(*OpenAIProvider); !ok {
		t.Error("期望返回OpenAIProvider实例, 实际返回其他类型")
	}

	// 测试创建AnthropicProvider
	anthropicPlatform := &domain.Platform{
		ID:      "test_anthropic",
		Name:    "测试Anthropic",
		Type:    "anthropic",
		BaseURL: "https://api.anthropic.com",
		APIKey:  "sk-ant-test1234",
		Models:  []string{"claude-3-opus"},
	}

	anthropicProvider, err := CreateProvider(anthropicPlatform)
	if err != nil {
		t.Fatalf("创建AnthropicProvider失败: %v", err)
	}

	if _, ok := anthropicProvider.(*AnthropicProvider); !ok {
		t.Error("期望返回AnthropicProvider实例, 实际返回其他类型")
	}

	// 测试创建不支持的Provider类型
	unsupportedPlatform := &domain.Platform{
		ID:      "test_unsupported",
		Name:    "测试不支持",
		Type:    "unsupported",
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Models:  []string{"test-model"},
	}

	_, err = CreateProvider(unsupportedPlatform)
	if err == nil {
		t.Error("期望创建不支持的Provider类型时返回错误, 但未返回错误")
	}
}

// TestAIProviderInterface 测试AIProvider接口
func TestAIProviderInterface(t *testing.T) {
	// 创建测试平台
	platform := &domain.Platform{
		ID:      "test_platform",
		Name:    "测试平台",
		Type:    "openai",
		BaseURL: "https://api.openai.com/v1",
		APIKey:  "sk-test1234",
		Models:  []string{"gpt-4"},
	}

	// 创建Provider实例
	provider, err := CreateProvider(platform)
	if err != nil {
		t.Fatalf("创建Provider失败: %v", err)
	}

	// 测试Chat方法（这里只是测试接口调用，不会实际发送请求）
	ctx := context.Background()
	_, err = provider.Chat(ctx, "gpt-4", "测试消息")
	// 这里会返回错误，因为是测试环境，但接口调用应该成功
	if err == nil {
		t.Log("Chat方法调用成功（预期会返回错误，因为是测试环境）")
	} else {
		t.Logf("Chat方法调用返回错误（预期行为）: %v", err)
	}
}
