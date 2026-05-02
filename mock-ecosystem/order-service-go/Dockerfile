# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Thay đổi cực kỳ quan trọng ở đây: Khởi tạo module go trước khi build
RUN go mod init order-service

# Lưu ý: Ở đây tôi thấy bạn vẫn đang dùng tên order.go thay vì main.go
COPY order.go .

# Build code (Đã cập nhật để trỏ đúng tên file order.go)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order-service order.go

# Stage 2: Run (Dùng scratch để image chỉ có vài MB)
FROM scratch
WORKDIR /root/
COPY --from=builder /app/order-service .
EXPOSE 8081
CMD ["./order-service"]
