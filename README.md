# QuillHub 🪶

A modern, scalable backend API for a content sharing platform built with Go, PostgreSQL, and cloud-based image storage.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Development](#development)
- [Docker Deployment](#docker-deployment)
- [Contributing](#contributing)

## 📖 Overview

QuillHub is a RESTful API backend designed for a modern content sharing platform. It supports user authentication, post creation and management, image uploads to Cloudinary, and role-based access control.

## ✨ Features

- **User Authentication**: JWT-based authentication with secure password hashing
- **User Management**: User registration, login, and profile management
- **Content Management**: Create, read, update, and delete posts
- **Image Uploads**: Cloudinary integration for image storage and optimization
- **Role-Based Access Control**: Admin and user roles with middleware protection
- **CORS Support**: Pre-configured for frontend integration
- **PostgreSQL Database**: Reliable data persistence with UUID support
- **Redis Caching**: Ready for cache layer implementation (optional)
- **Graceful Shutdown**: Proper server shutdown handling
- **Development Hot Reload**: Air integration for auto-reload during development

## 🛠 Tech Stack

- **Language**: Go 1.25.4
- **Web Framework**: Gin Gonic
- **Database**: PostgreSQL 18
- **ORM/Database Driver**: pgx v5 (PostgreSQL driver)
- **Authentication**: JWT (golang-jwt)
- **Password Hashing**: golang.org/x/crypto
- **Image Storage**: Cloudinary
- **Cache**: Redis (optional)
- **CORS**: gin-contrib/cors
- **Environment Management**: godotenv
- **Containerization**: Docker & Docker Compose
- **Development Tool**: Air (live reload)

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: Version 1.25.4 or higher ([Download](https://golang.org/dl/))
- **PostgreSQL**: Version 16 or higher ([Download](https://www.postgresql.org/download/))
- **Redis**: Optional but recommended ([Download](https://redis.io/download))
- **Docker & Docker Compose**: For containerized development ([Download](https://www.docker.com/products/docker-desktop))
- **Git**: For version control ([Download](https://git-scm.com/))

### Optional Tools

- **Postman** or **cURL**: For testing API endpoints
- **pgAdmin**: For database management UI

## ⚙️ Setup

### Step 1: Clone the Repository

```bash
git clone https://github.com/britinogn/quillhub.git
cd quillhub/server
```

### Step 2: Install Dependencies

```bash
go mod download
go mod tidy
```

### Step 3: Configure Environment

Copy the environment template and configure with your settings:

```bash
cp config/.env.example .env
```

### Step 4: Set Up Database

**Option A: Using Docker Compose (Recommended)**

```bash
docker-compose up -d db
```

**Option B: Local PostgreSQL**

Create database and run migrations:

```bash
psql -U postgres -d quill_hub -f migrations/001_init.sql
```

## 🚀 Running the Application

### Option 1: Local Development (Recommended for Development)

```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest

# Run with Air
air
```

The server will start on `http://localhost:8080` and automatically reload on file changes.

### Option 2: Standard Go Run

```bash
go run cmd/server/main.go
```

### Option 3: Build and Run Binary

```bash
# Build the binary
go build -o build/quillhub cmd/server/main.go

# Run the binary
./build/quillhub
```

## 🐳 Docker Deployment

### Run All Services with Docker Compose

Start the complete stack (PostgreSQL, Redis, and Go server):

```bash
docker-compose up
```

For development with hot reload:

```bash
docker-compose up -d
```

For production build:

```bash
docker-compose -f docker-compose.yml up --build
```

### Stop Services

```bash
docker-compose down

# Remove volumes (⚠️ this deletes data)
docker-compose down -v
```

### View Logs

```bash
# All services
docker-compose logs

# Specific service
docker-compose logs server
docker-compose logs db
docker-compose logs redis
```

## 📡 API Endpoints

### Health Check

```
GET /api/health
```

Returns the API health status and version.

### Authentication

#### Sign Up
```
POST /api/auth/signup
Content-Type: application/json

{
  "name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

#### Login
```
POST /api/auth/login
Content-Type: application/json

{
  "identifier": "john@example.com",
  "password": "securepassword123"
}
```

Returns JWT token for authenticated requests.

### Posts (Protected Routes - Require Auth)

#### Create Post
```
POST /api/posts
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json

{
  "title": "My First Post",
  "content": "This is the content of my first post",
  "tags": ["writing", "thoughts"],
  "author_id": "user-uuid"
}
```

**Note**: More post endpoints (GET, PUT, DELETE) are commented in `routes.go` and ready to be implemented.

### Root Endpoint

```
GET /api/
```

Returns welcome message.

## 📁 Project Structure

```
server/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── config/
│   ├── config.go                  # Configuration loader
│   ├── config.yaml               # YAML configuration
│   └── .env.example              # Environment variables template
├── internal/
│   ├── database/
│   │   ├── db.go                 # Database initialization
│   │   ├── postgres.go           # PostgreSQL connection
│   │   └── cloudinary.go         # Cloudinary setup
│   ├── handlers/
│   │   ├── auth_handler.go       # Authentication endpoints
│   │   └── post_handler.go       # Post management endpoints
│   ├── middleware/
│   │   ├── authMiddleware.go     # JWT verification
│   │   ├── authRoleMiddleware.go # Role-based access
│   │   └── uploadMiddleware.go   # File upload handling
│   ├── model/
│   │   ├── User.go               # User data model
│   │   ├── Post.go               # Post data model
│   │   └── Comment.go            # Comment data model
│   ├── repository/
│   │   ├── user_repository.go    # User data access
│   │   ├── post_repository.go    # Post data access
│   │   └── post_repository_impl.go # Post implementation
│   ├── routes/
│   │   └── routes.go             # Route registration
│   ├── services/
│   │   ├── auth_service.go       # Authentication logic
│   │   ├── post_service.go       # Post business logic
│   │   └── cloudinary_service.go # Image upload logic
│   └── templates/
│       └── index.html            # HTML template
├── pkg/
│   ├── logger/
│   │   └── logger.go             # Logging utilities
│   ├── response/
│   │   └── response.go           # Response formatting
│   └── utils/
│       ├── jwt.go                # JWT utilities
│       ├── hashPassword.go       # Password hashing
│       └── response.go           # Response helpers
├── migrations/
│   └── 001_init.sql              # Initial database schema
├── deployments/                  # Deployment configurations
├── docker-compose.yml            # Docker Compose configuration
├── Dockerfile                    # Production Docker image
├── Dockerfile.dev               # Development Docker image
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── .env                         # Environment variables (not in git)
├── .env.example                 # Environment template
├── .air.toml                    # Air configuration
├── .gitignore                   # Git ignore rules
└── README.md                    # This file
```



## 🔧 Development

### Prerequisites for Development

Install development dependencies:

```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest

# Install Go tools (optional)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Running Tests

```bash
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Code Formatting

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run
```

### Adding New Dependencies

```bash
go get github.com/package/module

# Update dependencies
go mod tidy
```

## Database Workflow

### First Time Setup

1. Ensure PostgreSQL is running
2. Run migrations: `psql -U postgres -d quill_hub -f migrations/001_init.sql`
3. Start the application

### Backup Database

```bash
pg_dump -U postgres quill_hub > backup_$(date +%Y%m%d).sql
```

### Restore Database

```bash
psql -U postgres quill_hub < backup_20240209.sql
```

## 🐳 Docker Troubleshooting

### Database Connection Issues

```bash
# Check if database container is running
docker-compose ps

# Restart containers
docker-compose restart

# View database logs
docker-compose logs db
```

### Port Already in Use

If port 8080 is already in use, modify the `.env` file:

```env
PORT=8081
```

Then update `docker-compose.yml` port mapping if using Docker.

### Permission Denied on Windows

Run Docker Desktop with Administrator privileges.

## 🚢 Building for Production

### Build Docker Image

```bash
docker build -f Dockerfile -t quillhub:latest .
```

### Environment for Production

Create a `.env.production` file with production values:

```env
ENVIRONMENT=production
GIN_MODE=release
DB_SSLMODE=require
# ... other production values
```

## 📊 API Response Format

All API responses follow a standard format:

### Success Response (2xx)
```json
{
  "message": "Success message",
  "data": {
    "key": "value"
  }
}
```

### Error Response (4xx/5xx)
```json
{
  "error": "Error message",
  "status": 400
}
```

## 🤝 Contributing

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Make your changes and commit: `git commit -am 'Add your feature'`
3. Push to the branch: `git push origin feature/your-feature`
4. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For issues and questions:

1. Check existing issues on GitHub
2. Create a new issue with detailed description
3. Include error logs and environment details

## 🔗 Useful Links

- [Go Documentation](https://golang.org/doc/)
- [Gin Documentation](https://gin-gonic.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [JWT Introduction](https://jwt.io/introduction)
- [Cloudinary Documentation](https://cloudinary.com/documentation)
- [Docker Documentation](https://docs.docker.com/)

---

**Built with ❤️ by the QuillHub Team**
