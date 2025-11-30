# Dir-Monitor-Go 常见问题

> **版本**: v3.2.1  
> **最后更新**: 2025年10月16日

## 📋 目录

1. [安装与配置问题](#-安装与配置问题)
2. [监控问题](#-监控问题)
3. [命令执行问题](#-命令执行问题)
4. [性能问题](#-性能问题)
5. [日志问题](#-日志问题)
6. [权限问题](#-权限问题)
7. [网络问题](#-网络问题)
8. [高可用问题](#-高可用问题)
9. [其他问题](#-其他问题)

---

## 🛠️ 安装与配置问题

### Q: 如何安装Dir-Monitor-Go？

**A**: 有多种安装方式：

1. **二进制安装**（推荐）
   ```bash
   # 下载适合您系统的二进制文件
   wget https://github.com/your-repo/dir-monitor-go/releases/download/v3.2.1/dir-monitor-go-linux-amd64-v3.2.1.tar.gz
   
   # 解压
   tar -xzf dir-monitor-go-linux-amd64-v3.2.1.tar.gz
   
   # 复制到系统路径
   sudo cp dir-monitor-go /usr/local/bin/
   sudo chmod +x /usr/local/bin/dir-monitor-go
   ```

2. **源码编译**
   ```bash
   git clone https://github.com/your-repo/dir-monitor-go.git
   cd dir-monitor-go
   go build -o dir-monitor-go cmd/dir-monitor/main.go
   ```

3. **Docker安装**
   ```bash
   docker pull dirmonitor/go:v3.2.1
   ```

### Q: 配置文件在哪里？

**A**: 默认情况下，Dir-Monitor-Go会在以下位置查找配置文件：

1. 当前目录下的`config.json`
2. `/etc/dir-monitor-go/config.json`
3. `~/.dir-monitor-go/config.json`

您也可以使用`-config`参数指定配置文件路径：
```bash
dir-monitor-go -config /path/to/your/config.json
```

### Q: 如何验证配置文件是否正确？

**A**: 您可以使用内置的配置验证功能：

```bash
dir-monitor-go -validate -config /path/to/config.json
```

或者使用在线配置验证工具：https://dir-monitor.example.com/validator

### Q: 如何创建一个基本的配置文件？

**A**: 您可以使用以下命令生成示例配置文件：

```bash
dir-monitor-go -example-config > config.json
```

或者手动创建一个简单的配置文件：

```json
{
  "version": "3.2.1",
  "monitors": [
    {
      "name": "file-monitor",
      "path": "/path/to/monitor",
      "command": "echo 'File changed: {FILE_PATH}'",
      "patterns": ["*"],
      "recursive": true
    }
  ],
  "settings": {
    "log_level": "info",
    "max_concurrent": 5
  }
}
```

---

## 👀 监控问题

### Q: 为什么我的目录没有被监控？

**A**: 请检查以下几点：

1. **权限问题**
   ```bash
   # 检查目录权限
   ls -la /path/to/monitor
   
   # 确保用户有读取权限
   sudo usermod -a -G $USER $(stat -c "%G" /path/to/monitor)
   ```

2. **路径是否正确**
   ```bash
   # 检查路径是否存在
   ls /path/to/monitor
   
   # 检查路径是否为绝对路径
   pwd
   ```

3. **inotify限制**
   ```bash
   # 检查inotify限制
   cat /proc/sys/fs/inotify/max_user_watches
   
   # 增加inotify限制（临时）
   echo 8192 | sudo tee /proc/sys/fs/inotify/max_user_watches
   
   # 增加inotify限制（永久）
   echo fs.inotify.max_user_watches=8192 | sudo tee -a /etc/sysctl.conf
   sudo sysctl -p
   ```

### Q: 如何监控多个目录？

**A**: 在配置文件中添加多个监控器：

```json
{
  "version": "3.2.1",
  "monitors": [
    {
      "name": "documents",
      "path": "/home/user/documents",
      "command": "echo 'Document changed: {FILE_PATH}'",
      "patterns": ["*.doc", "*.pdf"],
      "recursive": true
    },
    {
      "name": "downloads",
      "path": "/home/user/downloads",
      "command": "echo 'Download changed: {FILE_PATH}'",
      "patterns": ["*"],
      "recursive": false
    }
  ],
  "settings": {
    "log_level": "info",
    "max_concurrent": 5
  }
}
```

### Q: 如何只监控特定类型的文件？

**A**: 使用文件模式匹配：

```json
{
  "monitors": [
    {
      "name": "image-monitor",
      "path": "/path/to/images",
      "command": "process-image.sh {FILE_PATH}",
      "patterns": ["*.jpg", "*.png", "*.gif"],
      "recursive": true
    }
  ]
}
```

您也可以使用排除模式：

```json
{
  "monitors": [
    {
      "name": "log-monitor",
      "path": "/var/log",
      "command": "process-log.sh {FILE_PATH}",
      "include_patterns": ["*.log"],
      "exclude_patterns": ["*.tmp", "*.bak"],
      "recursive": false
    }
  ]
}
```

### Q: 如何避免重复触发？

**A**: 使用防抖设置：

```json
{
  "monitors": [
    {
      "name": "stable-monitor",
      "path": "/path/to/monitor",
      "command": "process-file.sh {FILE_PATH}",
      "patterns": ["*"],
      "recursive": true,
      "debounce": {
        "enabled": true,
        "delay": "5s",
        "max_wait": "30s"
      }
    }
  ]
}
```

---

## ⚙️ 命令执行问题

### Q: 为什么我的命令没有执行？

**A**: 请检查以下几点：

1. **命令路径是否正确**
   ```bash
   # 使用绝对路径
   "command": "/usr/bin/python3 /path/to/script.py {FILE_PATH}"
   
   # 或者确保命令在PATH中
   "command": "python3 /path/to/script.py {FILE_PATH}"
   ```

2. **命令是否有执行权限**
   ```bash
   # 检查脚本权限
   ls -la /path/to/script.sh
   
   # 添加执行权限
   chmod +x /path/to/script.sh
   ```

3. **命令是否需要特殊环境**
   ```json
   {
     "monitors": [
       {
         "name": "env-monitor",
         "path": "/path/to/monitor",
         "command": "process.sh {FILE_PATH}",
         "env": {
           "PYTHONPATH": "/usr/lib/python3.8",
           "LD_LIBRARY_PATH": "/usr/local/lib"
         }
       }
     ]
   }
   ```

### Q: 如何传递多个参数给命令？

**A**: 使用占位符和引号：

```json
{
  "monitors": [
    {
      "name": "multi-arg-monitor",
      "path": "/path/to/monitor",
      "command": "process.sh \"{FILE_PATH}\" \"{FILE_NAME}\" \"{FILE_DIR}\"",
      "patterns": ["*"],
      "recursive": true
    }
  ]
}
```

### Q: 如何设置命令执行超时？

**A**: 在配置中设置超时：

```json
{
  "monitors": [
    {
      "name": "timeout-monitor",
      "path": "/path/to/monitor",
      "command": "long-running-task.sh {FILE_PATH}",
      "timeout": "60s",
      "patterns": ["*"],
      "recursive": true
    }
  ]
}
```

### Q: 如何限制并发执行的命令数量？

**A**: 在全局设置中设置最大并发数：

```json
{
  "settings": {
    "max_concurrent": 3
  },
  "monitors": [
    {
      "name": "concurrent-monitor",
      "path": "/path/to/monitor",
      "command": "process.sh {FILE_PATH}",
      "patterns": ["*"],
      "recursive": true
    }
  ]
}
```

---

## 🚀 性能问题

### Q: Dir-Monitor-Go占用太多内存怎么办？

**A**: 尝试以下优化：

1. **减少监控目录深度**
   ```json
   {
     "monitors": [
       {
         "recursive": false,
         "path": "/data/level1"
       }
     ]
   }
   ```

2. **优化文件过滤**
   ```json
   {
     "monitors": [
       {
         "include_patterns": ["*.txt"],
         "exclude_patterns": ["temp_*", "*.tmp"]
       }
     ]
   }
   ```

3. **调整并发执行数**
   ```json
   {
     "settings": {
       "max_concurrent": 2
     }
   }
   ```

4. **启用内存限制**
   ```json
   {
     "settings": {
       "memory_limit": "256MB"
     }
   }
   ```

### Q: 如何监控大量文件而不影响性能？

**A**: 使用以下策略：

1. **分批监控**
   ```json
   {
     "monitors": [
       {
         "name": "batch-1",
         "path": "/data/part1",
         "command": "process.sh {FILE_PATH}",
         "max_events": 100
       },
       {
         "name": "batch-2",
         "path": "/data/part2",
         "command": "process.sh {FILE_PATH}",
         "max_events": 100
       }
     ]
   }
   ```

2. **使用事件批处理**
   ```json
   {
     "monitors": [
       {
         "name": "batch-monitor",
         "path": "/path/to/monitor",
         "command": "process-batch.sh",
         "batch": {
           "enabled": true,
           "size": 10,
           "timeout": "5s"
         }
       }
     ]
   }
   ```

3. **调整事件处理间隔**
   ```json
   {
     "monitors": [
       {
         "name": "throttled-monitor",
         "path": "/path/to/monitor",
         "command": "process.sh {FILE_PATH}",
         "throttle": {
           "enabled": true,
           "interval": "1s",
           "burst": 5
         }
       }
     ]
   }
   ```

---

## 📝 日志问题

### Q: 如何查看详细日志？

**A**: 调整日志级别：

```json
{
  "settings": {
    "log_level": "debug"
  }
}
```

或者使用命令行参数：
```bash
dir-monitor-go -log-level debug -config /path/to/config.json
```

### Q: 如何将日志输出到文件？

**A**: 在配置中设置日志文件：

```json
{
  "settings": {
    "log_file": "/var/log/dir-monitor-go/app.log",
    "log_max_size": 100,
    "log_max_backups": 5,
    "log_max_age": 30
  }
}
```

### Q: 如何启用结构化日志？

**A**: 设置日志格式为JSON：

```json
{
  "settings": {
    "log_format": "json",
    "log_fields": ["timestamp", "level", "message", "monitor", "file_path"]
  }
}
```

---

## 🔐 权限问题

### Q: 如何以非root用户运行Dir-Monitor-Go？

**A**: 使用以下方法：

1. **创建专用用户**
   ```bash
   sudo useradd -r -s /bin/false dirmonitor
   sudo usermod -a -G dirmonitor $USER
   ```

2. **设置正确的文件权限**
   ```bash
   sudo chown -R dirmonitor:dirmonitor /etc/dir-monitor-go
   sudo chown -R dirmonitor:dirmonitor /var/log/dir-monitor-go
   ```

3. **修改systemd服务文件**
   ```ini
   [Service]
   User=dirmonitor
   Group=dirmonitor
   ```

### Q: 如何监控需要root权限的目录？

**A**: 使用以下方法：

1. **使用sudo（不推荐）**
   ```bash
   sudo dir-monitor-go -config /path/to/config.json
   ```

2. **设置capabilities（推荐）**
   ```bash
   sudo setcap cap_dac_read_search+ep /usr/local/bin/dir-monitor-go
   ```

3. **使用ACL**
   ```bash
   sudo setfacl -R -m u:dirmonitor:rx /path/to/monitor
   ```

---

## 🌐 网络问题

### Q: 如何通过API远程控制Dir-Monitor-Go？

**A**: 启用API服务：

```json
{
  "api": {
    "enabled": true,
    "address": "0.0.0.0:8080",
    "auth": {
      "enabled": true,
      "username": "admin",
      "password": "password"
    }
  }
}
```

### Q: 如何设置API访问认证？

**A**: 配置API认证：

```json
{
  "api": {
    "enabled": true,
    "address": "0.0.0.0:8080",
    "auth": {
      "enabled": true,
      "type": "basic",
      "username": "admin",
      "password": "secure_password"
    }
  }
}
```

或者使用JWT令牌：
```json
{
  "api": {
    "enabled": true,
    "address": "0.0.0.0:8080",
    "auth": {
      "enabled": true,
      "type": "jwt",
      "secret": "your_jwt_secret",
      "expiration": "24h"
    }
  }
}
```

---

## 🔄 高可用问题

### Q: 如何实现主备模式？

**A**: 配置主备节点：

主节点配置：
```json
{
  "role": "primary",
  "ha": {
    "enabled": true,
    "node_id": "node-1",
    "peer_nodes": ["node-2"],
    "heartbeat_interval": "5s",
    "failover_timeout": "15s"
  }
}
```

备节点配置：
```json
{
  "role": "secondary",
  "ha": {
    "enabled": true,
    "node_id": "node-2",
    "peer_nodes": ["node-1"],
    "heartbeat_interval": "5s",
    "failover_timeout": "15s"
  }
}
```

### Q: 如何实现负载均衡？

**A**: 使用负载均衡器：

Nginx配置：
```nginx
upstream dir_monitor_go {
    server 10.0.1.10:8080;
    server 10.0.1.11:8080;
    server 10.0.1.12:8080;
}

server {
    listen 80;
    server_name dir-monitor.example.com;
    
    location / {
        proxy_pass http://dir_monitor_go;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

---

## ❓ 其他问题

### Q: 如何获取帮助？

**A**: 有多种获取帮助的方式：

1. **查看帮助文档**
   ```bash
   dir-monitor-go -help
   ```

2. **查看在线文档**
   https://dir-monitor.example.com/docs

3. **提交问题**
   https://github.com/your-repo/dir-monitor-go/issues

4. **社区讨论**
   https://github.com/your-repo/dir-monitor-go/discussions

### Q: 如何贡献代码？

**A**: 请参考[贡献指南](CONTRIBUTING.md)。

### Q: 如何报告安全漏洞？

**A**: 请发送邮件至security@dir-monitor.example.com，不要在公开的问题跟踪器中报告安全漏洞。

### Q: Dir-Monitor-Go是否支持Windows？

**A**: 是的，Dir-Monitor-Go支持Windows、Linux和macOS。请注意，不同平台的文件系统监控机制可能有所不同。

### Q: 如何升级到新版本？

**A**: 升级步骤：

1. **备份配置**
   ```bash
   cp /etc/dir-monitor-go/config.json /etc/dir-monitor-go/config.json.bak
   ```

2. **停止服务**
   ```bash
   sudo systemctl stop dir-monitor-go
   ```

3. **下载新版本**
   ```bash
   wget https://github.com/your-repo/dir-monitor-go/releases/download/v3.2.1/dir-monitor-go-linux-amd64-v3.2.1.tar.gz
   tar -xzf dir-monitor-go-linux-amd64-v3.2.1.tar.gz
   sudo cp dir-monitor-go-linux-amd64-v3.2.1/dir-monitor-go /usr/local/bin/
   ```

4. **验证配置**
   ```bash
   dir-monitor-go -validate -config /etc/dir-monitor-go/config.json
   ```

5. **启动服务**
   ```bash
   sudo systemctl start dir-monitor-go
   ```

---

## 📚 更多资源

- [用户使用指南](USER_GUIDE.md)
- [配置参考](CONFIG.md)
- [API文档](API.md)
- [开发指南](DEVELOPMENT.md)
- [部署指南](DEPLOYMENT.md)