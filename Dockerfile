FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o main ./cmd/main.go

# 使用轻量级alpine镜像
FROM alpine:latest

WORKDIR /root/

# 从builder阶段复制编译好的应用
COPY --from=builder /app/main .

# 创建日志目录
RUN mkdir -p logs

# 暴露端口
EXPOSE 3006

# 启动应用
CMD ["./main"] 