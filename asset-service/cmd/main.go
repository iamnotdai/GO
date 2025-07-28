package main

import (
	"asset-service/config"
	"asset-service/database"
	"asset-service/routers"
	"log"

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

	// Mount các route app
	r := routers.SetupRouter(db)

	err = r.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
