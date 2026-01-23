# -------- Builder stage --------
FROM golang:1.25.3 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the app from cmd/server
RUN go build -tags netgo -ldflags="-s -w" -o app ./cmd/server

# -------- Runtime stage --------
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
