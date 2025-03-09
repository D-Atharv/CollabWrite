package docs_handlers

import (
	"context"
	"log"
	"server/internal/services"
	"server/internal/models"
	proto "server/proto"

	"github.com/google/uuid"
)

type DocumentHandler struct {
	proto.UnimplementedDocumentServiceServer
	Service docs_services.DocumentService
}

func (h *DocumentHandler) GetDocument(ctx context.Context, req *proto.GetDocumentRequest) (*proto.GetDocumentResponse, error) {
	docID := req.Id
	log.Printf("Received GetDocument request for ID: %s", docID)

	doc, err := h.Service.GetDocument(ctx, docID)
	if err != nil {
		return nil, err
	}

	return &proto.GetDocumentResponse{Document: h.Service.ConvertToProtoDocument(*doc)}, nil
}

func (h *DocumentHandler) CreateDocument(ctx context.Context, req *proto.CreateDocumentRequest) (*proto.CreateDocumentResponse, error) {
	log.Println("Received CreateDocument request")

	doc := models.Document{
		OwnerID: uuid.MustParse(req.OwnerId),
		Title:   req.Title,
		Content: req.Content,
	}

	docFull, err := h.Service.CreateDocument(ctx, &doc)
	if err != nil {
		return nil, err
	}

	return &proto.CreateDocumentResponse{Document: h.Service.ConvertToProtoDocument(*docFull)}, nil
}

func (h *DocumentHandler) UpdateDocument(ctx context.Context, req *proto.UpdateDocumentRequest) (*proto.UpdateDocumentResponse, error) {
	docID := req.Id
	log.Printf("Received UpdateDocument request for ID: %s", docID)

	doc, err := h.Service.UpdateDocument(ctx, docID, req.Title, req.Content)
	if err != nil {
		return nil, err
	}

	return &proto.UpdateDocumentResponse{Document: h.Service.ConvertToProtoDocument(*doc)}, nil
}

func (h *DocumentHandler) DeleteDocument(ctx context.Context, req *proto.DeleteDocumentRequest) (*proto.DeleteDocumentResponse, error) {
	docID := req.Id
	log.Printf("Received DeleteDocument request for ID: %s", docID)

	err := h.Service.DeleteDocument(ctx, docID)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteDocumentResponse{Message: "Document deleted successfully"}, nil
}
