#!/bin/bash
# 定义项目数组、对应的 main 文件和日志文件位置
projects=(
  "/root/simple-douyin/main.go:/root/douyinlogs/gin.log"
  "/root/simple-douyin/rpcServer/interaction/main.go:/root/douyinlogs/interaction.log"
  "/root/simple-douyin/rpcServer/message/main.go:/root/douyinlogs/message.log"
  "/root/simple-douyin/rpcServer/relation/main.go:/root/douyinlogs/relation.log"
  "/root/simple-douyin/rpcServer/user/main.go:/root/douyinlogs/user.log"
  "/root/simple-douyin/rpcServer/video/main.go:/root/douyinlogs/video.log"
)

# 启动项目并将日志输出到指定文件
function start_project {
  local project=$1
  local log_file=$2
  local project_dir=$(dirname "$project")
  local project_name=$(basename "$project" .go)

  echo "Starting $project_name..."
  cd "$project_dir"
  nohup go run "$project" > "$log_file" 2>&1 &
}

# 循环启动每个项目
for item in "${projects[@]}"; do
  IFS=':' read -ra parts <<< "$item"
  project="${parts[0]}"
  log_file="${parts[1]}"
  start_project "$project" "$log_file"
done