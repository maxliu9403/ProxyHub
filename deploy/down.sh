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
  echo "⚠️ 未找到 .env 文件（路径: $ENV_FILE），将继续执行但可能缺少变量"
fi

# 停止并移除容器、网络、卷
echo "🛑 正在停止 docker-compose 服务..."
docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" down

echo "✅ 服务已成功关闭并移除"
