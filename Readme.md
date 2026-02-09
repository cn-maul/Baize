这份文档是按照 **Go 语言标准项目布局 (Standard Go Project Layout)** 编写的。这种结构在 Go 社区非常流行，既能满足你目前“个人开发”的简单需求，又为未来扩展（比如增加新的 AI 厂商、增加数据库、增加 Web 界面）留足了空间。

---

# 白泽 (Baize) - AI 核心引擎设计文档 v1.0

## 1. 项目概述

**白泽 (Baize)** 是一款轻量级、可扩展的 AI 模型聚合网关。

* **核心目标**：统一管理不同 AI 厂商（OpenAI, Anthropic）的接口差异，对外提供统一的调用方式。
* **设计理念**：配置驱动（Configuration Driven），通过本地 YAML 文件管理所有的平台接入凭证和模型列表。

---

## 2. 标准目录结构 (Directory Structure)

这是基于 Go 社区标准 (`golang-standards/project-layout`) 裁剪出的最适合本项目的结构：

```text
baize/
├── cmd/
│   └── baize/
│       └── main.go           # 【入口】程序的唯一入口，只负责初始化和启动服务
├── configs/
│   └── config.yaml           # 【配置】默认的配置文件模板
├── internal/                 # 【私有代码】核心逻辑放在这里，防止被外部引用
│   ├── config/               # 配置加载模块 (读取 YAML)
│   ├── domain/               # 领域模型 (定义核心 Structs，如 Platform, Model)
│   ├── provider/             # 核心业务 (OpenAI/Anthropic 的具体实现)
│   │   ├── openai.go
│   │   ├── anthropic.go
│   │   └── factory.go        # 工厂模式，用于生产具体的 Provider
│   └── server/               # HTTP 服务模块 (处理路由、请求响应)
│       └── handler.go
├── pkg/                      # 【公共代码】(可选) 未来如果想把某些工具库开源，放这里
│   └── utils/                # 通用工具 (如 HTTP 请求封装、日志工具)
├── web/                      # 【前端界面】Web 前端界面，用于测试和使用 API
│   └── index.html            # 前端主页面
├── go.mod                    # 依赖管理
├── go.sum                    # 依赖校验
└── README.md                 # 项目说明文档

```

### 目录详细说明：

1. **`cmd/baize/main.go`**: 这里的代码应该非常少。它只做三件事：加载配置 -> 初始化服务 -> 监听端口。
2. **`internal/domain`**: 存放纯粹的数据结构定义（Structs），不包含复杂的逻辑。
3. **`internal/provider`**: 这是“白泽”的大脑。利用 **接口 (Interface)** 来抹平 OpenAI 和 Anthropic 的差异。
4. **`internal/server`**: 处理外界发来的 HTTP 请求，解析参数，调用 `provider`，然后返回 JSON。
5. **`web/index.html`**: 前端界面，用于测试和使用 API，支持选择平台、选择模型和发送消息。

---

## 3. 核心设计模式与理念

为了实现“接入不同平台”且“随时可扩展”，我们采用以下设计模式：

### 3.1 策略模式 (Strategy Pattern) & 接口抽象

不要在代码里写大量的 `if type == "openai"`。我们定义一个统一的接口。

**设计思路：**

* 定义一个接口 `AIProvider`，包含 `Chat()` 方法。
* `OpenAIProvider` 实现这个接口。
* `AnthropicProvider` 也实现这个接口。
* 上层业务只需要调用 `provider.Chat()`，不需要关心底层是谁。

### 3.2 工厂模式 (Factory Pattern)

用于根据配置文件中的字符串（"openai" 或 "anthropic"）自动创建对应的 Provider 实例。

### 3.3 配置驱动 (Config Driven)

所有的变动（增加平台、换 Key、加模型）都只修改 YAML 文件，不需要重新编译代码（除非重启程序以加载新配置）。

---

## 4. 详细模块设计

### 4.1 领域模型设计 (`internal/domain`)

这里定义数据在内存中的样子。

* **Config Struct**: 映射整个 YAML 文件。
* **Platform Struct**: 包含 `Name`, `Type` (枚举), `BaseURL`, `APIKey`。
* **Model Struct**: 包含 `Name`, `Alias`。

### 4.2 接口定义 (`internal/provider`)

这是解耦的关键。

```go
// 伪代码展示接口设计
type AIProvider interface {
    // 统一的对话方法
    // msg: 用户输入
    // model: 指定的模型ID
    Chat(ctx context.Context, model string, msg string) (string, error)
}

```

### 4.3 配置管理流程 (`internal/config`)

1. **加载策略**：程序启动时读取 `configs/config.yaml`。
2. **热加载 (可选)**：为了简单，第一版可以不做。修改配置后重启服务即可。
3. **安全性**：设计中要考虑到 API Key 在 Log 中需要脱敏显示。

### 4.4 HTTP 服务层 (`internal/server`)

对外暴露简单的 RESTful 接口：

* `POST /api/v1/chat`
* **Request**: `{ "platform": "MyOpenAI", "model": "gpt-4", "message": "你好" }`
* **Response**: `{ "reply": "你好！我是白泽...", "error": "" }`


* `GET /api/v1/platforms`
* **Response**: 返回当前配置中所有可用的平台和模型列表（方便前端下拉选择）。



---

## 5. 开发路线图 (Roadmap)

作为业余开发者，建议按照以下顺序实现，保证每一步都有产出：

1. **Phase 1: 骨架搭建**
* 创建上述目录结构。
* 编写 `config.yaml` 和对应的 `domain` 结构体。
* 实现 `internal/config`，确保能把 YAML 打印在控制台。


2. **Phase 2: 核心打通**
* 定义 `AIProvider` 接口。
* 实现 `OpenAIProvider` (先只做这就行)。
* 写一个单元测试，直接调用 Provider 发送请求，确保网络通畅。


3. **Phase 3: 服务化**
* 实现 `internal/server`，用 `net/http` 启动一个 Web Server。
* 把 HTTP 请求参数传给 Provider。


4. **Phase 4: 完整接入**
* 实现 `AnthropicProvider`。
* 利用工厂模式，根据请求参数动态切换 Provider。



---

## 6. 配置文件设计规范 (configs/config.yaml)

最终的配置文件应该长这样，清晰易读：

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

---

## 7. 使用指南

### 7.1 启动服务

1. **准备配置文件**：
   - 编辑 `configs/config.yaml` 文件，添加您的AI平台配置和API密钥

2. **启动服务**：
   ```bash
   # 在项目根目录运行
   go run cmd/baize/main.go
   ```
   或者
   ```bash
   # 在cmd/baize目录运行
   go run main.go
   ```

3. **服务默认配置**：
   - 默认端口：8080
   - 默认配置文件路径：configs/config.yaml

4. **环境变量**：
   - `BAIZE_CONFIG`：指定配置文件路径
   - `BAIZE_PORT`：指定服务端口

### 7.2 使用前端界面

1. **打开前端界面**：
   - 在浏览器中打开 `file:///path/to/baize/web/index.html`

2. **使用方法**：
   - 在"配置"区域选择平台和模型
   - 在聊天输入框中输入消息
   - 点击"发送"按钮或按回车键发送消息
   - 查看AI的回复

### 7.3 API文档

#### 7.3.1 聊天接口

- **URL**：`POST /api/v1/chat`
- **请求体**：
  ```json
  {
    "platform": "openai_main",
    "model": "gpt-4-turbo",
    "message": "你好"
  }
  ```
- **响应**：
  ```json
  {
    "reply": "你好！我是白泽...",
    "error": ""
  }
  ```

#### 7.3.2 平台列表接口

- **URL**：`GET /api/v1/platforms`
- **响应**：
  ```json
  {
    "platforms": [
      {
        "id": "openai_main",
        "name": "主OpenAI账号",
        "type": "openai",
        "models": ["gpt-4-turbo", "gpt-3.5-turbo"]
      },
      {
        "id": "claude_backup",
        "name": "Claude备用",
        "type": "anthropic",
        "models": ["claude-3-opus-20240229", "claude-3-sonnet-20240229"]
      }
    ]
  }
  ```

### 7.4 健康检查接口

- **URL**：`GET /api/v1/health`
- **响应**：
  ```json
  {
    "status": "ok",
    "timestamp": "2026-02-09T12:00:00Z"
  }
  ```

---

## 8. 技术栈

- **后端**：Go 1.20+
- **前端**：HTML5, Tailwind CSS v3, JavaScript
- **依赖**：
  - gopkg.in/yaml.v3：YAML配置解析
  - 标准库：net/http, encoding/json, context等

---

## 9. 项目状态

- ✅ 项目结构搭建
- ✅ 配置管理模块
- ✅ 核心接口定义
- ✅ 提供商实现（OpenAI, Anthropic）
- ✅ HTTP服务层
- ✅ 前端界面
- ✅ 健康检查接口
- ✅ 配置缓存机制
- ✅ 日志记录
- ✅ 错误处理
- ✅ 单元测试

项目已完成所有核心功能，可直接使用。

