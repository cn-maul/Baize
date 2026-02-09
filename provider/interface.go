package provider

import "context"

// AIProvider 定义了AI提供商的统一接口
type AIProvider interface {
	// Chat 发送聊天请求并获取回复
	// ctx: 上下文，用于控制请求超时等
	// model: 模型名称
	// msg: 用户输入的消息
	// 返回值: 模型的回复和可能的错误
	Chat(ctx context.Context, model string, msg string) (string, error)
}
