# 使用官方Golang镜像作为基础镜像
FROM golang:1.20.5
ENV GOPROXY https://goproxy.cn

# 设置工作目录
WORKDIR /app

# 复制应用程序源代码到容器中
COPY  ./ /app
# 进入 rpcServer/video 目录

RUN go mod tidy


