build:
	docker build -t sensezae/go-debugger:1.0.0 . --platform linux/amd64

run:
	go run cmd/main.go