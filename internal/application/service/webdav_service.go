package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yeying-community/webdav/internal/domain/permission"
	"github.com/yeying-community/webdav/internal/domain/user"
	"github.com/yeying-community/webdav/internal/infrastructure/config"
	"github.com/yeying-community/webdav/internal/interface/http/middleware"
	"go.uber.org/zap"
	"golang.org/x/net/webdav"
)

// WebDAVService WebDAV 服务
type WebDAVService struct {
	config          *config.Config
	permissionCheck permission.Checker
	logger          *zap.Logger
	lockSystem      webdav.LockSystem
}

// NewWebDAVService 创建 WebDAV 服务
func NewWebDAVService(
	cfg *config.Config,
	permissionCheck permission.Checker,
	logger *zap.Logger,
) *WebDAVService {
	return &WebDAVService{
		config:          cfg,
		permissionCheck: permissionCheck,
		logger:          logger,
		lockSystem:      webdav.NewMemLS(),
	}
}

// ServeHTTP 处理 WebDAV 请求
func (s *WebDAVService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户
	u, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		s.logger.Error("user not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 获取用户目录
	userDir := s.getUserDirectory(u)

	// 确保目录存在
	if err := s.ensureDirectory(userDir); err != nil {
		s.logger.Error("failed to ensure directory",
			zap.String("directory", userDir),
			zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 创建 WebDAV 处理器
	handler := &webdav.Handler{
		Prefix:     s.config.WebDAV.Prefix,
		FileSystem: webdav.Dir(userDir),
		LockSystem: s.lockSystem,
		Logger:     s.createLogger(u.Username),
	}

	// 检查权限
	if err := s.checkPermission(r.Context(), u, r); err != nil {
		s.logger.Warn("permission denied",
			zap.String("username", u.Username),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Error(err))
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// 设置响应头
	if s.config.WebDAV.NoSniff {
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	// 处理请求
	handler.ServeHTTP(w, r)
}

// createLogger 创建 WebDAV 日志记录器
func (s *WebDAVService) createLogger(username string) func(*http.Request, error) {
	return func(r *http.Request, err error) {
		if err == nil {
			return
		}

		// 分类错误
		fields := []zap.Field{
			zap.String("username", username),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		}

		// 判断错误类型
		if isNotFoundError(err) {
			// 文件不存在 - WARN 级别，不打印堆栈
			s.logger.Warn("resource not found",
				append(fields, zap.String("error", err.Error()))...)
		} else if isPermissionError(err) {
			// 权限错误 - WARN 级别
			s.logger.Warn("permission denied",
				append(fields, zap.String("error", err.Error()))...)
		} else if isClientError(err) {
			// 客户端错误 - INFO 级别
			s.logger.Info("client error",
				append(fields, zap.String("error", err.Error()))...)
		} else {
			// 系统错误 - ERROR 级别，打印堆栈
			s.logger.Error("webdav error",
				append(fields, zap.Error(err))...)
		}
	}
}

// isNotFoundError 判断是否为文件不存在错误
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否为 os.ErrNotExist
	if errors.Is(err, os.ErrNotExist) {
		return true
	}

	// 检查错误消息
	errMsg := err.Error()
	return contains(errMsg, "no such file") ||
		contains(errMsg, "not found") ||
		contains(errMsg, "does not exist")
}

// isPermissionError 判断是否为权限错误
func isPermissionError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否为 os.ErrPermission
	if errors.Is(err, os.ErrPermission) {
		return true
	}

	// 检查错误消息
	errMsg := err.Error()
	return contains(errMsg, "permission denied") ||
		contains(errMsg, "access denied") ||
		contains(errMsg, "forbidden")
}

// isClientError 判断是否为客户端错误
func isClientError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	return contains(errMsg, "invalid") ||
		contains(errMsg, "bad request") ||
		contains(errMsg, "malformed")
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(substr) == 0 ||
			containsCaseInsensitive(s, substr))
}

// containsCaseInsensitive 不区分大小写的字符串包含检查
func containsCaseInsensitive(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// toLower 转换为小写
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

// getUserDirectory 获取用户目录
func (s *WebDAVService) getUserDirectory(u *user.User) string {
	// 如果用户有自定义目录，使用用户目录
	if u.Directory != "" {
		// 如果是绝对路径，直接使用
		if filepath.IsAbs(u.Directory) {
			return u.Directory
		}
		// 否则拼接到基础目录
		return filepath.Join(s.config.WebDAV.Directory, u.Directory)
	}

	// 使用基础目录
	return s.config.WebDAV.Directory
}

// ensureDirectory 确保目录存在
func (s *WebDAVService) ensureDirectory(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// 创建目录
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			s.logger.Info("directory created", zap.String("directory", dir))
			return nil
		}
		return fmt.Errorf("failed to stat directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", dir)
	}

	return nil
}

// checkPermission 检查权限
func (s *WebDAVService) checkPermission(ctx context.Context, u *user.User, r *http.Request) error {
	// 映射 HTTP 方法到操作
	operation := permission.MapHTTPMethodToOperation(r.Method)

	// 检查权限
	return s.permissionCheck.Check(ctx, u, r.URL.Path, operation)
}
