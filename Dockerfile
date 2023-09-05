# 使用官方Golang镜像作为基础镜像
FROM golang:1.20.5
ENV GOPROXY https://goproxy.cn

# 设置工作目录
WORKDIR /app

# 复制应用程序源代码到容器中
COPY  ./ /app
# 进入 rpcServer/video 目录

RUN apt-get update && apt-get install -y --no-install-recommends ffmpeg

RUN go mod tidy && go build -o main .

WORKDIR /app/rpcServer/video
RUN go build -o main .

WORKDIR /app/rpcServer/user
RUN go build -o main .

WORKDIR /app/rpcServer/message
RUN go build -o main .

WORKDIR /app/rpcServer/interaction
RUN go build -o main .

WORKDIR /app/rpcServer/relation
RUN go build -o main .






