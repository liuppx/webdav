package config

import (
	"time"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	WebDAV   WebDAVConfig   `yaml:"webdav"`
	Web3     Web3Config     `yaml:"web3"`
	Security SecurityConfig `yaml:"security"`
	CORS     CORSConfig     `yaml:"cors"`
	Log      LogConfig      `yaml:"log"`
	Users    []UserConfig   `yaml:"users"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Address         string        `yaml:"address"`
	Port            int           `yaml:"port"`
	TLS             bool          `yaml:"tls"`
	CertFile        string        `yaml:"cert_file"`
	KeyFile         string        `yaml:"key_file"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// WebDAVConfig WebDAV 配置
type WebDAVConfig struct {
	Prefix      string `yaml:"prefix"`
	Directory   string `yaml:"directory"`
	NoSniff     bool   `yaml:"no_sniff"`
	Permissions string `yaml:"permissions"`
}

// Web3Config Web3 配置
type Web3Config struct {
	Enabled         bool          `yaml:"enabled"`
	JWTSecret       string        `yaml:"jwt_secret"`
	TokenExpiration time.Duration `yaml:"token_expiration"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	NoPassword   bool `yaml:"no_password"`
	BehindProxy  bool `yaml:"behind_proxy"`
}

// CORSConfig CORS 配置
type CORSConfig struct {
	Enabled        bool     `yaml:"enabled"`
	Credentials    bool     `yaml:"credentials"`
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
	ExposedHeaders []string `yaml:"exposed_headers"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level   string   `yaml:"level"`
	Format  string   `yaml:"format"`
	Colors  bool     `yaml:"colors"`
	Outputs []string `yaml:"outputs"`
}

// UserConfig 用户配置
type UserConfig struct {
	Username      string       `yaml:"username"`
	Password      string       `yaml:"password"`
	WalletAddress string       `yaml:"wallet_address"`
	Directory     string       `yaml:"directory"`
	Permissions   string       `yaml:"permissions"`
	Rules         []RuleConfig `yaml:"rules"`
}

// RuleConfig 规则配置
type RuleConfig struct {
	Path        string `yaml:"path"`
	Permissions string `yaml:"permissions"`
	Regex       bool   `yaml:"regex"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Address:         "0.0.0.0",
			Port:            6065,
			TLS:             false,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			IdleTimeout:     60 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		WebDAV: WebDAVConfig{
			Prefix:      "/",
			Directory:   "/data",
			NoSniff:     true,
			Permissions: "R",
		},
		Web3: Web3Config{
			Enabled:         false,
			TokenExpiration: 24 * time.Hour,
		},
		Security: SecurityConfig{
			NoPassword:  false,
			BehindProxy: false,
		},
		CORS: CORSConfig{
			Enabled:     false,
			Credentials: false,
		},
		Log: LogConfig{
			Level:   "info",
			Format:  "console",
			Colors:  true,
			Outputs: []string{"stderr"},
		},
		Users: []UserConfig{},
	}
}

