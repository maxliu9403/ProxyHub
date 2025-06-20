#!/bin/bash

set -e

# 交互输入镜像tar包路径
read -p "请输入镜像tar包路径: " IMAGE_TAR_PATH
if [ -z "$IMAGE_TAR_PATH" ]; then
  echo "错误：镜像tar包路径不能为空"
  exit 1
fi

# 交互输入docker-compose路径
read -p "请输入docker-compose文件路径: " DOCKER_COMPOSE_FILE
if [ -z "$DOCKER_COMPOSE_FILE" ]; then
  echo "错误：docker-compose路径不能为空"
  exit 1
fi

echo "将使用以下配置："
echo "镜像tar包路径: $IMAGE_TAR_PATH"
echo "docker-compose文件: $DOCKER_COMPOSE_FILE"

# 1. 加载镜像
if [ -f "$IMAGE_TAR_PATH" ]; then
  echo "正在加载镜像 $IMAGE_TAR_PATH ..."
  docker load -i "$IMAGE_TAR_PATH"
else
  echo "错误：镜像tar包不存在：$IMAGE_TAR_PATH"
  exit 1
fi

# 2. 启动服务
echo "重启 docker-compose 服务..."
docker-compose -f "$DOCKER_COMPOSE_FILE" down
docker-compose -f "$DOCKER_COMPOSE_FILE" up -d

echo "✅ 更新完成！"
