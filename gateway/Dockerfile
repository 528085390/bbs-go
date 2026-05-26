FROM golang:1.25-alpine AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache tzdata

RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -x -ldflags="-s -w" -o /app/bin ./gateway

FROM alpine:3.19

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

COPY --from=builder /app/bin /app/bin
COPY --from=builder /src/gateway/etc /app/etc

EXPOSE 8888

ENTRYPOINT ["/app/bin"]
