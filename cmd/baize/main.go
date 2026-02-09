package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"baize/internal/server"
	"baize/pkg/utils"
)

const (
	// DefaultConfigPath 默认配置文件路径
	DefaultConfigPath = "configs/config.yaml"
	// DefaultPort 默认端口
	DefaultPort = "8080"
)

var logger = utils.NewLogger(utils.InfoLevel)

func main() {
	// 获取配置文件路径
	configPath := getConfigPath()
	logger.Info("使用配置文件: %s", configPath)

	// 创建HTTP处理器
	handler := server.NewHandler(configPath)

	// 创建HTTP服务器
	mux := http.NewServeMux()

	// 注册路由
	handler.RegisterRoutes(mux)

	// 获取端口
	port := getPort()
	addr := fmt.Sprintf(":%s", port)

	// 启动HTTP服务器
	logger.Info("白泽(Baize) AI网关服务启动中...")
	logger.Info("监听地址: http://localhost%s", addr)
	logger.Info("API接口: http://localhost%s/api/v1/chat", addr)
	logger.Info("平台列表: http://localhost%s/api/v1/platforms", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Fatal("启动服务失败: %v", err)
	}
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	// 检查环境变量
	if path := os.Getenv("BAIZE_CONFIG"); path != "" {
		return path
	}

	// 尝试从多个位置查找配置文件
	possiblePaths := []string{
		"configs/config.yaml",                   // 当前目录
		"../configs/config.yaml",                // 父目录
		"../../configs/config.yaml",             // 祖父目录
		"../../../configs/config.yaml",          // 曾祖父目录
	}

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		return "configs/config.yaml"
	}

	// 检查每个可能的路径
	for _, relPath := range possiblePaths {
		absPath := filepath.Join(cwd, relPath)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	// 尝试从main.go文件所在目录的相对路径查找
	// 获取main.go文件所在目录
	mainDir := filepath.Dir(os.Args[0])
	if mainDir == "" {
		mainDir = cwd
	}

	// 从main.go所在目录向上查找
	for i := 0; i < 4; i++ {
		configPath := filepath.Join(mainDir, strings.Repeat("../", i), "configs/config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// 如果都找不到，返回默认路径
	return "configs/config.yaml"
}

// getPort 获取服务端口
func getPort() string {
	// 检查环境变量
	if port := os.Getenv("BAIZE_PORT"); port != "" {
		return port
	}

	// 使用默认端口
	return DefaultPort
}
