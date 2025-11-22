#!/bin/bash

# Dir-Monitor-Go 服务卸载脚本
# 完全卸载服务及相关配置

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
LOG_DIR="/var/log/dir-monitor-go"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

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
    echo "   Dir-Monitor-Go 服务卸载脚本"
    echo "======================================"
    echo -e "${NC}"
}

# 函数：检查服务是否存在
check_service_exists() {
    print_step "检查服务是否存在..."
    
    # 检查服务文件是否存在
    if [ ! -f "$SERVICE_FILE" ]; then
        print_warning "服务文件不存在: $SERVICE_FILE"
        return 1
    fi
    
    # 检查systemd是否已识别该服务
    if ! systemctl list-unit-files | grep -q "^${SERVICE_NAME}.service"; then
        print_warning "systemd未识别服务: ${SERVICE_NAME}"
        return 1
    fi
    
    print_message "服务存在: ${SERVICE_NAME}"
    return 0
}

# 函数：检查root权限
check_root() {
    if [[ "$EUID" -ne 0 ]]; then
        print_error "请使用root权限运行此脚本 (sudo)"
        exit 1
    fi
}

# 函数：显示卸载确认
show_uninstall_confirmation() {
    print_header
    
    echo -e "${YELLOW}⚠️  警告：此操作将完全卸载 Dir-Monitor-Go 服务${NC}"
    echo ""
    echo "将要执行的操作："
    echo "  - 停止并禁用服务"
    echo "  - 删除systemd服务文件"
    echo "  - 删除日志轮转配置"
    echo "  - 删除内核参数配置"
    echo "  - 重新加载systemd"
    echo ""
    echo -e "${YELLOW}注意：此操作不会删除项目文件和日志文件${NC}"
    echo "      如需删除项目文件，请手动删除: $PROJECT_DIR"
    echo "      如需删除日志文件，请手动删除: $LOG_DIR"
    echo ""
    
    read -p "确定要继续卸载吗? (yes/no): " -r confirmation
    
    if [[ "$confirmation" != "yes" ]]; then
        print_message "卸载操作已取消"
        exit 0
    fi
}

# 函数：停止服务
stop_service() {
    print_step "停止服务..."
    
    # 检查服务是否运行
    if systemctl is-active --quiet "${SERVICE_NAME}"; then
        systemctl stop "${SERVICE_NAME}"
        
        # 验证服务是否已停止
        if systemctl is-active --quiet "${SERVICE_NAME}"; then
            print_error "服务停止失败"
            systemctl status "${SERVICE_NAME}" --no-pager -l
            exit 1
        else
            print_message "服务已停止"
        fi
    else
        print_warning "服务未运行"
    fi
}

# 函数：禁用服务
disable_service() {
    print_step "禁用服务..."
    
    # 检查服务是否已启用
    if systemctl is-enabled --quiet "${SERVICE_NAME}" 2>/dev/null; then
        systemctl disable "${SERVICE_NAME}"
        
        # 验证服务是否已禁用
        if systemctl is-enabled --quiet "${SERVICE_NAME}" 2>/dev/null; then
            print_error "服务禁用失败"
            exit 1
        else
            print_message "服务已禁用"
        fi
    else
        print_warning "服务未启用"
    fi
}

# 函数：删除服务文件
remove_service_file() {
    print_step "删除服务文件..."
    
    if [ -f "$SERVICE_FILE" ]; then
        rm "$SERVICE_FILE"
        print_message "服务文件已删除: $SERVICE_FILE"
    else
        print_warning "服务文件不存在: $SERVICE_FILE"
    fi
}

# 函数：删除日志轮转配置
remove_logrotate_config() {
    print_step "删除日志轮转配置..."
    
    if [ -f "/etc/logrotate.d/${SERVICE_NAME}" ]; then
        rm "/etc/logrotate.d/${SERVICE_NAME}"
        print_message "日志轮转配置已删除"
    else
        print_warning "日志轮转配置不存在"
    fi
}

# 函数：删除内核参数配置
remove_kernel_config() {
    print_step "删除内核参数配置..."
    
    if [ -f "/etc/sysctl.d/99-dir-monitor-go.conf" ]; then
        rm "/etc/sysctl.d/99-dir-monitor-go.conf"
        print_message "内核参数配置已删除"
    else
        print_warning "内核参数配置不存在"
    fi
}

# 函数：重新加载systemd
reload_systemd() {
    print_step "重新加载systemd..."
    
    systemctl daemon-reload
    
    print_message "systemd已重新加载"
}

# 函数：显示卸载结果
show_uninstall_result() {
    echo ""
    echo -e "${GREEN}✅ Dir-Monitor-Go 服务卸载完成！${NC}"
    echo ""
    echo -e "${YELLOW}卸载总结：${NC}"
    echo "  ✅ 服务已停止"
    echo "  ✅ 服务已禁用"
    echo "  ✅ 服务文件已删除"
    echo "  ✅ 日志轮转配置已删除"
    echo "  ✅ 内核参数配置已删除"
    echo "  ✅ systemd已重新加载"
    echo ""
    echo -e "${YELLOW}剩余文件（如需手动清理）：${NC}"
    echo "  项目目录: $PROJECT_DIR"
    echo "  日志目录: $LOG_DIR"
    echo ""
    echo -e "${YELLOW}清理命令：${NC}"
    echo "  删除项目文件: sudo rm -rf $PROJECT_DIR"
    echo "  删除日志文件: sudo rm -rf $LOG_DIR"
    echo ""
}

# 函数：显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "Dir-Monitor-Go 服务卸载脚本"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示此帮助信息"
    echo "  -y, --yes      自动确认卸载（不提示确认）"
    echo ""
    echo "注意："
    echo "  - 此脚本需要root权限运行"
    echo "  - 卸载操作不会删除项目文件和日志文件"
    echo "  - 如需完全清理，请手动删除项目目录和日志目录"
    echo ""
    echo "示例："
    echo "  sudo $0                    # 交互式卸载（需要确认）"
    echo "  sudo $0 --yes              # 自动确认卸载"
    echo "  sudo $0 --help             # 显示帮助信息"
    echo ""
}

# 函数：执行完整卸载
perform_uninstall() {
    stop_service
    disable_service
    remove_service_file
    remove_logrotate_config
    remove_kernel_config
    reload_systemd
}

# 主函数
main() {
    # 检查root权限
    check_root
    
    # 解析命令行参数
    local auto_confirm=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -y|--yes)
                auto_confirm=true
                shift
                ;;
            *)
                print_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 检查服务是否存在
    if ! check_service_exists; then
        print_warning "服务不存在或已卸载: ${SERVICE_NAME}"
        print_message "如需重新安装，请使用部署脚本"
        exit 0
    fi
    
    # 显示确认信息（除非自动确认）
    if [ "$auto_confirm" = false ]; then
        show_uninstall_confirmation
    fi
    
    # 执行卸载
    perform_uninstall
    
    # 显示结果
    show_uninstall_result
}

# 执行主函数
main "$@"