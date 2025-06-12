package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"smarthome/db"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Set up database connection
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smarthome")
	database, err := db.New(dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer database.Close()

	// Initialize router
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./static")

	// Serve API documentation files
	router.Static("/api-docs", "./api-docs")

	// Swagger UI endpoint (встроенный через gin-swagger)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Custom Swagger UI endpoint (использует наш openapi.yaml)
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/swagger.html")
	})

	// Main documentation page
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "smart-home-api",
			"version":   "1.0",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// API routes (mock endpoints for documentation)
	apiRoutes := router.Group("/api/v1")
	setupMockRoutes(apiRoutes)

	// Start server
	srv := &http.Server{
		Addr:    getEnv("PORT", ":8080"),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on %s\n", srv.Addr)
		log.Printf("📚 Documentation available at: http://localhost%s/\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}

// setupMockRoutes создает mock endpoints для демонстрации документации
func setupMockRoutes(router *gin.RouterGroup) {
	// Mock devices endpoint
	router.GET("/devices", func(c *gin.Context) {
		c.JSON(http.StatusOK, []gin.H{
			{
				"id":           1,
				"name":         "Датчик гостиной",
				"type":         "temperature",
				"location":     "living_room",
				"value":        22.5,
				"unit":         "°C",
				"status":       "active",
				"last_updated": "2024-01-15T10:30:00Z",
			},
			{
				"id":           2,
				"name":         "Отопление гостиной",
				"type":         "heating",
				"location":     "living_room",
				"value":        nil,
				"unit":         "",
				"status":       "active",
				"last_updated": "2024-01-15T10:25:00Z",
			},
		})
	})

	// Mock device command endpoint
	router.POST("/devices/:deviceId/commands", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"command_id": "cmd_abc123",
			"status":     "executed",
			"message":    "Команда выполнена успешно",
			"timestamp":  time.Now().Format(time.RFC3339),
		})
	})

	// Mock scenarios endpoint
	router.POST("/scenarios", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{
			"id":          1,
			"name":        "Автогрев гостиной",
			"description": "Включает отопление когда температура ниже 20°C",
			"enabled":     true,
			"created_at":  time.Now().Format(time.RFC3339),
		})
	})

	// Mock telemetry endpoint
	router.GET("/telemetry/sensors/:sensorId/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"sensor_id":   1,
			"sensor_name": "Датчик гостиной",
			"period": gin.H{
				"from": "2024-01-15T00:00:00Z",
				"to":   "2024-01-15T23:59:59Z",
			},
			"data_points": []gin.H{
				{
					"timestamp": "2024-01-15T10:00:00Z",
					"value":     22.1,
					"status":    "active",
				},
				{
					"timestamp": "2024-01-15T10:30:00Z",
					"value":     22.5,
					"status":    "active",
				},
			},
		})
	})

	log.Println("Mock API endpoints configured for documentation demo")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
