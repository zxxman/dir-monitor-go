package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"dir-monitor-go/internal/config"
	"dir-monitor-go/internal/logger"
	"dir-monitor-go/internal/monitor"
)

const (
	// 默认停止文件检查间隔
	DefaultStopFileCheckInterval = 200 * time.Millisecond

	// 默认日志文件大小（10MB）
	DefaultLogMaxSize = 10 * 1024 * 1024
)

// 版本信息（将由构建参数注入）
var (
	Version    = "{version}"     // 构建时注入版本号
	BuildTime  = "{build_time}"  // 构建时注入构建时间
	CommitHash = "{commit_hash}" // 构建时注入提交哈希
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "configs/config.json", "配置文件路径")
	stopFile := flag.String("stop-file", "", "当该文件出现时优雅退出（测试/集成用）")

	showVersion := flag.Bool("version", false, "显示版本信息")
	dryRun := flag.Bool("dry-run", false, "仅验证配置，不启动实际监控")
	flag.Parse()

	// 显示版本信息
	if *showVersion {
		fmt.Printf("dir-monitor-go 版本 %s\n", Version)
		fmt.Printf("构建时间: %s\n", BuildTime)
		if CommitHash != "" && CommitHash != "{commit_hash}" {
			fmt.Printf("提交哈希: %s\n", CommitHash)
		}
		return
	}

	// 初始化日志
	var log *logger.Logger
	var err error

	// 使用控制台日志（启动期固定 INFO，配置加载后会重建）
	log = logger.NewLogger(logger.INFO, os.Stdout)
	// 确保退出前关闭日志器（文件日志器可安全释放句柄）
	defer func() {
		if log != nil {
			_ = log.Close()
		}
	}()

	log.Info("Dir-Monitor-Go %s 启动中...", Version)

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Error("加载配置文件失败: %v", err)
		safeExit(log, 1)
	}

	// 配置加载后，按配置重建日志器（忽略命令行 log-level/log-file）
	{
		// 解析级别（优先 settings，其次顶层）
		cfgLevel := strings.ToLower(cfg.Settings.LogLevel)
		if cfgLevel == "" {
			cfgLevel = strings.ToLower(cfg.LogLevel)
		}
		level := logger.INFO
		switch cfgLevel {
		case "debug":
			level = logger.DEBUG
		case "warn":
			level = logger.WARN
		case "error":
			level = logger.ERROR
		default:
			level = logger.INFO
		}

		// 解析日志文件（优先 settings，其次顶层）
		cfgLogFile := cfg.Settings.LogFile
		if cfgLogFile == "" {
			cfgLogFile = cfg.LogFile
		}

		// 解析是否显示调用者信息
		logShowCaller := cfg.Settings.LogShowCaller

		// 解析日志文件最大大小
		logMaxSize := cfg.Settings.LogMaxSize
		if logMaxSize <= 0 {
			logMaxSize = DefaultLogMaxSize
		}

		// 重新创建日志器
		if strings.TrimSpace(cfgLogFile) != "" {
			// 创建文件日志器
			newLog, err2 := logger.NewFileLogger(level, cfgLogFile, logMaxSize)
			if err2 != nil {
				fmt.Fprintf(os.Stderr, "初始化文件日志失败: %v", err2)
				safeExit(log, 1)
			}
			log = newLog
		} else {
			// 创建控制台日志器
			log = logger.NewLogger(level, os.Stdout)
		}

		// 应用日志配置
		log.SetCaller(logShowCaller)

		// 提示最终日志配置
		levelName := map[logger.LogLevel]string{logger.DEBUG: "debug", logger.INFO: "info", logger.WARN: "warn", logger.ERROR: "error"}[level]
		log.Info("使用配置的日志级别: %s", levelName)
		if strings.TrimSpace(cfgLogFile) != "" {
			log.Info("使用配置的日志文件: %s", cfgLogFile)
		}
		// Debug 验证：确认 debug 级别是否已生效
		log.Debug("Debug 级别已生效（验证输出）")
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Error("配置验证失败: %v", err)
		os.Exit(1)
	}

	// 显示配置统计信息
	log.Info("配置统计: 监控项 %d 个", len(cfg.Monitors))

	// 如果是测试模式，验证后退出
	if *dryRun {
		log.Info("[测试模式] 配置验证通过，程序验证完成，退出")
		return
	}

	// 创建监控管理器
	monitorManager := monitor.NewMonitorManager(log)

	// 创建并添加监控器
	monitorInstance, err := monitor.NewMonitor(cfg, log)
	if err != nil {
		log.Error("创建监控器失败: %v", err)
		safeExit(log, 1)
	}
	monitorManager.AddMonitor(monitorInstance)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动监控器
	if err := monitorManager.Start(ctx); err != nil {
		log.Error("启动监控器失败: %v", err)
		return
	}
	log.Info("[生产模式] 监控器启动成功")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 若指定了停止标记文件，则后台轮询（仅测试/集成用）
	if strings.TrimSpace(*stopFile) != "" {
		sf := strings.TrimSpace(*stopFile)
		go func() {
			t := time.NewTicker(DefaultStopFileCheckInterval)
			defer t.Stop()
			for range t.C {
				if _, err := os.Stat(sf); err == nil {
					log.Info("检测到停止标记文件: %s，触发优雅退出", sf)
					quit <- syscall.SIGTERM
					return
				}
			}
		}()
	}

	<-quit

	log.Info("[生产模式] 收到退出信号，正在关闭监控服务...")

	// 优雅关闭监控服务
	monitorManager.Stop(ctx)

	log.Info("[生产模式] 监控服务已关闭，程序退出")
}

// safeExit 在退出前尝试关闭日志器，避免丢失或占用文件句柄
func safeExit(log *logger.Logger, code int) {
	if log != nil {
		_ = log.Close()
	}
	os.Exit(code)
}
