# Khởi tạo go
go mod init team-service

# Cài đặt Gin và GORM
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres  // hoặc mysql, sqlite tùy DB bạn dùng

# Sau khi cài đặt Gin và GORM. Nếu cái nào không dùng thì clear
go mod tidy

# Chạy localhost 
go run cmd/main.go

# Add Swagger
# Step 1:
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files

# Step 2: phải chứa comment này ở main.go
// @title           Team Service API
// @version         1.0
// @description     API cho quản lý team
// @host            localhost:8080
// @BasePath        /

# Step 3:
swag init --generalInfo cmd/main.go  // -- output docs

# Nếu bị lỗi câu lệnh swag thì install
go install github.com/swaggo/swag/cmd/swag@latest