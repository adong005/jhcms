package main

import (
	"adcms-backend/internal/bootstrap"
	"adcms-backend/internal/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully")
	log.Printf("Starting database migration and seed (mode=%s)...", cfg.Database.InitMode)

	if err := bootstrap.InitDatabase(db, cfg.Database.InitMode); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Database migration completed successfully")
}
