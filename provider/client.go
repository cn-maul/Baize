package provider

import (
	"net/http"
	"time"
)

// sharedClient 是一个共享的HTTP客户端，用于所有Provider实例
// 配置了连接池和合理的超时设置，以提高性能
var sharedClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	},
}
