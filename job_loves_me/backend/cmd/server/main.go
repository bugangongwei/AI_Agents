package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bugangongwei/job_loves_me/backend/internal/api/handler"
	"github.com/bugangongwei/job_loves_me/backend/internal/api/middleware"
	"github.com/bugangongwei/job_loves_me/backend/internal/repository"
	"github.com/bugangongwei/job_loves_me/backend/internal/service"
	"github.com/bugangongwei/job_loves_me/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default environment variables")
	}

	// Initialize DB
	repository.InitDB()

	// Initialize JWT
	utils.InitJWT()

	// Initialize handlers
	authHandler := handler.NewAuthHandler()
	resumeHandler := handler.NewResumeHandler()
	// Get the resume service from the handler to share with greeting handler
	// In a real app, you'd use a better DI container or manual injection
	resumeService := service.NewResumeService()
	greetingHandler := handler.NewGreetingHandler(resumeService)

	// Create Gin router
	r := gin.Default()

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		// Resume routes
		resumes := protected.Group("/resumes")
		{
			resumes.POST("/upload", resumeHandler.UploadResume)
		}

		// Greeting routes
		greetings := protected.Group("/greetings")
		{
			greetings.POST("/generate", greetingHandler.GenerateGreeting)
		}

		// TODO: Interview routes
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
