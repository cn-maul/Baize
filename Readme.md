# 白泽 (Baize) - AI 模型聚合网关

[![GitHub](https://img.shields.io/github/stars/cn-maul/Baize?style=social)](https://github.com/cn-maul/Baize)

## 项目概述

**白泽 (Baize)** 是一款轻量级、可扩展的 AI 模型聚合网关库。

* **核心目标**：统一管理不同 AI 厂商（OpenAI, Anthropic）的接口差异，对外提供统一的调用方式。
* **设计理念**：配置驱动（Configuration Driven），通过本地 YAML 文件管理所有的平台接入凭证和模型列表。

## 标准目录结构 (Directory Structure)

这是基于 Go 社区标准 (`golang-standards/project-layout`) 裁剪出的最适合本项目的结构：

```text
baize/
├── config/                    # 配置加载模块 (读取 YAML)
│   └── config.go
├── domain/                    # 领域模型 (定义核心 Structs，如 Platform, Model)
│   └── model.go
├── provider/                  # 核心业务 (OpenAI/Anthropic 的具体实现)
│   ├── openai.go
│   ├── anthropic.go
│   ├── factory.go             # 工厂模式，用于生产具体的 Provider
│   ├── interface.go           # 核心接口定义
│   ├── client.go              # 共享HTTP客户端
│   ├── common.go              # 公共逻辑
│   ├── errors.go              # 错误定义
│   └── options.go             # Option模式支持
├── pkg/                       # 【公共代码】通用工具库
│   └── utils/                 # 通用工具 (如 HTTP 请求封装、日志工具)
│       ├── http.go
│       └── logger.go

├── config.yaml                # 【配置】默认的配置文件
├── go.mod                     # 依赖管理
├── go.sum                     # 依赖校验
└── README.md                  # 项目说明文档

```

### 目录详细说明：

1. **`config/`**: 配置加载模块，用于加载和解析 YAML 配置文件
   - `config.go`: 配置加载和解析逻辑
2. **`domain/`**: 领域模型，定义核心数据结构如 Platform, Model
   - `model.go`: 核心数据结构定义
3. **`provider/`**: 核心业务逻辑，实现不同 AI 厂商的接口
   - `openai.go`: OpenAI 提供商实现
   - `anthropic.go`: Anthropic 提供商实现
   - `factory.go`: 工厂模式，用于生产具体的 Provider
   - `interface.go`: 核心接口定义
   - `client.go`: 共享 HTTP 客户端实现
   - `common.go`: 公共逻辑封装
   - `errors.go`: 错误定义
   - `options.go`: Option 模式支持
4. **`pkg/utils/`**: 通用工具库，如 HTTP 请求封装、日志工具
   - `http.go`: HTTP 工具函数
   - `logger.go`: 日志工具实现

## 核心设计模式与理念

为了实现"接入不同平台"且"随时可扩展"，我们采用以下设计模式：

### 策略模式 (Strategy Pattern) & 接口抽象

通过接口抽象抹平不同AI厂商的差异：

* 定义一个接口 `AIProvider`，包含以下方法：
  - `Chat()`: 基本的单次聊天方法
  - `ChatWithContext()`: 带上下文的聊天方法，支持维护对话历史
  - `ChatStream()`: 流式输出的聊天方法，实时显示AI的回复
  - `ChatStreamWithContext()`: 同时支持上下文和流式输出的聊天方法
* `OpenAIProvider` 和 `AnthropicProvider` 分别实现这个接口
* 上层业务只需要调用相应的方法，不需要关心底层是谁

### 工厂模式 (Factory Pattern)

用于根据配置文件中的字符串（"openai" 或 "anthropic"）自动创建对应的 Provider 实例。

### 配置驱动 (Config Driven)

所有的变动（增加平台、换 Key、加模型）都只修改 YAML 文件，不需要重新编译代码。

## 性能优化

为了提高库的性能和可靠性，我们进行了以下优化：

### 1. 性能提升
- **替换 LineScanner**：使用标准库 `bufio.Scanner` 替换了自定义的低性能实现，解决了错误处理问题
- **实现共享 HTTP 客户端**：创建了配置了连接池的共享 HTTP 客户端，提高了网络请求效率
- **优化字符串拼接**：使用 `strings.Builder` 替代了直接的字符串拼接，减少了内存分配
- **消除重复代码**：提取了公共逻辑到 `BaseProvider`，减少了代码重复率

### 2. 功能增强
- **添加 HTTP 状态码检查**：对所有 HTTP 请求添加了状态码检查，确保正确处理 4xx/5xx 错误
- **集成日志记录**：在关键操作中添加了详细的日志记录，包括请求/响应信息和错误处理
- **添加配置验证增强**：实现了更严格的配置验证，包括 URL 格式检查、API Key 长度验证和平台类型检查
- **添加 Option 模式支持**：实现了灵活的配置选项机制，支持设置超时、重试次数和日志级别

### 3. 质量保证
- **修复编译错误**：解决了所有编译错误，确保代码可以正常构建
- **提高代码可读性**：优化了代码结构和命名，添加了必要的注释

## 配置文件设计规范 (config.yaml)

配置文件示例：

```yaml
version: "1.0"

platforms:
  - id: "openai_main"          # 唯一标识，代码里用这个找平台
    name: "主OpenAI账号"        # 显示给用户看的名字
    type: "openai"             # 核心字段：决定了调用哪个代码逻辑
    base_url: "https://api.openai.com/v1"
    api_key: "sk-xxxxxxxx"
    models:                    # 该账号下可用的模型
      - "gpt-4-turbo"
      - "gpt-3.5-turbo"

  - id: "claude_backup"
    name: "Claude备用"
    type: "anthropic"
    base_url: "https://api.anthropic.com"
    api_key: "sk-ant-xxxxxxx"
    models:
      - "claude-3-opus-20240229"
      - "claude-3-sonnet-20240229"

```

## 使用指南

### 作为库在其他项目中使用

```go
import (
    "context"
    "fmt"
    "time"
    "github.com/cn-maul/Baize/config"
    "github.com/cn-maul/Baize/provider"
    "github.com/cn-maul/Baize/pkg/utils"
)

// 加载配置
cfg, err := config.LoadConfig("path/to/config.yaml")

// 创建Provider实例（使用默认选项）
prov, err := provider.CreateProvider(&cfg.Platforms["openai"])

// 使用Option模式创建Provider（自定义配置）
provWithOptions, err := provider.CreateProvider(&cfg.Platforms["openai"],
    provider.WithTimeout(60*time.Second),
    provider.WithMaxRetries(3),
    provider.WithLogLevel(utils.DebugLevel),
)

ctx := context.Background()

// 基本的聊天请求
reply, err := prov.Chat(ctx, "model-name", "Hello!")

// 带上下文的聊天请求
messages := []provider.Message{
    {Role: "user", Content: "你好，你是谁？"},
    {Role: "assistant", Content: "我是AI助手"},
    {Role: "user", Content: "你能做什么？"},
}
contextReply, err := prov.ChatWithContext(ctx, "model-name", messages)

// 流式输出的聊天请求
err := prov.ChatStream(ctx, "model-name", "你好，请简单介绍一下自己", func(chunk string) error {
    fmt.Print(chunk) // 实时显示AI的回复
    return nil
})

// 带上下文的流式输出聊天请求
streamMessages := []provider.Message{
    {Role: "user", Content: "你好，你是谁？"},
    {Role: "assistant", Content: "我是AI助手"},
    {Role: "user", Content: "你能做什么？请详细说明"},
}
err := prov.ChatStreamWithContext(ctx, "model-name", streamMessages, func(chunk string) error {
    fmt.Print(chunk) // 实时显示AI的回复
    return nil
})
```

## 技术栈

- **后端**：Go 1.25+

- **依赖**：
  - gopkg.in/yaml.v3：YAML配置解析
  - 标准库：net/http, encoding/json, context等

## 项目状态

- ✅ 项目结构搭建
- ✅ 配置管理模块
- ✅ 核心接口定义
- ✅ 提供商实现（OpenAI, Anthropic）
- ✅ 上下文支持功能
- ✅ 流式输出功能
- ✅ 性能优化
  - ✅ 替换LineScanner为标准库bufio.Scanner
  - ✅ 实现共享HTTP客户端
  - ✅ 优化字符串拼接
  - ✅ 消除重复代码
- ✅ 功能增强
  - ✅ 添加HTTP状态码检查
  - ✅ 集成日志记录
  - ✅ 添加配置验证增强
  - ✅ 添加Option模式支持
- ✅ 质量保证
  - ✅ 修复编译错误
  - ✅ 提高代码可读性

项目已完成所有核心功能和优化，可直接在生产环境中使用。

## 项目信息

- **GitHub URL**: https://github.com/cn-maul/Baize
- **作者**: cn-maul
- **许可证**: MIT
