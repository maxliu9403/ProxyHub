# -------- 构建阶段 --------
FROM golang:1.23.5 AS builder

# 版本信息
ARG COMMITID
ENV COMMITID=${COMMITID:-v1.0}

# 多平台编译目标（由 buildx 传入）
ARG TARGETOS
ARG TARGETARCH

# 编译环境变量
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV CGO_ENABLED=0

WORKDIR /app

# 复制代码并编译
COPY . .
RUN go mod vendor
RUN go build -mod=vendor -ldflags="-s -w -X main.Build=${COMMITID}" -o proxyhub ./cmd/app

# -------- 运行阶段 --------
FROM alpine:3.18

# 安装证书和时区数据
RUN apk add --no-cache ca-certificates tzdata

# 创建工作目录
WORKDIR /app

# 复制编译后的二进制和配置文件
COPY --from=builder /app/proxyhub /app/proxyhub
COPY --from=builder /app/docs /app/docs

# 设置默认时区
ENV TZ=Asia/Shanghai

# 暴露服务端口
EXPOSE 9901

# 启动命令（绝对路径）
CMD ["/app/proxyhub", "--config", "/app/configs/dev.yaml"]

#docker buildx build \
#  --platform linux/arm64 \
#  -t proxyhub:v1.0.0 \
#  -f ./build/Dockerfile \
#  --output type=docker .