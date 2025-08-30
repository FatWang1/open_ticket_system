# 多阶段构建Dockerfile
# 第一阶段：依赖下载阶段
FROM golang:1.25-alpine AS deps

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 设置Go代理
ENV GOPROXY=https://goproxy.cn,direct

# 下载依赖
RUN go mod 

# 第二阶段：构建阶段
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git ca-certificates tzdata build-base

# 设置Go环境变量
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# 从依赖阶段复制go mod缓存
COPY --from=deps /go/pkg/mod /go/pkg/mod

# 复制go mod文件
COPY go.mod go.sum ./

# 复制源代码
COPY . .

# 构建应用
RUN go build -ldflags="-w -s" -o open_ticket_system ./cmd/open_ticket_system

# 第三阶段：运行阶段
FROM alpine:latest

# 安装必要的系统工具和运行时依赖
RUN apk --no-cache add ca-certificates tzdata wget curl

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/open_ticket_system .

# 复制配置文件
COPY --from=builder /app/cmd/open_ticket_system/conf ./cmd/open_ticket_system/conf

# 创建必要的目录
RUN mkdir -p /app/logs && \
    chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 设置环境变量
ENV GIN_MODE=release
ENV ENV=production

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./open_ticket_system"]
