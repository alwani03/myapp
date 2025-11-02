run:
	go run ./cmd/server

build:
	go build -o bin/myapp ./cmd/server

tidy:
	go mod tidy

test:
	go test ./...