FROM golang:1.25-alpine AS builder

ARG SERVICE

RUN apk add --no-cache tzdata

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin ./${SERVICE}

FROM alpine:3.19

ARG SERVICE

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

COPY --from=builder /app/bin /app/bin
COPY --from=builder /src/${SERVICE}/etc /app/etc

EXPOSE 8000-8100

ENTRYPOINT ["/app/bin"]
