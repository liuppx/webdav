package auth

import "errors"

var (
	// ErrInvalidToken 无效的 token
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired token 已过期
	ErrTokenExpired = errors.New("token expired")

	// ErrInvalidCredentials 无效的凭证
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrChallengeExpired 挑战已过期
	ErrChallengeExpired = errors.New("challenge expired")

	// ErrInvalidSignature 无效的签名
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrInvalidChallenge 无效挑战信息
	ErrInvalidChallenge = errors.New("invalid challenge")
)
