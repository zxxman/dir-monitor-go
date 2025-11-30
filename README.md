# Dir-Monitor-Go

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/zxxman/dir-monitor-go)

一个高效、可靠的目录文件监控工具，使用Go语言开发，支持实时监控文件系统变化并执行自定义命令。

## 特性

- 🚀 **高性能**：基于fsnotify实现高效的文件系统事件监控
- ⏰ **调度支持**：内置cron表达式支持，可在指定时间窗口执行命令
- 🔄 **重试机制**：内置命令执行失败重试机制
- 📝 **详细日志**：支持多级别日志记录和文件日志轮转
- 🎯 **防抖动**：支持文件和目录稳定性检测，避免重复触发
- 🔧 **灵活配置**：JSON格式配置文件，支持多监控项
- 🛡️ **并发控制**：可配置最大并发操作数，防止资源耗尽
- 📦 **轻量部署**：单一二进制文件，支持系统服务部署

## 快速开始

### 安装

#### 从源码构建

```bash
git clone https://github.com/zxxman/dir-monitor-go.git
cd dir-monitor-go
make build
```

#### 直接下载二进制文件

从[Releases](https://github.com/zxxman/dir-monitor-go/releases)页面下载适合您系统的二进制文件。

### 配置

1. 复制示例配置文件：
```bash
cp config.json.example configs/config.json
```

2. 编辑配置文件：
```json
{
  "version": "3.2.1",
  "monitors": [
    {
      "id": "log-monitor",
      "name": "日志文件监控",
      "directory": "/var/log",
      "command": "echo '检测到日志文件变化: ${FILE_PATH}'",
      "file_patterns": ["*.log"],
      "timeout": 60,
      "enabled": true
    }
  ],
  "settings": {
    "log_file": "logs/dir-monitor-go.log",
    "log_level": "info",
    "max_concurrent_operations": 5
  }
}
```

### 运行

```bash
# 使用默认配置
./dir-monitor-go

# 指定配置文件
./dir-monitor-go -config /path/to/config.json

# 仅验证配置
./dir-monitor-go -dry-run

# 查看版本信息
./dir-monitor-go -version
```

## 配置说明

### 监控项配置

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| id | string | 否 | 监控项唯一标识符 |
| name | string | 否 | 监控项名称 |
| directory | string | 是 | 要监控的目录路径 |
| command | string | 是 | 文件变化时执行的命令 |
| file_patterns | []string | 是 | 监控的文件模式（如["*.log", "*.txt"]） |
| timeout | int | 是 | 命令执行超时时间（秒） |
| schedule | string | 否 | cron表达式，指定执行时间窗口 |
| debounce_seconds | int | 否 | 防抖动时间（秒） |
| enabled | bool | 否 | 是否启用此监控项，默认true |

### 全局设置

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| log_file | string | - | 日志文件路径 |
| log_level | string | "info" | 日志级别（debug/info/warn/error） |
| log_show_caller | bool | false | 是否在日志中显示调用者信息 |
| log_max_size | int | 10485760 | 日志文件最大大小（字节） |
| max_concurrent_operations | int | 5 | 最大并发操作数 |
| operation_timeout_seconds | int | 300 | 默认操作超时时间（秒） |
| min_stability_time_ms | int | 5000 | 文件最小稳定时间（毫秒） |
| directory_stability_quiet_ms | int | 2000 | 目录稳定静默时间（毫秒） |

## 命令行参数

| 参数 | 说明 |
|------|------|
| -config | 指定配置文件路径（默认：configs/config.json） |
| -version | 显示版本信息 |
| -dry-run | 仅验证配置，不启动实际监控 |
| -stop-file | 指定停止文件，当该文件存在时优雅退出 |

## 部署

### 作为系统服务

```bash
# 安装为系统服务
sudo make install-service

# 启动服务
sudo systemctl start dir-monitor-go

# 查看服务状态
sudo systemctl status dir-monitor-go

# 查看日志
sudo journalctl -u dir-monitor-go -f
```

### Docker部署

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/dir-monitor-go .
COPY --from=builder /app/config.json.example .
CMD ["./dir-monitor-go"]
```

## 使用示例

### 示例1：监控日志文件

```json
{
  "id": "log-monitor",
  "name": "日志文件监控",
  "directory": "/var/log/myapp",
  "command": "tail -n 10 ${FILE_PATH}",
  "file_patterns": ["*.log"],
  "timeout": 30,
  "enabled": true
}
```

### 示例2：上传文件处理

```json
{
  "id": "upload-processor",
  "name": "上传文件处理",
  "directory": "/uploads",
  "command": "python /scripts/process_upload.py ${FILE_PATH}",
  "file_patterns": ["*.csv", "*.xlsx", "*.json"],
  "timeout": 300,
  "schedule": "* 9-18 * * 1-5",
  "debounce_seconds": 10,
  "enabled": true
}
```

### 示例3：代码构建触发

```json
{
  "id": "build-trigger",
  "name": "代码变化构建",
  "directory": "/src/myproject",
  "command": "cd /src/myproject && make build",
  "file_patterns": ["*.go", "go.mod", "go.sum"],
  "timeout": 600,
  "debounce_seconds": 5,
  "enabled": true
}
```

## 命令变量

在执行命令时，可以使用以下变量：

| 变量 | 说明 |
|------|------|
| ${FILE_PATH} | 变化的文件完整路径 |
| ${FILE_NAME} | 文件名（不含路径） |
| ${FILE_DIR} | 文件所在目录 |
| ${FILE_EXT} | 文件扩展名 |
| ${EVENT_TYPE} | 事件类型（create/write/remove/rename） |
| ${TIMESTAMP} | 当前时间戳 |

## 故障排除

### 常见问题

1. **权限问题**：确保程序有权限访问监控目录和执行命令
2. **配置错误**：使用`-dry-run`参数验证配置文件
3. **命令执行失败**：检查命令路径和权限，查看日志获取详细错误信息

### 日志分析

```bash
# 查看实时日志
tail -f logs/dir-monitor-go.log

# 查看错误日志
grep "ERROR" logs/dir-monitor-go.log

# 查看特定监控项的日志
grep "monitor-id" logs/dir-monitor-go.log
```

## 开发

### 构建

```bash
# 开发构建
make build

# 生产构建
make build LDFLAGS="-ldflags '-s -w'"
```

### 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test -v ./internal/monitor
```

### 代码质量检查

```bash
# 代码格式化
make fmt

# 静态分析
make vet

# 代码质量检查
make quality
```

## 贡献

欢迎提交Issue和Pull Request！请确保：

1. 代码通过所有测试
2. 遵循Go代码规范
3. 添加必要的文档和注释
4. 更新相关文档

## 许可证

本项目采用MIT许可证 - 查看[LICENSE](LICENSE)文件了解详情。

## 更新日志

查看[CHANGELOG.md](CHANGELOG.md)了解版本更新历史。