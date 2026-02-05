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
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid req"})
		return
	}

    // Convert to User model
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
        Role: req.Role,
	}

	ctx := c.Request.Context()
	err := h.authService.Register(ctx, &req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyRegistered) {
			c.JSON(409, gin.H{"error": "email already registered"}) // ← friendly message
			return
		}
		// Other errors → generic 500
		c.JSON(500, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
        "message": "user registered successfully",
        "data": model.UserResponse{
            ID: user.ID,
            Name:user.Username,
            Email: user.Email,
            Role: user.Role,
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

	// Here: generate JWT, set cookie, etc.
	// For now just return user (without password)
	c.JSON(http.StatusOK, gin.H{
        "message": "login successful",
        "data": model.LoginResponse{
            Token: token,
            User: model.UserResponse{
                ID: user.ID,
                Name:user.Name,
                Username: user.Username,
                Email: user.Email,
                Role: user.Role,
            },
        },
    })
}