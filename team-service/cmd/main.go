package main

import (
	"log"
	"team-service/config"
	"team-service/database"
	"team-service/routers"
	"team-service/swagger"

	"github.com/joho/godotenv"
)

// @title           Team Service API
// @version         1.0
// @description     API quản lý team
// @host            localhost:8081
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Nhập token theo dạng: Bearer <token>

// @security BearerAuth

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load file .env")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Unable to load config")
	}

	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Connect database failed")
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Migrate database failed")
	}

	// Mount các route app
	r := routers.SetupRouter(db)

	// Mount Swagger
	swagger.SetupSwagger(r)

	err = r.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
