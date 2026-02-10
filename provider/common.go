package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cn-maul/Baize/pkg/utils"
)

// BaseProvider 是所有Provider的基础结构，包含公共逻辑
type BaseProvider struct {
	baseURL string
	apiKey  string
	client  *http.Client
	logger  *utils.Logger
}

// NewBaseProvider 创建一个新的BaseProvider实例
func NewBaseProvider(baseURL, apiKey string, options ...ProviderOption) *BaseProvider {
	// 获取默认选项
	opts := getDefaultOptions()
	// 应用用户提供的选项
	applyOptions(opts, options...)

	// 创建一个基于sharedClient配置但使用自定义超时的客户端
	client := &http.Client{
		Timeout:   opts.Timeout,
		Transport: sharedClient.Transport,
	}

	return &BaseProvider{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  client,
		logger:  utils.NewLogger(opts.LogLevel),
	}
}

// sendRequest 发送HTTP请求并处理响应
func (p *BaseProvider) sendRequest(ctx context.Context, method, endpoint string, reqBody interface{}, headers map[string]string) (*http.Response, error) {
	// 序列化请求体
	requestJSON, err := json.Marshal(reqBody)
	if err != nil {
		p.logger.Error("序列化请求体失败: %v", err)
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, method, p.baseURL+endpoint, bytes.NewBuffer(requestJSON))
	if err != nil {
		p.logger.Error("创建请求失败: %v", err)
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置默认请求头
	req.Header.Set("Content-Type", "application/json")

	// 设置自定义请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 记录请求信息（脱敏处理）
	sanitizedHeaders := make(map[string]string)
	for key, value := range headers {
		if key == "Authorization" || key == "x-api-key" {
			sanitizedHeaders[key] = "***"
		} else {
			sanitizedHeaders[key] = value
		}
	}
	p.logger.Info("发送HTTP请求: %s %s", method, req.URL.String())
	p.logger.Debug("请求头: %v", sanitizedHeaders)

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Error("发送请求失败: %v", err)
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}

	// 记录响应状态码
	p.logger.Info("收到HTTP响应: %d", resp.StatusCode)

	// 检查HTTP状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		errorMsg := fmt.Sprintf("HTTP请求失败: %d, body: %s", resp.StatusCode, string(body))
		p.logger.Error("%s", errorMsg)
		return nil, fmt.Errorf("%s", errorMsg)
	}

	return resp, nil
}
