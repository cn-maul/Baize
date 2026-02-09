package config

import (
	"os"
	"path/filepath"
	"testing"

	"baize/internal/domain"
)

// TestLoadConfig 测试加载配置文件
func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// 写入测试配置
	testConfig := `version: "1.0"
platforms:
  - id: "test_openai"
    name: "测试OpenAI"
    type: "openai"
    base_url: "https://api.openai.com/v1"
    api_key: "sk-test1234"
    models:
      - "gpt-4"
      - "gpt-3.5-turbo"
`

	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("写入测试配置文件失败: %v", err)
	}

	// 加载配置
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}

	// 验证配置
	if config.Version != "1.0" {
		t.Errorf("期望版本为 1.0, 实际为 %s", config.Version)
	}

	if len(config.Platforms) != 1 {
		t.Errorf("期望平台数量为 1, 实际为 %d", len(config.Platforms))
	}

	platform := config.Platforms[0]
	if platform.ID != "test_openai" {
		t.Errorf("期望平台ID为 test_openai, 实际为 %s", platform.ID)
	}

	if platform.Type != "openai" {
		t.Errorf("期望平台类型为 openai, 实际为 %s", platform.Type)
	}

	if len(platform.Models) != 2 {
		t.Errorf("期望模型数量为 2, 实际为 %d", len(platform.Models))
	}
}

// TestGetPlatformByID 测试根据ID获取平台配置
func TestGetPlatformByID(t *testing.T) {
	// 创建测试配置
	config := &domain.Config{
		Version: "1.0",
		Platforms: []domain.Platform{
			{
				ID:      "platform1",
				Name:    "平台1",
				Type:    "openai",
				BaseURL: "https://api.openai.com/v1",
				APIKey:  "sk-test1234",
				Models:  []string{"gpt-4"},
			},
			{
				ID:      "platform2",
				Name:    "平台2",
				Type:    "anthropic",
				BaseURL: "https://api.anthropic.com",
				APIKey:  "sk-ant-test1234",
				Models:  []string{"claude-3-opus"},
			},
		},
	}

	// 测试获取存在的平台
	platform, err := GetPlatformByID(config, "platform1")
	if err != nil {
		t.Fatalf("获取平台失败: %v", err)
	}

	if platform.ID != "platform1" {
		t.Errorf("期望平台ID为 platform1, 实际为 %s", platform.ID)
	}

	if platform.Type != "openai" {
		t.Errorf("期望平台类型为 openai, 实际为 %s", platform.Type)
	}

	// 测试获取不存在的平台
	_, err = GetPlatformByID(config, "non_existent")
	if err == nil {
		t.Error("期望获取不存在的平台时返回错误, 但未返回错误")
	}
}

// TestLoadConfig_InvalidConfig 测试加载无效配置
func TestLoadConfig_InvalidConfig(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid_config.yaml")

	// 写入无效配置（缺少必要字段）
	invalidConfig := `version: "1.0"
platforms:
  - id: "test_openai"
    name: "测试OpenAI"
    # 缺少type字段
    base_url: "https://api.openai.com/v1"
    api_key: "sk-test1234"
    models:
      - "gpt-4"
`

	if err := os.WriteFile(configPath, []byte(invalidConfig), 0644); err != nil {
		t.Fatalf("写入测试配置文件失败: %v", err)
	}

	// 加载配置，应该返回错误
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("期望加载无效配置时返回错误, 但未返回错误")
	}
}

// TestLoadConfig_EmptyConfig 测试加载空配置
func TestLoadConfig_EmptyConfig(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "empty_config.yaml")

	// 写入空配置
	emptyConfig := ""

	if err := os.WriteFile(configPath, []byte(emptyConfig), 0644); err != nil {
		t.Fatalf("写入测试配置文件失败: %v", err)
	}

	// 加载配置，应该返回错误
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("期望加载空配置时返回错误, 但未返回错误")
	}
}

// TestLoadConfig_Cache 测试配置缓存
func TestLoadConfig_Cache(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "cache_config.yaml")

	// 写入测试配置
	testConfig := `version: "1.0"
platforms:
  - id: "test_openai"
    name: "测试OpenAI"
    type: "openai"
    base_url: "https://api.openai.com/v1"
    api_key: "sk-test1234"
    models:
      - "gpt-4"
`

	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("写入测试配置文件失败: %v", err)
	}

	// 第一次加载配置
	config1, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("第一次加载配置文件失败: %v", err)
	}

	// 第二次加载配置（应该使用缓存）
	config2, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("第二次加载配置文件失败: %v", err)
	}

	// 验证两次加载的是同一个配置对象
	if config1 != config2 {
		t.Error("期望第二次加载使用缓存，但返回了不同的配置对象")
	}

	// 清除缓存
	ClearConfigCache()

	// 第三次加载配置（应该重新加载）
	config3, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("第三次加载配置文件失败: %v", err)
	}

	// 验证清除缓存后加载的是不同的配置对象
	if config1 == config3 {
		t.Error("期望清除缓存后重新加载配置，但返回了相同的配置对象")
	}
}
