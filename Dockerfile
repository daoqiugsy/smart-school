# Dockerfile

# --- 第一阶段：构建 (Builder) ---
# 使用官方的 Go 语言镜像作为构建环境
# 我们选择一个具体的版本以保证一致性
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件，并下载依赖
# 这一步是单独的，因为依赖通常不经常变动，可以利用 Docker 的缓存机制
COPY go.mod go.sum ./
RUN go mod download

# 复制项目所有源代码到工作目录
COPY . .

# 编译我们的 Go 应用
# -o /app/server 指定输出文件名为 server
# CGO_ENABLED=0 禁用 CGO，以构建一个静态链接的二进制文件，这在 Alpine 这种轻量级镜像中很重要
# -ldflags="-w -s" 减小编译后文件的大小
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server -ldflags="-w -s" ./cmd/server/main.go


# --- 第二阶段：运行 (Runner) ---
# 使用一个非常轻量级的 Alpine Linux 镜像作为最终运行环境
# 这使得我们的最终镜像非常小
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从第一阶段 (builder) 中复制编译好的二进制文件和配置文件
# 我们只需要这两个东西就可以运行我们的应用了！
COPY --from=builder /app/server /app/server
COPY config/config.yaml /app/config/config.yaml

# 暴露我们的应用监听的端口 (请确保和你配置文件中的端口一致)
EXPOSE 8080

# 容器启动时执行的命令
# 运行我们的服务器，并指定配置文件的路径
ENTRYPOINT ["/app/server"]