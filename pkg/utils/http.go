package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPClient HTTP客户端接口
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(timeout time.Duration) HTTPClient {
	return &http.Client{
		Timeout: timeout,
	}
}

// PostJSON 发送POST请求，使用JSON格式
func PostJSON(ctx context.Context, client HTTPClient, url string, headers map[string]string, data interface{}) (*http.Response, error) {
	// 序列化请求数据
	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置默认请求头
	req.Header.Set("Content-Type", "application/json")

	// 设置自定义请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	return client.Do(req)
}

// ParseJSONResponse 解析JSON响应
func ParseJSONResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	return nil
}

// GetStatusCode 获取HTTP响应状态码
func GetStatusCode(resp *http.Response) int {
	if resp == nil {
		return 0
	}
	return resp.StatusCode
}

// IsSuccessStatusCode 检查是否是成功的状态码
func IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
