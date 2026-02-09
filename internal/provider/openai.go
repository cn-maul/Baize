package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"baize/internal/domain"
)

// OpenAIProvider OpenAI提供商实现
type OpenAIProvider struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewOpenAIProvider 创建新的OpenAIProvider实例
func NewOpenAIProvider(platform *domain.Platform) (AIProvider, error) {
	return &OpenAIProvider{
		baseURL: platform.BaseURL,
		apiKey:  platform.APIKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// OpenAIRequest OpenAI API请求结构
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse OpenAI API响应结构
type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
	Error   *Error   `json:"error,omitempty"`
}

// Choice 选择结构
type Choice struct {
	Message Message `json:"message"`
}

// Error 错误结构
type Error struct {
	Message string `json:"message"`
}

// Chat 实现AIProvider接口的Chat方法
func (p *OpenAIProvider) Chat(ctx context.Context, model string, msg string) (string, error) {
	// 构建请求体
	requestBody := OpenAIRequest{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: msg},
		},
	}

	// 序列化请求体
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查错误
	if response.Error != nil {
		return "", fmt.Errorf("API错误: %s", response.Error.Message)
	}

	// 检查响应
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("响应中没有选择")
	}

	return response.Choices[0].Message.Content, nil
}
