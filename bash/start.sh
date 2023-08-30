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

# 定义项目数组、对应的 main 文件和日志文件位置
projects=(
  "/Users/wuwei/go/src/simple-douyin/main.go:/Users/wuwei/go/src/logs/gin.log"
  "/Users/wuwei/go/src/simple-douyin/rpcServer/interaction/main.go:/Users/wuwei/go/src/logs/interaction.log"
  "/Users/wuwei/go/src/simple-douyin/rpcServer/message/main.go:/Users/wuwei/go/src/logs/message.log"
  "/Users/wuwei/go/src/simple-douyin/rpcServer/relation/main.go:/Users/wuwei/go/src/logs/relation.log"
  "/Users/wuwei/go/src/simple-douyin/rpcServer/user/main.go:/Users/wuwei/go/src/logs/user.log"
  "/Users/wuwei/go/src/simple-douyin/rpcServer/video/main.go:/Users/wuwei/go/src/logs/video.log"
)

# 启动项目并将日志输出到指定文件
function start_project {
  local project=$1
  local log_file=$2

  go run "$project" 2>&1 | ts '[%Y-%m-%d %H:%M:%S]' >> "$log_file" &
}

# 循环启动每个项目
for item in "${projects[@]}"; do
  IFS=':' read -ra parts <<< "$item"
  project="${parts[0]}"
  log_file="${parts[1]}"
  start_project "$project" "$log_file"
done