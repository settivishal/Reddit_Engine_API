# Makefile

# Directories
PROTO_DIR := .
GO_OUT_DIR := .

# Protobuf Compiler and Plugin
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go

# Proto Files
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# Go Files (generated)
GEN_GO_FILES := $(patsubst $(PROTO_DIR)/%.proto, $(GO_OUT_DIR)/%.pb.go, $(PROTO_FILES))

# Targets
.PHONY: all proto build clean run-engine run-simulator help

# Default target: Generate protobuf files and build the project
all: proto build

# Help command to list all targets
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  proto           Generate Go code from Protobuf files"
	@echo "  build           Build the engine and simulator binaries"
	@echo "  run-engine      Run the engine process (starts on port 8080)"
	@echo "  all             Generate Protobuf files and build the project"
	@echo "  clean           Remove generated files and binaries"
	@echo "  help            Show this help message"

# Generate Go files from .proto files
proto: 
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative messages/messages.proto

$(GO_OUT_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	@echo "Generating $@ from $<"
	$(PROTOC) --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative $<

# Build the project
build:
	@echo "Building the project..."
	go build -o engine ./engine/main.go

# Run the engine
run-engine:
	@echo "Starting the engine..."
	@go run ./engine/main.go


# Clean up generated and built files
clean:
	@echo "Cleaning up..."
	@rm -f $(GEN_GO_FILES) engine/main simulator/main
