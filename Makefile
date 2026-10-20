APP = helwan-penguin

.PHONY: build test run clean

build:
	go build -o bin/$(APP) ./cmd/helwan-penguin

test:
	go test ./...

run:
	go run ./cmd/helwan-penguin

clean:
	rm -rf bin
