#!/bin/bash

# Dir-Monitor-Go 服务部署脚本
# 专注于服务部署，不包含构建功能
# 适用于 Ubuntu/Debian 系统

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 默认配置
PROJECT_DIR="/opt/dir-monitor-go"
SERVICE_NAME="dir-monitor-go"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

# 从配置文件中获取日志路径
get_log_path_from_config() {
    local config_file="$PROJECT_DIR/configs/config.json"
    if [ -f "$config_file" ]; then
        # 提取配置文件中的日志路径
        local log_file=$(grep -o '"log_file"[[:space:]]*:[[:space:]]*"[^"]*"' "$config_file" | sed 's/.*"log_file"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/')
        if [ -n "$log_file" ]; then
            # 如果是相对路径，转换为绝对路径
            if [[ "$log_file" != /* ]]; then
                echo "$PROJECT_DIR/$log_file"
            else
                echo "$log_file"
            fi
        else
            echo "$PROJECT_DIR/logs/dir-monitor-go.log"
        fi
    else
        echo "$PROJECT_DIR/logs/dir-monitor-go.log"
    fi
}

# 获取日志目录路径
LOG_FILE=$(get_log_path_from_config)
LOG_DIR=$(dirname "$LOG_FILE")

# 函数：打印彩色消息
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

print_header() {
    echo -e "${CYAN}"
    echo "======================================"
    echo "   Dir-Monitor-Go 服务部署脚本"
    echo "======================================"
    echo -e "${NC}"
}

# 函数：检查root权限
check_root() {
    if [[ "$EUID" -ne 0 ]]; then
        print_error "请使用root权限运行此脚本 (sudo)"
        exit 1
    fi
}

# 函数：检查项目文件
check_project_files() {
    print_step "检查项目文件..."
    
    if [ ! -d "$PROJECT_DIR" ]; then
        print_error "项目目录不存在: $PROJECT_DIR"
        print_message "请确保已将项目文件复制到 $PROJECT_DIR 目录"
        exit 1
    fi
    
    if [ ! -f "$PROJECT_DIR/${SERVICE_NAME}" ]; then
        print_error "服务可执行文件不存在: $PROJECT_DIR/${SERVICE_NAME}"
        print_message "请确保已构建项目或复制预编译的二进制文件"
        exit 1
    fi
    
    if [ ! -f "$PROJECT_DIR/configs/config.json" ]; then
        print_error "配置文件不存在: $PROJECT_DIR/configs/config.json"
        print_message "请确保配置文件已存在"
        exit 1
    fi
    
    print_message "项目文件检查通过"
}

# 函数：配置内核参数
configure_kernel() {
    print_step "配置内核参数..."
    
    # 配置inotify参数
    tee /etc/sysctl.d/99-dir-monitor-go.conf > /dev/null <<EOF
# Dir-Monitor-Go 内核参数配置
fs.inotify.max_user_watches = 524288
fs.inotify.max_user_instances = 512
fs.inotify.max_queued_events = 16384
EOF
    
    # 立即应用参数
    sysctl -p /etc/sysctl.d/99-dir-monitor-go.conf > /dev/null 2>&1 || true
    
    print_message "内核参数配置完成"
}

# 函数：创建日志目录
create_log_directory() {
    print_step "创建日志目录..."
    
    mkdir -p "$LOG_DIR"
    chmod 755 "$LOG_DIR"
    
    print_message "日志目录创建完成: $LOG_DIR"
}

# 函数：设置文件权限
set_permissions() {
    print_step "设置文件权限..."
    
    # 设置可执行权限
    chmod +x "$PROJECT_DIR/${SERVICE_NAME}"
    
    # 设置配置文件权限
    chmod 644 "$PROJECT_DIR/configs/config.json"
    
    print_message "文件权限设置完成"
}

# 函数：创建systemd服务
create_systemd_service() {
    print_step "创建systemd服务..."
    
    # 创建服务文件
    tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=Directory Monitor Go Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$PROJECT_DIR
ExecStart=$PROJECT_DIR/${SERVICE_NAME} -config $PROJECT_DIR/configs/config.json 
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

# 安全设置已移除

# 资源限制
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
    
    # 重新加载systemd
    systemctl daemon-reload
    
    print_message "systemd服务创建完成"
}

# 函数：配置日志轮转
configure_logrotate() {
    print_step "配置日志轮转..."
    
    tee "/etc/logrotate.d/${SERVICE_NAME}" > /dev/null <<EOF
$LOG_DIR/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 root root
    postrotate
        systemctl reload ${SERVICE_NAME} > /dev/null 2>&1 || true
    endscript
}
EOF
    
    print_message "日志轮转配置完成"
}

# 函数：启动服务
start_service() {
    print_step "启动服务..."
    
    # 启用服务
    systemctl enable "${SERVICE_NAME}"
    
    # 启动服务
    systemctl start "${SERVICE_NAME}"
    
    # 检查服务状态
    if systemctl is-active --quiet "${SERVICE_NAME}"; then
        print_message "服务启动成功"
    else
        print_error "服务启动失败，请检查日志"
        systemctl status "${SERVICE_NAME}" --no-pager -l
        exit 1
    fi
}

# 函数：显示部署结果
show_deployment_result() {
    echo ""
    echo -e "${GREEN}✅ Dir-Monitor-Go 服务部署成功！${NC}"
    echo ""
    echo -e "${YELLOW}服务管理命令：${NC}"
    echo "  启动服务:     systemctl start ${SERVICE_NAME}"
    echo "  停止服务:     systemctl stop ${SERVICE_NAME}"
    echo "  重启服务:     systemctl restart ${SERVICE_NAME}"
    echo "  查看状态:     systemctl status ${SERVICE_NAME}"
    echo "  开机自启:     systemctl enable ${SERVICE_NAME}"
    echo ""
    echo -e "${YELLOW}日志查看命令：${NC}"
    echo "  服务日志:     journalctl -u ${SERVICE_NAME} -f"
    echo "  应用日志:     tail -f ${LOG_DIR}/app.log"
    echo ""
    echo -e "${YELLOW}配置文件位置：${NC}"
    echo "  $PROJECT_DIR/configs/config.json"
    echo ""
    echo -e "${YELLOW}⚠️  重要：${NC}"
    echo "   请务必在配置文件中设置正确的监控目录路径！"
    echo ""
}

# 函数：显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "Dir-Monitor-Go 服务部署脚本 - 专注于服务部署，不包含构建功能"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示此帮助信息"
    echo "  -c, --check    仅检查项目文件，不执行部署"
    echo ""
    echo "部署前提："
    echo "  1. 项目文件已存在于 $PROJECT_DIR"
    echo "  2. 服务可执行文件已构建: $PROJECT_DIR/${SERVICE_NAME}"
    echo "  3. 配置文件已存在: $PROJECT_DIR/configs/config.json"
    echo ""
    echo "示例："
    echo "  sudo $0                    # 执行完整部署"
    echo "  sudo $0 --check            # 仅检查项目文件"
    echo "  sudo $0 --help             # 显示帮助信息"
    echo ""
}

# 主函数
main() {
    print_header
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -c|--check)
                check_root
                check_project_files
                echo -e "${GREEN}✅ 项目文件检查通过${NC}"
                exit 0
                ;;
            *)
                print_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
        shift
    done
    
    # 执行完整部署流程
    check_root
    check_project_files
    configure_kernel
    create_log_directory
    set_permissions
    create_systemd_service
    configure_logrotate
    start_service
    show_deployment_result
}

# 执行主函数
main "$@"