PROTO_DIR=D:\Projects\WEB\Final\CollabWrite\server\proto
PROTO_FILES=$(PROTO_DIR)/*.proto
OUT_DIR=$(PROTO_DIR)
PROTOC_INCLUDE="C:\Program Files\protobuf\include"

# Generate gRPC and Go files
generate:
	protoc --proto_path=$(PROTOC_INCLUDE) --proto_path=$(PROTO_DIR) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_FILES)

# Clean generated files
clean:
	rm -f $(OUT_DIR)/*.pb.go $(OUT_DIR)/*_grpc.pb.go

# Run the Go server (Modify this based on your server entry point)
run:
	go run main.go

# Build the Go project
build:
	go build -o server.exe main.go


# protoc --proto_path="C:\Program Files\protobuf\include" --proto_path="D:\Projects\WEB\Final\CollabWrite\server\proto" --go_out="D:\Projects\WEB\Final\CollabWrite\server\proto" --go-grpc_out="D:\Projects\WEB\Final\CollabWrite\server\proto" "D:\Projects\WEB\Final\CollabWrite\server\proto\document.proto"

