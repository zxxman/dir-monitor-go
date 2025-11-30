# Dir-Monitor-Go 部署指南

> **版本**: v3.2.1  
> **最后更新**: 2025年10月16日

## 📋 目录

1. [部署概述](#-部署概述)
2. [系统要求](#-系统要求)
3. [二进制部署](#-二进制部署)
4. [Docker部署](#-docker部署)
5. [Kubernetes部署](#-kubernetes部署)
6. [系统服务部署](#-系统服务部署)
7. [云平台部署](#-云平台部署)
8. [高可用部署](#-高可用部署)
9. [监控与日志](#-监控与日志)
10. [故障排除](#-故障排除)

---

## 🚀 部署概述

Dir-Monitor-Go 是一个轻量级的文件系统监控工具，可以以多种方式部署：

- **二进制部署**: 直接下载预编译的二进制文件
- **Docker部署**: 使用Docker容器运行
- **Kubernetes部署**: 在Kubernetes集群中部署
- **系统服务部署**: 作为系统服务运行
- **云平台部署**: 在各种云平台上部署

选择合适的部署方式取决于您的具体需求和环境限制。

---

## 💻 系统要求

### 最低要求
- **操作系统**: Linux, macOS, Windows
- **内存**: 64MB
- **磁盘**: 50MB
- **网络**: 可选（用于远程日志和API访问）

### 推荐配置
- **操作系统**: Linux (Ubuntu 20.04+, CentOS 8+, RHEL 8+)
- **内存**: 256MB
- **磁盘**: 500MB
- **CPU**: 2核心
- **网络**: 100Mbps（用于大量文件传输）

### 依赖要求
- **文件系统**: 支持inotify（Linux）或FSEvents（macOS）
- **权限**: 对监控目录的读取权限和执行命令的权限
- **Shell**: Bash或兼容的shell（用于命令执行）

---

## 📦 二进制部署

### 下载预编译二进制文件

1. **访问Releases页面**
   ```
   https://github.com/your-repo/dir-monitor-go/releases
   ```

2. **选择适合的版本**
   - Linux AMD64: `dir-monitor-go-linux-amd64-v3.2.1.tar.gz`
   - Linux ARM64: `dir-monitor-go-linux-arm64-v3.2.1.tar.gz`
   - macOS AMD64: `dir-monitor-go-darwin-amd64-v3.2.1.tar.gz`
   - Windows AMD64: `dir-monitor-go-windows-amd64-v3.2.1.zip`

3. **下载并解压**
   ```bash
   # 下载
   wget https://github.com/your-repo/dir-monitor-go/releases/download/v3.2.1/dir-monitor-go-linux-amd64-v3.2.1.tar.gz
   
   # 解压
   tar -xzf dir-monitor-go-linux-amd64-v3.2.1.tar.gz
   cd dir-monitor-go-linux-amd64-v3.2.1
   ```

### 安装配置

1. **复制二进制文件**
   ```bash
   sudo cp dir-monitor-go /usr/local/bin/
   sudo chmod +x /usr/local/bin/dir-monitor-go
   ```

2. **创建配置目录**
   ```bash
   sudo mkdir -p /etc/dir-monitor-go
   sudo mkdir -p /var/log/dir-monitor-go
   ```

3. **复制配置文件**
   ```bash
   sudo cp config.json.example /etc/dir-monitor-go/config.json
   ```

4. **编辑配置文件**
   ```bash
   sudo nano /etc/dir-monitor-go/config.json
   ```

### 运行

1. **直接运行**
   ```bash
   dir-monitor-go -config /etc/dir-monitor-go/config.json
   ```

2. **后台运行**
   ```bash
   nohup dir-monitor-go -config /etc/dir-monitor-go/config.json > /var/log/dir-monitor-go/output.log 2>&1 &
   ```

---

## 🐳 Docker部署

### 使用官方镜像

1. **拉取镜像**
   ```bash
   docker pull dirmonitor/go:v3.2.1
   ```

2. **运行容器**
   ```bash
   docker run -d \
     --name dir-monitor-go \
     -v /path/to/your/config.json:/app/config.json \
     -v /path/to/monitor:/data \
     -v /var/log/dir-monitor-go:/app/logs \
     dirmonitor/go:v3.2.1
   ```

### 构建自定义镜像

1. **创建Dockerfile**
   ```dockerfile
   FROM alpine:latest
   
   RUN apk --no-cache add ca-certificates
   WORKDIR /app
   
   COPY dir-monitor-go .
   COPY config.json .
   
   RUN mkdir -p logs
   VOLUME ["/app/logs", "/data"]
   
   EXPOSE 8080
   
   CMD ["./dir-monitor-go", "-config", "config.json"]
   ```

2. **构建镜像**
   ```bash
   docker build -t my-dir-monitor-go:v3.2.1 .
   ```

3. **运行自定义镜像**
   ```bash
   docker run -d \
     --name my-dir-monitor-go \
     -v /path/to/monitor:/data \
     -v /var/log/dir-monitor-go:/app/logs \
     my-dir-monitor-go:v3.2.1
   ```

### Docker Compose部署

1. **创建docker-compose.yml**
   ```yaml
   version: '3.8'
   
   services:
     dir-monitor-go:
       image: dirmonitor/go:v3.2.1
       container_name: dir-monitor-go
       restart: unless-stopped
       volumes:
         - ./config.json:/app/config.json:ro
         - /path/to/monitor:/data:ro
         - ./logs:/app/logs
       ports:
         - "8080:8080"
       environment:
         - TZ=Asia/Shanghai
   ```

2. **启动服务**
   ```bash
   docker-compose up -d
   ```

3. **查看日志**
   ```bash
   docker-compose logs -f dir-monitor-go
   ```

---

## ☸️ Kubernetes部署

### 创建部署配置

1. **创建Namespace**
   ```yaml
   # namespace.yaml
   apiVersion: v1
   kind: Namespace
   metadata:
     name: dir-monitor-go
   ```

2. **创建ConfigMap**
   ```yaml
   # configmap.yaml
   apiVersion: v1
   kind: ConfigMap
   metadata:
     name: dir-monitor-go-config
     namespace: dir-monitor-go
   data:
     config.json: |
       {
         "version": "3.2.1",
         "monitors": [
           {
             "name": "file-monitor",
             "path": "/data",
             "command": "echo 'File changed: {FILE_PATH}'",
             "patterns": ["*.txt", "*.log"],
             "recursive": true
           }
         ],
         "settings": {
           "log_level": "info",
           "max_concurrent": 5
         }
       }
   ```

3. **创建Deployment**
   ```yaml
   # deployment.yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: dir-monitor-go
     namespace: dir-monitor-go
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
           image: dirmonitor/go:v3.2.1
           imagePullPolicy: IfNotPresent
           ports:
           - containerPort: 8080
           volumeMounts:
           - name: config
             mountPath: /app/config.json
             subPath: config.json
           - name: data
             mountPath: /data
           - name: logs
             mountPath: /app/logs
           resources:
             requests:
               memory: "64Mi"
               cpu: "50m"
             limits:
               memory: "256Mi"
               cpu: "200m"
         volumes:
         - name: config
           configMap:
             name: dir-monitor-go-config
         - name: data
           hostPath:
             path: /path/to/monitor
             type: Directory
         - name: logs
           emptyDir: {}
   ```

4. **创建Service**
   ```yaml
   # service.yaml
   apiVersion: v1
   kind: Service
   metadata:
     name: dir-monitor-go-service
     namespace: dir-monitor-go
   spec:
     selector:
       app: dir-monitor-go
     ports:
     - protocol: TCP
       port: 8080
       targetPort: 8080
     type: ClusterIP
   ```

### 部署应用

1. **应用配置**
   ```bash
   kubectl apply -f namespace.yaml
   kubectl apply -f configmap.yaml
   kubectl apply -f deployment.yaml
   kubectl apply -f service.yaml
   ```

2. **检查部署状态**
   ```bash
   kubectl get pods -n dir-monitor-go
   kubectl logs -f deployment/dir-monitor-go -n dir-monitor-go
   ```

3. **端口转发（可选）**
   ```bash
   kubectl port-forward service/dir-monitor-go-service 8080:8080 -n dir-monitor-go
   ```

---

## 🔧 系统服务部署

### Systemd服务（Linux）

1. **创建服务文件**
   ```bash
   sudo nano /etc/systemd/system/dir-monitor-go.service
   ```

2. **添加服务配置**
   ```ini
   [Unit]
   Description=Dir Monitor Go Service
   After=network.target
   
   [Service]
   Type=simple
   User=root
   Group=root
   WorkingDirectory=/opt/dir-monitor-go
   ExecStart=/usr/local/bin/dir-monitor-go -config /etc/dir-monitor-go/config.json
   Restart=always
   RestartSec=5
   StandardOutput=journal
   StandardError=journal
   
   [Install]
   WantedBy=multi-user.target
   ```

3. **启用并启动服务**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable dir-monitor-go
   sudo systemctl start dir-monitor-go
   ```

4. **检查服务状态**
   ```bash
   sudo systemctl status dir-monitor-go
   sudo journalctl -u dir-monitor-go -f
   ```

### Windows服务

1. **使用NSSM安装服务**
   ```cmd
   nssm install DirMonitorGo "C:\path\to\dir-monitor-go.exe"
   nssm set DirMonitorGo Arguments "-config C:\path\to\config.json"
   nssm set DirMonitorGo DisplayName "Dir Monitor Go Service"
   nssm set DirMonitorGo Description "File system monitoring service"
   nssm start DirMonitorGo
   ```

2. **使用sc命令安装服务**
   ```cmd
   sc create DirMonitorGo binPath= "C:\path\to\dir-monitor-go.exe -config C:\path\to\config.json"
   sc start DirMonitorGo
   ```

---

## ☁️ 云平台部署

### AWS部署

1. **使用EC2实例**
   - 创建EC2实例（Ubuntu 20.04 LTS）
   - 配置安全组（开放8080端口）
   - 使用用户数据脚本自动安装

2. **用户数据脚本示例**
   ```bash
   #!/bin/bash
   apt-get update
   apt-get install -y wget
   
   # 下载并安装dir-monitor-go
   wget https://github.com/your-repo/dir-monitor-go/releases/download/v3.2.1/dir-monitor-go-linux-amd64-v3.2.1.tar.gz
   tar -xzf dir-monitor-go-linux-amd64-v3.2.1.tar.gz
   sudo cp dir-monitor-go-linux-amd64-v3.2.1/dir-monitor-go /usr/local/bin/
   sudo chmod +x /usr/local/bin/dir-monitor-go
   
   # 创建systemd服务
   cat > /etc/systemd/system/dir-monitor-go.service << EOF
   [Unit]
   Description=Dir Monitor Go Service
   After=network.target
   
   [Service]
   Type=simple
   ExecStart=/usr/local/bin/dir-monitor-go -config /etc/dir-monitor-go/config.json
   Restart=always
   
   [Install]
   WantedBy=multi-user.target
   EOF
   
   systemctl daemon-reload
   systemctl enable dir-monitor-go
   systemctl start dir-monitor-go
   ```

3. **使用ECS部署**
   ```json
   {
     "family": "dir-monitor-go",
     "networkMode": "awsvpc",
     "requiresCompatibilities": ["FARGATE"],
     "cpu": "256",
     "memory": "512",
     "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
     "containerDefinitions": [
       {
         "name": "dir-monitor-go",
         "image": "dirmonitor/go:v3.2.1",
         "portMappings": [
           {
             "containerPort": 8080,
             "protocol": "tcp"
           }
         ],
         "logConfiguration": {
           "logDriver": "awslogs",
           "options": {
             "awslogs-group": "/ecs/dir-monitor-go",
             "awslogs-region": "us-west-2",
             "awslogs-stream-prefix": "ecs"
           }
         }
       }
     ]
   }
   ```

### Google Cloud Platform部署

1. **使用Compute Engine**
   ```bash
   # 创建实例
   gcloud compute instances create dir-monitor-go \
     --image-family=ubuntu-2004-lts \
     --image-project=ubuntu-os-cloud \
     --machine-type=e2-micro \
     --zone=us-central1-a \
     --tags=http-server
   
   # 创建防火墙规则
   gcloud compute firewall-rules create allow-http \
     --allow tcp:8080 \
     --source-ranges 0.0.0.0/0 \
     --target-tags http-server
   ```

2. **使用Cloud Run**
   ```bash
   # 构建并推送镜像
   gcloud builds submit --tag gcr.io/PROJECT-ID/dir-monitor-go
   
   # 部署到Cloud Run
   gcloud run deploy dir-monitor-go \
     --image gcr.io/PROJECT-ID/dir-monitor-go \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated
   ```

### Azure部署

1. **使用Virtual Machine**
   ```bash
   # 创建资源组
   az group create --name dir-monitor-go-rg --location eastus
   
   # 创建虚拟机
   az vm create \
     --resource-group dir-monitor-go-rg \
     --name dir-monitor-go-vm \
     --image UbuntuLTS \
     --admin-username azureuser \
     --generate-ssh-keys
   
   # 开放端口
   az vm open-port \
     --resource-group dir-monitor-go-rg \
     --name dir-monitor-go-vm \
     --port 8080
   ```

2. **使用Container Instances**
   ```bash
   # 创建容器实例
   az container create \
     --resource-group dir-monitor-go-rg \
     --name dir-monitor-go \
     --image dirmonitor/go:v3.2.1 \
     --ports 8080 \
     --dns-name-label dir-monitor-go-unique
   ```

---

## 🔄 高可用部署

### 主备模式部署

1. **主节点配置**
   ```json
   {
     "version": "3.2.1",
     "role": "primary",
     "monitors": [
       {
         "name": "file-monitor",
         "path": "/data",
         "command": "process-file.sh {FILE_PATH}",
         "patterns": ["*"],
         "recursive": true
       }
     ],
     "settings": {
       "log_level": "info",
       "max_concurrent": 10
     },
     "ha": {
       "enabled": true,
       "node_id": "node-1",
       "peer_nodes": ["node-2"],
       "heartbeat_interval": "5s",
       "failover_timeout": "15s"
     }
   }
   ```

2. **备节点配置**
   ```json
   {
     "version": "3.2.1",
     "role": "secondary",
     "monitors": [
       {
         "name": "file-monitor",
         "path": "/data",
         "command": "process-file.sh {FILE_PATH}",
         "patterns": ["*"],
         "recursive": true
       }
     ],
     "settings": {
       "log_level": "info",
       "max_concurrent": 10
     },
     "ha": {
       "enabled": true,
       "node_id": "node-2",
       "peer_nodes": ["node-1"],
       "heartbeat_interval": "5s",
       "failover_timeout": "15s"
     }
   }
   ```

### 负载均衡部署

1. **使用Nginx负载均衡**
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

2. **使用HAProxy负载均衡**
   ```
   global
       log stdout format raw local0
   
   defaults
       log global
       mode http
       timeout connect 5000ms
       timeout client 50000ms
       timeout server 50000ms
   
   frontend dir_monitor_frontend
       bind *:80
       default_backend dir_monitor_backend
   
   backend dir_monitor_backend
       balance roundrobin
       server node1 10.0.1.10:8080 check
       server node2 10.0.1.11:8080 check
       server node3 10.0.1.12:8080 check
   ```

---

## 📊 监控与日志

### 日志配置

1. **配置日志级别**
   ```json
   {
     "settings": {
       "log_level": "info",
       "log_file": "/var/log/dir-monitor-go/app.log",
       "log_max_size": 100,
       "log_max_backups": 5,
       "log_max_age": 30
     }
   }
   ```

2. **结构化日志**
   ```json
   {
     "settings": {
       "log_format": "json",
       "log_fields": ["timestamp", "level", "message", "monitor", "file_path"]
     }
   }
   ```

### 监控指标

1. **内置指标**
   - 文件事件计数
   - 命令执行次数
   - 执行成功/失败次数
   - 平均执行时间
   - 当前并发执行数

2. **Prometheus集成**
   ```yaml
   # prometheus.yml
   global:
     scrape_interval: 15s
   
   scrape_configs:
     - job_name: 'dir-monitor-go'
       static_configs:
         - targets: ['localhost:8080']
   ```

3. **Grafana仪表板**
   - 事件率图表
   - 执行时间分布
   - 错误率统计
   - 资源使用情况

### 健康检查

1. **HTTP健康检查**
   ```bash
   curl http://localhost:8080/health
   ```

2. **响应示例**
   ```json
   {
     "status": "healthy",
     "timestamp": "2025-10-16T10:30:00Z",
     "uptime": "2h45m30s",
     "version": "3.2.1",
     "monitors": {
       "active": 5,
       "total": 5
     },
     "stats": {
       "events_processed": 1250,
       "commands_executed": 1200,
       "errors": 5
     }
   }
   ```

---

## 🔧 故障排除

### 常见问题

1. **权限问题**
   ```
   错误: permission denied
   解决: 确保对监控目录和执行命令有足够权限
   ```

2. **文件监控失败**
   ```
   错误: too many open files
   解决: 增加系统文件描述符限制
   ```

3. **命令执行超时**
   ```
   错误: command execution timeout
   解决: 增加超时设置或优化命令执行时间
   ```

4. **内存使用过高**
   ```
   错误: out of memory
   解决: 减少并发数或增加系统内存
   ```

### 调试技巧

1. **启用调试日志**
   ```json
   {
     "settings": {
       "log_level": "debug"
     }
   }
   ```

2. **使用strace跟踪系统调用**
   ```bash
   strace -p $(pidof dir-monitor-go)
   ```

3. **使用pprof分析性能**
   ```bash
   # CPU分析
   go tool pprof http://localhost:8080/debug/pprof/profile
   
   # 内存分析
   go tool pprof http://localhost:8080/debug/pprof/heap
   ```

### 性能优化

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
       "max_concurrent": 5
     }
   }
   ```

---

## 📚 更多资源

- [用户使用指南](USER_GUIDE.md)
- [配置参考](CONFIG.md)
- [API文档](API.md)
- [开发指南](DEVELOPMENT.md)