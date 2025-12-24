package container

import (
	"fmt"

	"github.com/yeying-community/webdav/internal/application/service"
	"github.com/yeying-community/webdav/internal/domain/auth"
	infraAuth "github.com/yeying-community/webdav/internal/infrastructure/auth"
	"github.com/yeying-community/webdav/internal/infrastructure/config"
	"github.com/yeying-community/webdav/internal/infrastructure/logger"
	"github.com/yeying-community/webdav/internal/infrastructure/permission"
	"github.com/yeying-community/webdav/internal/infrastructure/repository"
	"github.com/yeying-community/webdav/internal/interface/http"
	"github.com/yeying-community/webdav/internal/interface/http/handler"
	"go.uber.org/zap"
	"golang.org/x/net/webdav"
)

// Container 依赖注入容器
type Container struct {
	Config *config.Config
	Logger *zap.Logger

	// Repositories
	UserRepo *repository.MemoryUserRepository

	// Authenticators
	Authenticators []auth.Authenticator
	BasicAuth      *infraAuth.BasicAuthenticator
	Web3Auth       *infraAuth.Web3Authenticator

	// Services
	WebDAVService *service.WebDAVService

	// Handlers
	HealthHandler *handler.HealthHandler
	Web3Handler   *handler.Web3Handler
	WebDAVHandler *handler.WebDAVHandler

	// HTTP
	Router *http.Router
	Server *http.Server
}

// NewContainer 创建容器
func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{
		Config:         cfg,
		Authenticators: make([]auth.Authenticator, 0),
	}

	// 初始化组件
	if err := c.initLogger(); err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	if err := c.initRepositories(); err != nil {
		return nil, fmt.Errorf("failed to init repositories: %w", err)
	}

	if err := c.initAuthenticators(); err != nil {
		return nil, fmt.Errorf("failed to init authenticators: %w", err)
	}

	if err := c.initServices(); err != nil {
		return nil, fmt.Errorf("failed to init services: %w", err)
	}

	if err := c.initHandlers(); err != nil {
		return nil, fmt.Errorf("failed to init handlers: %w", err)
	}

	if err := c.initHTTP(); err != nil {
		return nil, fmt.Errorf("failed to init http: %w", err)
	}

	return c, nil
}

// initLogger 初始化日志器
func (c *Container) initLogger() error {
	l, err := logger.NewLogger(c.Config.Log)
	if err != nil {
		return err
	}

	c.Logger = l
	c.Logger.Info("logger initialized",
		zap.String("level", c.Config.Log.Level),
		zap.String("format", c.Config.Log.Format))

	return nil
}

// initRepositories 初始化仓储
func (c *Container) initRepositories() error {
	c.UserRepo = repository.NewMemoryUserRepository(c.Config.Users)

	c.Logger.Info("repositories initialized",
		zap.Int("users", len(c.Config.Users)))

	return nil
}

// initAuthenticators 初始化认证器
func (c *Container) initAuthenticators() error {
	// Basic 认证器
	c.BasicAuth = infraAuth.NewBasicAuthenticator(
		c.UserRepo,
		c.Config.Security.NoPassword,
		c.Logger,
	)
	c.Authenticators = append(c.Authenticators, c.BasicAuth)

	// Web3 认证器
	if c.Config.Web3.Enabled {
		c.Web3Auth = infraAuth.NewWeb3Authenticator(
			c.UserRepo,
			c.Config.Web3.JWTSecret,
			c.Config.Web3.TokenExpiration,
			c.Logger,
		)
		c.Authenticators = append(c.Authenticators, c.Web3Auth)

		c.Logger.Info("web3 authentication enabled",
			zap.Duration("token_expiration", c.Config.Web3.TokenExpiration))
	}

	c.Logger.Info("authenticators initialized",
		zap.Int("count", len(c.Authenticators)))

	return nil
}

// initServices 初始化服务
func (c *Container) initServices() error {
	// WebDAV 服务
	fileSystem := webdav.Dir(c.Config.WebDAV.Directory)
	permissionChecker := permission.NewWebDAVChecker(fileSystem, c.Logger)

	c.WebDAVService = service.NewWebDAVService(
		c.Config,
		permissionChecker,
		c.Logger,
	)

	c.Logger.Info("services initialized")

	return nil
}

// initHandlers 初始化处理器
func (c *Container) initHandlers() error {
	// 健康检查处理器
	c.HealthHandler = handler.NewHealthHandler(c.Logger)

	// Web3 处理器
	if c.Web3Auth != nil {
		c.Web3Handler = handler.NewWeb3Handler(
			c.Web3Auth,
			c.UserRepo,
			c.Logger,
		)
	}

	// WebDAV 处理器
	c.WebDAVHandler = handler.NewWebDAVHandler(c.WebDAVService, c.Logger)

	c.Logger.Info("handlers initialized")

	return nil
}

// initHTTP 初始化 HTTP
func (c *Container) initHTTP() error {
	// 路由器
	c.Router = http.NewRouter(
		c.Config,
		c.Authenticators,
		c.HealthHandler,
		c.Web3Handler,
		c.WebDAVHandler,
		c.Logger,
	)

	// 服务器
	c.Server = http.NewServer(c.Config, c.Router, c.Logger)

	c.Logger.Info("http components initialized")

	return nil
}

// Close 关闭容器
func (c *Container) Close() error {
	if c.Logger != nil {
		c.Logger.Info("closing container")
		_ = c.Logger.Sync()
	}

	return nil
}
