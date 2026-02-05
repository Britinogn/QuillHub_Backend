package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"os/signal"
	"syscall"

	"github.com/britinogn/quillhub/internal/database"
	"github.com/britinogn/quillhub/internal/handlers"
	"github.com/britinogn/quillhub/internal/repository"
	"github.com/britinogn/quillhub/internal/routes"
	"github.com/britinogn/quillhub/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	//Connect to PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := database.ConnectPostgres(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbPool.Close()

	// Create repository
	userRepo := repository.NewUserRepository(dbPool)

	// Create service
	authService := services.NewAuthService(userRepo)

	// Create handler
	authHandler := handlers.NewAuthHandler(authService)

	//Set up Gin router
	router := gin.Default() // or gin.New() if you want full control


	// Add CORS middleware (very important for frontend)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "https://your-frontend-domain.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register all routes 
	routes.RegisterRoutes(router, authHandler)

	// Graceful shutdown
	srv := &http.Server{
		Addr:   os.Getenv("PORT"),   // or ":8080"
		Handler: router,
	}

	// Start server in a goroutine so we can handle shutdown
	go func() {
		log.Printf("Starting QuillHub API on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C or kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give up to 5 seconds to finish ongoing requests
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}


	log.Println("Postgres connected successfully with pgx")
}