ALL: clean lint test build

build:
	go build -mod=vendor -o build/robinhood-export cmd/robinhood-export/*.go

lint: gofmt goimports
	docker run --rm -e LOG_LEVEL=error -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

test:
	go test ./...

gofmt:
	gofmt -w .

goimports:
	goimports -local "github.com/vitoordaz/robinhood-export" -w .

clean:
	rm -rf build
