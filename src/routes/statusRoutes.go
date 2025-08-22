package routes

import (
	statusController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

// StatusRoutes registers status and health check routes
func StatusRoutes(router *gin.Engine) {
	// Health check routes (commonly used paths)
	router.GET("/health", statusController.GetHealthCheck())
	router.GET("/healthz", statusController.GetHealthCheck())
	// Note: /healthCheck is handled separately in main.go to avoid conflicts
	
	// Kubernetes-style health checks
	router.GET("/readyz", statusController.GetReadinessCheck())
	router.GET("/livez", statusController.GetLivenessCheck())
	
	// Detailed status endpoint
	router.GET("/status", statusController.GetApplicationStatus())
	
	// System information endpoint
	router.GET("/system", statusController.GetSystemInfo())
	
	// API endpoints inventory
	router.GET("/endpoints", statusController.GetEndpoints())
	
	// Status API group
	statusGroup := router.Group("/api/v1/status")
	{
		statusGroup.GET("/", statusController.GetApplicationStatus())
		statusGroup.GET("/health", statusController.GetHealthCheck())
		statusGroup.GET("/system", statusController.GetSystemInfo())
		statusGroup.GET("/endpoints", statusController.GetEndpoints())
		statusGroup.GET("/readiness", statusController.GetReadinessCheck())
		statusGroup.GET("/liveness", statusController.GetLivenessCheck())
	}
	
	// Monitoring endpoints
	monitoringGroup := router.Group("/api/v1/monitoring")
	{
		monitoringGroup.GET("/status", statusController.GetApplicationStatus())
		monitoringGroup.GET("/metrics", statusController.GetSystemInfo())
		monitoringGroup.GET("/health", statusController.GetHealthCheck())
	}
}