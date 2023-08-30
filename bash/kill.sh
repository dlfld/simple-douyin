#!/bin/bash

# 定义要关闭的端口列表
port_list=(8080 8081 8082 8083 8084 8085)

# 关闭指定端口的服务
function stop_service {
  local port=$1
  echo "Stopping service on port $port..."
  lsof -ti :$port | xargs kill
}

# 循环关闭端口列表中的服务
for port in "${port_list[@]}"; do
  stop_service $port
done