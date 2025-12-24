package middleware

import (
	"net/http"
	"strings"
)

// WebDAVMiddleware WebDAV 兼容性中间件
type WebDAVMiddleware struct{}

// NewWebDAVMiddleware 创建 WebDAV 中间件
func NewWebDAVMiddleware() *WebDAVMiddleware {
	return &WebDAVMiddleware{}
}

// Handle 处理请求
func (m *WebDAVMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 处理 OPTIONS 请求
		if r.Method == "OPTIONS" {
			m.handleOptions(w, r)
			return
		}

		// 为所有 WebDAV 请求添加必需的响应头
		w.Header().Set("DAV", "1, 2")
		w.Header().Set("MS-Author-Via", "DAV")

		next.ServeHTTP(w, r)
	})
}

// handleOptions 处理 OPTIONS 请求
func (m *WebDAVMiddleware) handleOptions(w http.ResponseWriter, _ *http.Request) {
	// 允许的方法
	methods := []string{
		"OPTIONS",
		"GET", "HEAD", "POST", "PUT", "DELETE",
		"PROPFIND", "PROPPATCH",
		"MKCOL", "COPY", "MOVE",
		"LOCK", "UNLOCK",
	}

	// 设置响应头
	w.Header().Set("Allow", strings.Join(methods, ", "))
	w.Header().Set("DAV", "1, 2")
	w.Header().Set("MS-Author-Via", "DAV")
	w.Header().Set("Accept-Ranges", "bytes")

	// 返回 200 OK（不是 204）
	w.WriteHeader(http.StatusOK)
}
