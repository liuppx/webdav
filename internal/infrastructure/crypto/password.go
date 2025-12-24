package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPasswordFormat = errors.New("invalid password format")
	ErrPasswordMismatch      = errors.New("password mismatch")
)

// PasswordHasher 密码哈希器
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher 创建密码哈希器
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		cost: bcrypt.DefaultCost,
	}
}

// Hash 哈希密码
func (h *PasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return "{bcrypt}" + string(hash), nil
}

// Verify 验证密码
func (h *PasswordHasher) Verify(hashedPassword, password string) error {
	// 移除前缀
	if !strings.HasPrefix(hashedPassword, "{bcrypt}") {
		return ErrInvalidPasswordFormat
	}
	
	hash := strings.TrimPrefix(hashedPassword, "{bcrypt}")
	
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}
	
	return nil
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

