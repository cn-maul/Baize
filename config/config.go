package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/cn-maul/Baize/domain"
)

// LoadConfig 加载并解析配置文件
func LoadConfig(configPath string) (*domain.Config, error) {
	// 确保配置文件路径是绝对路径
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法获取配置文件绝对路径: %w", err)
	}

	// 读取配置文件内容
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件: %w", err)
	}

	// 检查配置文件是否为空
	if len(data) == 0 {
		return nil, fmt.Errorf("配置文件为空")
	}

	// 解析YAML配置
	var config domain.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("无法解析配置文件: %w", err)
	}

	// 验证配置有效性
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig 验证配置有效性
func validateConfig(config *domain.Config) error {
	// 检查版本
	if config.Version == "" {
		return fmt.Errorf("配置文件缺少版本信息")
	}

	// 检查平台列表
	if len(config.Platforms) == 0 {
		return fmt.Errorf("配置文件中没有定义平台")
	}

	// 检查每个平台的有效性
	for name, platform := range config.Platforms {
		if platform.ID == "" {
			return fmt.Errorf("平台 %s 缺少ID", name)
		}
		if platform.Name == "" {
			return fmt.Errorf("平台 %s 缺少名称", name)
		}
		if platform.Type == "" {
			return fmt.Errorf("平台 %s 缺少类型", name)
		}
		if platform.BaseURL == "" {
			return fmt.Errorf("平台 %s 缺少基础URL", name)
		}
		if platform.APIKey == "" {
			return fmt.Errorf("平台 %s 缺少API Key", name)
		}
		if len(platform.Models) == 0 {
			return fmt.Errorf("平台 %s 没有定义模型", name)
		}

		// 检查BaseURL格式
		if _, err := url.Parse(platform.BaseURL); err != nil {
			return fmt.Errorf("平台 %s 的BaseURL格式无效: %w", name, err)
		}

		// 检查API Key长度
		if len(platform.APIKey) < 10 {
			return fmt.Errorf("平台 %s 的API Key长度不足", name)
		}

		// 检查平台类型是否支持
		if platform.Type != "openai" && platform.Type != "anthropic" {
			return fmt.Errorf("平台 %s 的类型 %s 不支持", name, platform.Type)
		}
	}

	return nil
}

// GetPlatformByID 根据ID获取平台配置
func GetPlatformByID(config *domain.Config, platformID string) (*domain.Platform, error) {
	// 首先尝试通过map键查找
	if platform, ok := config.Platforms[platformID]; ok {
		return platform, nil
	}

	// 然后尝试通过ID字段查找
	for _, platform := range config.Platforms {
		if platform.ID == platformID {
			return platform, nil
		}
	}

	return nil, fmt.Errorf("未找到ID为 %s 的平台", platformID)
}
