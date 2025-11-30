# Dir-Monitor-Go

> **版本**: v3.2.1  
> **最后更新**: 2025年10月16日

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/your-repo/dir-monitor-go/actions)
[![Release](https://img.shields.io/badge/Release-v3.2.1-red.svg)](https://github.com/your-repo/dir-monitor-go/releases)

## 📖 简介

Dir-Monitor-Go 是一个高性能的文件系统监控工具，使用 Go 语言开发。它可以实时监控指定目录的变化，并在检测到文件事件时触发自定义命令。该工具具有以下特点：

- 🔍 **实时监控**: 基于操作系统原生事件通知机制，高效监控文件系统变化
- 🎯 **灵活配置**: 支持多种配置选项，满足不同场景需求
- ⚡ **高性能**: 优化的并发处理机制，支持大规模文件监控
- 🛡️ **稳定可靠**: 内置文件稳定性检测，避免处理未完全传输的文件
- 🔄 **高可用**: 支持主备模式和负载均衡，确保服务连续性
- 🐳 **容器化**: 提供Docker镜像，便于容器化部署
- ☸️ **云原生**: 支持Kubernetes部署，适应云原生环境

## 🚀 快速开始

### 安装

#### 二进制安装（推荐）

```bash
# 下载适合您系统的二进制文件
wget https://github.com/your-repo/dir-monitor-go/releases/download/v3.2.1/dir-monitor-go-linux-amd64-v3.2.1.tar.gz

# 解压
tar -xzf dir-monitor-go-linux-amd64-v3.2.1.tar.gz

# 复制到系统路径
sudo cp dir-monitor-go /usr/local/bin/
sudo chmod +x /usr/local/bin/dir-monitor-go
```

#### Docker安装

```bash
docker pull dirmonitor/go:v3.2.1
```

#### 源码编译

```bash
git clone https://github.com/your-repo/dir-monitor-go.git
cd dir-monitor-go
go build -o dir-monitor-go cmd/dir-monitor/main.go
```

### 基本使用

1. **创建配置文件**
   ```bash
   dir-monitor-go -example-config > config.json
   ```

2. **编辑配置文件**
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

3. **启动监控**
   ```bash
   dir-monitor-go -config config.json
   ```

### Docker运行

```bash
docker run -d \
  --name dir-monitor-go \
  -v /path/to/your/config.json:/app/config.json \
  -v /path/to/monitor:/data \
  dirmonitor/go:v3.2.1
```

## 📚 文档

完整文档请访问 [📋 文档中心](docs/)，包含以下内容：

- [📖 用户指南](docs/USER_GUIDE.md) - 安装、配置和使用指南
- [⚙️ 配置参考](docs/CONFIG.md) - 详细的配置选项说明
- [� 部署指南](docs/DEPLOYMENT.md) - 多种部署方式说明
- [🔌 API文档](docs/API.md) - REST API接口文档
- [📝 更新日志](docs/CHANGELOG.md) - 版本更新历史

## 🌟 主要特性

### 实时文件监控

- 基于操作系统原生事件通知机制（inotify/FSEvents）
- 支持递归目录监控
- 高效的事件过滤和处理

### 灵活的配置选项

- 支持多种文件模式匹配
- 可配置的事件类型过滤
- 灵活的命令执行环境设置

### 高性能并发处理

- 优化的并发控制机制
- 可配置的并发执行数量
- 智能的事件批处理

### 文件稳定性检测

- 智能检测文件传输完成状态
- 可配置的稳定性检查参数
- 避免处理未完全传输的文件

### 高可用性支持

- 主备模式部署
- 负载均衡支持
- 健康检查和故障转移

### 丰富的日志和监控

- 结构化日志输出
- 多级别日志控制
- Prometheus指标集成

## 🏗️ 架构设计

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   文件系统事件   │───▶│   事件处理器     │───▶│   命令执行器     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   过滤器        │    │   并发控制器     │
                       └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   稳定性检测    │    │   日志记录器     │
                       └─────────────────┘    └─────────────────┘
```

## 📊 性能指标

| 指标 | 值 |
|------|-----|
| 支持监控文件数 | 100万+ |
| 事件处理延迟 | <10ms |
| 内存占用 | <50MB |
| CPU占用 | <5% |

## 🌍 使用场景

### 文件处理自动化

- 监控上传目录，自动处理上传的文件
- 监控日志目录，自动分析日志文件
- 监控备份目录，自动验证备份完整性

### 数据同步

- 监控源目录，自动同步到目标目录
- 监控配置文件变更，自动重载服务配置
- 监控代码变更，自动触发构建流程

### 安全监控

- 监控敏感目录，记录文件访问
- 监控系统目录，检测异常文件变化
- 监控审计日志，自动触发安全响应

## 🤝 贡献

我们欢迎所有形式的贡献！请查看[贡献指南](docs/DEVELOPMENT.md#贡献指南)了解如何参与项目开发。

### 贡献方式

- 🐛 报告问题
- 💡 提出新功能建议
- 📝 改进文档
- 🔧 提交代码修复
- 🧪 添加测试用例

### 开发流程

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

感谢以下开源项目：

- [fsnotify](https://github.com/fsnotify/fsnotify) - 文件系统事件监控
- [gronx](https://github.com/adhocore/gronx) - Cron表达式解析
- [cobra](https://github.com/spf13/cobra) - 命令行界面

## 📞 联系我们

- 📧 邮箱: contact@dir-monitor.example.com
- 💬 讨论: [GitHub Discussions](https://github.com/your-repo/dir-monitor-go/discussions)
- 🐛 问题: [GitHub Issues](https://github.com/your-repo/dir-monitor-go/issues)
- 📖 文档: [在线文档](https://dir-monitor.example.com/docs)

---

⭐ 如果这个项目对您有帮助，请给我们一个星标！