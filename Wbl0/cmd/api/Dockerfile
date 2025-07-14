# Используем официальный образ Go 1.23
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./Wbl0/cmd/api

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

COPY --from=builder /app/main .

RUN chown appuser:appgroup main

USER appuser

EXPOSE 8081

CMD ["./main"]