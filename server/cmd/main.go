package main

import (
	"log"
	"net/http"
	"server/internal/db"
	"server/internal/redis"
	"server/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//TODO: to setup KONG
func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.SetTrustedProxies(nil) // Disable trusted proxies
	// r.SetTrustedProxies([]string{"YOUR_PROXY_IP"}) // Replace with the actual proxy IP

	db.InitDB()
	redis.InitRedis()

	routes.SetupRoutes(r)

	r.GET("/ping",func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK,gin.H{
			"message": "Hello World",
		})
	})

	r.GET("/health",func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK,gin.H{
			"message": "Health Check OK",
		})
	})

	r.Run(":8080")
}