# 目录监控系统 API 文档

## 概述

目录监控系统是一个基于Go语言开发的文件监控和处理系统，它可以监控指定目录的文件变化，并根据预定义的监控器配置执行相应的操作。

当前版本：v3.2.1

## 配置结构

### 主配置结构

```json
{
  "version": "3.2.1",
  "metadata": {
    "name": "配置名称",
    "description": "配置描述"
  },
  "monitors": [
    {
      "id": "monitor_id",
      "name": "监控器名称",
      "description": "监控器描述",
      "directory": "/path/to/directory",
      "command": "/path/to/command",
      "file_patterns": ["*.txt", "*.log"],
      "timeout": 300,
      "schedule": "* 15-21 * * 1-5",
      "debounce_seconds": 15,
      "enabled": true
    }
  ],
  "settings": {
    "log_level": "info",
    "log_file": "/path/to/log/file",
    "log_show_caller": true,
    "log_max_size": 10485760,
    "log_compress": true,
    "max_concurrent_operations": 10,
    "operation_timeout_seconds": 600,
    "file_stability_check_interval_ms": 1000,
    "file_stability_timeout_seconds": 180,
    "min_stability_time_ms": 10000,
    "default_debounce_seconds": 10,
    "directory_stability_quiet_ms": 5000,
    "directory_stability_timeout_seconds": 180,
    "small_file_threshold": 1048576,
    "medium_file_threshold": 10485760,
    "log_max_backups": 3,
    "file_watcher_buffer_size": 100,
    "event_channel_buffer_size": 100,
    "execution_dedup_interval_seconds": 10,
    "retry_attempts": 3,
    "retry_delay_seconds": 5,
    "health_check_interval_seconds": 60
  }
}
```

### 配置版本兼容性

- **v3.2.x**: 当前版本，支持所有新特性
- **v3.1.x**: 兼容版本，支持基本功能
- **v3.0.x**: 旧版本，建议使用最新版本

### 配置字段说明

#### Metadata（配置元数据）

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 配置名称，用于标识配置用途 |
| description | string | 否 | 配置描述，详细说明配置的使用场景 |

#### Monitor（监控器配置）

监控器是系统的核心组件，定义了目录监控的完整配置。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 监控器唯一标识符，用于日志和调试 |
| name | string | 是 | 监控器名称，便于识别和管理 |
| description | string | 否 | 监控器描述，说明监控器用途 |
| directory | string | 是 | 要监控的目录路径 |
| command | string | 是 | 检测到文件变化时执行的命令 |
| file_patterns | string[] | 是 | 文件模式数组，支持通配符（*.txt, *.csv） |
| timeout | int | 否 | 命令执行超时时间（秒），默认300 |
| schedule | string | 否 | Cron表达式，控制执行时间段（如：* 15-21 * * 1-5） |
| debounce_seconds | int | 否 | 防抖时间（秒），防止重复触发，默认15 |
| enabled | boolean | 否 | 是否启用该监控器，默认true |

**Cron表达式格式说明：**
```
┌───────────── 分钟 (0 - 59)
│ ┌───────────── 小时 (0 - 23)
│ │ ┌───────────── 日 (1 - 31)
│ │ │ ┌───────────── 月 (1 - 12)
│ │ │ │ ┌───────────── 星期 (0 - 7) (0和7都表示星期日)
│ │ │ │ │
* * * * *
```

**常用示例：**
- `* 15-21 * * 1-5`：周一至周五 15-21点
- `* 0-9 * * 1-6`：周一至周六 0-9点
- `0 2 * * *`：每天凌晨2点
- `*/30 * * * *`：每30分钟

#### Settings（系统设置）

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| log_level | string | 否 | 日志级别，可选值：debug, info, warn, error，默认info |
| log_file | string | 否 | 日志文件路径，默认logs/dir-monitor-go.log |
| log_show_caller | boolean | 否 | 是否显示调用者信息，默认true |
| log_max_size | int | 否 | 日志文件最大大小（字节），默认10485760（10MB） |
| log_max_backups | int | 否 | 日志文件备份数量，默认3 |
| log_compress | boolean | 否 | 是否压缩轮转的日志文件，默认true |
| max_concurrent_operations | int | 否 | 最大并发操作数，默认10 |
| operation_timeout_seconds | int | 否 | 操作超时时间（秒），默认600 |
| file_watcher_buffer_size | int | 否 | 文件监视器缓冲区大小，默认100 |
| event_channel_buffer_size | int | 否 | 事件通道缓冲区大小，默认100 |
| min_stability_time_ms | int | 否 | 最小稳定性时间（毫秒），默认10000 |
| execution_dedup_interval_seconds | int | 否 | 执行去重间隔（秒），默认10 |
| directory_stability_quiet_ms | int | 否 | 目录稳定性静默时间（毫秒），默认5000 |
| directory_stability_timeout_seconds | int | 否 | 目录稳定性超时时间（秒），默认180 |
| retry_attempts | int | 否 | 重试次数，默认3 |
| retry_delay_seconds | int | 否 | 重试延迟时间（秒），默认5 |
| health_check_interval_seconds | int | 否 | 健康检查间隔（秒），默认60 |

## 环境变量

在执行操作时，系统会提供以下环境变量：

| 变量名 | 说明 | 示例 |
|--------|------|------|
| FILE_PATH | 文件完整路径 | /sftp/user2/data/report.csv |
| FILE_NAME | 文件名 | report.csv |
| FILE_DIR | 文件所在目录 | /sftp/user2/data |
| EVENT_TYPE | 事件类型（created, modified, deleted, renamed） | created |
| EVENT_TIME | 事件时间（Unix时间戳） | 1698768000 |

### 使用示例

在命令中使用环境变量：
```bash
echo "文件 ${FILE_NAME} 在 ${FILE_DIR} 目录中被 ${EVENT_TYPE}，时间：${EVENT_TIME}"
```

Python脚本中使用环境变量：
```python
import os
file_path = os.environ.get('FILE_PATH')
event_type = os.environ.get('EVENT_TYPE')
print(f"文件 {file_path} 被 {event_type}")
```

## 命令行接口

### 启动监控程序

```bash
# 基本启动（使用默认配置文件）
./dir-monitor-go

# 指定配置文件启动
./dir-monitor-go --config /path/to/config.json

# 验证配置
./dir-monitor-go --config /path/to/config.json --dry-run

# 显示版本信息
./dir-monitor-go --version

# 测试模式（使用停止文件）
./dir-monitor-go --config /path/to/config.json --stop-file /tmp/stop_marker
```

### 命令行参数

| 参数 | 必填 | 说明 | 示例 |
|------|------|------|------|
| --config, -c | 否 | 配置文件路径，默认configs/config.json | /etc/dir-monitor/config.json |
| --stop-file | 否 | 当该文件出现时优雅退出（测试/集成用） | /tmp/stop_marker |
| --version, -v | 否 | 显示版本信息 | - |
| --dry-run | 否 | 仅验证配置，不启动实际监控 | - |
| --help, -h | 否 | 显示帮助信息 | - |

## 示例配置

### 基本文件监控

适用于简单的文件监控场景，当文件创建或修改时执行指定命令。

```json
{
  "version": "3.2.1",
  "metadata": {
    "name": "基本文件监控配置",
    "description": "监控目录中的文件变化并执行命令"
  },
  "monitors": [
    {
      "id": "file_watcher",
      "name": "文件监控器",
      "directory": "/path/to/watch",
      "command": "echo '文件 $FILE_PATH 被 $EVENT_TYPE'",
      "file_patterns": ["*.txt", "*.log"],
      "timeout": 300,
      "debounce_seconds": 10,
      "enabled": true
    }
  ],
  "settings": {
    "log_level": "info",
    "log_file": "logs/dir-monitor-go.log",
    "max_concurrent_operations": 5,
    "operation_timeout_seconds": 300
  }
}
```

### Python脚本处理

适用于需要使用Python脚本处理文件的场景。

```json
{
  "version": "3.2.1",
  "metadata": {
    "name": "Python脚本处理配置",
    "description": "使用Python脚本处理监控到的文件"
  },
  "monitors": [
    {
      "id": "python_processor",
      "name": "Python处理监控器",
      "directory": "/sftp/incoming",
      "command": "/usr/bin/python3 /opt/scripts/process_file.py",
      "file_patterns": ["*.csv", "*.xlsx"],
      "timeout": 600,
      "schedule": "* 9-17 * * 1-5",
      "debounce_seconds": 30,
      "enabled": true
    }
  ],
  "settings": {
    "log_level": "info",
    "log_file": "/var/log/dir-monitor-go/app.log",
    "log_max_size": 20971520,
    "log_max_backups": 5,
    "max_concurrent_operations": 3,
    "operation_timeout_seconds": 600,
    "retry_attempts": 3,
    "retry_delay_seconds": 10
  }
}
```

## 错误处理

### 常见错误码

| 错误码 | 说明 | 建议处理方式 |
|--------|------|--------------|
| ERR001 | 配置文件不存在或无法读取 | 检查配置文件路径和权限 |
| ERR002 | 配置文件格式错误 | 使用JSON验证工具检查配置文件 |
| ERR003 | 监控目录不存在或无权限访问 | 检查目录路径和权限设置 |
| ERR004 | 命令执行失败 | 检查命令路径和参数设置 |
| ERR005 | 文件处理超时 | 增加timeout设置或优化处理逻辑 |
| ERR006 | 并发操作数超过限制 | 调整max_concurrent_operations设置 |

### 处理建议

1. **配置错误**：启动时会验证配置文件，如有错误会显示详细信息并退出
2. **目录权限**：确保监控目录具有读取和执行权限
3. **命令执行**：确保命令路径正确且具有执行权限
4. **日志分析**：通过日志文件分析问题原因
5. **调试模式**：使用--dry-run参数验证配置

### 调试技巧

1. 使用`--dry-run`参数验证配置而不启动监控
2. 设置`log_level`为`debug`获取详细日志信息
3. 检查系统日志获取更多错误信息
4. 使用`--stop-file`参数进行测试

## 性能优化

### 监控大量文件

当需要监控大量文件时，建议使用以下配置：

```json
{
  "settings": {
    "file_watcher_buffer_size": 2048,
    "event_channel_buffer_size": 500,
    "max_concurrent_operations": 20,
    "directory_stability_quiet_ms": 10000,
    "directory_stability_timeout_seconds": 300
  }
}
```

### 处理大文件

对于大文件处理，建议使用以下配置：

```json
{
  "settings": {
    "small_file_threshold": 1048576,
    "medium_file_threshold": 10485760,
    "operation_timeout_seconds": 1800,
    "file_stability_check_interval_ms": 2000,
    "file_stability_timeout_seconds": 600
  }
}
```

### 减少资源占用

为减少系统资源占用，可使用以下配置：

```json
{
  "settings": {
    "max_concurrent_operations": 2,
    "file_watcher_buffer_size": 50,
    "event_channel_buffer_size": 50,
    "default_debounce_seconds": 30,
    "execution_dedup_interval_seconds": 60
  }
}
```

## 故障排除

### 监控不工作

1. 检查配置文件中的目录路径是否正确
2. 确认目录存在且具有读取权限
3. 检查日志文件获取错误信息
4. 验证命令路径和参数设置

### 操作不执行

1. 检查命令路径是否正确
2. 确认命令具有执行权限
3. 查看日志了解命令执行失败原因
4. 验证文件匹配模式是否正确

### 性能问题

1. 检查是否监控了过多文件
2. 调整并发操作数设置
3. 增加防抖时间避免频繁触发
4. 优化命令执行逻辑

### 内存泄漏

1. 检查是否有未释放的资源
2. 监控系统内存使用情况
3. 更新到最新版本修复已知问题
4. 联系开发者报告问题

## 开发指南

### 编译程序

```bash
# 克隆代码仓库
git clone https://github.com/yourusername/dir-monitor-go.git
cd dir-monitor-go

# 安装依赖
go mod download

# 编译程序
go build -o dir-monitor-go cmd/dir-monitor-go/main.go

# 编译带版本信息
go build -ldflags "-X main.version=3.2.1 -X main.buildTime=$(date -u +%Y-%m-%d_%H:%M:%S)" -o dir-monitor-go cmd/dir-monitor-go/main.go

# 运行测试
go test ./...

# 安装程序
go install

# 交叉编译（例如编译Linux ARM版本）
GOOS=linux GOARCH=arm64 go build -o dir-monitor-go-linux-arm64 cmd/dir-monitor-go/main.go
```

### 代码结构

```
dir-monitor-go/
├── cmd/
│   └── dir-monitor-go/     # 主程序入口
├── internal/               # 内部包
│   ├── config/             # 配置处理
│   ├── logger/             # 日志处理
│   ├── model/              # 数据模型
│   └── monitor/            # 监控器
├── configs/                # 配置文件
├── docs/                   # 文档
├── deploy/                 # 部署脚本
├── logs/                   # 日志文件
├── build/                  # 构建输出
├── bin/                    # 二进制文件
├── Makefile                # 构建脚本
├── go.mod                  # Go模块定义
└── go.sum                  # Go模块校验和
```