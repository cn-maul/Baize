package main

import (
	"context"
	"fmt"
	"github.com/cn-maul/Baize/config"
	"github.com/cn-maul/Baize/provider"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	// 检查是否有配置的平台
	if len(cfg.Platforms) == 0 {
		fmt.Println("配置文件中没有定义平台")
		return
	}

	// 创建Provider实例
	prov, err := provider.CreateProvider(&cfg.Platforms[0])
	if err != nil {
		fmt.Printf("创建Provider失败: %v\n", err)
		return
	}

	ctx := context.Background()

	// 测试基本的Chat方法
	fmt.Println("测试基本的Chat方法:")
	reply, err := prov.Chat(ctx, cfg.Platforms[0].Models[0], "你好，你是谁？")
	if err != nil {
		fmt.Printf("Chat失败: %v\n", err)
	} else {
		fmt.Printf("回复: %s\n", reply)
	}

	// 测试带上下文的ChatWithContext方法
	fmt.Println("\n测试带上下文的ChatWithContext方法:")
	messages := []provider.Message{
		{Role: "user", Content: "你好，你是谁？"},
		{Role: "assistant", Content: reply},
		{Role: "user", Content: "你能做什么？"},
	}
	contextReply, err := prov.ChatWithContext(ctx, cfg.Platforms[0].Models[0], messages)
	if err != nil {
		fmt.Printf("ChatWithContext失败: %v\n", err)
	} else {
		fmt.Printf("回复: %s\n", contextReply)
	}

	// 测试流式输出的ChatStream方法
	fmt.Println("\n测试流式输出的ChatStream方法:")
	err = prov.ChatStream(ctx, cfg.Platforms[0].Models[0], "你好，请简单介绍一下自己", func(chunk string) error {
		fmt.Print(chunk)
		return nil
	})
	if err != nil {
		fmt.Printf("\nChatStream失败: %v\n", err)
	} else {
		fmt.Println("\n流式输出完成")
	}

	// 测试带上下文的流式输出ChatStreamWithContext方法
	fmt.Println("\n测试带上下文的流式输出ChatStreamWithContext方法:")
	streamMessages := []provider.Message{
		{Role: "user", Content: "你好，你是谁？"},
		{Role: "assistant", Content: reply},
		{Role: "user", Content: "你能做什么？请详细说明"},
	}
	err = prov.ChatStreamWithContext(ctx, cfg.Platforms[0].Models[0], streamMessages, func(chunk string) error {
		fmt.Print(chunk)
		return nil
	})
	if err != nil {
		fmt.Printf("\nChatStreamWithContext失败: %v\n", err)
	} else {
		fmt.Println("\n流式输出完成")
	}

	fmt.Println("\n所有测试完成")
}
