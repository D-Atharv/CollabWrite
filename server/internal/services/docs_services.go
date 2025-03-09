package docs_services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"server/internal/db"
	"server/internal/models"
	"server/internal/redis"
	"time"
	proto "server/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

//TODO: fix time format in proto file

type DocumentService struct{}

func (s *DocumentService) ConvertToProtoDocument(doc models.Document) *proto.Document {
	protoDoc := &proto.Document{
		Id:        doc.ID.String(),
		OwnerId:   doc.OwnerID.String(),
		Title:     doc.Title,
		Content:   doc.Content,
		CreatedAt: timestamppb.New(doc.CreatedAt),
		UpdatedAt: timestamppb.New(doc.UpdatedAt),
	}

	if doc.Owner.ID != uuid.Nil {
		protoDoc.OwnerEmail = doc.Owner.Email
		protoDoc.OwnerProvider = doc.Owner.Provider
	}

	return protoDoc
}

func (s *DocumentService) FetchDocument(docID string) (*models.Document, error) {
	var doc models.Document
	err := db.DB.Preload("Owner").First(&doc, "id = ?", docID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "Document not found")
		}
		return nil, status.Errorf(codes.Internal, "Database error: %v", err)
	}
	return &doc, nil
}

func (s *DocumentService) CacheDocument(ctx context.Context, doc models.Document) {
	if redis.RedisClient != nil {
		docJSON, _ := json.Marshal(doc)
		err := redis.RedisClient.HSet(ctx, "documents", doc.ID.String(), docJSON).Err()
		if err != nil {
			log.Println("Redis caching error:", err)
		}
	}
}

// Get document (from cache or DB)
func (s *DocumentService) GetDocument(ctx context.Context, docID string) (*models.Document, error) {
	if redis.RedisClient != nil {
		cachedDoc, err := redis.RedisClient.Get(ctx, "document:"+docID).Result()
		if err == nil && cachedDoc != "" {
			var doc models.Document
			if json.Unmarshal([]byte(cachedDoc), &doc) == nil {
				return &doc, nil
			}
		}
	}

	// Fetch from DB if not in cache
	return s.FetchDocument(docID)
}

func (s *DocumentService) CreateDocument(ctx context.Context, doc *models.Document) (*models.Document, error) {
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()

	if err := db.DB.Create(&doc).Error; err != nil {
		return nil, err
	}

	docFull, err := s.FetchDocument(doc.ID.String())
	if err != nil {
		return nil, err
	}

	s.CacheDocument(ctx, *docFull)
	return docFull, nil
}

func (s *DocumentService) UpdateDocument(ctx context.Context, docID string, title, content string) (*models.Document, error) {
	doc, err := s.FetchDocument(docID)
	if err != nil {
		return nil, err
	}

	doc.Title = title
	doc.Content = content
	doc.UpdatedAt = time.Now()

	if err := db.DB.Save(doc).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update document: %v", err)
	}

	s.CacheDocument(ctx, *doc)
	return doc, nil
}

func (s *DocumentService) DeleteDocument(ctx context.Context, docID string) error {
	doc, err := s.FetchDocument(docID)
	if err != nil {
		return err
	}

	if err := db.DB.Delete(&doc).Error; err != nil {
		return status.Errorf(codes.Internal, "Failed to delete document: %v", err)
	}

	if redis.RedisClient != nil {
		err := redis.RedisClient.HDel(ctx, "documents", docID).Err()
		if err != nil {
			log.Println("Redis cache delete error:", err)
		}
	}

	return nil
}
