# 目录监控服务用户指南

## 目录

1. [简介](#简介)
2. [安装与配置](#安装与配置)
3. [基本使用](#基本使用)
4. [配置文件详解](#配置文件详解)
5. [命令行参数](#命令行参数)
6. [常见使用场景](#常见使用场景)
7. [故障排除](#故障排除)
8. [高级功能](#高级功能)

## 简介

目录监控服务（dir-monitor-go）是一个轻量级的文件系统监控工具，可以监控指定目录的变化并触发自定义命令。它支持多种文件模式匹配、并发操作控制、命令变量替换等高级功能。

## 安装与配置

### 安装方式

#### 从源码编译

```bash
git clone https://github.com/zxxman/dir-monitor-go.git
cd dir-monitor-go
make build
```

#### 使用预编译二进制文件

从 [Releases](https://github.com/zxxman/dir-monitor-go/releases) 页面下载适合您系统的预编译二进制文件。

#### 使用Docker

```bash
docker pull zxxman/dir-monitor-go:latest
```

### 初始配置

首次运行时，创建配置文件：

```bash
# 复制示例配置文件
cp config.json.example config.json

# 编辑配置文件
nano config.json
```

## 基本使用

### 快速启动

1. 创建一个简单的配置文件：

```json
{
  "version": "3.2.1",
  "monitors": [
    {
      "directory": "/tmp/test",
      "command": "echo '文件已更改: {{.FilePath}}'",
      "file_patterns": ["*.txt", "*.log"]
    }
  ],
  "settings": {
    "log_level": "info",
    "max_concurrent_operations": 5,
    "file_stability_check": {
      "enabled": true,
      "min_stable_duration": "1s"
    }
  }
}
```

2. 启动监控服务：

```bash
./dir-monitor-go -c config.json
```

3. 在另一个终端中测试：

```bash
touch /tmp/test/example.txt
```

您应该会看到监控服务输出"文件已更改: /tmp/test/example.txt"。

### 作为系统服务运行

```bash
# 安装为系统服务
sudo make install-service

# 启动服务
sudo systemctl start dir-monitor-go

# 设置开机自启
sudo systemctl enable dir-monitor-go
```

## 配置文件详解

### 配置文件结构

配置文件采用JSON格式，包含三个主要部分：

- `version`: 配置文件版本
- `monitors`: 监控器配置数组
- `settings`: 全局设置

### 监控器配置

每个监控器包含以下字段：

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `directory` | string | 是 | 要监控的目录路径 |
| `command` | string | 是 | 触发时执行的命令 |
| `file_patterns` | []string | 否 | 文件模式匹配列表（默认匹配所有文件） |
| `recursive` | bool | 否 | 是否递归监控子目录（默认true） |
| `events` | []string | 否 | 要监控的事件类型（默认监控所有事件） |
| `debounce` | object | 否 | 防抖配置 |
| `filters` | object | 否 | 文件过滤器配置 |

#### 事件类型

支持的事件类型：
- `create`: 文件创建
- `write`: 文件写入
- `remove`: 文件删除
- `rename`: 文件重命名
- `chmod`: 文件权限变更

#### 防抖配置

```json
"debounce": {
  "enabled": true,
  "delay": "500ms"
}
```

#### 过滤器配置

```json
"filters": {
  "min_size": "1KB",
  "max_size": "10MB",
  "exclude_patterns": ["*.tmp", "*.bak"]
}
```

### 全局设置

| 字段 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `log_level` | string | "info" | 日志级别（debug, info, warn, error） |
| `log_file` | string | - | 日志文件路径（默认输出到控制台） |
| `max_concurrent_operations` | int | 5 | 最大并发操作数 |
| `file_stability_check` | object | - | 文件稳定性检查配置 |
| `performance_monitoring` | object | - | 性能监控配置 |

#### 文件稳定性检查

```json
"file_stability_check": {
  "enabled": true,
  "min_stable_duration": "1s",
  "check_interval": "100ms"
}
```

#### 性能监控

```json
"performance_monitoring": {
  "enabled": true,
  "report_interval": "1m"
}
```

## 命令行参数

| 参数 | 简写 | 描述 | 默认值 |
|------|------|------|--------|
| `--config` | `-c` | 配置文件路径 | `./config.json` |
| `--log-level` | `-l` | 日志级别 | `info` |
| `--version` | `-v` | 显示版本信息 | - |
| `--help` | `-h` | 显示帮助信息 | - |

### 示例

```bash
# 使用自定义配置文件
./dir-monitor-go -c /etc/dir-monitor/config.json

# 设置日志级别为debug
./dir-monitor-go -l debug

# 显示版本信息
./dir-monitor-go -v
```

## 常见使用场景

### 1. 日志文件监控

```json
{
  "monitors": [
    {
      "directory": "/var/log",
      "command": "tail -n 10 {{.FilePath}}",
      "file_patterns": ["*.log"],
      "events": ["write"]
    }
  ]
}
```

### 2. 上传文件处理

```json
{
  "monitors": [
    {
      "directory": "/uploads",
      "command": "python3 process_upload.py {{.FilePath}}",
      "file_patterns": ["*.jpg", "*.png", "*.pdf"],
      "events": ["create"],
      "debounce": {
        "enabled": true,
        "delay": "2s"
      }
    }
  ]
}
```

### 3. 配置文件热重载

```json
{
  "monitors": [
    {
      "directory": "/etc/myapp",
      "command": "systemctl reload myapp",
      "file_patterns": ["*.conf"],
      "events": ["write"]
    }
  ]
}
```

### 4. 备份新文件

```json
{
  "monitors": [
    {
      "directory": "/data",
      "command": "rsync -av {{.FilePath}} /backup/$(date +%Y%m%d)/",
      "file_patterns": ["*"],
      "events": ["create"],
      "filters": {
        "min_size": "1MB"
      }
    }
  ]
}
```

## 故障排除

### 常见问题

#### 1. 监控不工作

**可能原因**：
- 配置文件路径错误
- 监控目录不存在
- 权限不足

**解决方法**：
```bash
# 检查配置文件
./dir-monitor-go -c config.json -l debug

# 检查目录权限
ls -la /path/to/monitor

# 确保有读取目录的权限
sudo usermod -a -G $(stat -c '%G' /path/to/monitor) $USER
```

#### 2. 命令执行失败

**可能原因**：
- 命令路径错误
- 命令权限不足
- 变量替换错误

**解决方法**：
```bash
# 测试命令
echo "文件已更改: /tmp/test.txt" > /tmp/test.txt

# 检查命令权限
which your-command

# 使用绝对路径
"/usr/bin/python3" "/path/to/script.py" "{{.FilePath}}"
```

#### 3. 性能问题

**可能原因**：
- 监控目录过大
- 文件变化频繁
- 并发操作过多

**解决方法**：
```json
{
  "settings": {
    "max_concurrent_operations": 2,
    "file_stability_check": {
      "enabled": true,
      "min_stable_duration": "5s"
    }
  }
}
```

#### 4. 日志文件过大

**解决方法**：
```json
{
  "settings": {
    "log_level": "warn",
    "log_file": "/var/log/dir-monitor-go.log"
  }
}
```

并配置日志轮转：

```bash
# 创建logrotate配置
sudo nano /etc/logrotate.d/dir-monitor-go
```

内容：
```
/var/log/dir-monitor-go.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 644 dir-monitor dir-monitor
}
```

### 调试技巧

1. **使用调试日志级别**：
   ```bash
   ./dir-monitor-go -l debug
   ```

2. **测试配置文件**：
   ```bash
   ./dir-monitor-go -c config.json -t
   ```

3. **手动触发命令**：
   ```bash
   # 替换变量
   FilePath="/tmp/test.txt"
   echo "文件已更改: $FilePath"
   ```

## 高级功能

### 命令变量替换

在命令中可以使用以下变量：

| 变量 | 描述 |
|------|------|
| `{{.FilePath}}` | 完整文件路径 |
| `{{.FileName}}` | 文件名（不含路径） |
| `{{.FileExt}}` | 文件扩展名 |
| `{{.DirPath}}` | 目录路径 |
| `{{.EventName}}` | 事件名称 |
| `{{.Timestamp}}` | 时间戳 |
| `{{.PID}}` | 进程ID |

### 复杂命令示例

```json
{
  "command": "bash -c 'echo \"{{.Timestamp}}: {{.EventName}} on {{.FilePath}}\" >> /var/log/file-events.log'"
}
```

### 多监控器配置

```json
{
  "monitors": [
    {
      "directory": "/var/log",
      "command": "tail -n 5 {{.FilePath}}",
      "file_patterns": ["*.log"],
      "events": ["write"]
    },
    {
      "directory": "/uploads",
      "command": "python3 /scripts/process_upload.py {{.FilePath}}",
      "file_patterns": ["*.jpg", "*.png"],
      "events": ["create"]
    },
    {
      "directory": "/etc/myapp",
      "command": "systemctl reload myapp",
      "file_patterns": ["*.conf"],
      "events": ["write"]
    }
  ]
}
```

### Docker部署

创建Dockerfile：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o dir-monitor-go ./cmd/dir-monitor-go

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /root/

COPY --from=builder /app/dir-monitor-go .
COPY config.json.example config.json

CMD ["./dir-monitor-go"]
```

构建和运行：

```bash
docker build -t dir-monitor-go .
docker run -v /path/to/monitor:/data -v /path/to/config.json:/root/config.json dir-monitor-go
```

### Docker Compose部署

创建docker-compose.yml：

```yaml
version: '3.8'

services:
  dir-monitor:
    image: zxxman/dir-monitor-go:latest
    container_name: dir-monitor
    restart: unless-stopped
    volumes:
      - ./config.json:/app/config.json:ro
      - /path/to/monitor:/data:ro
      - /path/to/logs:/app/logs
    environment:
      - LOG_LEVEL=info
```

运行：

```bash
docker-compose up -d
```

---

如需更多帮助，请参考[开发文档](DEVELOPMENT.md)或提交[Issue](https://github.com/zxxman/dir-monitor-go/issues)。