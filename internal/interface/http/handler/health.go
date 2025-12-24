package handler

import (
	"encoding/json"
	"net/http"
	"time"
	
	"go.uber.org/zap"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	logger    *zap.Logger
	startTime time.Time
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		logger:    logger,
		startTime: time.Now(),
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status  string        `json:"status"`
	Uptime  time.Duration `json:"uptime"`
	Version string        `json:"version"`
}

// Handle 处理健康检查请求
func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "healthy",
		Uptime:  time.Since(h.startTime),
		Version: "2.0.0",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode health response", zap.Error(err))
	}
}

