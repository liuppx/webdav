package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// Loader 配置加载器
type Loader struct {
	defaultConfig *Config
}

// NewLoader 创建配置加载器
func NewLoader() *Loader {
	return &Loader{
		defaultConfig: DefaultConfig(),
	}
}

// Load 加载配置
func (l *Loader) Load(configFile string, flags *pflag.FlagSet) (*Config, error) {
	config := l.defaultConfig

	// 1. 从文件加载
	if configFile != "" {
		if err := l.loadFromFile(configFile, config); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	// 2. 从命令行参数覆盖
	if flags != nil {
		l.overrideFromFlags(config, flags)
	}

	// 3. 从环境变量覆盖
	l.overrideFromEnv(config)

	// 4. 验证配置
	if err := l.validate(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}

// loadFromFile 从文件加载配置
func (l *Loader) loadFromFile(filename string, config *Config) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

// overrideFromFlags 从命令行参数覆盖配置
func (l *Loader) overrideFromFlags(config *Config, flags *pflag.FlagSet) {
	if flags.Changed("address") {
		config.Server.Address, _ = flags.GetString("address")
	}
	if flags.Changed("port") {
		config.Server.Port, _ = flags.GetInt("port")
	}
	if flags.Changed("tls") {
		config.Server.TLS, _ = flags.GetBool("tls")
	}
	if flags.Changed("cert") {
		config.Server.CertFile, _ = flags.GetString("cert")
	}
	if flags.Changed("key") {
		config.Server.KeyFile, _ = flags.GetString("key")
	}
	if flags.Changed("prefix") {
		config.WebDAV.Prefix, _ = flags.GetString("prefix")
	}
	if flags.Changed("directory") {
		config.WebDAV.Directory, _ = flags.GetString("directory")
	}
}

// overrideFromEnv 从环境变量覆盖配置
func (l *Loader) overrideFromEnv(config *Config) {
	if v := os.Getenv("WEBDAV_ADDRESS"); v != "" {
		config.Server.Address = v
	}
	if v := os.Getenv("WEBDAV_PORT"); v != "" {
		// 解析端口...
	}
	if v := os.Getenv("WEBDAV_JWT_SECRET"); v != "" {
		config.Web3.JWTSecret = v
	}
}

// validate 验证配置
func (l *Loader) validate(config *Config) error {
	validator := NewValidator()
	return validator.Validate(config)
}
