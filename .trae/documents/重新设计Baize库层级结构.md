# 重新设计Baize库层级结构

## 问题分析

当前项目存在以下问题：

1. 导入路径中出现 `baize/baize/` 这样的重复层级，不够美观
2. 项目同时包含独立服务运行能力和库功能，需要删除独立服务能力
3. 层级结构设计不合理，需要重新组织

## 解决方案

### 新的层级结构设计

将库代码直接放在根目录，移除中间的 `baize/` 目录，使导入路径更简洁：

```text
baize/
├── config/            # 配置加载模块
│   └── config.go
├── domain/            # 领域模型
│   └── model.go
├── provider/          # AI提供商实现
│   ├── anthropic.go
│   ├── factory.go
│   ├── interface.go
│   └── openai.go
├── pkg/               # 通用工具库
│   └── utils/
│       ├── http.go
│       └── logger.go
├── configs/           # 配置文件
│   └── config.yaml
├── go.mod             # 依赖管理
├── go.sum             # 依赖校验
└── README.md          # 项目说明文档
```

### 实施步骤

1. **删除独立服务相关文件**

   * 删除 `cmd/` 目录及其内容

   * 删除 `baize_server` 可执行文件

   * 删除 `test_baize.go` 测试文件

2. **移动库文件到新位置**

   * 将 `baize/config/config.go` 移动到 `config/config.go`

   * 将 `baize/domain/model.go` 移动到 `domain/model.go`

   * 将 `baize/provider/` 目录下的所有文件移动到 `provider/` 目录

   * 删除空的 `baize/` 目录

3. **更新导入路径**

   * 更新 `config/config.go` 中的导入路径，从 `baize/baize/domain` 改为 `baize/domain`

   * 更新 `provider/` 目录下所有文件的导入路径，从 `baize/baize/domain` 改为 `baize/domain`

4. **清理不需要的文件**

   * 删除 `internal/` 目录及其内容（冗余代码）

   * 删除 `.trae/` 目录及其内容（文档已整合到README.md）

5. **更新README.md**

   * 更新目录结构说明

   * 更新使用示例中的导入路径

   * 移除与独立服务相关的内容

## 预期效果

* 导入路径从 `baize/baize/config` 变为 `baize/config`，更加简洁美观

* 移除了独立服务运行能力，专注于库功能

* 层级结构更加清晰合理

* 代码组织符合Go库的最佳实践

## 注意事项

* 确保所有导入路径都正确更新，避免编译错误

* 保持核心功能不变，只改变代码的组织方式

* 确保配置文件路径和其他依赖路径也相应更新

