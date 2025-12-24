package config

import (
	"errors"
	"fmt"
	"os"
)

// Validator 配置验证器
type Validator struct{}

// NewValidator 创建验证器
func NewValidator() *Validator {
	return &Validator{}
}

// Validate 验证配置
func (v *Validator) Validate(config *Config) error {
	if err := v.validateServer(config); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	if err := v.validateWebDAV(config); err != nil {
		return fmt.Errorf("webdav config: %w", err)
	}

	if err := v.validateWeb3(config); err != nil {
		return fmt.Errorf("web3 config: %w", err)
	}

	if err := v.validateUsers(config); err != nil {
		return fmt.Errorf("users config: %w", err)
	}

	return nil
}

// validateServer 验证服务器配置
func (v *Validator) validateServer(config *Config) error {
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return errors.New("invalid port number")
	}

	if config.Server.TLS {
		if config.Server.CertFile == "" {
			return errors.New("cert_file is required when TLS is enabled")
		}
		if config.Server.KeyFile == "" {
			return errors.New("key_file is required when TLS is enabled")
		}

		// 检查证书文件是否存在
		if _, err := os.Stat(config.Server.CertFile); err != nil {
			return fmt.Errorf("cert file not found: %w", err)
		}
		if _, err := os.Stat(config.Server.KeyFile); err != nil {
			return fmt.Errorf("key file not found: %w", err)
		}
	}

	return nil
}

// validateWebDAV 验证 WebDAV 配置
func (v *Validator) validateWebDAV(config *Config) error {
	if config.WebDAV.Directory == "" {
		return errors.New("directory is required")
	}

	// 检查目录是否存在
	info, err := os.Stat(config.WebDAV.Directory)
	if err != nil {
		return fmt.Errorf("directory not found: %w", err)
	}
	if !info.IsDir() {
		return errors.New("directory is not a directory")
	}

	return nil
}

// validateWeb3 验证 Web3 配置
func (v *Validator) validateWeb3(config *Config) error {
	if config.Web3.Enabled {
		if config.Web3.JWTSecret == "" {
			return errors.New("jwt_secret is required when web3 is enabled")
		}
		if len(config.Web3.JWTSecret) < 32 {
			return errors.New("jwt_secret must be at least 32 characters")
		}
	}

	return nil
}

// validateUsers 验证用户配置
func (v *Validator) validateUsers(config *Config) error {
	if len(config.Users) == 0 {
		return errors.New("at least one user is required")
	}

	usernames := make(map[string]bool)
	addresses := make(map[string]bool)

	for i, user := range config.Users {
		// 检查用户名
		if user.Username == "" {
			return fmt.Errorf("user[%d]: username is required", i)
		}
		if usernames[user.Username] {
			return fmt.Errorf("user[%d]: duplicate username: %s", i, user.Username)
		}
		usernames[user.Username] = true

		// 检查认证方式
		hasPassword := user.Password != ""
		hasWallet := user.WalletAddress != ""

		if !hasPassword && !hasWallet && !config.Security.NoPassword {
			return fmt.Errorf("user[%d]: must have password or wallet_address", i)
		}

		// 检查钱包地址唯一性
		if hasWallet {
			if addresses[user.WalletAddress] {
				return fmt.Errorf("user[%d]: duplicate wallet_address: %s", i, user.WalletAddress)
			}
			addresses[user.WalletAddress] = true
		}

		// 检查目录
		if user.Directory == "" {
			return fmt.Errorf("user[%d]: directory is required", i)
		}
	}

	return nil
}
