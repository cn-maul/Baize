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
├── domain/                    # 领域模型 (定义核心 Structs，如 Platform, Model)
├── provider/                  # 核心业务 (OpenAI/Anthropic 的具体实现)
│   ├── openai.go
│   ├── anthropic.go
│   ├── factory.go             # 工厂模式，用于生产具体的 Provider
│   └── interface.go           # 核心接口定义
├── pkg/                       # 【公共代码】通用工具库
│   └── utils/                 # 通用工具 (如 HTTP 请求封装、日志工具)

├── config.yaml                # 【配置】默认的配置文件
├── go.mod                     # 依赖管理
├── go.sum                     # 依赖校验
└── README.md                  # 项目说明文档

```

### 目录详细说明：

1. **`config/`**: 配置加载模块，用于加载和解析 YAML 配置文件
2. **`domain/`**: 领域模型，定义核心数据结构如 Platform, Model
3. **`provider/`**: 核心业务逻辑，实现不同 AI 厂商的接口
4. **`pkg/utils/`**: 通用工具库，如 HTTP 请求封装、日志工具

## 核心设计模式与理念

为了实现"接入不同平台"且"随时可扩展"，我们采用以下设计模式：

### 策略模式 (Strategy Pattern) & 接口抽象

通过接口抽象抹平不同AI厂商的差异：

* 定义一个接口 `AIProvider`，包含 `Chat()` 方法
* `OpenAIProvider` 和 `AnthropicProvider` 分别实现这个接口
* 上层业务只需要调用 `provider.Chat()`，不需要关心底层是谁

### 工厂模式 (Factory Pattern)

用于根据配置文件中的字符串（"openai" 或 "anthropic"）自动创建对应的 Provider 实例。

### 配置驱动 (Config Driven)

所有的变动（增加平台、换 Key、加模型）都只修改 YAML 文件，不需要重新编译代码。

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
    "baize/config"
    "baize/provider"
)

// 加载配置
cfg, err := config.LoadConfig("path/to/config.yaml")

// 创建Provider实例
prov, err := provider.CreateProvider(&cfg.Platforms[0])

// 发送聊天请求
reply, err := prov.Chat(ctx, "model-name", "Hello!")
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
- ✅ 配置缓存机制
- ✅ 项目已改造为可重用的Go库

项目已完成所有核心功能，可直接使用。

## 项目信息

- **GitHub URL**: https://github.com/cn-maul/Baize
- **作者**: cn-maul
- **许可证**: MIT
