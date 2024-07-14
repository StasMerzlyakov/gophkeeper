.PHONY: proto

proto:
	mkdir -p ./internal/proto
	protoc --go_out=./internal/proto --go_opt=paths=source_relative \
          --go-grpc_out=./internal/proto --go-grpc_opt=paths=source_relative \
          --proto_path proto\
          proto/gophkeeper.proto
