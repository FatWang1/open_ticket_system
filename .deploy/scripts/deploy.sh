#!/bin/bash

# 部署脚本 - Open Ticket System
# 使用方法: ./deploy.sh [environment] [action]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认值
ENVIRONMENT=${1:-"dev"}
ACTION=${2:-"deploy"}
NAMESPACE="default"
APP_NAME="open-ticket-system"

# 打印带颜色的消息
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

# 检查依赖
check_dependencies() {
    print_step "检查部署依赖..."
    
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl 未安装或不在PATH中"
        exit 1
    fi
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装或不在PATH中"
        exit 1
    fi
    
    print_message "依赖检查通过"
}

# 检查Kubernetes连接
check_k8s_connection() {
    print_step "检查Kubernetes连接..."
    
    if kubectl cluster-info &> /dev/null; then
        print_message "Kubernetes连接正常"
    else
        print_error "无法连接到Kubernetes集群"
        exit 1
    fi
}

# 创建命名空间
create_namespace() {
    if [ "$NAMESPACE" != "default" ]; then
        print_step "创建命名空间: $NAMESPACE"
        
        if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
            kubectl create namespace "$NAMESPACE"
            print_message "命名空间创建成功"
        else
            print_message "命名空间已存在"
        fi
    fi
}

# 构建镜像
build_image() {
    print_step "构建Docker镜像..."
    
    if [ -f ".deploy/scripts/build.sh" ]; then
        chmod +x .deploy/scripts/build.sh
        .deploy/scripts/build.sh "$ENVIRONMENT"
    else
        print_error "构建脚本不存在: .deploy/scripts/build.sh"
        exit 1
    fi
}

# 加载镜像到kind集群（如果是本地开发）
load_image_to_kind() {
    if command -v kind &> /dev/null && kind get clusters | grep -q "kind"; then
        print_step "加载镜像到kind集群..."
        kind load docker-image open-ticket-system:latest
        print_message "镜像加载完成"
    fi
}

# 部署应用
deploy_application() {
    print_step "部署应用到Kubernetes..."
    
    # 检查部署配置是否存在
    if [ -f ".deploy/k8s/generated/deployment.yaml" ]; then
        DEPLOYMENT_FILE=".deploy/k8s/generated/deployment.yaml"
    elif [ -f ".deploy/k8s/deployment.yaml" ]; then
        DEPLOYMENT_FILE=".deploy/k8s/deployment.yaml"
    else
        print_error "部署配置文件不存在"
        exit 1
    fi
    
    # 应用配置
    if kubectl apply -f "$DEPLOYMENT_FILE" -n "$NAMESPACE"; then
        print_message "部署配置应用成功"
    else
        print_error "部署配置应用失败"
        exit 1
    fi
    
    # 等待部署完成
    print_step "等待部署完成..."
    kubectl rollout status deployment/"$APP_NAME" -n "$NAMESPACE" --timeout=300s
    
    print_message "应用部署完成"
}

# 升级应用
upgrade_application() {
    print_step "升级应用..."
    
    # 重新构建镜像
    build_image
    
    # 加载镜像到kind（如果是本地开发）
    load_image_to_kind
    
    # 重启部署
    kubectl rollout restart deployment/"$APP_NAME" -n "$NAMESPACE"
    
    # 等待升级完成
    print_step "等待升级完成..."
    kubectl rollout status deployment/"$APP_NAME" -n "$NAMESPACE" --timeout=300s
    
    print_message "应用升级完成"
}

# 回滚应用
rollback_application() {
    print_step "回滚应用..."
    
    # 查看部署历史
    print_message "部署历史:"
    kubectl rollout history deployment/"$APP_NAME" -n "$NAMESPACE"
    
    # 执行回滚
    kubectl rollout undo deployment/"$APP_NAME" -n "$NAMESPACE"
    
    # 等待回滚完成
    print_step "等待回滚完成..."
    kubectl rollout status deployment/"$APP_NAME" -n "$NAMESPACE" --timeout=300s
    
    print_message "应用回滚完成"
}

# 删除应用
delete_application() {
    print_step "删除应用..."
    
    if [ -f ".deploy/k8s/generated/deployment.yaml" ]; then
        DEPLOYMENT_FILE=".deploy/k8s/generated/deployment.yaml"
    elif [ -f ".deploy/k8s/deployment.yaml" ]; then
        DEPLOYMENT_FILE=".deploy/k8s/deployment.yaml"
    else
        print_error "部署配置文件不存在"
        exit 1
    fi
    
    if kubectl delete -f "$DEPLOYMENT_FILE" -n "$NAMESPACE"; then
        print_message "应用删除成功"
    else
        print_error "应用删除失败"
        exit 1
    fi
}

# 查看应用状态
show_status() {
    print_step "应用状态:"
    
    echo
    echo "Pod状态:"
    kubectl get pods -l app="$APP_NAME" -n "$NAMESPACE"
    
    echo
    echo "服务状态:"
    kubectl get svc -l app="$APP_NAME" -n "$NAMESPACE"
    
    echo
    echo "部署状态:"
    kubectl get deployment "$APP_NAME" -n "$NAMESPACE"
    
    echo
    echo "最近事件:"
    kubectl get events -n "$NAMESPACE" --sort-by='.lastTimestamp' | tail -10
}

# 查看日志
show_logs() {
    print_step "查看应用日志..."
    
    POD_NAME=$(kubectl get pods -l app="$APP_NAME" -n "$NAMESPACE" -o jsonpath='{.items[0].metadata.name}')
    
    if [ -n "$POD_NAME" ]; then
        print_message "Pod: $POD_NAME"
        kubectl logs -f "$POD_NAME" -n "$NAMESPACE"
    else
        print_error "未找到运行中的Pod"
    fi
}

# 端口转发
port_forward() {
    print_step "启动端口转发..."
    
    POD_NAME=$(kubectl get pods -l app="$APP_NAME" -n "$NAMESPACE" -o jsonpath='{.items[0].metadata.name}')
    
    if [ -n "$POD_NAME" ]; then
        print_message "转发本地8080端口到Pod的8080端口"
        print_message "按Ctrl+C停止端口转发"
        kubectl port-forward "$POD_NAME" 8080:8080 -n "$NAMESPACE"
    else
        print_error "未找到运行中的Pod"
    fi
}

# 显示帮助信息
show_help() {
    echo "使用方法: $0 [environment] [action]"
    echo
    echo "环境 (environment):"
    echo "  dev      - 开发环境 (默认)"
    echo "  staging  - 预发布环境"
    echo "  prod     - 生产环境"
    echo
    echo "操作 (action):"
    echo "  deploy   - 部署应用 (默认)"
    echo "  upgrade  - 升级应用"
    echo "  rollback - 回滚应用"
    echo "  delete   - 删除应用"
    echo "  status   - 查看状态"
    echo "  logs     - 查看日志"
    echo "  forward  - 端口转发"
    echo "  help     - 显示帮助"
    echo
    echo "示例:"
    echo "  $0 dev deploy    # 部署到开发环境"
    echo "  $0 prod upgrade  # 升级生产环境"
    echo "  $0 dev status    # 查看开发环境状态"
}

# 主函数
main() {
    case "$ACTION" in
        "deploy")
            print_message "开始部署 Open Ticket System..."
            print_message "环境: $ENVIRONMENT"
            print_message "命名空间: $NAMESPACE"
            echo
            
            check_dependencies
            check_k8s_connection
            create_namespace
            build_image
            load_image_to_kind
            deploy_application
            show_status
            ;;
        "upgrade")
            print_message "开始升级应用..."
            upgrade_application
            show_status
            ;;
        "rollback")
            print_message "开始回滚应用..."
            rollback_application
            show_status
            ;;
        "delete")
            print_message "开始删除应用..."
            delete_application
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs
            ;;
        "forward")
            port_forward
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "未知操作: $ACTION"
            show_help
            exit 1
            ;;
    esac
    
    print_message "操作完成！"
}

# 脚本入口
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
