package provider

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cn-maul/Baize/domain"
)

// OpenAIProvider OpenAI提供商实现
type OpenAIProvider struct {
	*BaseProvider
}

// NewOpenAIProvider 创建新的OpenAIProvider实例
func NewOpenAIProvider(platform *domain.Platform, options ...ProviderOption) (AIProvider, error) {
	return &OpenAIProvider{
		BaseProvider: NewBaseProvider(platform.BaseURL, platform.APIKey, options...),
	}, nil
}

// OpenAIRequest OpenAI API请求结构
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
}

// OpenAIResponse OpenAI API响应结构
type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
	Error   *Error   `json:"error,omitempty"`
}

// OpenAIStreamResponse OpenAI API流式响应结构
type OpenAIStreamResponse struct {
	Choices []StreamChoice `json:"choices"`
	Error   *Error         `json:"error,omitempty"`
}

// StreamChoice 流式选择结构
type StreamChoice struct {
	Delta        Message `json:"delta"`
	FinishReason string  `json:"finish_reason"`
}

// Choice 选择结构
type Choice struct {
	Message Message `json:"message"`
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

	// 设置请求头
	headers := map[string]string{
		"Authorization": "Bearer " + p.apiKey,
	}

	// 发送请求
	resp, err := p.sendRequest(ctx, "POST", "/chat/completions", requestBody, headers)
	if err != nil {
		return "", err
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

// ChatWithContext 实现AIProvider接口的ChatWithContext方法
func (p *OpenAIProvider) ChatWithContext(ctx context.Context, model string, messages []Message) (string, error) {
	// 构建请求体
	requestBody := OpenAIRequest{
		Model:    model,
		Messages: messages,
	}

	// 设置请求头
	headers := map[string]string{
		"Authorization": "Bearer " + p.apiKey,
	}

	// 发送请求
	resp, err := p.sendRequest(ctx, "POST", "/chat/completions", requestBody, headers)
	if err != nil {
		return "", err
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

// ChatStream 实现AIProvider接口的ChatStream方法
func (p *OpenAIProvider) ChatStream(ctx context.Context, model string, msg string, callback func(chunk string) error) error {
	// 构建请求体
	requestBody := OpenAIRequest{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: msg},
		},
		Stream: true,
	}

	// 设置请求头
	headers := map[string]string{
		"Authorization": "Bearer " + p.apiKey,
	}

	// 发送请求
	resp, err := p.sendRequest(ctx, "POST", "/chat/completions", requestBody, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 处理流式响应
	return p.handleStreamResponse(resp.Body, callback)
}

// ChatStreamWithContext 实现AIProvider接口的ChatStreamWithContext方法
func (p *OpenAIProvider) ChatStreamWithContext(ctx context.Context, model string, messages []Message, callback func(chunk string) error) error {
	// 构建请求体
	requestBody := OpenAIRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	// 设置请求头
	headers := map[string]string{
		"Authorization": "Bearer " + p.apiKey,
	}

	// 发送请求
	resp, err := p.sendRequest(ctx, "POST", "/chat/completions", requestBody, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 处理流式响应
	return p.handleStreamResponse(resp.Body, callback)
}

// handleStreamResponse 处理流式响应
func (p *OpenAIProvider) handleStreamResponse(body io.Reader, callback func(chunk string) error) error {
	p.logger.Info("开始处理流式响应")
	
	// 创建一个扫描器来逐行读取响应
	scanner := bufio.NewScanner(body)
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
			p.logger.Info("收到流式响应结束信号")
			break
		}
		// 提取JSON数据
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			// 解析JSON
			var response OpenAIStreamResponse
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				p.logger.Error("解析流式响应失败: %v", err)
				return fmt.Errorf("解析流式响应失败: %w", err)
			}
			// 检查错误
			if response.Error != nil {
				errorMsg := fmt.Sprintf("API错误: %s", response.Error.Message)
				p.logger.Error("%s", errorMsg)
				return fmt.Errorf("%s", errorMsg)
			}
			// 处理响应
			if len(response.Choices) > 0 {
				chunk := response.Choices[0].Delta.Content
				if chunk != "" {
					p.logger.Debug("收到流式响应 chunk: %s", chunk)
					// 调用回调函数
					if err := callback(chunk); err != nil {
						p.logger.Error("回调函数执行失败: %v", err)
						return err
					}
				}
			}
		}
	}

	// 检查扫描器错误
	if err := scanner.Err(); err != nil {
		p.logger.Error("读取流式响应失败: %v", err)
		return fmt.Errorf("读取流式响应失败: %w", err)
	}

	p.logger.Info("流式响应处理完成")
	return nil
}
