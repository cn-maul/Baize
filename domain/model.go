package domain

// Config 配置结构体，映射整个YAML文件
type Config struct {
	Version   string                 `yaml:"version"`
	Platforms map[string]*Platform `yaml:"platforms"`
}

// Platform 平台结构体，包含平台的基本信息
type Platform struct {
	ID      string   `yaml:"id"`
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	BaseURL string   `yaml:"base_url"`
	APIKey  string   `yaml:"api_key"`
	Models  []string `yaml:"models"`
}

// Model 模型结构体，定义模型的基本信息
type Model struct {
	Name  string `yaml:"name"`
	Alias string `yaml:"alias"`
}
