#!/bin/bash

# Dir-Monitor-Go 服务管理脚本
# 提供统一的服务管理接口

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 默认配置
SERVICE_NAME="dir-monitor-go"
PROJECT_DIR="/opt/dir-monitor-go"

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

# 函数：检查root权限
check_root() {
    if [[ "$EUID" -ne 0 ]]; then
        print_error "请使用root权限运行此脚本 (sudo)"
        exit 1
    fi
}

# 函数：检查服务是否存在
check_service_exists() {
    print_step "检查服务是否存在..."
    
    if ! systemctl list-unit-files | grep -q "^${SERVICE_NAME}.service"; then
        print_error "服务不存在: ${SERVICE_NAME}"
        print_message "请先运行部署脚本: ./deploy/deploy-service.sh"
        exit 1
    fi
    
    print_message "服务检查通过"
}

# 函数：检查服务是否运行
is_service_running() {
    systemctl is-active --quiet "${SERVICE_NAME}"
    return $?
}

# 函数：检查服务是否启用
is_service_enabled() {
    systemctl is-enabled --quiet "${SERVICE_NAME}" 2>/dev/null
    return $?
}

# 函数：启动服务
start_service() {
    print_step "启动服务..."
    
    if is_service_running; then
        print_warning "服务已在运行中"
        return 0
    fi
    
    systemctl start "${SERVICE_NAME}"
    
    # 等待服务启动
    sleep 2
    
    if is_service_running; then
        print_message "服务启动成功"
        
        # 显示服务状态
        if is_service_enabled; then
            print_message "服务已设置为开机自启"
        else
            print_warning "服务未设置为开机自启"
        fi
    else
        print_error "服务启动失败"
        systemctl status "${SERVICE_NAME}" --no-pager -l
        exit 1
    fi
}

# 函数：停止服务
stop_service() {
    print_step "停止服务..."
    
    if ! is_service_running; then
        print_warning "服务未运行"
        return 0
    fi
    
    systemctl stop "${SERVICE_NAME}"
    
    if ! is_service_running; then
        print_message "服务停止成功"
    else
        print_error "服务停止失败"
        systemctl status "${SERVICE_NAME}" --no-pager -l
        exit 1
    fi
}

# 函数：重启服务
restart_service() {
    print_step "重启服务..."
    
    # 检查服务是否运行
    local was_running=false
    if is_service_running; then
        was_running=true
    fi
    
    systemctl restart "${SERVICE_NAME}"
    
    # 等待服务重启
    sleep 2
    
    if is_service_running; then
        print_message "服务重启成功"
        
        # 显示服务状态
        if is_service_enabled; then
            print_message "服务已设置为开机自启"
        else
            print_warning "服务未设置为开机自启"
        fi
    else
        print_error "服务重启失败"
        systemctl status "${SERVICE_NAME}" --no-pager -l
        exit 1
    fi
}

# 函数：查看服务状态
show_status() {
    print_step "查看服务状态..."
    
    echo ""
    systemctl status "${SERVICE_NAME}" --no-pager -l
    echo ""
    
    # 显示服务是否启用
    if is_service_enabled; then
        print_message "服务已设置为开机自启"
    else
        print_warning "服务未设置为开机自启"
    fi
}

# 函数：查看服务日志
show_logs() {
    print_step "查看服务日志..."
    
    echo ""
    echo -e "${CYAN}最近100行服务日志：${NC}"
    journalctl -u "${SERVICE_NAME}" -n 100 --no-pager
    echo ""
    
    # 如果应用日志文件存在，也显示其内容
    if [ -f "${LOG_DIR}/app.log" ]; then
        echo -e "${CYAN}应用日志文件：${NC}"
        echo "  ${LOG_DIR}/app.log"
        echo "  查看最新日志: tail -f ${LOG_DIR}/app.log"
        echo ""
    fi
}

# 函数：实时查看日志
follow_logs() {
    print_step "实时查看服务日志 (按 Ctrl+C 退出)..."
    
    journalctl -u "${SERVICE_NAME}" -f
}

# 函数：启用开机自启
enable_service() {
    print_step "启用开机自启..."
    
    systemctl enable "${SERVICE_NAME}"
    
    if systemctl is-enabled --quiet "${SERVICE_NAME}"; then
        print_message "开机自启已启用"
    else
        print_error "开机自启启用失败"
        exit 1
    fi
}

# 函数：禁用开机自启
disable_service() {
    print_step "禁用开机自启..."
    
    systemctl disable "${SERVICE_NAME}"
    
    if ! systemctl is-enabled --quiet "${SERVICE_NAME}" 2>/dev/null; then
        print_message "开机自启已禁用"
    else
        print_error "开机自启禁用失败"
        exit 1
    fi
}

# 函数：显示服务信息
show_info() {
    print_step "服务信息..."
    
    echo ""
    echo -e "${CYAN}Dir-Monitor-Go 服务信息${NC}"
    echo "=================================="
    echo ""
    echo -e "${YELLOW}服务状态：${NC}"
    if is_service_running; then
        echo -e "  运行状态: ${GREEN}运行中${NC}"
    else
        echo -e "  运行状态: ${RED}已停止${NC}"
    fi
    
    if is_service_enabled; then
        echo -e "  开机自启: ${GREEN}已启用${NC}"
    else
        echo -e "  开机自启: ${RED}已禁用${NC}"
    fi
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
    
    echo -e "${YELLOW}文件位置：${NC}"
    echo "  配置文件:     $PROJECT_DIR/configs/config.json"
    echo "  服务文件:     /etc/systemd/system/${SERVICE_NAME}.service"
    echo "  项目目录:     $PROJECT_DIR"
    echo "  日志目录:     $LOG_DIR"
    echo ""
    
    echo -e "${YELLOW}快速操作：${NC}"
    echo "  使用本脚本:   $0 [start|stop|restart|status|logs|follow|enable|disable|info]"
    echo ""
}

# 函数：显示帮助信息
show_help() {
    echo "用法: $0 [命令]"
    echo ""
    echo "Dir-Monitor-Go 服务管理脚本"
    echo ""
    echo "命令:"
    echo "  start          启动服务"
    echo "  stop           停止服务"
    echo "  restart        重启服务"
    echo "  status         查看服务状态"
    echo "  logs           查看最近日志"
    echo "  follow         实时查看日志"
    echo "  enable         启用开机自启"
    echo "  disable        禁用开机自启"
    echo "  info           显示服务信息"
    echo "  help           显示此帮助信息"
    echo ""
    echo "示例："
    echo "  sudo $0 start              # 启动服务"
    echo "  sudo $0 status             # 查看状态"
    echo "  sudo $0 logs               # 查看日志"
    echo "  sudo $0 follow             # 实时查看日志"
    echo "  sudo $0 restart            # 重启服务"
    echo ""
}

# 主函数
main() {
    # 检查是否有参数
    if [ $# -eq 0 ]; then
        show_help
        exit 0
    fi
    
    # 解析命令
    case "$1" in
        start)
            check_root
            check_service_exists
            start_service
            ;;
        stop)
            check_root
            check_service_exists
            stop_service
            ;;
        restart)
            check_root
            check_service_exists
            restart_service
            ;;
        status)
            check_root
            check_service_exists
            show_status
            ;;
        logs)
            check_root
            check_service_exists
            show_logs
            ;;
        follow)
            check_root
            check_service_exists
            follow_logs
            ;;
        enable)
            check_root
            check_service_exists
            enable_service
            ;;
        disable)
            check_root
            check_service_exists
            disable_service
            ;;
        info)
            check_root
            check_service_exists
            show_info
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"