package http

import (
	"net/http"
	"strings"

	"github.com/yeying-community/webdav/internal/domain/auth"
	"github.com/yeying-community/webdav/internal/infrastructure/config"
	"github.com/yeying-community/webdav/internal/interface/http/handler"
	"github.com/yeying-community/webdav/internal/interface/http/middleware"
	"go.uber.org/zap"
)

// Router HTTP 路由器
type Router struct {
	config         *config.Config
	authenticators []auth.Authenticator
	healthHandler  *handler.HealthHandler
	web3Handler    *handler.Web3Handler
	webdavHandler  *handler.WebDAVHandler
	logger         *zap.Logger
}

// NewRouter 创建路由器
func NewRouter(
	cfg *config.Config,
	authenticators []auth.Authenticator,
	healthHandler *handler.HealthHandler,
	web3Handler *handler.Web3Handler,
	webdavHandler *handler.WebDAVHandler,
	logger *zap.Logger,
) *Router {
	return &Router{
		config:         cfg,
		authenticators: authenticators,
		healthHandler:  healthHandler,
		web3Handler:    web3Handler,
		webdavHandler:  webdavHandler,
		logger:         logger,
	}
}

// Setup 设置路由
func (r *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	// 健康检查路由（无需认证）
	mux.HandleFunc("/health", r.healthHandler.Handle)

	// Web3 认证路由（无需认证）
	if r.config.Web3.Enabled {
		mux.HandleFunc("/api/auth/challenge", r.web3Handler.HandleChallenge)
		mux.HandleFunc("/api/auth/verify", r.web3Handler.HandleVerify)
	}

	// WebDAV 路由（需要认证）
	webdavPrefix := r.normalizePrefix(r.config.WebDAV.Prefix)
	mux.Handle(webdavPrefix, r.createWebDAVHandler())

	// 应用全局中间件
	handler := r.applyMiddlewares(mux)

	return handler
}

// createWebDAVHandler 创建 WebDAV 处理器（带中间件）
func (r *Router) createWebDAVHandler() http.Handler {
	// 原始的 WebDAV 处理器
	var handler http.Handler = http.HandlerFunc(r.webdavHandler.Handle)

	// 1. 先应用 WebDAV 中间件（最内层）
	//    - 处理 OPTIONS 请求
	//    - 添加 DAV 响应头
	webdavMiddleware := middleware.NewWebDAVMiddleware()
	handler = webdavMiddleware.Handle(handler)

	// 2. 再应用认证中间件（外层）
	//    - 验证用户身份
	//    - OPTIONS 请求需要在这里放行
	authMiddleware := middleware.NewAuthMiddleware(r.authenticators, true, r.logger)
	handler = authMiddleware.Handle(handler)

	return handler
}

// applyMiddlewares 应用全局中间件
func (r *Router) applyMiddlewares(handler http.Handler) http.Handler {
	// 1. 恢复中间件（最外层）
	recoveryMiddleware := middleware.NewRecoveryMiddleware(r.logger)
	handler = recoveryMiddleware.Handle(handler)

	// 2. 日志中间件
	loggerMiddleware := middleware.NewLoggerMiddleware(r.logger, r.config.Security.BehindProxy)
	handler = loggerMiddleware.Handle(handler)

	// 3. CORS 中间件
	if r.config.CORS.Enabled {
		corsConfig := &middleware.CORSConfig{
			Enabled:        r.config.CORS.Enabled,
			Credentials:    r.config.CORS.Credentials,
			AllowedOrigins: r.config.CORS.AllowedOrigins,
			AllowedMethods: r.config.CORS.AllowedMethods,
			AllowedHeaders: r.config.CORS.AllowedHeaders,
			ExposedHeaders: r.config.CORS.ExposedHeaders,
		}
		corsMiddleware := middleware.NewCORSMiddleware(corsConfig)
		handler = corsMiddleware.Handle(handler)
	}

	return handler
}

// normalizePrefix 规范化前缀
func (r *Router) normalizePrefix(prefix string) string {
	if prefix == "" {
		return "/"
	}

	// 确保以 / 开头
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	// 确保以 / 结尾
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	return prefix
}
