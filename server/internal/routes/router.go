package routes

import (
	"context"
	"server/internal/middleware"
	"server/internal/handlers/auth_handlers"
	proto "server/proto"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, grpcClient proto.DocumentServiceClient) {
	api := router.Group("/")
	{
		api.GET("/auth/google", auth_handlers.GoogleLogin)
		api.GET("/auth/google/callback", auth_handlers.GoogleCallback)

		protected := api.Group("/docs")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("", func(ctx *gin.Context) {
				var req proto.CreateDocumentRequest
				if err := ctx.ShouldBindJSON(&req); err != nil {
					ctx.JSON(400, gin.H{"error": "Invalid request"})
					return
				}
				resp, err := grpcClient.CreateDocument(context.Background(), &req)
				if err != nil {
					ctx.JSON(500, gin.H{"error": err.Error()})
					return
				}
				ctx.JSON(201, resp)
			})

			protected.GET(":id", func(ctx *gin.Context) {
				req := &proto.GetDocumentRequest{Id: ctx.Param("id")}
				resp, err := grpcClient.GetDocument(context.Background(), req)
				if err != nil {
					ctx.JSON(404, gin.H{"error": "Document not found"})
					return
				}
				ctx.JSON(200, resp)
			})

			protected.PUT(":id", func(ctx *gin.Context) {
				var req proto.UpdateDocumentRequest
				req.Id = ctx.Param("id")
				if err := ctx.ShouldBindJSON(&req); err != nil {
					ctx.JSON(400, gin.H{"error": "Invalid request"})
					return
				}
				resp, err := grpcClient.UpdateDocument(context.Background(), &req)
				if err != nil {
					ctx.JSON(500, gin.H{"error": err.Error()})
					return
				}
				ctx.JSON(200, resp)
			})

			protected.DELETE(":id", func(ctx *gin.Context) {
				req := &proto.DeleteDocumentRequest{Id: ctx.Param("id")}
				resp, err := grpcClient.DeleteDocument(context.Background(), req)
				if err != nil {
					ctx.JSON(500, gin.H{"error": err.Error()})
					return
				}
				ctx.JSON(200, resp)
			})
		}
	}
}
