package db

import (
	"fmt"
	"log"
	"os"
	"server/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
		Logger:      logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	fmt.Println("Connected to PostgreSQL")

	// Migrate(DB)
}

// package db

// import (
// 	"context"
// 	"log"

// 	"github.com/jackc/pgx/v5/pgxpool" // ✅ Use pgxpool for better connection handling
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/logger"
// )

// var DB *gorm.DB
// var PgPool *pgxpool.Pool

// func InitDB() {
// 	dsn := "postgresql://postgres.bcwnyyepiqycblvwlbtq:SupaBase@321@aws-0-ap-south-1.pooler.supabase.com:6543/postgres"

// 	// ✅ Initialize pgxpool
// 	var err error
// 	PgPool, err = pgxpool.New(context.Background(), dsn)
// 	if err != nil {
// 		log.Fatal("Failed to create pgxpool:", err)
// 	}

// 	// ✅ Use GORM with pgx driver
// 	DB, err = gorm.Open(postgres.New(postgres.Config{
// 		DSN:                  dsn,
// 		PreferSimpleProtocol: true, // ✅ Disables prepared statement caching
// 	}), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	})
// 	if err != nil {
// 		log.Fatal("Failed to connect to PostgreSQL:", err)
// 	}
// 	log.Println("Connected to PostgreSQL")

// 	// Migrate(DB)
// }


// // {
// //     "error": "rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:63848->127.0.0.1:50051: wsarecv: An existing connection was forcibly closed by the remote host."
// // }