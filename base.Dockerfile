From golang:1.20.5
ENV GOPROXY https://goproxy.cn

COPY ./ /app
RUN cd /app && go mod tidy && rm -rf /app
RUN apt-get update && apt-get install -y --no-install-recommends \
    ffmpeg \
    && rm -rf /var/lib/apt/lists/*

