# 目录监控服务开发文档

## 目录

1. [项目概述](#项目概述)
2. [开发环境搭建](#开发环境搭建)
3. [项目结构](#项目结构)
4. [核心模块](#核心模块)
5. [开发流程](#开发流程)
6. [代码规范](#代码规范)
7. [测试指南](#测试指南)
8. [构建与发布](#构建与发布)
9. [贡献指南](#贡献指南)

## 项目概述

目录监控服务（dir-monitor-go）是一个使用Go语言开发的文件系统监控工具，它可以监控指定目录的变化并触发自定义命令。项目采用模块化设计，具有良好的可扩展性和可维护性。

### 主要特性

- 实时文件系统监控
- 灵活的配置管理
- 并发操作控制
- 文件稳定性检查
- 详细的日志记录
- 多种部署方式

### 技术栈

- **语言**: Go 1.25+
- **主要依赖**: 
  - `github.com/fsnotify/fsnotify v1.9.0` - 文件系统监控
  - `github.com/adhocore/gronx v1.19.6` - Cron表达式解析
  - `golang.org/x/sys v0.37.0` - 系统调用

## 开发环境搭建

### 前置要求

- Go 1.25 或更高版本
- Git
- Make（可选，用于构建）

### 克隆项目

```bash
git clone https://github.com/zxxman/dir-monitor-go.git
cd dir-monitor-go
```

### 安装依赖

```bash
go mod download
```

### 验证环境

```bash
go version
go mod verify
```

### 开发工具推荐

- **IDE**: VS Code, GoLand, Vim
- **插件**: 
  - Go官方插件
  - golangci-lint（代码检查）
  - Delve（调试器）

## 项目结构

```
dir-monitor-go/
├── cmd/
│   └── dir-monitor-go/        # 主程序入口
│       └── main.go
├── internal/                  # 内部包
│   ├── config/               # 配置管理
│   │   └── config.go
│   ├── logger/               # 日志管理
│   │   └── logger.go
│   ├── model/                # 数据模型
│   │   ├── event.go
│   │   └── settings.go
│   └── monitor/              # 监控核心逻辑
│       └── monitor.go
├── configs/                  # 配置文件
│   └── config.json.example
├── deploy/                   # 部署相关
│   └── dir-monitor-go.service
├── docs/                     # 文档
├── scripts/                  # 脚本文件
├── .gitignore
├── CHANGELOG.md
├── LICENSE
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

## 核心模块

### 1. 配置管理 (internal/config)

配置管理模块负责加载和验证配置文件，提供配置访问接口。

主要结构体：
```go
type Config struct {
    Version  string    `json:"version"`
    Monitors []Monitor `json:"monitors"`
    Settings Settings  `json:"settings"`
}

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

主要方法：
- `LoadConfig(path string) (*Config, error)` - 加载配置文件
- `Validate() error` - 验证配置有效性
- `applyDefaults()` - 应用默认值

### 2. 日志管理 (internal/logger)

日志管理模块提供结构化日志记录功能，支持多种日志级别和输出方式。

主要功能：
- 多级别日志（debug, info, warn, error）
- 日志轮转
- 调用者信息记录
- 自定义格式化

主要方法：
- `NewLogger(config LoggerConfig) *Logger` - 创建日志器
- `Info(msg string, fields ...Field)` - 记录信息日志
- `Error(msg string, fields ...Field)` - 记录错误日志

### 3. 监控核心 (internal/monitor)

监控核心模块是项目的核心，负责文件系统监控和事件处理。

主要结构体：
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

主要方法：
- `NewMonitor(cfg *config.Config, log *logger.Logger) *Monitor` - 创建监控器
- `Start() error` - 启动监控
- `Stop()` - 停止监控
- `startWatching(monitor config.Monitor) error` - 启动目录监控
- `handleEvent(event fsnotify.Event, monitor config.Monitor)` - 处理文件事件

### 4. 数据模型 (internal/model)

数据模型模块定义了项目中使用的数据结构。

主要结构体：
- `Event` - 文件事件
- `Settings` - 系统设置
- `Debounce` - 防抖配置
- `Filters` - 过滤器配置

## 开发流程

### 1. 创建功能分支

```bash
git checkout -b feature/new-feature
```

### 2. 开发与测试

- 编写代码
- 添加单元测试
- 运行测试确保通过

### 3. 代码检查

```bash
# 运行代码格式化
go fmt ./...

# 运行代码检查
golangci-lint run

# 运行测试
go test -v ./...
```

### 4. 提交代码

```bash
git add .
git commit -m "feat: add new feature"
```

### 5. 推送与创建PR

```bash
git push origin feature/new-feature
```

在GitHub上创建Pull Request。

## 代码规范

### 1. 命名规范

- 包名：小写，简短，有意义
- 变量名：驼峰命名法
- 常量名：大写，下划线分隔
- 函数名：驼峰命名法，导出函数首字母大写

### 2. 注释规范

- 公开函数必须有注释
- 注释应以函数名开头
- 注释应说明函数的目的、参数和返回值

示例：
```go
// NewMonitor creates a new file system monitor with the given configuration.
// It initializes the logger, event channels, and semaphore for concurrent control.
func NewMonitor(cfg *config.Config, log *logger.Logger) *Monitor {
    // Implementation...
}
```

### 3. 错误处理

- 使用明确的错误变量
- 错误信息应包含上下文
- 避免忽略错误

示例：
```go
// 不好的做法
file, _ := os.Open(filename)

// 好的做法
file, err := os.Open(filename)
if err != nil {
    return fmt.Errorf("failed to open file %s: %w", filename, err)
}
```

### 4. 并发安全

- 使用互斥锁保护共享数据
- 使用通道进行通信
- 避免数据竞争

## 测试指南

### 1. 单元测试

单元测试应覆盖所有主要功能。测试文件应以`_test.go`结尾。

示例：
```go
func TestLoadConfig(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        want    *config.Config
        wantErr bool
    }{
        {
            name:    "valid config",
            path:    "testdata/valid_config.json",
            want:    &config.Config{Version: "3.2.1"},
            wantErr: false,
        },
        {
            name:    "invalid config",
            path:    "testdata/invalid_config.json",
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := config.LoadConfig(tt.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 2. 集成测试

集成测试应测试组件之间的交互。

### 3. 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/config

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 4. 基准测试

```go
func BenchmarkMonitorStart(b *testing.B) {
    cfg := &config.Config{
        Monitors: []config.Monitor{
            {
                Directory: "/tmp",
                Command:   "echo test",
            },
        },
    }
    log := logger.NewLogger(logger.Config{Level: "info"})
    monitor := monitor.NewMonitor(cfg, log)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        monitor.Start()
        monitor.Stop()
    }
}
```

## 构建与发布

### 1. 本地构建

```bash
# 使用Makefile
make build

# 手动构建
go build -o dir-monitor-go ./cmd/dir-monitor-go
```

### 2. 交叉编译

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o dir-monitor-go-linux-amd64 ./cmd/dir-monitor-go

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o dir-monitor-go-windows-amd64.exe ./cmd/dir-monitor-go

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o dir-monitor-go-darwin-amd64 ./cmd/dir-monitor-go
```

### 3. 构建标签

使用构建标签控制不同平台的代码：

```go
// +build linux

package monitor

import "syscall"

func getInode(path string) (uint64, error) {
    var stat syscall.Stat_t
    err := syscall.Stat(path, &stat)
    if err != nil {
        return 0, err
    }
    return stat.Ino, nil
}
```

### 4. 发布流程

1. 更新版本号
2. 更新CHANGELOG.md
3. 创建Git标签
4. 构建多平台二进制文件
5. 创建GitHub Release

## 贡献指南

### 1. 报告问题

- 使用GitHub Issues
- 提供详细的问题描述
- 包含复现步骤
- 附上相关日志和配置

### 2. 提交代码

1. Fork项目
2. 创建功能分支
3. 编写代码和测试
4. 确保所有测试通过
5. 提交Pull Request

### 3. 代码审查

- 所有代码必须经过审查
- 确保代码符合项目规范
- 检查测试覆盖率
- 验证功能正确性

### 4. 文档更新

- 更新相关文档
- 添加新功能的示例
- 更新CHANGELOG.md

## 性能优化

### 1. 监控性能

使用内置性能监控功能：

```go
// 在配置中启用性能监控
"performance_monitoring": {
    "enabled": true,
    "report_interval": "1m"
}
```

### 2. 优化建议

- 使用缓冲通道减少锁竞争
- 合理设置并发操作数
- 使用文件稳定性检查避免频繁触发
- 优化文件模式匹配算法

### 3. 内存管理

- 及时释放不再使用的资源
- 避免内存泄漏
- 使用对象池减少GC压力

## 故障排除

### 1. 常见问题

- **文件监控不工作**: 检查目录权限和文件系统类型
- **高CPU使用**: 检查监控目录大小和文件变化频率
- **内存泄漏**: 使用pprof工具分析内存使用

### 2. 调试工具

- 使用delve进行调试
- 使用pprof进行性能分析
- 启用详细日志记录

```bash
# 启用调试模式
./dir-monitor-go -l debug

# 使用pprof
go tool pprof http://localhost:6060/debug/pprof/profile
```

---

如有更多问题，请查看[用户指南](USER_GUIDE.md)或提交[Issue](https://github.com/zxxman/dir-monitor-go/issues)。