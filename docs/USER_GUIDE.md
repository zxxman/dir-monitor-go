# 目录监控系统用户指南

## 系统简介

目录监控系统是一个基于Go语言开发的文件监控和处理系统，它可以监控指定目录的文件变化，并根据预定义的监控器配置执行相应的操作。系统支持多种文件操作触发条件、灵活的调度配置和可靠的错误处理机制。

## 主要特性

- **实时文件监控**：监控目录中的文件创建、修改、删除等操作
- **灵活的文件匹配**：支持通配符和正则表达式匹配文件
- **定时调度控制**：通过Cron表达式精确控制操作执行时间
- **并发处理能力**：支持多文件同时处理，提高处理效率
- **稳定性和可靠性**：文件稳定性检测、防抖机制、重试机制
- **丰富的配置选项**：全面的系统设置和监控器配置
- **详细的日志记录**：结构化日志记录和日志轮转
- **易于部署和维护**：支持多种部署方式和系统服务管理

## 适用场景

- 文件传输后自动处理（FTP/SFTP接收文件处理）
- 日志文件监控和分析
- 数据文件导入和转换
- 自动化备份和同步
- 实时数据处理流水线

## 系统要求

### 操作系统

- Ubuntu 20.04 LTS 或更高版本
- Debian 11 或更高版本
- CentOS 7 或更高版本（实验性支持）
- 其他Linux发行版（可能需要额外配置）

### 硬件要求

- CPU：1核或以上
- 内存：512MB或以上
- 磁盘空间：根据监控文件大小确定
- 网络：如需网络操作，需要相应网络连接

### 软件依赖

- Go 1.25或更高版本（编译时需要）
- systemd（系统服务管理，可选）

## 安装指南

### 二进制安装

1. 下载预编译的二进制文件：

```bash
# 下载最新版本
wget https://github.com/yourusername/dir-monitor-go/releases/latest/download/dir-monitor-go-linux-amd64.tar.gz

# 解压文件
tar -xzf dir-monitor-go-linux-amd64.tar.gz

# 移动到系统路径
sudo mv dir-monitor-go /usr/local/bin/

# 验证安装
dir-monitor-go --version
```

### 源码安装

```bash
# 克隆代码仓库
git clone https://github.com/yourusername/dir-monitor-go.git
cd dir-monitor-go

# 编译程序
go build -o dir-monitor-go cmd/dir-monitor-go/main.go

# 安装到系统路径
sudo cp dir-monitor-go /usr/local/bin/

# 验证安装
dir-monitor-go --version
```

### 系统服务安装

```bash
# 创建系统用户（可选）
sudo useradd --system --home /opt/dir-monitor-go --shell /bin/false dir-monitor

# 创建配置目录
sudo mkdir -p /opt/dir-monitor-go/configs
sudo mkdir -p /opt/dir-monitor-go/logs

# 复制配置文件模板
sudo cp configs/config.json.example /opt/dir-monitor-go/configs/config.json

# 编辑配置文件
sudo nano /opt/dir-monitor-go/configs/config.json

# 复制系统服务文件
sudo cp deploy/dir-monitor-go.service /etc/systemd/system/

# 重新加载systemd配置
sudo systemctl daemon-reload

# 启用服务开机自启
sudo systemctl enable dir-monitor-go

# 启动服务
sudo systemctl start dir-monitor-go

# 检查服务状态
sudo systemctl status dir-monitor-go
```

## 配置详解

### 配置文件结构

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

### 系统设置字段说明

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

### 调度配置

系统支持通过Cron表达式来控制操作的执行时间。

#### Cron表达式格式

```
┌───────────── 分钟 (0 - 59)
│ ┌───────────── 小时 (0 - 23)
│ │ ┌───────────── 日 (1 - 31)
│ │ │ ┌───────────── 月 (1 - 12)
│ │ │ │ ┌───────────── 星期 (0 - 7) (0和7都表示星期日)
│ │ │ │ │
* * * * *
```

#### 常用调度示例

| 表达式 | 说明 |
|--------|------|
| `* 15-21 * * 1-5` | 周一至周五 15-21点 |
| `* 0-9 * * 1-6` | 周一至周六 0-9点 |
| `0 2 * * *` | 每天凌晨2点 |
| `*/30 * * * *` | 每30分钟 |
| `0 0 * * 0` | 每周日凌晨 |

### 完整配置示例

```json
{
  "version": "3.2.1",
  "metadata": {
    "name": "生产环境配置",
    "description": "用于生产环境的标准配置"
  },
  "monitors": [
    {
      "id": "sftp_processor",
      "name": "SFTP文件处理器",
      "description": "处理SFTP目录中的新文件",
      "directory": "/sftp/incoming",
      "command": "/opt/scripts/process_sftp_file.sh",
      "file_patterns": ["*.csv", "*.xlsx", "*.pdf"],
      "timeout": 600,
      "schedule": "* 9-17 * * 1-5",
      "debounce_seconds": 30,
      "enabled": true
    },
    {
      "id": "log_analyzer",
      "name": "日志分析器",
      "description": "分析应用日志文件",
      "directory": "/var/log/application",
      "command": "/opt/scripts/analyze_logs.py",
      "file_patterns": ["*.log"],
      "timeout": 300,
      "schedule": "*/15 * * * *",
      "debounce_seconds": 10,
      "enabled": true
    }
  ],
  "settings": {
    "log_level": "info",
    "log_file": "/var/log/dir-monitor-go/app.log",
    "log_show_caller": true,
    "log_max_size": 20971520,
    "log_compress": true,
    "log_max_backups": 5,
    "max_concurrent_operations": 5,
    "operation_timeout_seconds": 900,
    "file_stability_check_interval_ms": 1000,
    "file_stability_timeout_seconds": 180,
    "min_stability_time_ms": 10000,
    "default_debounce_seconds": 15,
    "directory_stability_quiet_ms": 5000,
    "directory_stability_timeout_seconds": 180,
    "small_file_threshold": 1048576,
    "medium_file_threshold": 10485760,
    "file_watcher_buffer_size": 100,
    "event_channel_buffer_size": 100,
    "execution_dedup_interval_seconds": 10,
    "retry_attempts": 3,
    "retry_delay_seconds": 5,
    "health_check_interval_seconds": 60
  }
}
```

## 命令行参数

### 启动参数

| 参数 | 必填 | 说明 | 示例 |
|------|------|------|------|
| --config, -c | 否 | 配置文件路径，默认configs/config.json | /etc/dir-monitor/config.json |
| --stop-file | 否 | 当该文件出现时优雅退出（测试/集成用） | /tmp/stop_marker |
| --version, -v | 否 | 显示版本信息 | - |
| --dry-run | 否 | 仅验证配置，不启动实际监控 | - |
| --help, -h | 否 | 显示帮助信息 | - |

### 使用示例

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

## 系统服务管理

### 启动服务

```bash
sudo systemctl start dir-monitor-go
```

### 停止服务

```bash
sudo systemctl stop dir-monitor-go
```

### 重启服务

```bash
sudo systemctl restart dir-monitor-go
```

### 查看服务状态

```bash
sudo systemctl status dir-monitor-go
```

### 查看服务日志

```bash
# 查看实时日志
sudo journalctl -u dir-monitor-go -f

# 查看最近的日志
sudo journalctl -u dir-monitor-go --since "1 hour ago"

# 查看特定日期的日志
sudo journalctl -u dir-monitor-go --since "2024-01-01" --until "2024-01-02"
```

### 禁用服务开机自启

```bash
sudo systemctl disable dir-monitor-go
```

### 启用服务开机自启

```bash
sudo systemctl enable dir-monitor-go
```

## Docker容器运行

### 构建Docker镜像

```bash
docker build -t dir-monitor-go .
```

### 运行Docker容器

```bash
# 基本运行
docker run -d \
  --name dir-monitor-go \
  -v /path/to/configs:/opt/dir-monitor-go/configs \
  -v /path/to/monitor:/path/to/monitor \
  dir-monitor-go

# 运行并映射日志目录
docker run -d \
  --name dir-monitor-go \
  -v /path/to/configs:/opt/dir-monitor-go/configs \
  -v /path/to/monitor:/path/to/monitor \
  -v /path/to/logs:/opt/dir-monitor-go/logs \
  dir-monitor-go

# 运行并指定配置文件
docker run -d \
  --name dir-monitor-go \
  -v /path/to/configs:/opt/dir-monitor-go/configs \
  -v /path/to/monitor:/path/to/monitor \
  dir-monitor-go --config /opt/dir-monitor-go/configs/custom-config.json
```

## 常见场景配置示例

### 监控日志文件

适用于监控应用程序日志文件并在文件更新时执行分析操作。

```json
{
  "version": "3.2.1",
  "metadata": {
    "name": "日志监控配置",
    "description": "监控应用程序日志文件并执行分析"
  },
  "monitors": [
    {
      "id": "log_watcher",
      "name": "日志监控器",
      "directory": "/var/log/application",
      "command": "/opt/scripts/analyze_log.sh",
      "file_patterns": ["*.log"],
      "timeout": 300,
      "debounce_seconds": 10,
      "enabled": true
    }
  ],
  "settings": {
    "log_level": "info",
    "log_file": "/var/log/dir-monitor-go/logs.log",
    "max_concurrent_operations": 3,
    "operation_timeout_seconds": 300
  }
}
```

### 处理上传文件

适用于处理通过SFTP或其他方式上传到指定目录的文件。

```json
{
  "version": "3.2.1",
  "metadata": {
    "name": "文件上传处理配置",
    "description": "处理上传到指定目录的文件"
  },
  "monitors": [
    {
      "id": "upload_processor",
      "name": "上传文件处理器",
      "directory": "/sftp/incoming",
      "command": "/opt/scripts/process_upload.sh",
      "file_patterns": ["*.csv", "*.xlsx", "*.pdf"],
      "timeout": 600,
      "schedule": "* 9-17 * * 1-5",
      "debounce_seconds": 30,
      "enabled": true
    }
  ],
  "settings": {
    "log_level": "info",
    "log_file": "/var/log/dir-monitor-go/upload.log",
    "log_max_size": 20971520,
    "log_max_backups": 5,
    "max_concurrent_operations": 3,
    "operation_timeout_seconds": 900,
    "retry_attempts": 3,
    "retry_delay_seconds": 10
  }
}
```

## 故障排除

### 服务无法启动

1. 检查配置文件是否存在且格式正确
2. 确认配置文件路径正确
3. 检查系统日志获取详细错误信息
4. 验证监控目录权限

### 监控不工作

1. 检查监控器是否启用
2. 确认目录路径和文件模式匹配
3. 查看日志文件了解具体错误
4. 验证调度时间配置

### 性能问题

1. 检查是否监控了过多文件
2. 调整并发操作数设置
3. 增加防抖时间避免频繁触发
4. 优化命令执行逻辑

### 日志相关问题

1. 检查日志文件路径和权限
2. 确认日志级别设置合适
3. 验证日志轮转配置
4. 检查磁盘空间是否充足

## 更新和维护

### 软件更新

```bash
# 停止服务
sudo systemctl stop dir-monitor-go

# 下载新版本（二进制安装方式）
wget https://github.com/yourusername/dir-monitor-go/releases/latest/download/dir-monitor-go-linux-amd64.tar.gz
tar -xzf dir-monitor-go-linux-amd64.tar.gz
sudo mv dir-monitor-go /usr/local/bin/

# 重启服务
sudo systemctl start dir-monitor-go
```

### 配置更新

```bash
# 编辑配置文件
sudo nano /opt/dir-monitor-go/configs/config.json

# 重启服务使配置生效
sudo systemctl restart dir-monitor-go
```

### 备份配置

```bash
# 备份配置文件
sudo cp /opt/dir-monitor-go/configs/config.json /opt/dir-monitor-go/configs/config.json.backup.$(date +%Y%m%d)

# 备份整个配置目录
sudo tar -czf dir-monitor-config-backup-$(date +%Y%m%d).tar.gz /opt/dir-monitor-go/configs/
```