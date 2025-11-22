# Dir-Monitor-Go

一个高性能的目录监控工具，基于fsnotify库实现实时文件系统监控，能够智能检测文件传输完成并执行自动化脚本。

## 特性

- 实时文件系统监控，零延迟事件驱动
- 传输完成检测，智能识别文件传输完成
- 单一配置结构，简化配置管理
- 支持正则表达式匹配文件和目录模式
- 统一配置管理，简化的JSON配置
- 完整日志系统，自动轮转和清理
- 错误恢复机制，健壮的错误处理和自动恢复
- 高级调度功能，支持cron表达式的时间监控
- 多用户支持，独立的白天和夜间处理规则
- 并发操作控制，防止系统过载
- 文件稳定性检测，可靠处理大文件
- 防抖机制，避免短时间内重复触发
- 性能监控，实时监控系统运行状态

## 安装

### 二进制安装

```bash
# 下载预编译的二进制文件
wget https://github.com/your-repo/dir-monitor-go/releases/download/v3.2.1/dir-monitor-go
chmod +x dir-monitor-go

# 移动到系统路径
sudo mv dir-monitor-go /usr/local/bin/
```

### 源码安装

```bash
# 克隆项目
git clone https://github.com/your-repo/dir-monitor-go.git
cd dir-monitor-go

# 构建项目（需要Go 1.21+）
go build -o dir-monitor-go ./cmd/dir-monitor-go

# 或者使用make
make build

# 运行应用
./dir-monitor-go
```

### 系统服务安装

```bash
# 使用提供的安装脚本
sudo ./deploy/install-service.sh

# 启动服务
sudo systemctl start dir-monitor-go

# 设置开机自启
sudo systemctl enable dir-monitor-go
```

## 配置

配置文件使用简化的JSON格式，主要包含监控器配置和系统设置两部分。系统支持多用户场景，可以为不同用户配置不同的处理规则和时间段。

### 配置文件结构

```json
{
  "version": "3.2.1",
  "monitors": [
    {
      "directory": "/path/to/monitor",
      "command": "echo 'File processed: ${FILE_PATH}'",
      "file_patterns": ["*.csv", "*.xlsx"],
      "timeout": 300,
      "schedule": "* 11-21 * * 1-5",
      "enabled": true,
      "debounce_seconds": 2
    }
  ],
  "settings": {
    "log_file": "logs/dir-monitor-go.log",
    "log_level": "info",
    "log_show_caller": false,
    "log_max_size": 10485760,
    "max_concurrent_operations": 10,
    "operation_timeout_seconds": 600,
    "file_watcher_buffer_size": 1024,
    "file_stability_check_interval_ms": 200,
    "file_stability_timeout_seconds": 5,
    "min_stability_time_ms": 500,
    "event_channel_buffer_size": 100,
    "default_debounce_seconds": 1,
    "execution_dedup_interval_seconds": 5,
    "directory_stability_quiet_ms": 1000,
    "directory_stability_timeout_seconds": 10,
    "small_file_threshold": 1048576,
    "medium_file_threshold": 10485760,
    "large_file_threshold": 104857600,
    "default_rule_debounce_seconds": 60,
    "default_rule_priority": 5
  }
}
```

### 配置说明

#### monitors
监控器配置定义了目录监控的完整设置：

- `id`: 监控器唯一标识符
- `name`: 监控器名称
- `directory`: 需要监控的目录路径
- `command`: 文件匹配时执行的命令
- `file_patterns`: 文件匹配模式列表（支持glob模式）
- `timeout`: 命令执行超时时间（秒）
- `schedule`: Cron表达式，指定时间执行（可选，如"* 15-21 * * 1-5"表示周一至周五15-21点）
- `enabled`: 是否启用该监控项
- `debounce_seconds`: 防抖时间（秒），防止短时间内重复触发

#### settings
系统设置配置：

- `log_file`: 日志文件路径
- `log_level`: 日志级别（debug, info, warn, error）
- `log_show_caller`: 是否显示调用者信息
- `log_max_size`: 日志文件最大大小（字节）
- `log_compress`: 是否压缩日志文件
- `max_concurrent_operations`: 最大并发操作数，防止系统过载
- `operation_timeout_seconds`: 操作超时时间（秒）
- `file_stability_check_interval_ms`: 文件稳定性检查间隔（毫秒）
- `file_stability_timeout_seconds`: 文件稳定性超时时间（秒）
- `min_stability_time_ms`: 最小稳定性时间（毫秒）
- `event_channel_buffer_size`: 事件通道缓冲区大小
- `default_debounce_seconds`: 默认防抖时间（秒）
- `execution_dedup_interval_seconds`: 执行去重时间窗口（秒）
- `directory_stability_quiet_ms`: 目录稳定性静默时间（毫秒）
- `directory_stability_timeout_seconds`: 目录稳定性超时时间（秒）
- `small_file_threshold`: 小文件阈值（字节，默认1MB）
- `medium_file_threshold`: 中等文件阈值（字节，默认10MB）
- `large_file_threshold`: 大文件阈值（字节，默认100MB）
- `default_rule_debounce_seconds`: 默认规则防抖时间（秒）
- `performance_monitoring`: 是否启用性能监控

### 变量替换

在命令中可以使用以下变量：

- `${FILE_PATH}`: 文件完整路径
- `${FILE_NAME}`: 文件名
- `${FILE_DIR}`: 文件所在目录
- `${EVENT_TYPE}`: 事件类型（create, modify, delete等）
- `${EVENT_TIME}`: 事件时间戳
- `${RULE_ID}`: 规则ID
- `${RULE_NAME}`: 规则名称

## 使用方法

### 基本使用

```bash
# 使用默认配置文件运行
./dir-monitor-go

# 使用自定义配置文件运行
./dir-monitor-go --config /path/to/config.json

# 调试模式运行
./dir-monitor-go --config /path/to/config.json --log-level DEBUG

# 验证配置文件
./dir-monitor-go --config /path/to/config.json --dry-run

# 使用停止标记文件（测试用）
./dir-monitor-go --config /path/to/config.json --stop-file /tmp/stop_marker
```

### 高级配置示例

多用户SFTP目录监控配置，支持白天和夜间不同处理规则：

```json
{
  "version": "3.2.1",
  "monitors": [
    {
      "id": "user2_daytime",
      "name": "用户2白天CSV处理",
      "directory": "/sftp/user2/data",
      "command": "/opt/python-envs/bin/python /sftp/user2.py",
      "file_patterns": ["*.csv"],
      "timeout": 300,
      "schedule": "* 15-21 * * 1-5",
      "debounce_seconds": 15,
      "enabled": true
    },
    {
      "id": "user2_nighttime",
      "name": "用户2夜间CSV处理",
      "directory": "/sftp/user2/data",
      "command": "/opt/python-envs/bin/python /sftp/user2N.py",
      "file_patterns": ["*.csv"],
      "timeout": 300,
      "schedule": "* 0-9 * * 1-6",
      "debounce_seconds": 15,
      "enabled": true
    },
    {
      "id": "user3_daytime",
      "name": "用户3白天Excel处理",
      "directory": "/sftp/user3/data",
      "command": "/opt/python-envs/bin/python /sftp/user3.py",
      "file_patterns": ["*.xls", "*.xlsx"],
      "timeout": 300,
      "schedule": "* 15-21 * * 1-5",
      "debounce_seconds": 15,
      "enabled": true
    },
    {
      "id": "outbound_daytime",
      "name": "白天数据发送",
      "directory": "/sftp/out",
      "command": "/opt/mail-tool/mail-tool -t DAILY_DATA1",
      "file_patterns": ["*.csv"],
      "timeout": 300,
      "schedule": "* 15-21 * * 1-5",
      "debounce_seconds": 15,
      "enabled": true
    }
  ],
  "settings": {
    "log_file": "logs/dir-monitor-go.log",
    "log_level": "info",
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
    "large_file_threshold": 104857600,
    "default_rule_debounce_seconds": 60,
    "performance_monitoring": true
  }
}
```

**配置说明：**
- 用户2：白天（15-21点，周一至周五）使用`user2.py`处理CSV文件，夜间（0-9点，周一至周六）使用`user2N.py`处理
- 用户3：白天处理Excel文件（.xls和.xlsx格式）
-  outbound：白天时段自动发送数据邮件
- 所有监控项都设置了15秒防抖时间，避免重复处理
- 启用了性能监控和日志压缩功能

## 架构设计

### 核心组件

1. **FsnotifyWatcher**: 基于fsnotify库的文件系统监控器，实时监听目录变化
2. **Monitor**: 目录监控器实现，管理多个监控任务
3. **ScriptRunner**: 脚本执行器，负责执行配置的命令
4. **Config**: 配置管理器，处理配置文件加载和验证
5. **Logger**: 日志系统，支持多级别日志和自动轮转
6. **Scheduler**: 调度器，基于cron表达式的时间调度

### 工作流程

1. **初始化阶段**: 加载配置文件，验证配置项，初始化监控器
2. **监控阶段**: 
   - 监控指定目录的文件系统事件（创建、修改、删除等）
   - 对事件进行过滤和去重处理
   - 检查文件是否稳定（传输完成检测）
3. **调度阶段**: 检查cron表达式，确定当前时间是否允许执行
4. **匹配阶段**: 匹配监控规则和文件模式
5. **执行阶段**: 执行相应的命令，处理执行结果和错误
6. **日志阶段**: 记录操作日志，支持性能监控

### 性能优化

- **防抖机制**: 避免短时间内重复触发相同事件
- **并发控制**: 限制最大并发操作数，防止系统过载
- **文件稳定性检测**: 确保大文件传输完成后再处理
- **内存优化**: 合理设置缓冲区大小，避免内存泄漏

## 命令行参数

- `--config`: 配置文件路径（默认: configs/config.json）
- `--stop-file`: 停止标记文件路径（测试用）
- `--version`: 显示版本信息
- `--dry-run`: 仅验证配置，不启动实际监控
- `--log-level`: 设置日志级别（DEBUG, INFO, WARN, ERROR）
- `--log-file`: 设置日志文件路径
- `--validate-config`: 验证配置文件格式和参数
- `--show-config`: 显示当前配置信息
- `--performance-monitor`: 启用性能监控模式
- `--max-concurrent`: 设置最大并发操作数（覆盖配置）
- `--daemon`: 以守护进程模式运行
- `--pid-file`: 指定PID文件路径（守护进程模式）

## 日志系统

### 日志特性

- **多级别日志**: 支持DEBUG, INFO, WARN, ERROR四个级别
- **结构化日志**: 支持JSON格式输出，便于日志分析
- **自动轮转**: 基于文件大小的自动日志轮转
- **日志压缩**: 支持轮转后日志文件的自动压缩
- **调用链追踪**: 可选的函数调用者信息显示
- **性能监控**: 集成性能指标记录

### 日志配置示例

```json
{
  "settings": {
    "log_file": "logs/dir-monitor-go.log",
    "log_level": "info",
    "log_show_caller": true,
    "log_max_size": 10485760,
    "log_compress": true,
    "performance_monitoring": true
  }
}
```

### 日志输出格式

```
2025-10-31 18:19:41 [INFO] [monitor.go:245] 监控器 user2_daytime 检测到文件变化: /sftp/user2/data/report.csv
2025-10-31 18:19:41 [DEBUG] [scheduler.go:89] 调度器检查: user2_daytime 调度表达式 "* 15-21 * * 1-5" 匹配结果: true
2025-10-31 18:19:56 [INFO] [runner.go:123] 命令执行完成，耗时: 15.2s，返回码: 0
```

## 版本历史

### v3.2.1 (2025-10-31)
- 修复调度器时间匹配逻辑
- 优化文件稳定性检测算法
- 增强多用户场景支持
- 改进性能监控功能
- 修复已知问题和边界条件

### v3.2.0 (2025-10-31)
- 新增配置文件元数据支持
- 增强监控项配置（ID、名称、描述、优先级）
- 优化目录稳定性处理机制
- 扩展日志配置选项（压缩、调用链追踪）
- 增加重试机制和健康检查
- 改进配置验证逻辑
- 新增性能监控和指标收集

### v3.1.2 (2025-10-31)
- 修复文件稳定性检查逻辑
- 增强日志记录
- 简化配置结构
- 移除冗余代码和配置项

### v3.1.1 (2025-10-25)
- 优化性能
- 修复已知问题

### v3.0.0
- 重构配置结构
- 移除复杂的多层级配置
- 简化操作执行逻辑

## 许可证

MIT License