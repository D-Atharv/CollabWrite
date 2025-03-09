package main

import (
	"log"
	"net/http"
	"server/internal/db"
	"server/internal/redis"
	"server/internal/routes"
	proto "server/proto"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: to setup KONG and fix the gRPC connection which is deprecated
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

	// Set up gRPC client connection with updated credentials
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	grpcClient := proto.NewDocumentServiceClient(conn)
	routes.SetupRoutes(r, grpcClient)
	// routes.SetupRoutes(r)

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	log.Println("Starting HTTP server on :8080")
	r.Run(":8080")
}