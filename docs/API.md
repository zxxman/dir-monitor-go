# 目录监控服务 API 文档

## 目录

1. [概述](#概述)
2. [配置API](#配置api)
3. [监控API](#监控api)
4. [日志API](#日志api)
5. [事件API](#事件api)
6. [错误处理](#错误处理)
7. [代码示例](#代码示例)

## 概述

目录监控服务（dir-monitor-go）提供了一套完整的API接口，用于配置管理、文件监控、日志记录和事件处理。本文档详细介绍了这些API的使用方法、参数说明和返回值。

### API版本

当前API版本：v3.2.1

### 基本约定

- 所有API都遵循Go语言的命名约定
- 公开API以大写字母开头
- 错误返回遵循Go的错误处理模式
- 配置项使用JSON格式

## 配置API

配置API负责加载、验证和管理应用程序配置。

### Config 结构体

```go
type Config struct {
    Version  string    `json:"version"`
    Monitors []Monitor `json:"monitors"`
    Settings Settings  `json:"settings"`
}
```

#### 字段说明

- `Version`: 配置文件版本号
- `Monitors`: 监控器配置列表
- `Settings`: 系统设置

### LoadConfig 函数

```go
func LoadConfig(path string) (*Config, error)
```

加载配置文件并解析为Config结构体。

#### 参数

- `path`: 配置文件路径

#### 返回值

- `*Config`: 解析后的配置对象
- `error`: 错误信息，如果加载失败则返回非nil值

#### 示例

```go
config, err := config.LoadConfig("/path/to/config.json")
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

### Validate 方法

```go
func (c *Config) Validate() error
```

验证配置的有效性。

#### 返回值

- `error`: 验证错误，如果配置无效则返回非nil值

#### 验证规则

1. 至少配置一个监控器
2. 每个监控器的目录和命令不能为空
3. 文件模式必须有效
4. 设置值必须在有效范围内

#### 示例

```go
if err := config.Validate(); err != nil {
    log.Fatalf("Invalid config: %v", err)
}
```

### Monitor 结构体

```go
type Monitor struct {
    Directory    string   `json:"directory"`
    Command      string   `json:"command"`
    FilePatterns []string `json:"file_patterns"`
    Recursive    bool     `json:"recursive"`
    Events       []string `json:"events"`
    Debounce     Debounce `json:"debounce"`
    Filters      Filters  `json:"filters"`
}
```

#### 字段说明

- `Directory`: 要监控的目录路径
- `Command`: 文件变化时执行的命令
- `FilePatterns`: 文件模式匹配列表
- `Recursive`: 是否递归监控子目录
- `Events`: 要监听的事件类型列表
- `Debounce`: 防抖配置
- `Filters`: 过滤器配置

### Settings 结构体

```go
type Settings struct {
    LogLevel                        string `json:"log_level"`
    LogFile                         string `json:"log_file"`
    MaxConcurrentOperations         int    `json:"max_concurrent_operations"`
    OperationTimeoutSeconds         int    `json:"operation_timeout_seconds"`
    MinStabilityTimeMs              int    `json:"min_stability_time_ms"`
    ExecutionDedupIntervalSeconds   int    `json:"execution_dedup_interval_seconds"`
    DirectoryStabilityQuietMs       int    `json:"directory_stability_quiet_ms"`
    DirectoryStabilityTimeoutSeconds int   `json:"directory_stability_timeout_seconds"`
    RetryAttempts                   int    `json:"retry_attempts"`
    RetryDelaySeconds               int    `json:"retry_delay_seconds"`
    HealthCheckIntervalSeconds      int    `json:"health_check_interval_seconds"`
}
```

#### 字段说明

- `LogLevel`: 日志级别 (debug, info, warn, error)
- `LogFile`: 日志文件路径
- `MaxConcurrentOperations`: 最大并发操作数
- `OperationTimeoutSeconds`: 操作超时时间（秒）
- `MinStabilityTimeMs`: 最小稳定性时间（毫秒）
- `ExecutionDedupIntervalSeconds`: 执行去重间隔（秒）
- `DirectoryStabilityQuietMs`: 目录稳定性静默时间（毫秒）
- `DirectoryStabilityTimeoutSeconds`: 目录稳定性超时时间（秒）
- `RetryAttempts`: 重试次数
- `RetryDelaySeconds`: 重试延迟（秒）
- `HealthCheckIntervalSeconds`: 健康检查间隔（秒）

## 监控API

监控API负责文件系统监控的核心功能。

### Monitor 结构体

```go
type Monitor struct {
    config       *config.Config
    logger       *logger.Logger
    watchers     map[string]*fsnotify.Watcher
    eventChan    chan fsnotify.Event
    stopChan     chan struct{}
    wg           sync.WaitGroup
    sem          chan struct{}
}
```

#### 字段说明

- `config`: 配置对象
- `logger`: 日志记录器
- `watchers`: 文件监视器映射
- `eventChan`: 事件通道
- `stopChan`: 停止信号通道
- `wg`: 等待组，用于goroutine同步
- `sem`: 信号量，用于并发控制

### NewMonitor 函数

```go
func NewMonitor(cfg *config.Config, log *logger.Logger) *Monitor
```

创建一个新的文件系统监控器。

#### 参数

- `cfg`: 配置对象
- `log`: 日志记录器

#### 返回值

- `*Monitor`: 新创建的监控器对象

#### 示例

```go
monitor := monitor.NewMonitor(config, logger)
```

### Start 方法

```go
func (m *Monitor) Start() error
```

启动文件系统监控。

#### 返回值

- `error`: 启动错误，如果启动失败则返回非nil值

#### 示例

```go
if err := monitor.Start(); err != nil {
    log.Fatalf("Failed to start monitor: %v", err)
}
```

### Stop 方法

```go
func (m *Monitor) Stop()
```

停止文件系统监控。

#### 示例

```go
monitor.Stop()
```

### startWatching 方法

```go
func (m *Monitor) startWatching(monitorConfig config.Monitor) error
```

启动对指定目录的监控。

#### 参数

- `monitorConfig`: 监控器配置

#### 返回值

- `error`: 启动错误，如果启动失败则返回非nil值

### handleEvent 方法

```go
func (m *Monitor) handleEvent(event fsnotify.Event, monitorConfig config.Monitor)
```

处理文件系统事件。

#### 参数

- `event`: 文件系统事件
- `monitorConfig`: 监控器配置

## 日志API

日志API提供结构化日志记录功能。

### Logger 结构体

```go
type Logger struct {
    logger *zap.Logger
}
```

#### 字段说明

- `logger`: Zap日志记录器

### NewLogger 函数

```go
func NewLogger(config LoggerConfig) *Logger
```

创建一个新的日志记录器。

#### 参数

- `config`: 日志配置

#### 返回值

- `*Logger`: 新创建的日志记录器

#### 示例

```go
logConfig := logger.LoggerConfig{
    Level:      "info",
    OutputPath: "/var/log/dir-monitor-go.log",
    MaxSize:    100, // MB
    MaxBackups: 3,
    MaxAge:     28, // days
    Compress:   true,
}
logger := logger.NewLogger(logConfig)
```

### LoggerConfig 结构体

```go
type LoggerConfig struct {
    Level      string
    OutputPath string
    MaxSize    int
    MaxBackups int
    MaxAge     int
    Compress   bool
}
```

#### 字段说明

- `Level`: 日志级别
- `OutputPath`: 日志文件路径
- `MaxSize`: 日志文件最大大小（MB）
- `MaxBackups`: 保留的旧日志文件数量
- `MaxAge`: 日志文件保留天数
- `Compress`: 是否压缩旧日志文件

### Info 方法

```go
func (l *Logger) Info(msg string, fields ...Field)
```

记录信息级别日志。

#### 参数

- `msg`: 日志消息
- `fields`: 日志字段

#### 示例

```go
logger.Info("File changed", 
    logger.String("file", "/path/to/file.txt"),
    logger.String("operation", "write"))
```

### Error 方法

```go
func (l *Logger) Error(msg string, fields ...Field)
```

记录错误级别日志。

#### 参数

- `msg`: 日志消息
- `fields`: 日志字段

#### 示例

```go
logger.Error("Failed to execute command",
    logger.String("command", "process-file.sh"),
    logger.Error(err))
```

### Field 类型

```go
type Field struct {
    Key   string
    Value interface{}
}
```

#### 字段说明

- `Key`: 字段名
- `Value`: 字段值

### 辅助函数

```go
func String(key, val string) Field
func Int(key string, val int) Field
func Bool(key string, val bool) Field
func Error(err error) Field
```

创建特定类型的日志字段。

## 事件API

事件API定义了文件系统事件的数据结构。

### Event 结构体

```go
type Event struct {
    Name      string    `json:"name"`
    Operation string    `json:"operation"`
    Timestamp time.Time `json:"timestamp"`
    Size      int64     `json:"size"`
    IsDir     bool      `json:"is_dir"`
}
```

#### 字段说明

- `Name`: 文件或目录名
- `Operation`: 操作类型 (create, write, remove, rename, chmod)
- `Timestamp`: 事件时间戳
- `Size`: 文件大小（字节）
- `IsDir`: 是否为目录

### EventType 类型

```go
type EventType string
```

#### 常量

```go
const (
    EventCreate = EventType("create")
    EventWrite  = EventType("write")
    EventRemove = EventType("remove")
    EventRename = EventType("rename")
    EventChmod  = EventType("chmod")
)
```

## 错误处理

API遵循Go的错误处理模式，大多数函数返回一个error类型的值作为最后一个返回值。

### 错误类型

```go
type ConfigError struct {
    Field   string
    Message string
}

type MonitorError struct {
    Operation string
    Path      string
    Err       error
}
```

### 错误检查

```go
config, err := config.LoadConfig("config.json")
if err != nil {
    // 处理错误
    switch e := err.(type) {
    case *config.ConfigError:
        log.Printf("Config error in field %s: %s", e.Field, e.Message)
    case *os.PathError:
        log.Printf("File error: %s", e.Err)
    default:
        log.Printf("Unknown error: %v", err)
    }
    return
}
```

## 代码示例

### 完整示例

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/zxxman/dir-monitor-go/internal/config"
    "github.com/zxxman/dir-monitor-go/internal/logger"
    "github.com/zxxman/dir-monitor-go/internal/monitor"
)

func main() {
    // 加载配置
    cfg, err := config.LoadConfig("config.json")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 验证配置
    if err := cfg.Validate(); err != nil {
        log.Fatalf("Invalid config: %v", err)
    }
    
    // 创建日志记录器
    logConfig := logger.LoggerConfig{
        Level:      cfg.Settings.LogLevel,
        OutputPath: cfg.Settings.LogFile,
        MaxSize:    100,
        MaxBackups: 3,
        MaxAge:     28,
        Compress:   true,
    }
    logger := logger.NewLogger(logConfig)
    
    // 创建监控器
    mon := monitor.NewMonitor(cfg, logger)
    
    // 启动监控
    if err := mon.Start(); err != nil {
        logger.Error("Failed to start monitor", logger.Error(err))
        os.Exit(1)
    }
    
    logger.Info("Monitor started successfully")
    
    // 等待中断信号
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    // 停止监控
    logger.Info("Shutting down monitor...")
    mon.Stop()
    logger.Info("Monitor stopped")
}
```

### 自定义事件处理

```go
// 创建自定义事件处理器
func customEventHandler(event fsnotify.Event, monitorConfig config.Monitor) {
    // 检查文件模式
    if !matchPatterns(event.Name, monitorConfig.FilePatterns) {
        return
    }
    
    // 检查事件类型
    if !containsEvent(event.Op.String(), monitorConfig.Events) {
        return
    }
    
    // 执行自定义命令
    cmd := exec.Command(monitorConfig.Command, event.Name)
    if err := cmd.Run(); err != nil {
        log.Printf("Failed to execute command: %v", err)
    }
}

// 辅助函数：检查文件是否匹配模式
func matchPatterns(filename string, patterns []string) bool {
    for _, pattern := range patterns {
        matched, err := filepath.Match(pattern, filepath.Base(filename))
        if err != nil {
            continue
        }
        if matched {
            return true
        }
    }
    return false
}

// 辅助函数：检查事件类型是否在列表中
func containsEvent(eventType string, events []string) bool {
    for _, e := range events {
        if e == eventType {
            return true
        }
    }
    return false
}
```

### 动态配置更新

```go
// 监听配置文件变化并重新加载配置
func watchConfig(configPath string, monitor *monitor.Monitor, logger *logger.Logger) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        logger.Error("Failed to create config watcher", logger.Error(err))
        return
    }
    defer watcher.Close()
    
    if err := watcher.Add(configPath); err != nil {
        logger.Error("Failed to watch config file", logger.Error(err))
        return
    }
    
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            
            if event.Op&fsnotify.Write == fsnotify.Write {
                logger.Info("Config file modified, reloading...")
                
                // 重新加载配置
                newCfg, err := config.LoadConfig(configPath)
                if err != nil {
                    logger.Error("Failed to reload config", logger.Error(err))
                    continue
                }
                
                if err := newCfg.Validate(); err != nil {
                    logger.Error("Invalid config after reload", logger.Error(err))
                    continue
                }
                
                // 停止当前监控
                monitor.Stop()
                
                // 创建新监控器
                newMonitor := monitor.NewMonitor(newCfg, logger)
                
                // 启动新监控
                if err := newMonitor.Start(); err != nil {
                    logger.Error("Failed to start new monitor", logger.Error(err))
                    continue
                }
                
                // 更新监控器引用
                *monitor = *newMonitor
                
                logger.Info("Config reloaded successfully")
            }
            
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            logger.Error("Config watcher error", logger.Error(err))
        }
    }
}
```

---

如有更多问题，请查看[用户指南](USER_GUIDE.md)或[开发文档](DEVELOPMENT.md)。