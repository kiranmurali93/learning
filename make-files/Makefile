hello:
	echo "Hello"

build:
	go build -o bin/main main.go

run:
	go run main.go

compile:
	echo "Compiling for different os"
	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm main.go
	GOOS=linux   GOARCH=arm64 go build -o bin/main-linux-arm64 main.go

all: run compile