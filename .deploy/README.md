# Open Ticket System 部署指南

## 概述

本文档描述了如何部署 Open Ticket System 到 Kubernetes 集群。

## 目录结构

```
.deploy/
├── k8s/                    # Kubernetes 配置文件
│   ├── docker/            # Docker 相关文件
│   │   ├── .dockerignore  # Docker 忽略文件
│   │   └── Dockerfile     # 多阶段构建 Dockerfile
│   ├── deployment.yaml    # 主应用部署配置
│   ├── config-prod.yaml   # 生产环境配置
│   └── config-staging.yaml # 预发布环境配置
├── scripts/                # 部署脚本
│   ├── build.sh           # 构建脚本
│   └── deploy.sh          # 部署脚本
└── README.md              # 本文档
```

## 前置要求

### 必需工具
- Docker
- kubectl
- Go 1.24+

### 可选工具
- kind (本地 Kubernetes 集群)
- helm (包管理器)

## 快速开始

### 1. 构建应用

```bash
# 构建最新版本
./deploy/scripts/build.sh

# 构建指定版本
./deploy/scripts/build.sh v1.0.0

# 构建指定平台
./deploy/scripts/build.sh v1.0.0 linux/arm64
```

### 2. 部署应用

```bash
# 部署到开发环境
./deploy/scripts/deploy.sh dev deploy

# 部署到预发布环境
./deploy/scripts/deploy.sh staging deploy

# 部署到生产环境
./deploy/scripts/deploy.sh prod deploy
```

### 3. 管理应用

```bash
# 查看状态
./deploy/scripts/deploy.sh dev status

# 查看日志
./deploy/scripts/deploy.sh dev logs

# 升级应用
./deploy/scripts/deploy.sh dev upgrade

# 回滚应用
./deploy/scripts/deploy.sh dev rollback

# 删除应用
./deploy/scripts/deploy.sh dev delete

# 端口转发（本地访问）
./deploy/scripts/deploy.sh dev forward
```

## 配置说明

### 环境变量

应用通过以下环境变量进行配置：

| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| MYSQL_HOST | MySQL 主机地址 | localhost | 是 |
| MYSQL_PORT | MySQL 端口 | 3306 | 否 |
| MYSQL_USERNAME | MySQL 用户名 | root | 是 |
| MYSQL_DATABASE | MySQL 数据库名 | open_ticket_system | 是 |
| MYSQL_PASSWORD | MySQL 密码 | - | 是 |
| JWT_SECRET | JWT 密钥 | - | 是 |

### 资源配置

默认资源配置：

```yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

### 健康检查

应用提供健康检查端点：

- **路径**: `/health`
- **方法**: GET
- **响应**: JSON 格式的健康状态

## 部署架构

### 组件

1. **Deployment**: 应用部署
2. **Service**: 内部服务暴露
3. **ConfigMap**: 配置管理
4. **Secret**: 敏感信息管理
5. **Ingress**: 外部访问入口
6. **HPA**: 自动扩缩容

### 网络

- **内部端口**: 8080
- **服务端口**: 80
- **外部访问**: 通过 Ingress

## 监控和日志

### 日志

- 日志输出到 stdout
- 支持结构化日志 (JSON)
- 可配置日志级别和轮转

### 监控

- 健康检查端点
- Kubernetes 原生监控
- 资源使用指标

## 故障排除

### 常见问题

1. **Pod 启动失败**
   - 检查配置是否正确
   - 查看 Pod 事件和日志
   - 验证数据库连接

2. **健康检查失败**
   - 确认应用正常启动
   - 检查端口配置
   - 验证健康检查路径

3. **数据库连接失败**
   - 检查数据库服务状态
   - 验证连接参数
   - 确认网络策略

### 调试命令

```bash
# 查看 Pod 状态
kubectl get pods -l app=open-ticket-system

# 查看 Pod 详情
kubectl describe pod <pod-name>

# 查看 Pod 日志
kubectl logs <pod-name>

# 查看服务状态
kubectl get svc -l app=open-ticket-system

# 查看事件
kubectl get events --sort-by='.lastTimestamp'
```

## 安全考虑

### 生产环境

1. **密钥管理**
   - 使用 Kubernetes Secrets 或外部密钥管理
   - 定期轮换密钥
   - 避免硬编码密钥

2. **网络策略**
   - 限制 Pod 间通信
   - 配置防火墙规则
   - 使用 TLS 加密

3. **资源限制**
   - 设置合理的资源限制
   - 监控资源使用
   - 配置自动扩缩容

## 扩展和定制

### 自定义配置

可以通过修改以下文件进行定制：

- `deployment.yaml`: 部署配置
- `config-*.yaml`: 环境配置
- `Dockerfile`: 容器构建

### 添加组件

1. 创建新的 YAML 配置文件
2. 更新部署脚本
3. 测试配置有效性

## 支持

如有问题，请：

1. 查看本文档
2. 检查应用日志
3. 查看 Kubernetes 事件
4. 联系开发团队

## 更新日志

- v1.0.0: 初始版本
  - 基础部署配置
  - 构建和部署脚本
  - 多环境支持
