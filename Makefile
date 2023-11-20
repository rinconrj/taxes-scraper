BINARY_NAME=golang-scraper

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin ./cmd/server/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux ./cmd/server/main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows ./cmd/server/main.go

run:
	make build
	./${BINARY_NAME}-linux

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all
