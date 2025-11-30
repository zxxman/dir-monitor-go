# Dir-Monitor-Go 配置参考

> **版本**: v3.2.1  
> **最后更新**: 2025年10月16日

> 💡 **提示**: 查看完整的配置示例，请参考 [配置示例文档](CONFIG_EXAMPLE.md)。

## 📋 目录

1. [配置文件结构](#-配置文件结构)
2. [全局配置](#-全局配置)
3. [监控器配置](#-监控器配置)
4. [文件模式匹配](#-文件模式匹配)
5. [脚本执行配置](#-脚本执行配置)
6. [调度配置](#-调度配置)
7. [日志配置](#-日志配置)
8. [性能配置](#-性能配置)
9. [配置示例](#-配置示例)
10. [配置验证](#-配置验证)

---

## 📁 配置文件结构

### 基本结构
```json
{
  "version": "3.2.1",
  "global": {
    // 全局配置
  },
  "monitors": [
    // 监控器配置数组
  ]
}
```

### 配置文件版本
- **当前版本**: 3.2.1
- **兼容性**: 向后兼容3.x版本
- **升级注意**: 从2.x升级需要手动调整配置格式

---

## 🌍 全局配置

### 日志配置
```json
{
  "global": {
    "log_level": "info",
    "log_file": "/var/log/dir-monitor-go.log",
    "log_max_size": 100,
    "log_max_backups": 5,
    "log_max_age": 30
  }
}
```

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| log_level | string | "info" | 日志级别: debug, info, warn, error |
| log_file | string | "" | 日志文件路径，空则输出到控制台 |
| log_max_size | int | 100 | 日志文件最大大小(MB) |
| log_max_backups | int | 5 | 保留的备份日志文件数 |
| log_max_age | int | 30 | 日志文件保留天数 |

### 执行控制
```json
{
  "global": {
    "max_concurrent_executions": 10,
    "global_execution_lock": false,
    "lock_timeout": 300,
    "default_execution_timeout": 300
  }
}
```

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| max_concurrent_executions | int | 10 | 最大并发执行数 |
| global_execution_lock | bool | false | 是否启用全局执行锁 |
| lock_timeout | int | 300 | 锁超时时间(秒) |
| default_execution_timeout | int | 300 | 默认执行超时时间(秒) |

### 性能配置
```json
{
  "global": {
    "performance_monitoring": {
      "enabled": true,
      "report_interval": "1m",
      "metrics_retention": "24h"
    },
    "file_stability_check": {
      "enabled": true,
      "default_check_interval": "500ms",
      "default_stable_duration": "1s",
      "max_file_size_for_check": "1GB"
    }
  }
}
```

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| performance_monitoring.enabled | bool | true | 是否启用性能监控 |
| performance_monitoring.report_interval | string | "1m" | 性能报告间隔 |
| performance_monitoring.metrics_retention | string | "24h" | 指标保留时间 |
| file_stability_check.enabled | bool | true | 是否启用文件稳定性检查 |
| file_stability_check.default_check_interval | string | "500ms" | 默认检查间隔 |
| file_stability_check.default_stable_duration | string | "1s" | 默认稳定持续时间 |
| file_stability_check.max_file_size_for_check | string | "1GB" | 检查的最大文件大小 |

---

## 🔍 监控器配置

### 基本配置
```json
{
  "name": "example_monitor",
  "path": "/path/to/directory",
  "patterns": ["*.txt", "*.pdf"],
  "command": "process.sh {FILE_PATH}",
  "enabled": true
}
```

| 选项 | 类型 | 必需 | 默认值 | 描述 |
|------|------|------|--------|------|
| name | string | 是 | - | 监控器名称，必须唯一 |
| path | string | 是 | - | 监控目录路径 |
| patterns | array | 否 | ["*"] | 文件匹配模式数组 |
| command | string | 是 | - | 触发执行的命令 |
| enabled | bool | 否 | true | 是否启用此监控器 |

### 高级配置
```json
{
  "recursive": true,
  "ignore_patterns": ["*.tmp", ".*"],
  "debounce_time": 1000,
  "transfer_complete_check": true,
  "execution_timeout": 300,
  "events": ["create", "modify"]
}
```

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| recursive | bool | true | 是否递归监控子目录 |
| ignore_patterns | array | [] | 忽略的文件模式数组 |
| debounce_time | int | 500 | 防抖时间(毫秒) |
| transfer_complete_check | bool | true | 是否启用传输完成检测 |
| execution_timeout | int | 300 | 命令执行超时时间(秒) |
| events | array | ["create", "modify", "delete", "rename"] | 监控的事件类型 |

---

## 🎯 文件模式匹配

### 基本模式
```json
{
  "patterns": [
    "*.txt",           // 扩展名匹配
    "report_*.pdf",    // 前缀匹配
    "data_???.csv",    // 通配符匹配
    "image?.jpg"       // 单字符通配符
  ]
}
```

### 正则表达式
```json
{
  "patterns": [
    "\\d{4}-\\d{2}-\\d{2}.*",  // 日期格式
    "(?i)\\.(jpg|png)$",       // 不区分大小写
    "^[A-Z][a-z]+.*",          // 以大写字母开头
    "^test_.*_\\d{8}$"         // 特定命名格式
  ]
}
```

### 忽略模式
```json
{
  "ignore_patterns": [
    "*.tmp",           // 临时文件
    ".*",              // 隐藏文件
    "__MACOSX",        // Mac系统文件
    "Thumbs.db",       // Windows缩略图
    "*.bak", "*.swp"   // 备份文件
  ]
}
```

---

## 🔧 脚本执行配置

### 基本命令配置
```json
{
  "command": "/path/to/script.sh {FILE_PATH}",
  "execution_timeout": 300,
  "working_directory": "/tmp"
}
```

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| command | string | 必需 | 要执行的命令 |
| execution_timeout | int | 300 | 执行超时时间(秒) |
| working_directory | string | "" | 命令执行的工作目录 |

### 环境变量
```json
{
  "environment": {
    "CUSTOM_VAR": "value",
    "PATH": "/usr/local/bin:$PATH"
  }
}
```

### 执行模式
```json
{
  "execution_mode": "async",
  "retry_on_failure": true,
  "max_retries": 3,
  "retry_delay": 5000
}
```

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| execution_mode | string | "async" | 执行模式: async, sync |
| retry_on_failure | bool | false | 失败时是否重试 |
| max_retries | int | 3 | 最大重试次数 |
| retry_delay | int | 5000 | 重试延迟(毫秒) |

---

## ⏰ 调度配置

### 时间窗口
```json
{
  "schedule": {
    "time_windows": [
      {
        "start": "02:00",
        "end": "04:00",
        "days": ["monday", "tuesday", "wednesday", "thursday", "friday"]
      },
      {
        "start": "14:00",
        "end": "16:00",
        "days": ["saturday", "sunday"]
      }
    ],
    "timezone": "UTC"
  }
}
```

| 选项 | 类型 | 必需 | 描述 |
|------|------|------|------|
| time_windows | array | 是 | 时间窗口数组 |
| timezone | string | 否 | 时区设置，默认UTC |

### 时间窗口选项
| 选项 | 类型 | 必需 | 描述 |
|------|------|------|------|
| start | string | 是 | 开始时间(HH:MM格式) |
| end | string | 是 | 结束时间(HH:MM格式) |
| days | array | 否 | 执行日期，默认所有天 |

### 限制配置
```json
{
  "schedule": {
    "max_executions_per_hour": 10,
    "min_interval_between_executions": 300,
    "skip_holidays": true,
    "holiday_countries": ["US", "GB"]
  }
}
```

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| max_executions_per_hour | int | - | 每小时最大执行次数 |
| min_interval_between_executions | int | - | 执行间最小间隔(秒) |
| skip_holidays | bool | false | 是否跳过节假日 |
| holiday_countries | array | [] | 节假日国家代码 |

---

## 📝 日志配置

### 日志级别
```json
{
  "log_level": "info",
  "log_file": "/var/log/dir-monitor-go.log"
}
```

### 日志轮转
```json
{
  "log_max_size": 100,
  "log_max_backups": 5,
  "log_max_age": 30
}
```

### 结构化日志
```json
{
  "structured_logging": {
    "enabled": true,
    "format": "json",
    "include_timestamp": true,
    "include_level": true,
    "include_monitor": true
  }
}
```

---

## ⚡ 性能配置

### 并发控制
```json
{
  "max_concurrent_executions": 10,
  "execution_queue_size": 100,
  "worker_pool_size": 5
}
```

### 内存管理
```json
{
  "memory_limit": "512MB",
  "gc_percentage": 100,
  "max_buffer_size": "10MB"
}
```

### 文件检查优化
```json
{
  "file_stability_check": {
    "enabled": true,
    "check_interval": "500ms",
    "stable_duration": "1s",
    "max_file_size_for_check": "1GB",
    "skip_check_for_small_files": true,
    "small_file_threshold": "1MB"
  }
}
```

---

## 📋 配置示例

### 示例1：简单文件监控
```json
{
  "version": "3.2.1",
  "global": {
    "log_level": "info"
  },
  "monitors": [
    {
      "name": "downloads",
      "path": "/home/user/Downloads",
      "patterns": ["*.pdf", "*.docx"],
      "command": "echo 'New file: {FILE_PATH}'"
    }
  ]
}
```

### 示例2：高级配置
```json
{
  "version": "3.2.1",
  "global": {
    "log_level": "debug",
    "log_file": "/var/log/dir-monitor-go.log",
    "max_concurrent_executions": 5,
    "performance_monitoring": {
      "enabled": true,
      "report_interval": "5m"
    }
  },
  "monitors": [
    {
      "name": "uploads",
      "path": "/var/uploads",
      "patterns": ["*"],
      "ignore_patterns": ["*.tmp", ".*"],
      "command": "/usr/local/bin/process-upload.sh {FILE_PATH}",
      "recursive": true,
      "debounce_time": 2000,
      "transfer_complete_check": true,
      "execution_timeout": 600,
      "execution_mode": "async",
      "retry_on_failure": true,
      "max_retries": 3,
      "schedule": {
        "time_windows": [
          {
            "start": "02:00",
            "end": "06:00"
          }
        ],
        "timezone": "UTC"
      }
    }
  ]
}
```

### 示例3：多监控器配置
```json
{
  "version": "3.2.1",
  "global": {
    "log_level": "info",
    "max_concurrent_executions": 10
  },
  "monitors": [
    {
      "name": "documents",
      "path": "/home/user/Documents",
      "patterns": ["*.pdf", "*.docx", "*.xlsx"],
      "command": "organize-docs.sh {FILE_PATH}",
      "debounce_time": 1000
    },
    {
      "name": "images",
      "path": "/home/user/Pictures",
      "patterns": ["*.jpg", "*.png", "*.gif"],
      "command": "process-image.sh {FILE_PATH}",
      "debounce_time": 2000
    },
    {
      "name": "logs",
      "path": "/var/log/app",
      "patterns": ["*.log"],
      "command": "log-processor.sh {FILE_PATH}",
      "ignore_patterns": ["*.tmp"],
      "schedule": {
        "time_windows": [
          {
            "start": "01:00",
            "end": "03:00"
          }
        ]
      }
    }
  ]
}
```

---

## ✅ 配置验证

### 验证命令
```bash
# 验证配置文件语法
dir-monitor-go -config config.json -validate

# 验证并显示详细错误
dir-monitor-go -config config.json -validate -verbose
```

### 常见验证错误

#### 错误1：监控器名称重复
```
错误: 监控器名称"uploads"已存在
解决: 确保每个监控器有唯一的名称
```

#### 错误2：路径不存在
```
错误: 监控路径"/nonexistent/path"不存在
解决: 确保监控路径存在且有访问权限
```

#### 错误3：无效的正则表达式
```
错误: 文件模式"[invalid"不是有效的正则表达式
解决: 检查正则表达式语法
```

#### 错误4：无效的时间格式
```
错误: 时间窗口开始时间"25:00"格式无效
解决: 使用HH:MM格式，小时范围00-23
```

---

## 📚 更多资源

- [用户使用指南](USER_GUIDE.md)
- [API文档](API.md)
- [开发指南](DEVELOPMENT.md)
- [部署指南](DEPLOYMENT.md)