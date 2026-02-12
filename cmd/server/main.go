package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/britinogn/quillhub/internal/database"
	"github.com/britinogn/quillhub/internal/handlers"
	"github.com/britinogn/quillhub/internal/repository"
	"github.com/britinogn/quillhub/internal/routes"
	"github.com/britinogn/quillhub/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create root context cancelled on OS signals (SIGINT, SIGTERM, SIGQUIT)
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Connect to PostgreSQL with timeout
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := database.ConnectPostgres(dbCtx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbPool.Close()

	log.Println("✓ Database connected successfully")

	// Initialize Cloudinary client
	cld, err := database.NewCloudinary()
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}
	log.Println("✓ Cloudinary initialized successfully")

	
	
	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)
	postRepo := repository.NewPostRepository(dbPool)
	commentRepo := repository.NewCommentRepository(dbPool)
	dashboardRepo := repository.NewDashboardRepository(dbPool)

	// Get or create AI bot user
	ctx = context.Background()
	botUserID, err := userRepo.GetOrCreateAIBot(ctx)
	if err != nil {
		log.Fatalf("Failed to setup AI bot user: %v", err)
	}

	// Initialize services
	authService := services.NewAuthService(userRepo)
	postService := services.NewPostService(postRepo, cld)
	commentService := services.NewCommentService(commentRepo, postRepo)
	aiService := services.NewAIService()

	// Create auto-poster service
	autoPoster := services.NewAutoPosterService(aiService, postRepo, botUserID)
	
	// Start auto-poster
	autoPoster.Start()
	defer autoPoster.Stop()
	defer aiService.Close() // ✅ Clean up client

	
	// Initialize dashboard service 
	dashboardService := services.NewDashboardService(dashboardRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Configure Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Add CORS middleware with explicit config
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "https://quill-hub-blog.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register all application routes
	routes.RegisterRoutes(router, authHandler, postHandler, commentHandler, dashboardHandler)

	// Determine server port (env or default)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configure HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background goroutine
	go func() {
		log.Printf("🚀 QuillHub API server starting on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("⏳ Shutdown signal received")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced shutdown: %v", err)
	}

	log.Println("✓ Server shutdown complete")
}