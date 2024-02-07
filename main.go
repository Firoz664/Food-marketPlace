package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/foodmngtapp/food-management-apps/src/routes"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// logger
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	// Initialize Zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	router := gin.New()
	router.Use(gin.Recovery())

	// Set up CORS middleware to allow all origins
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Custom logging middleware using Zap
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// Log using Zap
		logger.Info("Request",
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3007"
	}

	// Health check
	router.GET("/healthCheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server running successfully",
		})
	})
	// Catch-all route for any request not matched by the above routes
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"status": 404, "message": "Not found❗️🤷‍♂️ 404  Check the URL of the page."})
	})

	// importing all routes here
	routes.RegisterRoutes(router)

	log.Printf("Server is running on port http://localhost:%s", port)
	log.Fatal(router.Run(":" + port))
}
