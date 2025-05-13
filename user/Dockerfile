# ---------- Build 阶段 ----------
FROM golang:alpine AS builder

# 设置 Go 代理为七牛云的代理
ENV GOPROXY=https://goproxy.cn,direct

# 安装需要的依赖
RUN apk update && apk add --no-cache git

# 切换到 app 目录，构建二进制文件
WORKDIR /app

# 拷贝整个 user 项目
COPY . .

RUN go mod tidy && go build -o user-app ./cmd/main.go

# ---------- 运行阶段 ----------
FROM alpine:latest

# 安装 tzdata 来设置时区
RUN apk add --no-cache tzdata

# 设置时区为 Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 设置工作目录为
WORKDIR /app

# 拷贝编译好的二进制文件和配置文件
COPY --from=builder /app/user-app .

# 开放 gRPC 端口（按你实际用的）
EXPOSE 30001

CMD ["./user-app"]
