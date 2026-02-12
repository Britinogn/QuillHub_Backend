package handlers

import (
	"errors"
	"net/http"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/britinogn/quillhub/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct{
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler{
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req model.CreateUserRequest  // ← Use CreateUserRequest, not model.User
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }


    // Convert CreateUserRequest to User model
    user := &model.User{
        Name:     req.Name,
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
        // Role:     "user",
    }

    // Set role: use provided role or default to "user"
    if req.Role != nil && *req.Role != "" {
        user.Role = *req.Role
    } else {
        user.Role = "user"
    }

    ctx := c.Request.Context()
    err := h.authService.Register(ctx, user)
    if err != nil {
        if errors.Is(err, services.ErrEmailAlreadyRegistered) {
            c.JSON(409, gin.H{"error": "email already registered"})
            return
        }
        if errors.Is(err, services.ErrUsernameTaken) {
            c.JSON(409, gin.H{"error": "username already taken"})
            return
        }
        c.JSON(500, gin.H{
            "error": "failed to create user",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "user registered successfully",
        "data": model.UserResponse{
            ID:        user.ID.String(),
            Name:      user.Name,
            Username:  user.Username,
            Email:     user.Email,
            Role:      user.Role,
            CreatedAt: user.CreatedAt,
        },
    })
}


func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid req"})
		return
	}

	user, token ,err := h.authService.Login(c.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(500, gin.H{"error": "something went wrong"})
		return
	}

	// Just return user (without password)
	c.JSON(http.StatusOK, gin.H{
        "message": "login successful",
        "data": model.LoginResponse{
            Token: token,
            User: model.UserResponse{
				ID : user.ID.String(),
                Name:user.Name,
                Username: user.Username,
                Email: user.Email,
                Role: user.Role,
                CreatedAt: user.CreatedAt,
            },
        },
    })
}

func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
    // Get the requesting user's role from JWT token
    requestingUserRole, exists := c.Get("userRole")
    if !exists {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    var req model.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid input", "details": err.Error()})
        return
    }

    user := &model.User{
        Name:     req.Name,
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }

    ctx := c.Request.Context()
    err := h.authService.RegisterAdmin(ctx, user, requestingUserRole.(string))
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "admin user created successfully",
        "data": model.UserResponse{
            ID:        user.ID.String(),
            Name:      user.Name,
            Username:  user.Username,
            Email:     user.Email,
            Role:      user.Role,
            CreatedAt: user.CreatedAt,
        },
    })
}