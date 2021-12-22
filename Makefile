test:
	go test -v -race ./... -cover

lint:
	golangci-lint run ./...

build:
	go build ./...