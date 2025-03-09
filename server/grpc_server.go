package main

import (
	"log"
	"net"
	"server/internal/db"
	handlers "server/internal/handlers/docs_handlers"
	proto "server/proto"

	"google.golang.org/grpc"
)

func main() {
	db.InitDB()

	grpcServer := grpc.NewServer()
	proto.RegisterDocumentServiceServer(grpcServer, &handlers.DocumentHandler{})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("gRPC Server is running on port 50051...")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
