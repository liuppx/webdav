package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	
	"go.uber.org/zap"
)

// RecoveryMiddleware 恢复中间件
type RecoveryMiddleware struct {
	logger *zap.Logger
}

// NewRecoveryMiddleware 创建恢复中间件
func NewRecoveryMiddleware(logger *zap.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		logger: logger,
	}
}

// Handle 处理恢复
func (m *RecoveryMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 信息
				m.logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())))
				
				// 返回 500 错误
				http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

