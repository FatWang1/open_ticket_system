# 部署脚本说明

本目录包含了用于构建和部署 Open Ticket System 的各种脚本。

## 脚本列表

### build.sh - Docker镜像构建脚本

用于构建项目的Docker镜像。

#### 使用方法

```bash
# 基本构建
./build.sh

# 指定镜像名称和标签
./build.sh -n myapp -t v1.0.0

# 指定目标平台
./build.sh -p linux/arm64

# 构建并推送到镜像仓库
./build.sh --push

# 显示帮助信息
./build.sh -h
```

#### 参数说明

- `-n, --name IMAGE_NAME`: 镜像名称 (默认: open_ticket_system)
- `-t, --tag TAG`: 标签 (默认: latest)
- `-p, --platform PLATFORM`: 目标平台 (默认: linux/amd64)
- `--push`: 构建完成后推送到镜像仓库
- `-h, --help`: 显示帮助信息

#### 示例

```bash
# 构建开发版本
./build.sh -n open-ticket-system -t dev

# 构建ARM64版本
./build.sh -p linux/arm64 -t v1.0.0

# 构建并推送生产版本
./build.sh -n open-ticket-system -t v1.0.0 --push
```

## 目录结构

```
.deploy/
├── scripts/
│   ├── build.sh          # Docker构建脚本
│   └── README.md         # 本说明文档
├── k8s/
│   ├── docker/
│   │   ├── Dockerfile    # Docker镜像定义
│   │   └── .dockerignore # Docker构建忽略文件
│   └── ...               # Kubernetes部署文件
└── ...
```

## 注意事项

1. 确保Docker已安装并正在运行
2. 构建前确保项目代码已编译通过
3. 如需推送到镜像仓库，请先配置Docker登录
4. 构建脚本会自动使用项目根目录下的Dockerfile

## 故障排除

### 常见问题

1. **Docker未运行**
   - 错误: "Docker未运行或无法访问"
   - 解决: 启动Docker服务

2. **构建失败**
   - 检查项目代码是否编译通过
   - 确认Dockerfile路径正确
   - 查看Docker构建日志

3. **推送失败**
   - 确认已登录到目标镜像仓库
   - 检查网络连接
   - 验证镜像标签格式

### 获取帮助

```bash
# 查看脚本帮助
./build.sh -h

# 查看Docker构建日志
docker build --progress=plain -f .deploy/k8s/docker/Dockerfile .
```
