all: lint test

test:
	go test ./...
lint:
	golangci-lint run -E gofmt -E golint --exclude-use-default=false
