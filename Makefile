BINARY_NAME=golang-scraper

build:
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin ./cmd/server/main.go
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux ./cmd/server/main.go
	GOARCH=amd64 GOOS=windows go build -o ./bin/${BINARY_NAME}-windows.exe ./cmd/server/main.go

run:
	go get ./...
	make build
	./bin/${BINARY_NAME}-linux

clean:
	go clean
	rm rf ./bin

test:
	go test ./... -v --race

test_coverage:
	go test ./... -coverprofile=coverage.out

lint:
	golangci-lint run --enable-all
