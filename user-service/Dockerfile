# Stage duy nhất: build và chạy luôn
FROM golang:1.22-slim

# Tạo thư mục làm việc trong container
WORKDIR /app

# Copy go.mod và go.sum để tải dependencies trước
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy toàn bộ mã nguồn vào image
COPY . .

# Build file nhị phân tên là user-servicecd/

RUN go build -o user-service .

# Expose cổng mặc định nếu bạn chạy HTTP server ở đó
EXPOSE 8080

# Lệnh khởi chạy container
CMD ["./user-service"]