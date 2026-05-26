FROM golang:1.25-alpine AS builder

# 固定服务名（按需修改）
ARG SERVICE_NAME=userservice

# 换 alpine 源加速
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache tzdata

# 换 Go 代理
RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

# 固定编译 ./${SERVICE_NAME} 目录下的 main.go（或按实际入口调整）
RUN CGO_ENABLED=0 GOOS=linux go build -v -x -ldflags="-s -w" -o /app/bin ./${SERVICE_NAME}

FROM alpine:3.19

# 固定服务名（仅用于复制配置文件路径）
ARG SERVICE_NAME=userservice

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

COPY --from=builder /app/bin /app/bin
# 假设每个服务有自己的 etc 目录，路径如 userservice/etc
COPY --from=builder /src/${SERVICE_NAME}/etc /app/etc

EXPOSE 8001

ENTRYPOINT ["/app/bin"]