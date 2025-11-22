# Dir-Monitor-Go 部署指南

本目录包含用于部署和管理 Dir-Monitor-Go 服务的脚本，专注于服务部署和卸载，不包含构建功能。

## 目录结构

```
deploy/
├── deploy-service.sh      # 服务部署脚本
├── service-manager.sh     # 服务管理脚本
├── uninstall-service.sh   # 服务卸载脚本
└── README.md             # 本文件
```

## 脚本说明

### deploy-service.sh
服务部署脚本，用于将 Dir-Monitor-Go 部署为系统服务。

**功能：**
- 检查项目文件完整性
- 配置内核参数（inotify）
- 创建日志目录
- 创建systemd服务
- 配置日志轮转
- 启动服务

**使用方法：**
```bash
# 基本部署（需要root权限）
sudo bash deploy-service.sh

# 检查项目文件（不部署）
sudo bash deploy-service.sh --check

# 显示帮助信息
sudo bash deploy-service.sh --help
```

**部署前提：**
1. 项目文件已存在于 `/opt/dir-monitor-go`
2. 服务可执行文件已构建：`/opt/dir-monitor-go/dir-monitor-go`
3. 配置文件已存在：`/opt/dir-monitor-go/configs/config.json`

### service-manager.sh
服务管理脚本，用于管理已部署的服务。

**功能：**
- 启动/停止/重启服务
- 查看服务状态
- 查看服务日志
- 启用/禁用开机自启
- 显示服务信息

**使用方法：**
```bash
# 启动服务
sudo bash service-manager.sh start

# 停止服务
sudo bash service-manager.sh stop

# 重启服务
sudo bash service-manager.sh restart

# 查看服务状态
sudo bash service-manager.sh status

# 查看最近日志
sudo bash service-manager.sh logs

# 实时查看日志
sudo bash service-manager.sh follow

# 启用开机自启
sudo bash service-manager.sh enable

# 显示服务信息
sudo bash service-manager.sh info

# 显示帮助信息
sudo bash service-manager.sh help
```

### uninstall-service.sh
服务卸载脚本，用于完全卸载服务及相关配置。

**功能：**
- 停止并禁用服务
- 删除systemd服务文件
- 删除日志轮转配置
- 删除内核参数配置
- 重新加载systemd

**使用方法：**
```bash
# 交互式卸载（需要确认）
sudo bash uninstall-service.sh

# 自动确认卸载（不提示）
sudo bash uninstall-service.sh --yes

# 显示帮助信息
sudo bash uninstall-service.sh --help
```

**注意：**
- 卸载操作不会删除项目文件和日志文件
- 如需完全清理，请手动删除项目目录和日志目录

## 快速开始

### 1. 准备项目文件
确保项目文件已准备就绪：
```bash
# 创建项目目录
sudo mkdir -p /opt/dir-monitor-go

# 复制项目文件（包括已构建的二进制文件）
sudo cp -r /path/to/dir-monitor-go/* /opt/dir-monitor-go/

# 确保配置文件存在
sudo cp /opt/dir-monitor-go/config.json.example /opt/dir-monitor-go/configs/config.json

# 编辑配置文件
sudo nano /opt/dir-monitor-go/configs/config.json
```

### 2. 部署服务
```bash
# 进入部署目录
cd /opt/dir-monitor-go/deploy

# 执行部署
sudo bash deploy-service.sh
```

### 3. 管理服务
```bash
# 查看服务状态
sudo bash service-manager.sh status

# 查看日志
sudo bash service-manager.sh logs

# 重启服务（配置更改后）
sudo bash service-manager.sh restart
```

### 4. 卸载服务（如需要）
```bash
# 卸载服务
sudo bash uninstall-service.sh

# 手动清理项目文件（可选）
sudo rm -rf /opt/dir-monitor-go
sudo rm -rf /var/log/dir-monitor-go
```

## 配置文件

部署完成后，主要配置文件位置：
- **配置文件**：`/opt/dir-monitor-go/configs/config.json`
- **服务文件**：`/etc/systemd/system/dir-monitor-go.service`
- **日志轮转**：`/etc/logrotate.d/dir-monitor-go`
- **内核参数**：`/etc/sysctl.d/99-dir-monitor-go.conf`

### 日志路径配置

系统现在会自动从配置文件中读取日志路径设置：
- 如果配置文件中设置了 `log_file` 路径，系统会自动使用该路径
- 相对路径会自动转换为相对于项目目录的绝对路径
- 如果配置文件中未设置日志路径，默认使用 `/opt/dir-monitor-go/logs/dir-monitor-go.log`

示例配置：
```json
{
  "log_file": "logs/dir-monitor-go.log",
  "log_level": "info"
}
```

## 系统要求

- **操作系统**：Ubuntu 18.04+ 或 Debian 9+
- **权限要求**：所有脚本需要root权限运行
- **依赖要求**：systemd, logrotate

## 服务管理命令

除了使用脚本，也可以直接使用systemd命令：

```bash
# 启动服务
sudo systemctl start dir-monitor-go

# 停止服务
sudo systemctl stop dir-monitor-go

# 重启服务
sudo systemctl restart dir-monitor-go

# 查看服务状态
sudo systemctl status dir-monitor-go

# 设置开机自启
sudo systemctl enable dir-monitor-go

# 查看服务日志
sudo journalctl -u dir-monitor-go -f
```

## 故障排除

### 服务启动失败
1. 检查服务状态：`sudo bash service-manager.sh status`
2. 查看服务日志：`sudo bash service-manager.sh logs`
3. 检查配置文件语法：`cat /opt/dir-monitor-go/configs/config.json`

### 部署失败
1. 检查项目文件：`sudo bash deploy-service.sh --check`
2. 确保二进制文件存在且可执行
3. 确保配置文件存在且格式正确

### 权限问题
确保所有操作都使用root权限：
```bash
sudo bash deploy-service.sh
sudo bash service-manager.sh start
```

## 注意事项

1. **备份配置**：在修改配置文件前建议备份
2. **服务重启**：配置更改后需要重启服务生效
3. **日志管理**：日志会自动轮转，无需手动清理
4. **内核参数**：部署时会自动优化inotify参数
5. **安全设置**：服务以root用户运行，确保系统安全

## 版本历史

### v1.1.0 (当前版本)
- 日志路径现在从配置文件自动获取
- 部署脚本和服务管理脚本动态解析配置
- 支持相对路径和绝对路径的日志配置
- 保持向后兼容性

### v1.0.0
- 重构部署脚本，专注于服务部署和卸载
- 移除构建功能，假设二进制文件已准备就绪
- 简化部署流程，提高可靠性
- 增强错误处理和用户提示
- 提供独立的服务管理和卸载脚本