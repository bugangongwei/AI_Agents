# 阶段 1: 构建阶段 (Builder Stage)
# 这个基础镜像已经设置了 GOROOT 和 GOPATH
FROM golang:1.21-alpine AS builder

# 1. 设置工作目录
WORKDIR /app

# 2. 复制 go.mod 和 go.sum 并下载依赖 (利用缓存)
COPY go.mod go.sum ./
RUN go mod download

# 3. 复制源代码
COPY . .

# 4. 构建 Go 应用，并将其编译为静态二进制文件
# -v: 开启版本信息嵌入
# -ldflags: 移除调试信息，减小镜像体积
RUN go build -o /outfit-agent -ldflags '-s -w' ./main.go

# 阶段 2: 运行阶段 (Minimal Runtime Stage)
# 使用极小基础镜像，降低最终镜像大小
FROM alpine:latest
WORKDIR /root/
# 复制构建阶段生成的二进制文件
COPY --from=builder /outfit-agent .

# 暴露端口，运行应用程序
EXPOSE 8080
CMD ["./outfit-agent"]