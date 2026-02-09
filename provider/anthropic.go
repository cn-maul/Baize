package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"baize/domain"
)

// AnthropicProvider Anthropic提供商实现
type AnthropicProvider struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewAnthropicProvider 创建新的AnthropicProvider实例
func NewAnthropicProvider(platform *domain.Platform) (AIProvider, error) {
	return &AnthropicProvider{
		baseURL: platform.BaseURL,
		apiKey:  platform.APIKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// AnthropicRequest Anthropic API请求结构
type AnthropicRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

// AnthropicResponse Anthropic API响应结构
type AnthropicResponse struct {
	Content []ContentBlock `json:"content"`
	Error   *Error         `json:"error,omitempty"`
}

// ContentBlock 内容块结构
type ContentBlock struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
}

// Chat 实现AIProvider接口的Chat方法
func (p *AnthropicProvider) Chat(ctx context.Context, model string, msg string) (string, error) {
	// 构建请求体
	requestBody := AnthropicRequest{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: msg},
		},
		MaxTokens: 1000,
	}

	// 序列化请求体
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/v1/messages", bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response AnthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查错误
	if response.Error != nil {
		return "", fmt.Errorf("API错误: %s", response.Error.Message)
	}

	// 检查响应
	if len(response.Content) == 0 {
		return "", fmt.Errorf("响应中没有内容")
	}

	// 提取文本内容
	var reply string
	for _, block := range response.Content {
		if block.Type == "text" {
			reply += block.Text
		}
	}

	return reply, nil
}
