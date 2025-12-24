package auth

import (
	"time"
)

// Challenge Web3 认证挑战
type Challenge struct {
	Nonce     string
	Message   string
	Address   string
	ExpiresAt time.Time
}

// IsExpired 是否过期
func (c *Challenge) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// Validate 验证挑战
func (c *Challenge) Validate() error {
	if c.Nonce == "" {
		return ErrInvalidChallenge
	}
	if c.Address == "" {
		return ErrInvalidChallenge
	}
	if c.IsExpired() {
		return ErrChallengeExpired
	}
	return nil
}

