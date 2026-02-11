# QuillHub 🪶

[![Go Version](https://img.shields.io/badge/Go-1.25.4-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

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
- [Testing](#testing)
- [Security Best Practices](#security-best-practices)
- [Monitoring & Logging](#monitoring--logging)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)
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

### Base URL
```
http://localhost:8080/api
```

### Health Check

```http
GET /api/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2026-02-11T10:30:00Z",
  "version": "1.0.0"
}
```

---

### Authentication Endpoints

#### 1. User Registration (Sign Up)

```http
POST /api/auth/signup
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123",
  "role": "user"
}
```

**Response (201 Created):**
```json
{
  "message": "user registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2026-02-11T10:30:00Z"
  }
}
```

**Error Responses:**
- `409 Conflict` - Email already registered or username taken
- `400 Bad Request` - Invalid request format

---

#### 2. User Login

```http
POST /api/auth/login
Content-Type: application/json
```

**Request Body:**
```json
{
  "identifier": "john@example.com",
  "password": "123456"
}
```

**Note:** `identifier` can be either email or username.

**Response (200 OK):**
```json
{
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "username": "johndoe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2026-02-11T10:30:00Z"
    }
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid credentials
- `400 Bad Request` - Invalid request format

---

#### 3. Admin Registration (Admin Only)

```http
POST /api/auth/admin/register
Authorization: Bearer {ADMIN_JWT_TOKEN}
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Admin User",
  "username": "adminuser",
  "email": "admin@example.com",
  "password": "123456",
  "role": "admin"
}
```

**Response (201 Created):**
```json
{
  "message": "admin user created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Admin User",
    "username": "adminuser",
    "email": "admin@example.com",
    "role": "admin",
    "created_at": "2026-02-11T10:30:00Z"
  }
}
```

---

### Post Endpoints

#### 4. Create Post (Protected)

```http
POST /api/posts
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json
```

**Request Body (JSON):**
```json
{
  "title": "My First Post",
  "content": "This is the content of my first post",
  "tags": ["writing", "thoughts"],
  "category": "Technology"
}
```

**Or with Form Data (for image uploads):**
```http
POST /api/posts
Authorization: Bearer {JWT_TOKEN}
Content-Type: multipart/form-data

title: My First Post
content: This is the content of my first post
tags: writing,thoughts
category: Technology
images: [file1.jpg, file2.png]
```

**Response (201 Created):**
```json
{
  "message": "Post created successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "author_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "My First Post",
    "content": "This is the content of my first post",
    "image_url": ["https://cloudinary.com/image1.jpg"],
    "tags": ["writing", "thoughts"],
    "category": "Technology",
    "is_published": true,
    "view_count": 0,
    "created_at": "2026-02-11T10:30:00Z",
    "updated_at": "2026-02-11T10:30:00Z"
  }
}
```

---

#### 5. Get All Posts (Public)

```http
GET /api/posts?page=1&limit=10
```

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `limit` (optional) - Items per page (default: 10, max: 100)

**Response (200 OK):**
```json
{
  "posts": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "author_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "My First Post",
      "content": "This is the content...",
      "image_url": ["https://cloudinary.com/image1.jpg"],
      "tags": ["writing", "thoughts"],
      "category": "Technology",
      "is_published": true,
      "view_count": 42,
      "created_at": "2026-02-11T10:30:00Z",
      "updated_at": "2026-02-11T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

---

#### 6. Get Post by ID (Public)

```http
GET /api/posts/:id
```

**Response (200 OK):**
```json
{
  "post": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "author_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "My First Post",
    "content": "This is the content of my first post",
    "image_url": ["https://cloudinary.com/image1.jpg"],
    "tags": ["writing", "thoughts"],
    "category": "Technology",
    "is_published": true,
    "view_count": 42,
    "created_at": "2026-02-11T10:30:00Z",
    "updated_at": "2026-02-11T10:30:00Z"
  }
}
```

**Error Responses:**
- `404 Not Found` - Post not found

---

#### 7. Get Posts by Author (Public)

```http
GET /api/posts/author/:authorId
```

**Response (200 OK):**
```json
{
  "posts": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "author_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "My First Post",
      "content": "This is the content...",
      "tags": ["writing"],
      "category": "Technology",
      "created_at": "2026-02-11T10:30:00Z"
    }
  ],
  "count": 5
}
```

---

#### 8. Update Post (Protected - Author Only)

```http
PUT /api/posts/:id
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json
```

**Request Body (all fields optional):**
```json
{
  "title": "Updated Title",
  "content": "Updated content",
  "tags": ["updated", "tags"],
  "category": "Updated Category",
  "is_published": true
}
```

**Or with Form Data:**
```http
PUT /api/posts/:id
Authorization: Bearer {JWT_TOKEN}
Content-Type: multipart/form-data

title: Updated Title
content: Updated content
tags: updated,tags
category: Updated Category
is_published: true
images: [new_image.jpg]
```

**Response (200 OK):**
```json
{
  "message": "Post updated successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "title": "Updated Title",
    "content": "Updated content",
    "updated_at": "2026-02-11T11:00:00Z"
  }
}
```

**Error Responses:**
- `404 Not Found` - Post not found
- `403 Forbidden` - Not authorized to update this post

---

#### 9. Delete Post (Protected - Author Only)

```http
DELETE /api/posts/:id
Authorization: Bearer {JWT_TOKEN}
```

**Response (200 OK):**
```json
{
  "message": "Post deleted successfully"
}
```

**Error Responses:**
- `404 Not Found` - Post not found
- `403 Forbidden` - Not authorized to delete this post

---

### Comment Endpoints

#### 10. Create Comment (Protected)

```http
POST /api/posts/:postId/comments
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json
```

**Request Body:**
```json
{
  "content": "Great post! Very informative."
}
```

**Response (201 Created):**
```json
{
  "message": "Comment created successfully",
  "comment": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "post_id": "660e8400-e29b-41d4-a716-446655440000",
    "author_id": "550e8400-e29b-41d4-a716-446655440000",
    "content": "Great post! Very informative.",
    "created_at": "2026-02-11T10:30:00Z"
  }
}
```

**Error Responses:**
- `404 Not Found` - Post not found
- `401 Unauthorized` - User not authenticated

---

#### 11. Get Comments by Post ID (Public)

```http
GET /api/posts/:id/comments
```

**Response (200 OK):**
```json
{
  "comments": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440000",
      "post_id": "660e8400-e29b-41d4-a716-446655440000",
      "author_id": "550e8400-e29b-41d4-a716-446655440000",
      "author_name": "John Doe",
      "content": "Great post!",
      "created_at": "2026-02-11T10:30:00Z"
    }
  ],
  "count": 5
}
```

---

#### 12. Delete Comment (Protected - Author Only)

```http
DELETE /api/comments/:commentId
Authorization: Bearer {JWT_TOKEN}
```

**Response (200 OK):**
```json
{
  "message": "Comment deleted successfully"
}
```

**Error Responses:**
- `404 Not Found` - Comment not found
- `403 Forbidden` - Can only delete your own comments

---

### Dashboard Endpoints

#### 13. Get User Dashboard (Protected)

```http
GET /api/dashboard
Authorization: Bearer {JWT_TOKEN}
```

**Response (200 OK):**
```json
{
  "message": "User dashboard fetched successfully",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "total_posts": 15,
    "total_comments": 42,
    "total_views": 1250,
    "recent_posts": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440000",
        "title": "My Recent Post",
        "view_count": 100,
        "created_at": "2026-02-11T10:30:00Z"
      }
    ]
  }
}
```

---

#### 14. Get Admin Dashboard (Admin Only)

```http
GET /api/admin/dashboard
Authorization: Bearer {ADMIN_JWT_TOKEN}
```

**Response (200 OK):**
```json
{
  "message": "Admin dashboard fetched successfully",
  "data": {
    "total_users": 1500,
    "total_posts": 5000,
    "total_comments": 12000,
    "total_views": 250000,
    "recent_users": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "johndoe",
        "email": "john@example.com",
        "created_at": "2026-02-11T10:30:00Z"
      }
    ],
    "recent_posts": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440000",
        "title": "Latest Post",
        "author": "johndoe",
        "created_at": "2026-02-11T10:30:00Z"
      }
    ]
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Not authenticated
- `403 Forbidden` - Admin access required

---

### Root Endpoint

```http
GET /api/
```

**Response (200 OK):**
```json
{
  "message": "Welcome to QuillHub API",
  "version": "1.0.0",
  "documentation": "/api/docs"
}
```

## 📁 Project Structure

```
server/
├── cmd/                          # Application entry points
│   ├── app/                      # Additional app commands (if any)
│   └── server/
│       └── main.go               # Main application entry point
│
├── config/                       # Configuration files
│   ├── config.go                 # Configuration loader and parser
│   ├── config.yaml               # YAML configuration file
│   └── .env.example              # Environment variables template
│
├── internal/                     # Private application code
│   ├── database/                 # Database connections
│   │   ├── db.go                 # Database initialization
│   │   ├── postgres.go           # PostgreSQL connection setup
│   │   └── cloudinary.go         # Cloudinary client setup
│   │
│   ├── handlers/                 # HTTP request handlers
│   │   ├── auth_handler.go       # Authentication endpoints (signup, login, admin)
│   │   ├── post_handler.go       # Post CRUD endpoints
│   │   ├── comment_handler.go    # Comment CRUD endpoints
│   │   └── dashboard_handlers.go # Dashboard statistics endpoints
│   │
│   ├── middleware/               # HTTP middleware
│   │   ├── authMiddleware.go     # JWT token verification
│   │   ├── authRoleMiddleware.go # Role-based access control (admin/user)
│   │   └── uploadMiddleware.go   # File upload handling and validation
│   │
│   ├── model/                    # Data models and DTOs
│   │   ├── User.go               # User model and request/response types
│   │   ├── Post.go               # Post model and request/response types
│   │   ├── Comment.go            # Comment model and request/response types
│   │   └── Dashboard.go          # Dashboard statistics models
│   │
│   ├── repository/               # Data access layer
│   │   ├── user_repository.go    # User database operations
│   │   ├── post_repository.go    # Post database operations
│   │   ├── comment_repository.go # Comment database operations
│   │   └── dashboard_repository.go # Dashboard statistics queries
│   │
│   ├── routes/                   # Route definitions
│   │   └── routes.go             # API route registration and grouping
│   │
│   ├── services/                 # Business logic layer
│   │   ├── auth_service.go       # Authentication business logic
│   │   ├── post_service.go       # Post business logic and validation
│   │   ├── comment_service.go    # Comment business logic
│   │   └── dashboard_service.go  # Dashboard data aggregation
│   │
│   └── templates/                # HTML templates
│       └── index.html            # Welcome page template
│
├── pkg/                          # Public reusable packages
│   ├── logger/                   # Logging utilities
│   │   └── logger.go             # Structured logger implementation
│   │
│   ├── response/                 # HTTP response helpers
│   │   └── response.go           # Standard response formatting
│   │
│   └── utils/                    # Utility functions
│       ├── jwt.go                # JWT token generation and validation
│       ├── hashPassword.go       # Password hashing with bcrypt
│       └── response.go           # Additional response utilities
│
├── migrations/                   # Database migrations
│   └── 001_init.sql              # Initial schema (users, posts, comments)
│
├── deployments/                  # Deployment configurations
│   └── (kubernetes, terraform, etc.)
│
├── tmp/                          # Temporary build files (Air)
│   ├── build-errors.log          # Build error logs
│   └── main.exe                  # Compiled binary (development)
│
├── build/                        # Production build output
│   └── quillhub                  # Production binary
│
├── docker-compose.yml            # Docker Compose for local development
├── Dockerfile                    # Production Docker image
├── Dockerfile.dev                # Development Docker image with hot reload
├── .air.toml                     # Air configuration for hot reload
├── .dockerignore                 # Docker ignore patterns
├── .gitignore                    # Git ignore patterns
├── .env                          # Environment variables (not in git)
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
└── README.md                     # Project documentation
```

### Architecture Overview

The project follows a clean architecture pattern with clear separation of concerns:

- **cmd/**: Application entry points and initialization
- **internal/**: Private application code (not importable by other projects)
  - **handlers/**: HTTP layer - handles requests and responses
  - **services/**: Business logic layer - contains core application logic
  - **repository/**: Data access layer - database operations
  - **middleware/**: Cross-cutting concerns (auth, logging, etc.)
  - **model/**: Data structures and DTOs
- **pkg/**: Public reusable packages that could be imported by other projects
- **config/**: Configuration management
- **migrations/**: Database schema versioning



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

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Testing API Endpoints

Use the provided examples with cURL or Postman:

```bash
# Health check
curl http://localhost:8080/api/health

# Sign up
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","username":"testuser","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"test@example.com","password":"password123"}'
```

## 🔒 Security Best Practices

### Environment Variables

- Never commit `.env` files to version control
- Use strong, unique passwords for database and JWT secrets
- Rotate JWT secrets regularly in production
- Use different credentials for development and production

### Password Security

- Passwords are hashed using bcrypt with cost factor 10
- Minimum password length should be enforced (recommended: 8+ characters)
- Consider implementing password complexity requirements

### JWT Token Security

- Tokens expire after configured duration (default: 24 hours)
- Store tokens securely on client side (httpOnly cookies recommended)
- Implement token refresh mechanism for better UX
- Validate tokens on every protected route

### Database Security

- Use SSL/TLS for database connections in production (`DB_SSLMODE=require`)
- Implement prepared statements to prevent SQL injection (already done with pgx)
- Regular database backups
- Principle of least privilege for database users

### API Security

- CORS is configured - update allowed origins for production
- Rate limiting should be implemented for production
- Input validation on all endpoints
- Sanitize user inputs to prevent XSS attacks

## 📊 Monitoring & Logging

### Application Logging

The application uses structured logging with different levels:

```go
// Log levels: DEBUG, INFO, WARN, ERROR
logger.Info("Server started", "port", config.Port)
logger.Error("Database connection failed", "error", err)
```

### Log Files

Logs are written to:
- Console output (development)
- Log files in production (configure in `config.yaml`)

### Monitoring Endpoints

```bash
# Health check
GET /api/health

# Returns:
{
  "status": "healthy",
  "timestamp": "2026-02-11T10:30:00Z",
  "version": "1.0.0"
}
```

### Recommended Monitoring Tools

- **Prometheus**: Metrics collection
- **Grafana**: Metrics visualization
- **ELK Stack**: Log aggregation and analysis
- **Sentry**: Error tracking and monitoring

## 🚀 Deployment

### Deploying to Production

#### Prerequisites

- Server with Docker and Docker Compose installed
- Domain name configured
- SSL certificate (Let's Encrypt recommended)
- Environment variables configured

#### Deployment Steps

1. **Clone repository on server**
```bash
git clone https://github.com/britinogn/quillhub.git
cd quillhub/server
```

2. **Configure production environment**
```bash
cp config/.env.example .env
# Edit .env with production values
```

3. **Build and start services**
```bash
docker-compose -f docker-compose.yml up -d --build
```

4. **Run database migrations**
```bash
docker-compose exec server psql -U postgres -d quill_hub -f migrations/001_init.sql
```

5. **Verify deployment**
```bash
curl https://your-domain.com/api/health
```

### Deploying to Cloud Platforms

#### AWS (Elastic Beanstalk / ECS)

- Use `Dockerfile` for container deployment
- Configure RDS for PostgreSQL
- Use ElastiCache for Redis
- Store secrets in AWS Secrets Manager

#### Google Cloud Platform (Cloud Run)

```bash
gcloud builds submit --tag gcr.io/PROJECT_ID/quillhub
gcloud run deploy quillhub --image gcr.io/PROJECT_ID/quillhub --platform managed
```

#### Heroku

```bash
heroku create quillhub-api
heroku addons:create heroku-postgresql:hobby-dev
heroku addons:create heroku-redis:hobby-dev
git push heroku main
```

### Reverse Proxy Setup (Nginx)

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 🔧 Troubleshooting

### Common Issues

#### Database Connection Failed

**Problem**: `connection refused` or `database does not exist`

**Solutions**:
```bash
# Check if PostgreSQL is running
docker-compose ps

# Check database logs
docker-compose logs db

# Recreate database
docker-compose down -v
docker-compose up -d db
```

#### Port Already in Use

**Problem**: `bind: address already in use`

**Solutions**:
```bash
# Find process using port 8080
netstat -ano | findstr :8080

# Kill the process (Windows)
taskkill /PID <process_id> /F

# Or change port in .env
PORT=8081
```

#### JWT Token Invalid

**Problem**: `invalid token` or `token expired`

**Solutions**:
- Ensure JWT_SECRET matches between token generation and validation
- Check token expiration time
- Verify Authorization header format: `Bearer <token>`

#### Cloudinary Upload Failed

**Problem**: Image upload returns error

**Solutions**:
- Verify Cloudinary credentials in `.env`
- Check file size limits
- Ensure file format is supported (jpg, png, gif, webp)

#### Migration Errors

**Problem**: Database schema mismatch

**Solutions**:
```bash
# Reset database (⚠️ deletes all data)
docker-compose down -v
docker-compose up -d db
docker-compose exec db psql -U postgres -d quill_hub -f /migrations/001_init.sql
```

#### Docker Build Fails

**Problem**: Build errors or dependency issues

**Solutions**:
```bash
# Clear Docker cache
docker-compose build --no-cache

# Remove old images
docker system prune -a

# Rebuild from scratch
docker-compose down
docker-compose up --build
```

### Getting Help

If you encounter issues:

1. Check the logs: `docker-compose logs server`
2. Verify environment variables are set correctly
3. Ensure all prerequisites are installed
4. Check GitHub issues for similar problems
5. Create a new issue with:
   - Error message
   - Steps to reproduce
   - Environment details (OS, Go version, Docker version)

## � Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `ENVIRONMENT` | Environment mode | `development` | No |
| `GIN_MODE` | Gin framework mode | `debug` | No |
| `DB_HOST` | PostgreSQL host | `localhost` | Yes |
| `DB_PORT` | PostgreSQL port | `5432` | Yes |
| `DB_USER` | Database user | `postgres` | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | `quill_hub` | Yes |
| `DB_SSLMODE` | SSL mode | `disable` | No |
| `JWT_SECRET` | JWT signing secret | - | Yes |
| `JWT_EXPIRATION` | Token expiration | `24h` | No |
| `CLOUDINARY_CLOUD_NAME` | Cloudinary cloud name | - | Yes |
| `CLOUDINARY_API_KEY` | Cloudinary API key | - | Yes |
| `CLOUDINARY_API_SECRET` | Cloudinary API secret | - | Yes |
| `REDIS_HOST` | Redis host | `localhost` | No |
| `REDIS_PORT` | Redis port | `6379` | No |

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
