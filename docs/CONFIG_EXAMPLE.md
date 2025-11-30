# Dir-Monitor-Go 配置示例

> **版本**: v3.2.1  
> **最后更新**: 2025年10月16日

## 📋 说明

本文档提供了一个完整的Dir-Monitor-Go配置文件示例，展示了各种配置选项的实际用法。这个示例文件可以直接作为配置文件使用，只需根据实际需求修改相应参数。

> 💡 **提示**: 有关所有配置选项的详细说明，请参考 [配置参考文档](CONFIG.md)。

## 🚀 快速开始

1. 复制此示例文件为您的配置文件：
   ```bash
   cp CONFIG_EXAMPLE.md /etc/dir-monitor-go/config.json
   ```

2. 根据您的需求修改配置参数

3. 验证配置文件：
   ```bash
   dir-monitor-go -config /etc/dir-monitor-go/config.json -validate
   ```

4. 启动服务：
   ```bash
   dir-monitor-go -config /etc/dir-monitor-go/config.json
   ```

## 📝 配置文件
  // 配置文件版本，必须与当前软件版本匹配
  "version": "3.2.1",
  
  // 全局配置，适用于所有监控器
  "global": {
    // 日志配置
    "log_level": "info",
    "log_file": "/var/log/dir-monitor-go/app.log",
    "log_max_size": 100,
    "log_max_backups": 5,
    "log_max_age": 30,
    
    // 执行控制
    "max_concurrent_executions": 5,
    "global_execution_lock": false,
    "lock_timeout": 300,
    "default_execution_timeout": 300,
    
    // 性能配置
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
  },
  
  // 监控器配置数组，每个对象定义一个监控任务
  "monitors": [
    {
      // 监控器唯一标识符
      "id": "web-assets-watch",
      // 监控器名称，用于日志和UI显示
      "name": "Web Assets Monitor",
      // 监控器描述
      "description": "Monitor web assets directory for changes and rebuild frontend",
      // 监控目录路径
      "path": "/var/www/assets",
      // 文件变化时执行的命令
      "command": "cd /var/www && npm run build",
      // 匹配的文件模式数组
      "patterns": [
        "*.js",
        "*.css",
        "*.scss",
        "*.html"
      ],
      // 忽略的文件模式数组
      "ignore_patterns": [
        "*.tmp",
        "node_modules/**/*"
      ],
      // 命令执行超时时间（秒）
      "execution_timeout": 30,
      // 是否启用此监控器
      "enabled": true,
      // 防抖时间（毫秒），文件变化后等待时间
      "debounce_time": 2000,
      // 是否递归监控子目录
      "recursive": true,
      // 监控的文件事件类型
      "events": ["create", "modify", "delete"],
      // 调度配置，可选
      "schedule": {
        "time_windows": [
          {
            "start": "02:00",
            "end": "04:00",
            "days": ["monday", "tuesday", "wednesday", "thursday", "friday"]
          }
        ],
        "timezone": "UTC"
      }
    },
    {
      "id": "config-reload",
      "name": "Configuration Monitor",
      "description": "Monitor configuration files and reload services",
      "path": "/etc/myapp",
      "command": "systemctl reload myapp",
      "patterns": [
        "*.conf",
        "*.yaml",
        "*.yml",
        "*.json"
      ],
      "ignore_patterns": [
        "*.tmp",
        "*.bak"
      ],
      "execution_timeout": 10,
      "enabled": true,
      "debounce_time": 5000,
      "recursive": true,
      "events": ["modify"],
      "execution_mode": "sync",
      "retry_on_failure": true,
      "max_retries": 3,
      "retry_delay": 5000
    },
    {
      "id": "log-rotation",
      "name": "Log Rotation Monitor",
      "description": "Monitor log directory and trigger rotation when needed",
      "path": "/var/log/myapp",
      "command": "/usr/local/sbin/rotate-logs.sh",
      "patterns": [
        "*.log"
      ],
      "execution_timeout": 60,
      "enabled": true,
      "debounce_time": 10000,
      "recursive": false,
      "events": ["modify"],
      "schedule": "0 2 * * *"
    },
    {
      "id": "backup-trigger",
      "name": "Backup Trigger",
      "description": "Monitor data directory and trigger backup on changes",
      "path": "/data/myapp",
      "command": "/usr/local/bin/backup-data.sh",
      "patterns": [
        "*.db",
        "*.sqlite",
        "*.data"
      ],
      "ignore_patterns": [
        "*.tmp",
        "*.lock"
      ],
      "execution_timeout": 300,
      "enabled": true,
      "debounce_time": 30000,
      "recursive": true,
      "events": ["create", "modify"],
      "schedule": "0 3 * * 0",
      "working_directory": "/tmp",
      "environment": {
        "BACKUP_DIR": "/backup/myapp",
        "BACKUP_RETENTION": "7d"
      }
    }
  ]
}

## 📚 更多示例

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

## 🔗 相关文档

- [配置参考文档](CONFIG.md) - 详细的配置选项说明
- [用户使用指南](USER_GUIDE.md) - 完整的使用指南
- [API文档](API.md) - REST API参考
- [部署指南](DEPLOYMENT.md) - 生产环境部署指南

## ⚠️ 注意事项

1. **配置文件格式**: 配置文件必须是有效的JSON格式
2. **路径权限**: 确保Dir-Monitor-Go进程有权限访问监控目录和执行命令
3. **命令安全**: 避免在命令中直接使用用户输入，以防命令注入
4. **资源限制**: 合理设置并发执行数和超时时间，避免系统资源耗尽
5. **日志轮转**: 配置适当的日志轮转策略，避免日志文件过大

## 🐛 故障排除

### 常见问题

1. **配置验证失败**
   - 检查JSON格式是否正确
   - 确认所有必需字段都已填写
   - 验证路径和命令是否有效

2. **监控不工作**
   - 检查目录路径是否存在
   - 确认文件模式是否正确
   - 查看日志文件获取详细错误信息

3. **命令执行失败**
   - 检查命令路径和权限
   - 确认工作目录设置正确
   - 验证环境变量配置

更多故障排除信息，请参考 [常见问题文档](FAQ.md)。