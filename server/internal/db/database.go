package db

import (
	"fmt"
	"log"
	"os"
	"server/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Document{},
		&models.DocumentHistory{},
		&models.DocumentAccess{},
	)
	log.Println("Database migrated")
}

func InitDB() {
	dsn := os.Getenv("DB_STRING")

	if dsn == "" {
		log.Fatal("DB_STRING is not set in the environment")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	fmt.Println("Connected to PostgreSQL")

	Migrate(DB)
}
