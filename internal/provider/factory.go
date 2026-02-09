package provider

import (
	"fmt"

	"baize/internal/domain"
)

// ProviderFactory 工厂函数类型，用于创建Provider实例
type ProviderFactory func(platform *domain.Platform) (AIProvider, error)

// providerFactories 存储不同类型平台的工厂函数
var providerFactories = map[string]ProviderFactory{
	"openai":     NewOpenAIProvider,
	"anthropic":  NewAnthropicProvider,
}

// RegisterProviderFactory 注册新的Provider工厂函数
func RegisterProviderFactory(providerType string, factory ProviderFactory) {
	providerFactories[providerType] = factory
}

// CreateProvider 根据平台配置创建Provider实例
func CreateProvider(platform *domain.Platform) (AIProvider, error) {
	// 查找对应的工厂函数
	factory, exists := providerFactories[platform.Type]
	if !exists {
		return nil, fmt.Errorf("不支持的平台类型: %s", platform.Type)
	}

	// 使用工厂函数创建Provider实例
	return factory(platform)
}
