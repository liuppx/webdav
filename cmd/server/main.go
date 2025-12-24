package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/spf13/pflag"
	"github.com/yeying-community/webdav/internal/container"
	"github.com/yeying-community/webdav/internal/infrastructure/config"
	"go.uber.org/zap"
)

var (
	version   = "2.0.0"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// 解析命令行参数
	flags := parseFlags()
	
	// 显示版本信息
	if showVersion, _ := flags.GetBool("version"); showVersion {
		printVersion()
		os.Exit(0)
	}
	
	// 加载配置
	configFile, _ := flags.GetString("config")
	cfg, err := loadConfig(configFile, flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}
	
	// 创建容器
	c, err := container.NewContainer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create container: %v\n", err)
		os.Exit(1)
	}
	defer c.Close()
	
	// 打印启动信息
	printStartupInfo(c)
	
	// 启动服务器
	go func() {
		if err := c.Server.Start(); err != nil {
			c.Logger.Fatal("failed to start server", zap.Error(err))
		}
	}()
	
	// 等待中断信号
	waitForShutdown(c)
}

// parseFlags 解析命令行参数
func parseFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("webdav", pflag.ExitOnError)
	
	flags.StringP("config", "c", "", "Config file path")
	flags.String("address", "", "Server address")
	flags.IntP("port", "p", 0, "Server port")
	flags.Bool("tls", false, "Enable TLS")
	flags.String("cert", "", "TLS certificate file")
	flags.String("key", "", "TLS key file")
	flags.String("prefix", "", "WebDAV prefix")
	flags.StringP("directory", "d", "", "WebDAV directory")
	flags.BoolP("version", "v", false, "Show version")
	flags.BoolP("help", "h", false, "Show help")
	
	flags.Parse(os.Args[1:])
	
	if help, _ := flags.GetBool("help"); help {
		printHelp(flags)
		os.Exit(0)
	}
	
	return flags
}

// loadConfig 加载配置
func loadConfig(configFile string, flags *pflag.FlagSet) (*config.Config, error) {
	loader := config.NewLoader()
	return loader.Load(configFile, flags)
}

// printVersion 打印版本信息
func printVersion() {
	fmt.Printf("WebDAV Server\n")
	fmt.Printf("Version:    %s\n", version)
	fmt.Printf("Build Time: %s\n", buildTime)
	fmt.Printf("Git Commit: %s\n", gitCommit)
}

// printHelp 打印帮助信息
func printHelp(flags *pflag.FlagSet) {
	fmt.Println("WebDAV Server with Web3 Authentication")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  webdav [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	flags.PrintDefaults()
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Start with config file")
	fmt.Println("  webdav -c config.yaml")
	fmt.Println()
	fmt.Println("  # Start with command line flags")
	fmt.Println("  webdav -p 8080 -d /data")
	fmt.Println()
	fmt.Println("  # Start with TLS")
	fmt.Println("  webdav -c config.yaml --tls --cert cert.pem --key key.pem")
}

// printStartupInfo 打印启动信息
func printStartupInfo(c *container.Container) {
	c.Logger.Info("=================================")
	c.Logger.Info("WebDAV Server Starting")
	c.Logger.Info("=================================")
	c.Logger.Info("version", zap.String("version", version))
	c.Logger.Info("build_time", zap.String("build_time", buildTime))
	c.Logger.Info("git_commit", zap.String("git_commit", gitCommit))
	c.Logger.Info("=================================")
	c.Logger.Info("server",
		zap.String("address", c.Config.Server.Address),
		zap.Int("port", c.Config.Server.Port),
		zap.Bool("tls", c.Config.Server.TLS))
	c.Logger.Info("webdav",
		zap.String("prefix", c.Config.WebDAV.Prefix),
		zap.String("directory", c.Config.WebDAV.Directory))
	c.Logger.Info("web3",
		zap.Bool("enabled", c.Config.Web3.Enabled))
	c.Logger.Info("cors",
		zap.Bool("enabled", c.Config.CORS.Enabled))
	c.Logger.Info("=================================")
}

// waitForShutdown 等待关闭信号
func waitForShutdown(c *container.Container) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	sig := <-quit
	c.Logger.Info("received shutdown signal", zap.String("signal", sig.String()))
	
	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), c.Config.Server.ShutdownTimeout)
	defer cancel()
	
	if err := c.Server.Shutdown(ctx); err != nil {
		c.Logger.Error("failed to shutdown server", zap.Error(err))
	}
	
	c.Logger.Info("server stopped gracefully")
}

