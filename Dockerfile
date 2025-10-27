FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build ./cmd/main.go

CMD ["./main"]

EXPOSE 8080