FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY main.go .
RUN swag init
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o catpics-api .
FROM ubuntu:latest
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /root/
COPY --from=builder /app/catpics-api .
EXPOSE 8080
CMD ["./catpics-api"]
