build:
	@go build -o bin/golem

run: build
	@./bin/golem

test:
	@go test ./... -v