# 目录监控系统开发指南

## 系统架构

目录监控系统采用模块化设计，主要包括以下核心组件：

### 核心模块

1. **配置管理模块 (config)**
   - 负责配置文件的加载、解析和验证
   - 支持配置版本兼容性检查

2. **文件监控模块 (monitor)**
   - 基于fsnotify实现文件系统事件监听
   - 支持目录递归监控和文件模式匹配
   - 实现文件稳定性检测机制
   - 集成调度功能，控制操作执行的时间窗口

3. **日志模块 (logger)**
   - 基于zap实现结构化日志记录
   - 支持日志轮转和压缩
   - 提供不同级别的日志输出

4. **数据模型模块 (model)**
   - 定义系统的数据结构
   - 包括文件事件和配置设置等模型

### 数据流

```
[配置文件] --> [配置管理] --> [文件监控] --> [调度控制] --> [执行引擎] --> [操作执行]
     ^              |              |              |              |              |
     |              v              v              v              v              v
     |         [日志记录]    [事件处理]    [时间检查]    [并发控制]    [环境变量]
     |              |              |              |              |              |
     |              v              v              v              v              v
     |         [健康检查]    [稳定性检测]   [调度决策]    [超时管理]    [结果处理]
     |              |              |              |              |              |
     +---------------------------------------------------------------------------+
```

## 环境准备

### 开发环境要求

- Go 1.25或更高版本
- Git版本控制工具
- Docker (可选，用于容器化测试)
- Make (用于构建和测试)

### 安装Go环境

```bash
# 下载Go安装包
wget https://golang.org/dl/go1.25.3.linux-amd64.tar.gz

# 解压到/usr/local
sudo tar -C /usr/local -xzf go1.25.3.linux-amd64.tar.gz

# 添加到PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

### 克隆项目

```bash
git clone https://github.com/yourusername/dir-monitor-go.git
cd dir-monitor-go
```

### 安装依赖

```bash
go mod tidy
```

## 项目结构

```
dir-monitor-go/
├── cmd/
│   └── dir-monitor-go/
│       └── main.go              # 程序入口点
├── configs/
│   └── config.json.example      # 配置文件示例
├── deploy/
│   ├── dir-monitor-go.service   # systemd服务文件
│   ├── deploy-service.sh        # 服务部署脚本
│   └── service-manager.sh       # 服务管理脚本
├── docs/
│   ├── API.md                   # API文档
│   ├── USER_GUIDE.md            # 用户指南
│   └── DEVELOPMENT.md           # 开发指南
├── internal/
│   ├── config/                  # 配置管理模块
│   │   └── config.go            # 配置结构定义和加载
│   ├── logger/                  # 日志模块
│   │   └── logger.go            # Zap日志实现
│   ├── model/                   # 数据模型
│   │   ├── event.go             # 文件事件模型
│   │   └── settings.go          # 配置模型
│   └── monitor/                 # 文件监控模块
│       ├── fsnotify_watcher.go  # 基于fsnotify的文件监视器
│       ├── manager.go           # 监控管理器
│       ├── monitor.go           # 核心监控逻辑（包含调度功能）
│       ├── process_kill.go      # 进程终止功能
│       ├── runner.go            # 命令运行器
│       └── watcher.go           # 文件监视器接口
├── go.mod                       # Go模块定义
├── go.sum                       # Go模块校验和
├── Makefile                     # 构建文件
└── README.md                    # 项目说明
```

## 构建和测试

### 本地构建

```bash
# 构建项目
go build -o dir-monitor-go cmd/dir-monitor-go/main.go

# 运行程序
./dir-monitor-go --config configs/config.json.example
```

### 交叉编译

```bash
# 编译Linux版本
GOOS=linux GOARCH=amd64 go build -o dir-monitor-go-linux-amd64 cmd/dir-monitor-go/main.go

# 编译Windows版本
GOOS=windows GOARCH=amd64 go build -o dir-monitor-go-windows-amd64.exe cmd/dir-monitor-go/main.go

# 编译macOS版本
GOOS=darwin GOARCH=amd64 go build -o dir-monitor-go-darwin-amd64 cmd/dir-monitor-go/main.go
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./internal/config/...
go test ./internal/monitor/...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 使用Makefile

```bash
# 构建项目
make build

# 运行测试
make test

# 运行测试并生成覆盖率报告
make coverage

# 清理构建文件
make clean

# 安装到系统路径
make install

# 卸载程序
make uninstall

# 构建Docker镜像
make docker-build

# 运行Docker容器
make docker-run
```

## 核心模块详解

### 配置管理模块

配置管理模块负责加载、解析和验证配置文件。主要功能包括：

1. 配置结构定义 (`internal/config/config.go`)
2. 配置文件加载 (`internal/config/loader.go`)
3. 配置验证 (`internal/config/validator.go`)

#### 配置结构

```go
type Config struct {
    Version  string    `json:"version"`
    Metadata Metadata  `json:"metadata"`
    Monitors []Monitor `json:"monitors"`
    Settings Settings  `json:"settings"`
}

type Monitor struct {
    ID              string   `json:"id"`
    Name            string   `json:"name"`
    Description     string   `json:"description"`
    Directory       string   `json:"directory"`
    Command         string   `json:"command"`
    FilePatterns    []string `json:"file_patterns"`
    Timeout         int      `json:"timeout"`
    Schedule        string   `json:"schedule"`
    DebounceSeconds int      `json:"debounce_seconds"`
    Enabled         bool     `json:"enabled"`
}

type Settings struct {
    LogLevel                        string `json:"log_level"`
    LogFile                         string `json:"log_file"`
    LogShowCaller                   bool   `json:"log_show_caller"`
    LogMaxSize                      int    `json:"log_max_size"`
    LogMaxBackups                   int    `json:"log_max_backups"`
    LogCompress                     bool   `json:"log_compress"`
    MaxConcurrentOperations         int    `json:"max_concurrent_operations"`
    OperationTimeoutSeconds         int    `json:"operation_timeout_seconds"`
    FileWatcherBufferSize           int    `json:"file_watcher_buffer_size"`
    EventChannelBufferSize          int    `json:"event_channel_buffer_size"`
    MinStabilityTimeMs              int    `json:"min_stability_time_ms"`
    ExecutionDedupIntervalSeconds   int    `json:"execution_dedup_interval_seconds"`
    DirectoryStabilityQuietMs       int    `json:"directory_stability_quiet_ms"`
    DirectoryStabilityTimeoutSeconds int   `json:"directory_stability_timeout_seconds"`
    RetryAttempts                   int    `json:"retry_attempts"`
    RetryDelaySeconds               int    `json:"retry_delay_seconds"`
    HealthCheckIntervalSeconds      int    `json:"health_check_interval_seconds"`
}
```

### 文件监控模块

文件监控模块基于fsnotify实现文件系统事件监听，主要组件包括：

1. 文件监视器 (`internal/monitor/watcher.go`)
2. 命令运行器 (`internal/monitor/runner.go`)
3. 稳定性检测器 (`internal/monitor/stability.go`)

#### 文件监视器工作流程

1. 初始化fsnotify watcher
2. 添加监控目录到watcher
3. 启动事件监听goroutine
4. 处理文件系统事件
5. 过滤和去重事件
6. 触发命令执行

#### 稳定性检测机制

为避免对正在写入的文件过早执行操作，系统实现了稳定性检测机制：

1. 文件事件触发后，等待指定的稳定性时间
2. 在稳定性时间内，定期检查文件是否仍在变化
3. 如果文件在稳定性时间内保持不变，则认为文件已稳定
4. 文件稳定后，触发相应的操作执行

### 调度功能

调度功能集成在文件监控模块中，使用github.com/adhocore/gronx库实现Cron表达式解析和调度控制：

1. 调度检查逻辑 (`internal/monitor/monitor.go`中的`isScheduleActive`方法)
2. 调度决策逻辑

#### 调度工作流程

1. 解析监控器配置中的Cron表达式
2. 在文件事件触发时检查当前时间是否匹配调度表达式
3. 时间窗口匹配时允许操作执行
4. 时间窗口不匹配时跳过操作执行

### 命令执行功能

命令执行功能集成在文件监控模块中，负责管理命令执行队列和并发控制：

1. 命令执行器 (`internal/monitor/runner.go`中的`CommandExecutor`结构体)
2. 并发控制逻辑

#### 命令执行特性

1. 并发控制：限制同时执行的操作数量
2. 超时管理：为每个操作设置超时时间
3. 重试机制：操作失败时自动重试（通过配置实现）
4. 环境变量注入：为命令执行提供环境变量

### 日志模块

日志模块基于zap实现结构化日志记录：

1. Zap日志实现 (`internal/logger/zap_logger.go`)
2. 日志配置和管理

#### 日志特性

1. 多级别日志输出（debug, info, warn, error）
2. 日志轮转和压缩
3. 结构化日志格式
4. 调用者信息显示

## 性能优化策略

### 监控大量文件

当需要监控大量文件时，可以通过以下配置优化性能：

```json
{
  "settings": {
    "file_watcher_buffer_size": 500,
    "event_channel_buffer_size": 500,
    "max_concurrent_operations": 5,
    "default_debounce_seconds": 30
  }
}
```

### 处理大文件

对于大文件处理，建议调整以下配置：

```json
{
  "settings": {
    "min_stability_time_ms": 30000,
    "directory_stability_quiet_ms": 10000,
    "directory_stability_timeout_seconds": 300,
    "operation_timeout_seconds": 1800
  }
}
```

### 减少资源占用

为了减少系统资源占用，可以使用以下配置：

```json
{
  "settings": {
    "log_level": "warn",
    "log_max_size": 5242880,
    "log_max_backups": 2,
    "max_concurrent_operations": 3,
    "file_watcher_buffer_size": 50,
    "event_channel_buffer_size": 50
  }
}
```

## 部署指南

### 系统服务部署

1. 创建系统用户（可选）：
   ```bash
   sudo useradd --system --home /opt/dir-monitor-go --shell /bin/false dir-monitor
   ```

2. 创建配置目录：
   ```bash
   sudo mkdir -p /opt/dir-monitor-go/configs
   sudo mkdir -p /opt/dir-monitor-go/logs
   ```

3. 复制配置文件：
   ```bash
   sudo cp configs/config.json.example /opt/dir-monitor-go/configs/config.json
   ```

4. 编辑配置文件：
   ```bash
   sudo nano /opt/dir-monitor-go/configs/config.json
   ```

5. 复制系统服务文件：
   ```bash
   sudo cp deploy/dir-monitor-go.service /etc/systemd/system/
   ```

6. 重新加载systemd配置：
   ```bash
   sudo systemctl daemon-reload
   ```

7. 启用服务开机自启：
   ```bash
   sudo systemctl enable dir-monitor-go
   ```

8. 启动服务：
   ```bash
   sudo systemctl start dir-monitor-go
   ```

### Docker部署

1. 构建Docker镜像：
   ```bash
   docker build -t dir-monitor-go .
   ```

2. 运行Docker容器：
   ```bash
   docker run -d \
     --name dir-monitor-go \
     -v /path/to/configs:/opt/dir-monitor-go/configs \
     -v /path/to/monitor:/path/to/monitor \
     -v /path/to/logs:/opt/dir-monitor-go/logs \
     dir-monitor-go
   ```

3. 查看容器日志：
   ```bash
   docker logs -f dir-monitor-go
   ```

### Kubernetes部署

1. 创建ConfigMap：
   ```yaml
   apiVersion: v1
   kind: ConfigMap
   metadata:
     name: dir-monitor-go-config
   data:
     config.json: |
       {
         "version": "3.2.1",
         "metadata": {
           "name": "K8s部署配置",
           "description": "Kubernetes环境下的目录监控配置"
         },
         "monitors": [
           {
             "id": "k8s_monitor",
             "name": "K8s文件监控器",
             "directory": "/data/input",
             "command": "/scripts/process-file.sh",
             "file_patterns": ["*.txt", "*.csv"],
             "timeout": 300,
             "debounce_seconds": 15,
             "enabled": true
           }
         ],
         "settings": {
           "log_level": "info",
           "log_file": "/data/logs/app.log",
           "max_concurrent_operations": 5
         }
       }
   ```

2. 创建Deployment：
   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: dir-monitor-go
   spec:
     replicas: 1
     selector:
       matchLabels:
         app: dir-monitor-go
     template:
       metadata:
         labels:
           app: dir-monitor-go
       spec:
         containers:
         - name: dir-monitor-go
           image: dir-monitor-go:latest
           volumeMounts:
           - name: config-volume
             mountPath: /opt/dir-monitor-go/configs
           - name: data-volume
             mountPath: /data
         volumes:
         - name: config-volume
           configMap:
             name: dir-monitor-go-config
         - name: data-volume
           persistentVolumeClaim:
             claimName: dir-monitor-data-pvc
   ```

3. 部署到Kubernetes集群：
   ```bash
   kubectl apply -f configmap.yaml
   kubectl apply -f deployment.yaml
   ```

## 故障排除

### 编译问题

1. 确保Go版本符合要求（1.25+）
2. 运行`go mod tidy`更新依赖
3. 检查是否有语法错误

### 运行时问题

1. 检查配置文件格式和路径
2. 确认监控目录权限
3. 查看日志文件获取详细错误信息

### 性能问题

1. 检查是否监控了过多文件
2. 调整并发操作数设置
3. 增加防抖时间避免频繁触发

## 贡献指南

我们欢迎社区贡献！在提交贡献之前，请确保：

1. Fork项目并创建功能分支
2. 遵循项目编码规范
3. 添加相应的测试用例
4. 确保所有测试通过
5. 提交清晰的commit信息
6. 创建Pull Request并详细描述变更内容

## 许可证

本项目采用MIT许可证，详情请参见[LICENSE](../LICENSE)文件。