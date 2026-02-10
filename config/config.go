package config

import (
	"fmt"
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
	for i, platform := range config.Platforms {
		if platform.ID == "" {
			return fmt.Errorf("第 %d 个平台缺少ID", i+1)
		}
		if platform.Name == "" {
			return fmt.Errorf("第 %d 个平台缺少名称", i+1)
		}
		if platform.Type == "" {
			return fmt.Errorf("第 %d 个平台缺少类型", i+1)
		}
		if platform.BaseURL == "" {
			return fmt.Errorf("第 %d 个平台缺少基础URL", i+1)
		}
		if platform.APIKey == "" {
			return fmt.Errorf("第 %d 个平台缺少API Key", i+1)
		}
		if len(platform.Models) == 0 {
			return fmt.Errorf("第 %d 个平台没有定义模型", i+1)
		}
	}

	return nil
}

// GetPlatformByID 根据ID获取平台配置
func GetPlatformByID(config *domain.Config, platformID string) (*domain.Platform, error) {
	for _, platform := range config.Platforms {
		if platform.ID == platformID {
			return &platform, nil
		}
	}
	return nil, fmt.Errorf("未找到ID为 %s 的平台", platformID)
}
