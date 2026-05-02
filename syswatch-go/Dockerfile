# Stage 1: Build
FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY main.go .

RUN go mod init syswatch && \
    go get github.com/moby/moby/client@latest && \
    go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o syswatch .

# Stage 2: Run
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/syswatch .

CMD ["./syswatch"]
