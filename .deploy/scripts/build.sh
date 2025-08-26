#!/bin/bash

# 构建脚本 - Open Ticket System
# 使用方法: ./build.sh [version] [platform]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认值
VERSION=${1:-"latest"}
PLATFORM=${2:-"linux/amd64"}
IMAGE_NAME="open-ticket-system"
REGISTRY=""

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
    print_step "检查构建依赖..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装或不在PATH中"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        print_error "Go 未安装或不在PATH中"
        exit 1
    fi
    
    print_message "依赖检查通过"
}

# 清理旧的构建产物
cleanup() {
    print_step "清理旧的构建产物..."
    
    # 清理Docker镜像
    if docker images | grep -q "$IMAGE_NAME"; then
        print_message "删除旧的Docker镜像..."
        docker rmi "$IMAGE_NAME:$VERSION" 2>/dev/null || true
    fi
    
    # 清理本地构建产物
    if [ -d "bin" ]; then
        rm -rf bin/
    fi
    
    print_message "清理完成"
}

# 运行测试
run_tests() {
    print_step "运行单元测试..."
    
    if go test ./... -v; then
        print_message "测试通过"
    else
        print_error "测试失败"
        exit 1
    fi
}

# 构建Go二进制文件
build_binary() {
    print_step "构建Go二进制文件..."
    
    # 创建输出目录
    mkdir -p bin
    
    # 设置构建参数
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=amd64
    
    # 构建
    if go build -ldflags="-s -w" -o bin/open_ticket_system ./cmd/open_ticket_system; then
        print_message "二进制文件构建成功: bin/open_ticket_system"
    else
        print_error "二进制文件构建失败"
        exit 1
    fi
}

# 构建Docker镜像
build_docker_image() {
    print_step "构建Docker镜像..."
    
    # 构建镜像
    if docker build -t "$IMAGE_NAME:$VERSION" -f .deploy/k8s/docker/Dockerfile .; then
        print_message "Docker镜像构建成功: $IMAGE_NAME:$VERSION"
    else
        print_error "Docker镜像构建失败"
        exit 1
    fi
    
    # 显示镜像信息
    print_message "镜像信息:"
    docker images "$IMAGE_NAME:$VERSION"
}

# 推送镜像到仓库（可选）
push_image() {
    if [ -n "$REGISTRY" ]; then
        print_step "推送镜像到仓库..."
        
        # 标记镜像
        docker tag "$IMAGE_NAME:$VERSION" "$REGISTRY/$IMAGE_NAME:$VERSION"
        
        # 推送镜像
        if docker push "$REGISTRY/$IMAGE_NAME:$VERSION"; then
            print_message "镜像推送成功: $REGISTRY/$IMAGE_NAME:$VERSION"
        else
            print_error "镜像推送失败"
            exit 1
        fi
    else
        print_warning "未配置镜像仓库，跳过推送"
    fi
}

# 生成部署配置
generate_deployment_config() {
    print_step "生成部署配置..."
    
    # 创建部署配置目录
    mkdir -p .deploy/k8s/generated
    
    # 复制并替换版本号
    sed "s|image: open-ticket-system:latest|image: $IMAGE_NAME:$VERSION|g" \
        .deploy/k8s/deployment.yaml > .deploy/k8s/generated/deployment.yaml
    
    print_message "部署配置生成完成: .deploy/k8s/generated/deployment.yaml"
}

# 显示构建信息
show_build_info() {
    print_step "构建完成！"
    echo
    echo "构建信息:"
    echo "  镜像名称: $IMAGE_NAME"
    echo "  版本: $VERSION"
    echo "  平台: $PLATFORM"
    echo "  二进制文件: bin/open_ticket_system"
    echo "  部署配置: .deploy/k8s/generated/deployment.yaml"
    echo
    echo "下一步操作:"
    echo "  1. 部署到Kubernetes: kubectl apply -f .deploy/k8s/generated/"
    echo "  2. 查看部署状态: kubectl get pods -l app=open-ticket-system"
    echo "  3. 查看服务: kubectl get svc -l app=open-ticket-system"
    echo
}

# 主函数
main() {
    print_message "开始构建 Open Ticket System..."
    print_message "版本: $VERSION"
    print_message "平台: $PLATFORM"
    echo
    
    check_dependencies
    cleanup
    run_tests
    build_binary
    build_docker_image
    push_image
    generate_deployment_config
    show_build_info
    
    print_message "构建完成！"
}

# 脚本入口
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
