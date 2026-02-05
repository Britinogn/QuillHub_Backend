# Start from a lightweight Go base image
# golang:1.23-alpine is a minimal Alpine Linux image with Go 1.23
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
# All subsequent commands will run from this directory
WORKDIR /app

# Copy go.mod and go.sum first
# We do this separately to take advantage of Docker's layer caching
# If these files don't change, Docker won't re-download dependencies
COPY go.mod go.sum ./

# Download all Go module dependencies
# This downloads packages listed in go.mod
RUN go mod download

# Copy all remaining application files into the container
# This includes your Go source code (.go files)
COPY . .

# Build the Go binary
# CGO_ENABLED=0: Build a static binary without C dependencies
# GOOS=linux: Target Linux operating system
# GOARCH=amd64: Target 64-bit architecture
# -o server: Output binary named 'server'
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./cmd/server

# -------- Runtime stage --------
# Use a minimal Alpine Linux image for the final container
FROM alpine:3.19

# Set the working directory in the runtime container
WORKDIR /app

# Install CA certificates for HTTPS calls
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server .

# Tell Docker that the container will listen on the specified port
# This doesn't actually publish the port, just documents it
EXPOSE $PORT

# The command to run when the container starts
# This runs the compiled Go binary
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
