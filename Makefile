test:
	go test ./...

ALL: test
	go build -mod=vendor -i -o build/robinhood-export cmd/robinhood-export/*.go