.PHONY: proto clean test cover generate

test: clean
	GOOS=linux GOARCH=amd64 go build -buildvcs=false -o=cmd/staticlint ./cmd/staticlint/...
	cmd/staticlint/staticlint ./...
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

gen:
	go generate ./...

