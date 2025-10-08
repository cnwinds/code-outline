# 使用官方Go镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包（支持 CGO）
RUN apk add --no-cache git gcc musl-dev g++

# 复制go模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（启用 CGO）
RUN CGO_ENABLED=1 GOOS=linux go build -a -o contextgen ./cmd/contextgen

# 使用精简的alpine镜像作为运行环境
FROM alpine:latest

# 安装ca-certificates（用于HTTPS请求）
RUN apk --no-cache add ca-certificates

# 创建非root用户
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/contextgen .

# 创建语法目录
RUN mkdir -p grammars

# 设置权限
RUN chown -R appuser:appgroup /root

# 切换到非root用户
USER appuser

# 暴露端口（如果需要）
# EXPOSE 8080

# 设置入口点
ENTRYPOINT ["./contextgen"]

# 默认命令
CMD ["--help"]
