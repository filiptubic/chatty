all: clean build

run:
	go run . server

build:
	go build -o ./dist/chatty

clean:
	rm -rf ./dist

lint:
	golangci-lint run ./...
