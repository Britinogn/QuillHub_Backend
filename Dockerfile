# Start from a lightweight Go base image
# golang:1.23-alpine is a minimal Alpine Linux image with Go 1.23
FROM golang:1.25.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first (for better caching)
COPY go.mod go.sum ./

# Download dependencies (this layer will be cached)
RUN go mod download

# Copy vendor folder if it exists (much faster than go mod download)
COPY vendor ./vendor

# Copy source code
COPY . .

# Build with vendor (if available) or normal mode
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -o server ./cmd/server || \
    go build -o server ./cmd/server

# -------- Runtime stage --------
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port (use environment variable)
EXPOSE ${PORT:-8080}

# Run the application
CMD ["./server"]






# Stop containers
# docker-compose down

# Stop and remove volumes (clean slate)
# docker-compose down -v

# View logs
# docker-compose logs -f

# View logs for specific service
# docker-compose logs -f server

# Restart a specific service
# docker-compose restart server

# Rebuild and start
# docker-compose up --build --force-recreate

# Check running containers
# docker-compose ps

#2️⃣ Start with Hot Reload (Development Mode):

# docker-compose watch

# 1️⃣ Start Docker Compose (Production Mode):

# docker-compose up
# docker-compose up --build
