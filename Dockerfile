From amd64/centos

# 设置工作目录
WORKDIR /app

RUN rm -rf /app/*
COPY  ./ /app





