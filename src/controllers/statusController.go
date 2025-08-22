package controllers

import (
	"net/http"
	"time"

	"github.com/foodmngtapp/food-management-apps/src/common/status"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetApplicationStatus returns comprehensive application status
func GetApplicationStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusManager := status.GetStatusManager()
		if statusManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Status manager not initialized",
			})
			return
		}

		appStatus := statusManager.GetApplicationStatus()
		
		// Set appropriate HTTP status code based on application status
		httpStatus := http.StatusOK
		if appStatus.Status == "unhealthy" {
			httpStatus = http.StatusServiceUnavailable
		} else if appStatus.Status == "degraded" {
			httpStatus = http.StatusPartialContent
		}

		c.JSON(httpStatus, appStatus)
	}
}

// GetHealthCheck returns a simplified health check
func GetHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusManager := status.GetStatusManager()
		if statusManager == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"error":     "Status manager not initialized",
				"timestamp": time.Now().Format(time.RFC3339),
			})
			return
		}

		appStatus := statusManager.GetApplicationStatus()
		
		healthResponse := gin.H{
			"status":      appStatus.Status,
			"timestamp":   appStatus.Timestamp,
			"environment": appStatus.Environment,
			"version":     appStatus.Version,
			"uptime":      appStatus.Uptime,
		}

		// Add component summary
		components := gin.H{}
		if appStatus.Components.Database.Status != "healthy" {
			components["database"] = appStatus.Components.Database.Status
		}
		if appStatus.Components.Configuration.Status != "healthy" {
			components["configuration"] = appStatus.Components.Configuration.Status
		}
		if appStatus.Components.JWT.Status != "healthy" {
			components["jwt"] = appStatus.Components.JWT.Status
		}
		if len(components) > 0 {
			healthResponse["issues"] = components
		}

		// Set appropriate HTTP status code
		httpStatus := http.StatusOK
		if appStatus.Status == "unhealthy" {
			httpStatus = http.StatusServiceUnavailable
		} else if appStatus.Status == "degraded" {
			httpStatus = http.StatusPartialContent
		}

		c.JSON(httpStatus, healthResponse)
	}
}

// GetSystemInfo returns system information
func GetSystemInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusManager := status.GetStatusManager()
		if statusManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Status manager not initialized",
			})
			return
		}

		appStatus := statusManager.GetApplicationStatus()
		
		response := gin.H{
			"system":      appStatus.System,
			"metrics":     appStatus.Metrics,
			"environment": appStatus.Environment,
			"version":     appStatus.Version,
			"uptime":      appStatus.Uptime,
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetEndpoints returns API endpoint inventory
func GetEndpoints() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusManager := status.GetStatusManager()
		if statusManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Status manager not initialized",
			})
			return
		}

		appStatus := statusManager.GetApplicationStatus()
		
		c.JSON(http.StatusOK, appStatus.Endpoints)
	}
}

// GetReadinessCheck returns readiness check (for Kubernetes)
func GetReadinessCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusManager := status.GetStatusManager()
		if statusManager == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"ready": false,
				"error": "Status manager not initialized",
			})
			return
		}

		appStatus := statusManager.GetApplicationStatus()
		
		// Application is ready if core components are healthy
		isReady := appStatus.Components.Database.Status == "healthy" &&
			appStatus.Components.Configuration.Status == "healthy" &&
			appStatus.Components.Server.Status == "healthy"

		response := gin.H{
			"ready":       isReady,
			"timestamp":   time.Now().Format(time.RFC3339),
			"environment": appStatus.Environment,
		}

		if !isReady {
			// Add details about what's not ready
			issues := gin.H{}
			if appStatus.Components.Database.Status != "healthy" {
				issues["database"] = appStatus.Components.Database.Status
			}
			if appStatus.Components.Configuration.Status != "healthy" {
				issues["configuration"] = appStatus.Components.Configuration.Status
			}
			if appStatus.Components.Server.Status != "healthy" {
				issues["server"] = appStatus.Components.Server.Status
			}
			response["issues"] = issues
		}

		httpStatus := http.StatusOK
		if !isReady {
			httpStatus = http.StatusServiceUnavailable
		}

		c.JSON(httpStatus, response)
	}
}

// GetLivenessCheck returns liveness check (for Kubernetes)
func GetLivenessCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple liveness check - if we can respond, we're alive
		c.JSON(http.StatusOK, gin.H{
			"alive":     true,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
}

// RequestMetricsMiddleware tracks request metrics
func RequestMetricsMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate metrics
		latency := time.Since(start)
		isError := c.Writer.Status() >= 400

		// Track metrics
		statusManager := status.GetStatusManager()
		if statusManager != nil {
			statusManager.TrackRequest(latency, isError)
		}

		// Log request details
		if logger != nil {
			fields := []zap.Field{
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.Int("status", c.Writer.Status()),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
				zap.Bool("error", isError),
			}

			if isError {
				logger.Warn("HTTP Request Error", fields...)
			} else {
				logger.Info("HTTP Request", fields...)
			}
		}
	})
}