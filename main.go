package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/foodmngtapp/food-management-apps/src/common/status"
	statusController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/foodmngtapp/food-management-apps/src/config"
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/foodmngtapp/food-management-apps/src/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database with configuration
	_, err = database.InitializeDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Zap logger based on config
	var logger *zap.Logger
	if cfg.Server.Environment == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	// Initialize status manager
	_ = status.InitializeStatusManager(cfg, logger)
	logger.Info("Status manager initialized", zap.String("version", "1.0.0"))

	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.Server.Environment == "staging" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Set up CORS middleware with configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CORS.AllowedOrigins,
		AllowMethods:     cfg.Server.CORS.AllowedMethods,
		AllowHeaders:     cfg.Server.CORS.AllowedHeaders,
		ExposeHeaders:    cfg.Server.CORS.ExposedHeaders,
		AllowCredentials: cfg.Server.CORS.AllowCredentials,
		MaxAge:           cfg.Server.CORS.MaxAge,
	}))

	// Use request metrics middleware for tracking
	router.Use(statusController.RequestMetricsMiddleware(logger))

	port := cfg.Server.Port

	// router.LoadHTMLGlob("template/*")

	// Serve the HTML file at the default route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Food Marketplace API v1.0 - Environment: %s", cfg.Server.Environment),
			"version": "1.0.0",
			"environment": cfg.Server.Environment,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Register status routes first (before other routes)
	routes.StatusRoutes(router)

	// Catch-all route for any request not matched by the above routes
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":    404, 
			"message":   "Endpoint not found. Please check the URL and try again.",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// importing all routes here
	routes.RegisterRoutes(router)

	// Create HTTP server with configured timeouts
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	logger.Info("Starting Food Marketplace API Server",
		zap.String("environment", cfg.Server.Environment),
		zap.String("host", cfg.Server.Host),
		zap.String("port", port),
		zap.String("database", cfg.Database.Name),
		zap.Duration("readTimeout", cfg.Server.ReadTimeout),
		zap.Duration("writeTimeout", cfg.Server.WriteTimeout),
	)

	log.Printf("🚀 Food Marketplace API Server starting on %s:%s (Environment: %s)", 
		cfg.Server.Host, port, cfg.Server.Environment)
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
