# Use official Go image with version >= 1.25.2
FROM golang:1.25.3 AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN go build -tags netgo -ldflags="-s -w" -o app

# Use a small final image
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary from builder
COPY --from=builder /app/app .

# Expose port (change if your app uses another port)
EXPOSE 8080

# Run the binary
CMD ["./app"]
