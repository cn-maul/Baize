package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cn-maul/Baize/config"
	"github.com/cn-maul/Baize/domain"
	"github.com/cn-maul/Baize/provider"
	"github.com/cn-maul/Baize/pkg/utils"
)

func main() {
	// 加载配置
	fmt.Println("=== 加载配置文件 ===")
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

	// 打印配置的平台信息
	fmt.Println("\n=== 配置的平台信息 ===")
	for id, platform := range cfg.Platforms {
		fmt.Printf("平台ID: %s\n", id)
		fmt.Printf("  名称: %s\n", platform.Name)
		fmt.Printf("  类型: %s\n", platform.Type)
		fmt.Printf("  基础URL: %s\n", platform.BaseURL)
		fmt.Printf("  可用模型: %v\n", platform.Models)
		fmt.Println()
	}

	// 选择一个平台进行测试
	var testPlatform *domain.Platform
	var platformID string
	
	// 优先选择SiliconFlow平台（如果配置了）
	if platform, exists := cfg.Platforms["SiliconFlow"]; exists {
		testPlatform = platform
		platformID = "SiliconFlow"
	} else {
		// 否则选择第一个可用的平台
		for id, platform := range cfg.Platforms {
			testPlatform = platform
			platformID = id
			break
		}
	}

	fmt.Printf("=== 选择测试平台: %s ===\n", platformID)

	// 检查是否有可用模型
	if len(testPlatform.Models) == 0 {
		fmt.Println("所选平台没有配置可用模型")
		return
	}

	// 选择第一个模型进行测试
	testModel := testPlatform.Models[0]
	fmt.Printf("选择测试模型: %s\n\n", testModel)

	ctx := context.Background()

	// 测试1: 创建Provider实例（使用默认选项）
	fmt.Println("=== 测试1: 创建Provider实例（默认选项）===")
	prov, err := provider.CreateProvider(testPlatform, provider.WithTimeout(120*time.Second))
	if err != nil {
		fmt.Printf("创建Provider失败: %v\n", err)
		return
	}
	fmt.Println("Provider创建成功！")

	// 测试2: 使用Option模式创建Provider
	fmt.Println("\n=== 测试2: 使用Option模式创建Provider ===")
	_, err = provider.CreateProvider(testPlatform,
		provider.WithTimeout(60*time.Second),
		provider.WithMaxRetries(3),
		provider.WithLogLevel(utils.DebugLevel),
	)
	if err != nil {
		fmt.Printf("使用Option模式创建Provider失败: %v\n", err)
		return
	}
	fmt.Println("使用Option模式创建Provider成功！")

	// 测试3: 基本的Chat方法
	fmt.Println("\n=== 测试3: 基本的Chat方法 ===")
	fmt.Println("发送消息: 你好，你是谁？")
	start := time.Now()
	reply, err := prov.Chat(ctx, testModel, "你好，你是谁？")
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("Chat失败: %v\n", err)
	} else {
		fmt.Printf("回复: %s\n", reply)
		fmt.Printf("响应时间: %v\n", elapsed)
	}

	// 测试4: 带上下文的ChatWithContext方法
	fmt.Println("\n=== 测试4: 带上下文的ChatWithContext方法 ===")
	messages := []provider.Message{
		{Role: "user", Content: "你好，你是谁？"},
		{Role: "assistant", Content: reply},
		{Role: "user", Content: "你能做什么？请简单介绍一下你的功能。"},
	}
	
	fmt.Println("发送带上下文的消息:")
	for _, msg := range messages {
		fmt.Printf("  [%s]: %s\n", msg.Role, msg.Content)
	}
	
	start = time.Now()
	contextReply, err := prov.ChatWithContext(ctx, testModel, messages)
	elapsed = time.Since(start)
	if err != nil {
		fmt.Printf("ChatWithContext失败: %v\n", err)
	} else {
		fmt.Printf("\n回复: %s\n", contextReply)
		fmt.Printf("响应时间: %v\n", elapsed)
	}

	// 测试5: 流式输出的ChatStream方法
	fmt.Println("\n=== 测试5: 流式输出的ChatStream方法 ===")
	fmt.Println("发送消息: 你好，请简单介绍一下自己")
	fmt.Println("流式回复:")
	
	start = time.Now()
	err = prov.ChatStream(ctx, testModel, "你好，请简单介绍一下自己", func(chunk string) error {
		fmt.Print(chunk)
		return nil
	})
	elapsed = time.Since(start)
	if err != nil {
		fmt.Printf("\nChatStream失败: %v\n", err)
	} else {
		fmt.Printf("\n\n流式输出完成！响应时间: %v\n", elapsed)
	}

	// 测试6: 带上下文的流式输出ChatStreamWithContext方法
	fmt.Println("\n=== 测试6: 带上下文的流式输出ChatStreamWithContext方法 ===")
	streamMessages := []provider.Message{
		{Role: "user", Content: "你好，你是谁？"},
		{Role: "assistant", Content: reply},
		{Role: "user", Content: "你能做什么？请详细说明你的功能。"},
	}
	
	fmt.Println("发送带上下文的消息:")
	for _, msg := range streamMessages {
		fmt.Printf("  [%s]: %s\n", msg.Role, msg.Content)
	}
	
	fmt.Println("\n流式回复:")
	start = time.Now()
	err = prov.ChatStreamWithContext(ctx, testModel, streamMessages, func(chunk string) error {
		fmt.Print(chunk)
		return nil
	})
	elapsed = time.Since(start)
	if err != nil {
		fmt.Printf("\nChatStreamWithContext失败: %v\n", err)
	} else {
		fmt.Printf("\n\n流式输出完成！响应时间: %v\n", elapsed)
	}

	fmt.Println("\n=== 所有测试完成 ===")
	fmt.Println("Baize库功能测试成功！")
}
