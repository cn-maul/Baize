package provider

import (
	"time"

	"github.com/cn-maul/Baize/pkg/utils"
)

// ProviderOptions 定义了Provider的配置选项
type ProviderOptions struct {
	Timeout    time.Duration
	MaxRetries int
	LogLevel   utils.LogLevel
}

// ProviderOption 定义了Option模式的函数类型
type ProviderOption func(*ProviderOptions)

// WithTimeout 设置请求超时时间
func WithTimeout(timeout time.Duration) ProviderOption {
	return func(opts *ProviderOptions) {
		opts.Timeout = timeout
	}
}

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(maxRetries int) ProviderOption {
	return func(opts *ProviderOptions) {
		opts.MaxRetries = maxRetries
	}
}

// WithLogLevel 设置日志级别
func WithLogLevel(logLevel utils.LogLevel) ProviderOption {
	return func(opts *ProviderOptions) {
		opts.LogLevel = logLevel
	}
}

// getDefaultOptions 获取默认的ProviderOptions
func getDefaultOptions() *ProviderOptions {
	return &ProviderOptions{
		Timeout:    30 * time.Second,
		MaxRetries: 0,
		LogLevel:   utils.InfoLevel,
	}
}

// applyOptions 应用ProviderOption到ProviderOptions
func applyOptions(opts *ProviderOptions, options ...ProviderOption) {
	for _, option := range options {
		option(opts)
	}
}
