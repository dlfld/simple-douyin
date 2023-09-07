FROM 15150276667/douyin

# 设置工作目录
WORKDIR /app

RUN rm -rf /app/*
COPY  ./ /app

RUN go build -o main .

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






