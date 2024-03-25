# Build stage
FROM golang:1.22 as builder
WORKDIR /app

# Copy the Go Modules manifests and download the dependencies
COPY src/go.* ./
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY src/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o catpics-api .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage to the final stage
COPY --from=builder /app/catpics-api .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./catpics-api"]
