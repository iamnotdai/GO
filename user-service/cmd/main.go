package main

import (
	"log"
	"user-service/config"
	"user-service/database"
	"user-service/routers"

	"github.com/joho/godotenv"
)

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

	r := routers.SetupRouter(db)

	err = r.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
