.PHONY: proto clean test cover

test: clean
	go vet ./...
	golangci-lint run ./...
	go mod tidy
	go clean -testcache
	go test ./... -coverprofile cover.out

proto:
	mkdir -p ./internal/proto
	protoc --go_out=./internal/proto --go_opt=paths=source_relative \
          --go-grpc_out=./internal/proto --go-grpc_opt=paths=source_relative \
          --proto_path proto\
          proto/gophkeeper.proto

cover: test
	go tool cover -html=cover.out -o coverage.html
	firefox coverage.html &	
