ALL: clean lint test build

build:
	go build -mod=vendor -o build/robinhood-export cmd/robinhood-export/*.go

lint:
	docker run --rm -e LOG_LEVEL=error -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

test:
	go test ./...

clean:
	rm -rf build