#!/bin/bash

set -e

# 默认 docker-compose 文件路径
DOCKER_COMPOSE_FILE="./docker-compose.yaml"

# 自动寻找 .env 文件
ENV_FILE_DIR=$(dirname "$DOCKER_COMPOSE_FILE")
ENV_FILE="$ENV_FILE_DIR/.env"

if [ -f "$ENV_FILE" ]; then
  echo "📦 加载环境变量文件: $ENV_FILE"
  set -o allexport
  source "$ENV_FILE"
  set +o allexport
else
  echo "❌ 错误：未找到 .env 文件（路径: $ENV_FILE）"
  exit 1
fi

# 从 .env 中读取镜像 tar 包路径变量（例如：IMAGE_TAR_PATH）
if [ -z "$IMAGE_TAR_PATH" ]; then
  echo "❌ 错误：.env 中未定义 IMAGE_TAR_PATH"
  exit 1
fi

# 显示配置信息
echo "======= 使用配置 ======="
echo "镜像 tar 包路径: $IMAGE_TAR_PATH"
echo "docker-compose 文件: $DOCKER_COMPOSE_FILE"
echo "env 文件路径: $ENV_FILE"
echo "镜像名（从 .env 读取）: ${PROXYHUB_IMAGE:-未定义}"
echo "========================"

# 1. 加载镜像
if [ -f "$IMAGE_TAR_PATH" ]; then
  echo "🧊 正在加载镜像: $IMAGE_TAR_PATH ..."
  docker load -i "$IMAGE_TAR_PATH"
else
  echo "❌ 错误：镜像 tar 包不存在：$IMAGE_TAR_PATH"
  exit 1
fi

# 2. 重启服务
echo "🚀 重启 docker-compose 服务..."
docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" down
docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" up -d

echo "✅ 更新完成！"
