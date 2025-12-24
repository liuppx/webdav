package handler

import (
	"net/http"
	
	"github.com/yeying-community/webdav/internal/application/service"
	"go.uber.org/zap"
)

// WebDAVHandler WebDAV 处理器
type WebDAVHandler struct {
	webdavService *service.WebDAVService
	logger        *zap.Logger
}

// NewWebDAVHandler 创建 WebDAV 处理器
func NewWebDAVHandler(webdavService *service.WebDAVService, logger *zap.Logger) *WebDAVHandler {
	return &WebDAVHandler{
		webdavService: webdavService,
		logger:        logger,
	}
}

// Handle 处理 WebDAV 请求
func (h *WebDAVHandler) Handle(w http.ResponseWriter, r *http.Request) {
	h.webdavService.ServeHTTP(w, r)
}

