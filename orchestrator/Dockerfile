# 1. Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Build the binary
COPY . .
# CGO_ENABLED=0 creates a static binary (no dependencies on host libs)
RUN CGO_ENABLED=0 GOOS=linux go build -o custodian ./cmd/api/main.go

# 2. Run Stage (Tiny image)
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/custodian .
# Copy the .env file (or use docker-compose env vars)
COPY .env . 

EXPOSE 8000

CMD ["./custodian"]
