# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order-service .

# Stage 2: Run (Dùng scratch để image chỉ có vài MB)
FROM scratch
WORKDIR /root/
COPY --from=builder /app/order-service .
EXPOSE 8081
CMD ["./order-service"]
