# 白泽(Baize) AI模型聚合网关实施计划

## 项目概述
基于Go语言标准项目布局，实现一个轻量级、可扩展的AI模型聚合网关，统一管理OpenAI和Anthropic等不同AI厂商的接口差异。

## 实施步骤

### 1. 项目结构搭建
- 创建标准目录结构：`cmd/baize/`、`configs/`、`internal/`、`pkg/`等
- 初始化Go模块，创建`go.mod`文件
- 设置基本的项目结构和文件

### 2. 配置管理模块
- 创建`configs/config.yaml`配置文件模板
- 实现`internal/config`模块，用于加载和解析YAML配置
- 实现`internal/domain`模块，定义核心数据结构（Config、Platform、Model等）

### 3. 核心接口定义
- 在`internal/provider`目录下定义`AIProvider`接口
- 实现工厂模式（`factory.go`），用于根据配置创建具体的Provider实例

### 4. 提供商实现
- 实现`OpenAIProvider`（`openai.go`），对接OpenAI API
- 实现`AnthropicProvider`（`anthropic.go`），对接Anthropic API
- 确保两个实现都符合`AIProvider`接口规范

### 5. HTTP服务层
- 实现`internal/server/handler.go`，处理HTTP请求
- 提供`/api/v1/chat`接口用于发送聊天请求
- 提供`/api/v1/platforms`接口用于获取可用平台和模型列表

### 6. 主程序入口
- 实现`cmd/baize/main.go`，作为程序的唯一入口
- 实现配置加载、服务初始化和端口监听功能

### 7. 工具模块
- 实现`pkg/utils`目录下的通用工具，如HTTP请求封装、日志工具等

### 8. 测试和验证
- 编写单元测试，验证Provider的功能
- 测试HTTP接口的可用性
- 验证不同AI厂商接口的调用是否正常

## 技术要点

1. **接口抽象**：通过`AIProvider`接口统一不同AI厂商的调用方式
2. **工厂模式**：根据配置动态创建Provider实例
3. **配置驱动**：所有平台和模型配置通过YAML文件管理
4. **RESTful API**：提供简洁的HTTP接口供外部调用
5. **安全性**：确保API Key在日志中脱敏显示

## 预期成果

完成后，用户可以通过配置YAML文件添加不同的AI平台和模型，然后通过统一的HTTP接口调用这些模型，无需关心底层的实现细节。