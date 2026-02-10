package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cn-maul/Baize/domain"
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
	Stream    bool      `json:"stream,omitempty"`
}

// AnthropicResponse Anthropic API响应结构
type AnthropicResponse struct {
	Content []ContentBlock `json:"content"`
	Error   *Error         `json:"error,omitempty"`
}

// AnthropicStreamResponse Anthropic API流式响应结构
type AnthropicStreamResponse struct {
	Type    string         `json:"type"`
	Message *StreamMessage `json:"message,omitempty"`
	Error   *Error         `json:"error,omitempty"`
}

// StreamMessage 流式消息结构
type StreamMessage struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Role       string         `json:"role"`
	Content    []ContentBlock `json:"content"`
	StopReason string         `json:"stop_reason,omitempty"`
}

// ContentBlock 内容块结构
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
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

// ChatWithContext 实现AIProvider接口的ChatWithContext方法
func (p *AnthropicProvider) ChatWithContext(ctx context.Context, model string, messages []Message) (string, error) {
	// 构建请求体
	requestBody := AnthropicRequest{
		Model:     model,
		Messages:  messages,
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

// ChatStream 实现AIProvider接口的ChatStream方法
func (p *AnthropicProvider) ChatStream(ctx context.Context, model string, msg string, callback func(chunk string) error) error {
	// 构建请求体
	requestBody := AnthropicRequest{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: msg},
		},
		MaxTokens: 1000,
		Stream:    true,
	}

	// 序列化请求体
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/v1/messages", bytes.NewBuffer(requestJSON))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 处理流式响应
	return p.handleStreamResponse(resp.Body, callback)
}

// ChatStreamWithContext 实现AIProvider接口的ChatStreamWithContext方法
func (p *AnthropicProvider) ChatStreamWithContext(ctx context.Context, model string, messages []Message, callback func(chunk string) error) error {
	// 构建请求体
	requestBody := AnthropicRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: 1000,
		Stream:    true,
	}

	// 序列化请求体
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/v1/messages", bytes.NewBuffer(requestJSON))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 处理流式响应
	return p.handleStreamResponse(resp.Body, callback)
}

// handleStreamResponse 处理流式响应
func (p *AnthropicProvider) handleStreamResponse(body io.Reader, callback func(chunk string) error) error {
	// 创建一个扫描器来逐行读取响应
	scanner := NewLineScanner(body)
	for scanner.Scan() {
		line := scanner.Text()
		// 跳过空行
		if line == "" {
			continue
		}
		// 跳过SSE注释
		if strings.HasPrefix(line, ":") {
			continue
		}
		// 检查是否是结束信号
		if line == "data: [DONE]" {
			break
		}
		// 提取JSON数据
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			// 解析JSON
			var response AnthropicStreamResponse
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				return fmt.Errorf("解析流式响应失败: %w", err)
			}
			// 检查错误
			if response.Error != nil {
				return fmt.Errorf("API错误: %s", response.Error.Message)
			}
			// 处理响应
			if response.Type == "content_block_delta" && response.Message != nil {
				for _, block := range response.Message.Content {
					if block.Type == "text" && block.Text != "" {
						// 调用回调函数
						if err := callback(block.Text); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	// 检查扫描器错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取流式响应失败: %w", err)
	}

	return nil
}
