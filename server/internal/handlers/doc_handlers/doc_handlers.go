package handlers

import (
	"encoding/json"
	// "log"
	"net/http"
	"server/internal/db"
	"server/internal/models"
	"server/internal/redis"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//TODO: fix the database error and split this into services and handlers
func CreateDocument(ctx *gin.Context) {
	var doc models.Document

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := ctx.ShouldBindJSON(&doc); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	ownerID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}

	doc.OwnerID = ownerID

	if err := db.DB.Create(&doc).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
		return
	}

	docJSON, _ := json.Marshal(doc)
	err = redis.RedisClient.HSet(ctx, "documents", doc.ID.String(), docJSON).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache document"})
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Document created successfully", "data": doc})
}

func GetDocument(ctx *gin.Context) {
	doc_id := ctx.Param("id")
	if doc_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	// Check Redis cache
	cachedDoc, err := redis.RedisClient.Get(ctx, "document:"+doc_id).Result()
	if err == nil && cachedDoc != "" {
		var doc models.Document
		if json.Unmarshal([]byte(cachedDoc), &doc) == nil {
			ctx.JSON(http.StatusOK, gin.H{"message": "Document fetched from cache", "data": doc})
			return
		}
	}

	// Fetch from DB
	var doc models.Document
	if err := db.DB.First(&doc, "id = ?", doc_id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Cache document in Redis with expiry (e.g., 10 minutes)
	docJSON, _ := json.Marshal(doc)
	err = redis.RedisClient.Set(ctx, "document:"+doc_id, docJSON, 10*time.Minute).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache document"})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Document fetched successfully", "data": doc})
}


func UpdateDocument(ctx *gin.Context) {
	doc_id := ctx.Param("id")
	if doc_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var doc models.Document
	if err := db.DB.First(&doc, "id = ?", doc_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	ownerID := doc.OwnerID.String()
	if ownerID != userID.(string) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this document"})
		return
	}

	var updatedDoc models.Document
	if err := ctx.ShouldBindJSON(&updatedDoc); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	updatedDoc.ID = doc.ID
	updatedDoc.OwnerID = doc.OwnerID

	if err := db.DB.Save(&updatedDoc).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	docJSON, _ := json.Marshal(updatedDoc)
	err := redis.RedisClient.HSet(ctx, "documents", doc_id, docJSON).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache document"})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Document updated successfully", "data": updatedDoc})
}

func DeleteDocument(ctx *gin.Context) {
	doc_id := ctx.Param("id")
	if doc_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var doc models.Document
	if err := db.DB.First(&doc, "id = ?", doc_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if doc.OwnerID.String() != userID.(string) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this document"})
		return
	}

	if err := db.DB.Delete(&doc).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	err := redis.RedisClient.HDel(ctx, "documents", doc_id).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove from cache",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
