#!/bin/bash

set -e

# 交互输入 docker-compose 路径
read -p "请输入 docker-compose 文件路径: " DOCKER_COMPOSE_FILE
if [ -z "$DOCKER_COMPOSE_FILE" ]; then
  echo "错误：docker-compose 路径不能为空"
  exit 1
fi

# 停止并移除容器、网络等
echo "正在停止 docker-compose 服务..."
docker-compose -f "$DOCKER_COMPOSE_FILE" down

echo "✅ 服务已成功关闭并移除"
