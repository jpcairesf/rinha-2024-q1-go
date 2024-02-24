build:
	go build .
run:
	go run .
compile:
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-386 main.go

