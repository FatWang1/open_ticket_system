#!/bin/bash

# Docker构建脚本
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认值
IMAGE_NAME="open_ticket_system"
TAG="latest"
PLATFORM="linux/amd64"
PUSH=false

# 显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo "选项:"
    echo "  -n, --name IMAGE_NAME    镜像名称 (默认: open_ticket_system)"
    echo "  -t, --tag TAG            标签 (默认: latest)"
    echo "  -p, --platform PLATFORM  目标平台 (默认: linux/amd64)"
    echo "  --push                   构建完成后推送到镜像仓库"
    echo "  -h, --help               显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 -n myapp -t v1.0.0"
    echo "  $0 --platform linux/arm64 --push"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            IMAGE_NAME="$2"
            shift 2
            ;;
        -t|--tag)
            TAG="$2"
            shift 2
            ;;
        -p|--platform)
            PLATFORM="$2"
            shift 2
            ;;
        --push)
            PUSH=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}错误: 未知参数 $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 显示构建信息
echo -e "${GREEN}开始构建Docker镜像...${NC}"
echo "镜像名称: $IMAGE_NAME"
echo "标签: $TAG"
echo "目标平台: $PLATFORM"
echo ""

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}错误: Docker未运行或无法访问${NC}"
    exit 1
fi

# 构建镜像
echo -e "${YELLOW}构建镜像: $IMAGE_NAME:$TAG${NC}"
docker build \
    --platform $PLATFORM \
    --tag $IMAGE_NAME:$TAG \
    --file .deploy/k8s/docker/Dockerfile \
    .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}镜像构建成功: $IMAGE_NAME:$TAG${NC}"
else
    echo -e "${RED}镜像构建失败${NC}"
    exit 1
fi

# 显示镜像信息
echo ""
echo -e "${YELLOW}镜像信息:${NC}"
docker images $IMAGE_NAME:$TAG

# 如果指定了推送，则推送到镜像仓库
if [ "$PUSH" = true ]; then
    echo ""
    echo -e "${YELLOW}推送镜像到仓库...${NC}"
    docker push $IMAGE_NAME:$TAG
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}镜像推送成功${NC}"
    else
        echo -e "${RED}镜像推送失败${NC}"
        exit 1
    fi
fi

echo ""
echo -e "${GREEN}构建完成!${NC}"
