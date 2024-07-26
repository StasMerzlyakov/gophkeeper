.PHONY: build proto clean test cover generate

build:
	go mod tidy
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.buildVersion=$(shell test `git tag --points-at HEAD` && git tag --points-at HEAD || echo N/A) \
	    -X main.buildDate=$(shell date +'%Y-%m-%d') -X main.buildCommit=$(shell git rev-parse HEAD)" -buildvcs=false  -o=build/server ./cmd/server/...
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.buildVersion=$(shell test `git tag --points-at HEAD` && git tag --points-at HEAD || echo N/A) \
	    -X main.buildDate=$(shell date +'%Y-%m-%d') -X main.buildCommit=$(shell git rev-parse HEAD)" -buildvcs=false  -o=build/client_win ./cmd/client/...
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.buildVersion=$(shell test `git tag --points-at HEAD` && git tag --points-at HEAD || echo N/A) \
	    -X main.buildDate=$(shell date +'%Y-%m-%d') -X main.buildCommit=$(shell git rev-parse HEAD)" -buildvcs=false  -o=build/client_mac ./cmd/client/...
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.buildVersion=$(shell test `git tag --points-at HEAD` && git tag --points-at HEAD || echo N/A) \
	    -X main.buildDate=$(shell date +'%Y-%m-%d') -X main.buildCommit=$(shell git rev-parse HEAD)" -buildvcs=false  -o=build/client_linux ./cmd/client/...

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

