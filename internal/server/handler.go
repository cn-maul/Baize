package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"baize/internal/config"
	"baize/internal/provider"
	"baize/pkg/utils"
)

// Handler HTTP处理器结构体
type Handler struct {
	configPath string
	logger     *utils.Logger
	metrics    *Metrics
}

var logger = utils.NewLogger(utils.InfoLevel)

// Metrics 监控指标结构体
type Metrics struct {
	TotalRequests      int64           // 总请求数
	SuccessfulRequests int64           // 成功请求数
	FailedRequests     int64           // 失败请求数
	ResponseTimes      map[string]int64 // 各接口的响应时间总和
	RequestCounts      map[string]int64 // 各接口的请求次数
	mutex              sync.RWMutex    // 互斥锁，用于线程安全
}

// NewMetrics 创建新的监控指标
func NewMetrics() *Metrics {
	return &Metrics{
		ResponseTimes: make(map[string]int64),
		RequestCounts: make(map[string]int64),
	}
}

// NewHandler 创建新的HTTP处理器
func NewHandler(configPath string) *Handler {
	return &Handler{
		configPath: configPath,
		logger:     logger,
		metrics:    NewMetrics(),
	}
}

// ChatRequest 聊天请求结构
type ChatRequest struct {
	Platform string `json:"platform"`
	Model    string `json:"model"`
	Message  string `json:"message"`
}

// ChatResponse 聊天响应结构
type ChatResponse struct {
	Reply string `json:"reply"`
	Error string `json:"error"`
}

// PlatformResponse 平台响应结构
type PlatformResponse struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Models  []string `json:"models"`
}

// PlatformsResponse 平台列表响应结构
type PlatformsResponse struct {
	Platforms []PlatformResponse `json:"platforms"`
}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// ChatHandler 处理聊天请求
func (h *Handler) ChatHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodPost {
		h.logger.Warn("收到非POST方法的请求: %s", r.Method)
		h.sendError(w, http.StatusMethodNotAllowed, "只允许POST方法")
		return
	}

	// 解析请求体
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("解析请求失败: %v", err)
		h.sendError(w, http.StatusBadRequest, fmt.Sprintf("解析请求失败: %v", err))
		return
	}

	// 验证请求参数
	if req.Platform == "" {
		h.logger.Warn("缺少平台参数")
		h.sendError(w, http.StatusBadRequest, "缺少平台参数")
		return
	}
	if req.Model == "" {
		h.logger.Warn("缺少模型参数")
		h.sendError(w, http.StatusBadRequest, "缺少模型参数")
		return
	}
	if req.Message == "" {
		h.logger.Warn("缺少消息参数")
		h.sendError(w, http.StatusBadRequest, "缺少消息参数")
		return
	}

	h.logger.Info("收到聊天请求 - 平台: %s, 模型: %s", req.Platform, req.Model)

	// 加载配置
	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		h.logger.Error("加载配置失败: %v", err)
		h.sendError(w, http.StatusInternalServerError, fmt.Sprintf("加载配置失败: %v", err))
		return
	}

	// 获取平台配置
	platform, err := config.GetPlatformByID(cfg, req.Platform)
	if err != nil {
		h.logger.Warn("获取平台配置失败: %v", err)
		h.sendError(w, http.StatusBadRequest, fmt.Sprintf("获取平台配置失败: %v", err))
		return
	}

	// 创建Provider实例
	prov, err := provider.CreateProvider(platform)
	if err != nil {
		h.logger.Error("创建Provider失败: %v", err)
		h.sendError(w, http.StatusInternalServerError, fmt.Sprintf("创建Provider失败: %v", err))
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// 发送聊天请求
	reply, err := prov.Chat(ctx, req.Model, req.Message)
	if err != nil {
		h.logger.Error("发送聊天请求失败: %v", err)
		h.sendError(w, http.StatusInternalServerError, fmt.Sprintf("发送聊天请求失败: %v", err))
		return
	}

	h.logger.Info("聊天请求处理成功 - 平台: %s, 模型: %s", req.Platform, req.Model)

	// 发送响应
	response := ChatResponse{
		Reply: reply,
		Error: "",
	}
	h.sendJSON(w, http.StatusOK, response)
}

// PlatformsHandler 处理获取平台列表请求
func (h *Handler) PlatformsHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodGet {
		h.logger.Warn("收到非GET方法的请求: %s", r.Method)
		h.sendError(w, http.StatusMethodNotAllowed, "只允许GET方法")
		return
	}

	h.logger.Info("收到平台列表请求")

	// 加载配置
	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		h.logger.Error("加载配置失败: %v", err)
		h.sendError(w, http.StatusInternalServerError, fmt.Sprintf("加载配置失败: %v", err))
		return
	}

	// 构建平台列表响应
	platforms := make([]PlatformResponse, len(cfg.Platforms))
	for i, p := range cfg.Platforms {
		platforms[i] = PlatformResponse{
			ID:     p.ID,
			Name:   p.Name,
			Type:   p.Type,
			Models: p.Models,
		}
	}

	h.logger.Info("平台列表请求处理成功，返回 %d 个平台", len(platforms))

	// 发送响应
	response := PlatformsResponse{
		Platforms: platforms,
	}
	h.sendJSON(w, http.StatusOK, response)
}

// HealthHandler 处理健康检查请求
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodGet {
		h.logger.Warn("收到非GET方法的请求: %s", r.Method)
		h.sendError(w, http.StatusMethodNotAllowed, "只允许GET方法")
		return
	}

	h.logger.Info("收到健康检查请求")

	// 构建健康检查响应
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 发送响应
	h.sendJSON(w, http.StatusOK, response)
}

// sendJSON 发送JSON响应
func (h *Handler) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("编码响应失败: %v", err)
	}
}

// sendError 发送错误响应
func (h *Handler) sendError(w http.ResponseWriter, statusCode int, message string) {
	response := ChatResponse{
		Reply: "",
		Error: message,
	}
	h.sendJSON(w, statusCode, response)
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// 使用CORS中间件和监控中间件包装处理器
	mux.HandleFunc("/api/v1/chat", h.corsMiddleware(h.metricsMiddleware(h.ChatHandler)))
	mux.HandleFunc("/api/v1/platforms", h.corsMiddleware(h.metricsMiddleware(h.PlatformsHandler)))
	mux.HandleFunc("/api/v1/health", h.corsMiddleware(h.metricsMiddleware(h.HealthHandler)))
	mux.HandleFunc("/api/v1/metrics", h.corsMiddleware(h.MetricsHandler))
}

// corsMiddleware CORS中间件
func (h *Handler) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理OPTIONS请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用下一个处理器
		next(w, r)
	}
}

// metricsMiddleware 监控中间件
func (h *Handler) metricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		path := r.URL.Path

		// 创建响应包装器，用于捕获状态码
		wrapper := &responseWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// 调用下一个处理器
		next(wrapper, r)

		// 计算响应时间
		responseTime := time.Since(startTime).Milliseconds()

		// 更新监控指标
		h.metrics.mutex.Lock()
		defer h.metrics.mutex.Unlock()

		h.metrics.TotalRequests++
		h.metrics.RequestCounts[path]++
		h.metrics.ResponseTimes[path] += responseTime

		if wrapper.statusCode >= 200 && wrapper.statusCode < 400 {
			h.metrics.SuccessfulRequests++
		} else {
			h.metrics.FailedRequests++
		}

		// 记录响应时间
		h.logger.Info("请求处理完成 - 路径: %s, 状态码: %d, 响应时间: %dms", path, wrapper.statusCode, responseTime)
	}
}

// responseWrapper 响应包装器，用于捕获状态码
type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 重写WriteHeader方法，捕获状态码
func (w *responseWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// MetricsHandler 处理监控指标请求
func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodGet {
		h.logger.Warn("收到非GET方法的请求: %s", r.Method)
		h.sendError(w, http.StatusMethodNotAllowed, "只允许GET方法")
		return
	}

	h.logger.Info("收到监控指标请求")

	// 构建监控指标响应
	h.metrics.mutex.RLock()
	metricsCopy := struct {
		TotalRequests      int64            `json:"total_requests"`
		SuccessfulRequests int64            `json:"successful_requests"`
		FailedRequests     int64            `json:"failed_requests"`
		ResponseTimes      map[string]int64  `json:"response_times"`
		RequestCounts      map[string]int64  `json:"request_counts"`
		AverageResponseTimes map[string]float64 `json:"average_response_times"`
	}{
		TotalRequests:      h.metrics.TotalRequests,
		SuccessfulRequests: h.metrics.SuccessfulRequests,
		FailedRequests:     h.metrics.FailedRequests,
		ResponseTimes:      make(map[string]int64),
		RequestCounts:      make(map[string]int64),
		AverageResponseTimes: make(map[string]float64),
	}

	// 复制响应时间和请求次数
	for path, time := range h.metrics.ResponseTimes {
		metricsCopy.ResponseTimes[path] = time
	}

	for path, count := range h.metrics.RequestCounts {
		metricsCopy.RequestCounts[path] = count
		// 计算平均响应时间
		if count > 0 {
			metricsCopy.AverageResponseTimes[path] = float64(h.metrics.ResponseTimes[path]) / float64(count)
		}
	}
	h.metrics.mutex.RUnlock()

	// 发送响应
	h.sendJSON(w, http.StatusOK, metricsCopy)
}
