package auth

import (
	"time"
)

// Token 认证令牌
type Token struct {
	Value     string
	Address   string
	ExpiresAt time.Time
	IssuedAt  time.Time
}

// IsExpired 是否过期
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// Validate 验证令牌
func (t *Token) Validate() error {
	if t.Value == "" {
		return ErrInvalidToken
	}
	if t.IsExpired() {
		return ErrTokenExpired
	}
	return nil
}
