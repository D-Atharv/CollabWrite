package routes

import (
	"server/internal/auth"
	handlers "server/internal/handlers/doc_handlers"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/")
	{
		api.GET("/auth/google", auth.GoogleLogin)
		api.GET("/auth/google/callback", auth.GoogleCallback)

		protected := api.Group("/docs")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("", handlers.CreateDocument)
			protected.GET(":id", handlers.GetDocument)
			protected.PUT(":id", handlers.UpdateDocument)
			protected.DELETE(":id", handlers.DeleteDocument)
		}
	}
}
