APP=helwan-penguin

.PHONY: run build test clean
run:
	go run ./cmd/helwan-penguin
build:
	go build -o $(APP) ./cmd/helwan-penguin
test:
	go test ./...
clean:
	rm -f $(APP) $(APP).exe
