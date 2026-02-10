package provider

import "context"

// Message 消息结构
type Message struct {
	Role    string `json:"role"`    // user, assistant, system
	Content string `json:"content"` // 消息内容
}

// AIProvider 定义了AI提供商的统一接口
type AIProvider interface {
	// Chat 发送聊天请求并获取回复
	// ctx: 上下文，用于控制请求超时等
	// model: 模型名称
	// msg: 用户输入的消息
	// 返回值: 模型的回复和可能的错误
	Chat(ctx context.Context, model string, msg string) (string, error)
	
	// ChatWithContext 发送带上下文的聊天请求并获取回复
	// ctx: 上下文，用于控制请求超时等
	// model: 模型名称
	// messages: 消息历史，包含用户和助手的对话
	// 返回值: 模型的回复和可能的错误
	ChatWithContext(ctx context.Context, model string, messages []Message) (string, error)
	
	// ChatStream 发送聊天请求并流式获取回复
	// ctx: 上下文，用于控制请求超时等
	// model: 模型名称
	// msg: 用户输入的消息
	// callback: 回调函数，用于处理流式输出的每一个 chunk
	// 返回值: 可能的错误
	ChatStream(ctx context.Context, model string, msg string, callback func(chunk string) error) error
	
	// ChatStreamWithContext 发送带上下文的聊天请求并流式获取回复
	// ctx: 上下文，用于控制请求超时等
	// model: 模型名称
	// messages: 消息历史，包含用户和助手的对话
	// callback: 回调函数，用于处理流式输出的每一个 chunk
	// 返回值: 可能的错误
	ChatStreamWithContext(ctx context.Context, model string, messages []Message, callback func(chunk string) error) error
}
