package routes

import (
	"server/internal/handlers/auth_handlers"
	"server/internal/handlers/doc_handlers"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/")
	{
		api.GET("/auth/google", auth_handlers.GoogleLogin)
		api.GET("/auth/google/callback", auth_handlers.GoogleCallback)

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
