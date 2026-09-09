.PHONY: run test build clean

run:
	go run cmd/api/main.go

build:
	go build -o bin/contable-fix cmd/api/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
