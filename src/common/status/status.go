package status

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/foodmngtapp/food-management-apps/src/config"
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"go.uber.org/zap"
)

// ApplicationStatus represents the overall application status
type ApplicationStatus struct {
	Status      string              `json:"status"`
	Environment string              `json:"environment"`
	Version     string              `json:"version"`
	Timestamp   string              `json:"timestamp"`
	Uptime      string              `json:"uptime"`
	Components  ComponentStatuses   `json:"components"`
	System      SystemInfo          `json:"system"`
	Metrics     ApplicationMetrics  `json:"metrics"`
	Endpoints   EndpointInventory   `json:"endpoints"`
}

// ComponentStatuses represents the status of various application components
type ComponentStatuses struct {
	Database      ComponentStatus `json:"database"`
	Configuration ComponentStatus `json:"configuration"`
	JWT           ComponentStatus `json:"jwt"`
	Logger        ComponentStatus `json:"logger"`
	Server        ComponentStatus `json:"server"`
	Routes        ComponentStatus `json:"routes"`
}

// ComponentStatus represents the status of a single component
type ComponentStatus struct {
	Status      string                 `json:"status"`      // healthy, degraded, unhealthy
	LastChecked string                 `json:"lastChecked"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// SystemInfo contains system-level information
type SystemInfo struct {
	GoVersion     string `json:"goVersion"`
	Goroutines    int    `json:"goroutines"`
	MemoryUsage   string `json:"memoryUsage"`
	CPUCount      int    `json:"cpuCount"`
	Architecture  string `json:"architecture"`
	OS            string `json:"os"`
}

// ApplicationMetrics contains application metrics
type ApplicationMetrics struct {
	RequestsTotal     int64   `json:"requestsTotal"`
	RequestsPerMinute float64 `json:"requestsPerMinute"`
	AverageLatency    string  `json:"averageLatency"`
	ErrorRate         float64 `json:"errorRate"`
	ActiveConnections int     `json:"activeConnections"`
}

// EndpointInventory contains information about registered API endpoints
type EndpointInventory struct {
	Total     int                    `json:"total"`
	Endpoints []EndpointInfo         `json:"endpoints"`
	ByMethod  map[string]int         `json:"byMethod"`
	ByGroup   map[string][]string    `json:"byGroup"`
}

// EndpointInfo represents information about a single endpoint
type EndpointInfo struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Handler     string `json:"handler"`
	Middlewares int    `json:"middlewares"`
}

// StatusManager manages application status and health checks
type StatusManager struct {
	startTime      time.Time
	mutex          sync.RWMutex
	logger         *zap.Logger
	config         *config.Config
	lastStatus     *ApplicationStatus
	requestMetrics *RequestMetrics
}

// RequestMetrics tracks request-related metrics
type RequestMetrics struct {
	mutex           sync.RWMutex
	totalRequests   int64
	totalLatency    time.Duration
	errorCount      int64
	requestsPerMin  []int64
	minuteIndex     int
	lastMinute      time.Time
}

var (
	globalStatusManager *StatusManager
	once                sync.Once
)

// InitializeStatusManager initializes the global status manager
func InitializeStatusManager(cfg *config.Config, logger *zap.Logger) *StatusManager {
	once.Do(func() {
		globalStatusManager = &StatusManager{
			startTime: time.Now(),
			logger:    logger,
			config:    cfg,
			requestMetrics: &RequestMetrics{
				requestsPerMin: make([]int64, 60), // Track last 60 minutes
				lastMinute:     time.Now(),
			},
		}
	})
	return globalStatusManager
}

// GetStatusManager returns the global status manager instance
func GetStatusManager() *StatusManager {
	return globalStatusManager
}

// GetApplicationStatus returns the current application status
func (sm *StatusManager) GetApplicationStatus() *ApplicationStatus {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	now := time.Now()
	uptime := time.Since(sm.startTime)

	status := &ApplicationStatus{
		Status:      sm.determineOverallStatus(),
		Environment: sm.config.Server.Environment,
		Version:     "1.0.0",
		Timestamp:   now.Format(time.RFC3339),
		Uptime:      uptime.String(),
		Components:  sm.checkComponents(),
		System:      sm.getSystemInfo(),
		Metrics:     sm.getMetrics(),
		Endpoints:   sm.getEndpointInventory(),
	}

	sm.lastStatus = status
	return status
}

// determineOverallStatus determines the overall application status
func (sm *StatusManager) determineOverallStatus() string {
	components := sm.checkComponents()
	
	unhealthyCount := 0
	degradedCount := 0

	if components.Database.Status != "healthy" {
		if components.Database.Status == "unhealthy" {
			unhealthyCount++
		} else {
			degradedCount++
		}
	}
	if components.Configuration.Status != "healthy" {
		if components.Configuration.Status == "unhealthy" {
			unhealthyCount++
		} else {
			degradedCount++
		}
	}
	if components.JWT.Status != "healthy" {
		if components.JWT.Status == "unhealthy" {
			unhealthyCount++
		} else {
			degradedCount++
		}
	}
	if components.Logger.Status != "healthy" {
		if components.Logger.Status == "unhealthy" {
			unhealthyCount++
		} else {
			degradedCount++
		}
	}
	if components.Server.Status != "healthy" {
		if components.Server.Status == "unhealthy" {
			unhealthyCount++
		} else {
			degradedCount++
		}
	}
	if components.Routes.Status != "healthy" {
		if components.Routes.Status == "unhealthy" {
			unhealthyCount++
		} else {
			degradedCount++
		}
	}

	if unhealthyCount > 0 {
		return "unhealthy"
	} else if degradedCount > 0 {
		return "degraded"
	}
	return "healthy"
}

// checkComponents checks the status of all application components
func (sm *StatusManager) checkComponents() ComponentStatuses {
	now := time.Now().Format(time.RFC3339)
	
	return ComponentStatuses{
		Database:      sm.checkDatabase(now),
		Configuration: sm.checkConfiguration(now),
		JWT:           sm.checkJWT(now),
		Logger:        sm.checkLogger(now),
		Server:        sm.checkServer(now),
		Routes:        sm.checkRoutes(now),
	}
}

// checkDatabase checks the database connection status
func (sm *StatusManager) checkDatabase(timestamp string) ComponentStatus {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.HealthCheck(ctx)
	if err != nil {
		return ComponentStatus{
			Status:      "unhealthy",
			LastChecked: timestamp,
			Error:       err.Error(),
		}
	}

	details := map[string]interface{}{
		"uri":         config.GetDatabaseURI(),
		"database":    sm.config.Database.Name,
		"maxPoolSize": sm.config.Database.MaxPoolSize,
		"minPoolSize": sm.config.Database.MinPoolSize,
		"timeout":     sm.config.Database.ConnectionTimeout.String(),
	}

	return ComponentStatus{
		Status:      "healthy",
		LastChecked: timestamp,
		Details:     details,
	}
}

// checkConfiguration checks the configuration status
func (sm *StatusManager) checkConfiguration(timestamp string) ComponentStatus {
	if sm.config == nil {
		return ComponentStatus{
			Status:      "unhealthy",
			LastChecked: timestamp,
			Error:       "configuration not loaded",
		}
	}

	details := map[string]interface{}{
		"environment":    sm.config.Server.Environment,
		"port":          sm.config.Server.Port,
		"database":      sm.config.Database.Name,
		"logLevel":      sm.config.Logging.Level,
		"corsOrigins":   len(sm.config.Server.CORS.AllowedOrigins),
	}

	return ComponentStatus{
		Status:      "healthy",
		LastChecked: timestamp,
		Details:     details,
	}
}

// checkJWT checks the JWT configuration status
func (sm *StatusManager) checkJWT(timestamp string) ComponentStatus {
	secret := config.GetJWTSecret()
	if secret == "" {
		return ComponentStatus{
			Status:      "unhealthy",
			LastChecked: timestamp,
			Error:       "JWT secret not configured",
		}
	}

	details := map[string]interface{}{
		"issuer":               sm.config.JWT.Issuer,
		"audience":             sm.config.JWT.Audience,
		"accessTokenDuration":  sm.config.JWT.AccessTokenDuration.String(),
		"refreshTokenDuration": sm.config.JWT.RefreshTokenDuration.String(),
		"secretConfigured":     true,
	}

	return ComponentStatus{
		Status:      "healthy",
		LastChecked: timestamp,
		Details:     details,
	}
}

// checkLogger checks the logger status
func (sm *StatusManager) checkLogger(timestamp string) ComponentStatus {
	if sm.logger == nil {
		return ComponentStatus{
			Status:      "unhealthy",
			LastChecked: timestamp,
			Error:       "logger not initialized",
		}
	}

	details := map[string]interface{}{
		"level":      sm.config.Logging.Level,
		"format":     sm.config.Logging.Format,
		"outputPath": sm.config.Logging.OutputPath,
	}

	return ComponentStatus{
		Status:      "healthy",
		LastChecked: timestamp,
		Details:     details,
	}
}

// checkServer checks the server status
func (sm *StatusManager) checkServer(timestamp string) ComponentStatus {
	details := map[string]interface{}{
		"host":         sm.config.Server.Host,
		"port":         sm.config.Server.Port,
		"readTimeout":  sm.config.Server.ReadTimeout.String(),
		"writeTimeout": sm.config.Server.WriteTimeout.String(),
		"idleTimeout":  sm.config.Server.IdleTimeout.String(),
		"environment":  sm.config.Server.Environment,
	}

	return ComponentStatus{
		Status:      "healthy",
		LastChecked: timestamp,
		Details:     details,
	}
}

// checkRoutes checks the routes registration status
func (sm *StatusManager) checkRoutes(timestamp string) ComponentStatus {
	endpoints := sm.getEndpointInventory()
	
	details := map[string]interface{}{
		"totalEndpoints": endpoints.Total,
		"byMethod":      endpoints.ByMethod,
		"groups":        len(endpoints.ByGroup),
	}

	status := "healthy"
	if endpoints.Total == 0 {
		status = "unhealthy"
	} else if endpoints.Total < 20 {
		status = "degraded"
	}

	return ComponentStatus{
		Status:      status,
		LastChecked: timestamp,
		Details:     details,
	}
}

// getSystemInfo returns system information
func (sm *StatusManager) getSystemInfo() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemInfo{
		GoVersion:    runtime.Version(),
		Goroutines:   runtime.NumGoroutine(),
		MemoryUsage:  formatBytes(m.Alloc),
		CPUCount:     runtime.NumCPU(),
		Architecture: runtime.GOARCH,
		OS:           runtime.GOOS,
	}
}

// getMetrics returns application metrics
func (sm *StatusManager) getMetrics() ApplicationMetrics {
	sm.requestMetrics.mutex.RLock()
	defer sm.requestMetrics.mutex.RUnlock()

	var avgLatency time.Duration
	if sm.requestMetrics.totalRequests > 0 {
		avgLatency = sm.requestMetrics.totalLatency / time.Duration(sm.requestMetrics.totalRequests)
	}

	var errorRate float64
	if sm.requestMetrics.totalRequests > 0 {
		errorRate = float64(sm.requestMetrics.errorCount) / float64(sm.requestMetrics.totalRequests) * 100
	}

	// Calculate requests per minute
	var totalRequestsLastHour int64
	for _, count := range sm.requestMetrics.requestsPerMin {
		totalRequestsLastHour += count
	}
	requestsPerMinute := float64(totalRequestsLastHour) / 60.0

	return ApplicationMetrics{
		RequestsTotal:     sm.requestMetrics.totalRequests,
		RequestsPerMinute: requestsPerMinute,
		AverageLatency:    avgLatency.String(),
		ErrorRate:         errorRate,
		ActiveConnections: 0, // Would need to track this separately
	}
}

// getEndpointInventory returns information about registered endpoints
func (sm *StatusManager) getEndpointInventory() EndpointInventory {
	// This is a simplified version - in a real implementation,
	// you would track registered routes dynamically
	endpoints := []EndpointInfo{
		{Method: "GET", Path: "/", Handler: "main.main.func2", Middlewares: 4},
		{Method: "GET", Path: "/healthCheck", Handler: "main.main.func3", Middlewares: 4},
		{Method: "POST", Path: "/api/v1/admin/verifyRestaurant", Handler: "AdminRoutes.VerifyRestaurant", Middlewares: 4},
		{Method: "POST", Path: "/api/v1/files/singleFileUpload", Handler: "FileUploadRoutes.S3FileUploadFile", Middlewares: 4},
		{Method: "DELETE", Path: "/api/v1/files/deleteFile", Handler: "FileUploadRoutes.DeleteFileFromS3Bucket", Middlewares: 5},
		{Method: "POST", Path: "/api/v1/files/multipleUploadFiles", Handler: "FileUploadRoutes.MultipleUploadFiles", Middlewares: 4},
		{Method: "POST", Path: "/api/v1/food/createFood", Handler: "FoodRoutes.CreateFood", Middlewares: 5},
		{Method: "GET", Path: "/api/v1/food/getAllFood", Handler: "FoodRoutes.GetFoods", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/invoice/getAllInvoice", Handler: "InvoiceRoutes.GetInvoice", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/menu/getAllMenu", Handler: "MenuRoutes.GetMenu", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/notes/getAllNotes", Handler: "NotesRoutes.GetNotes", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/notification/getAllNotification", Handler: "NotificationsRoutes.GetNotification", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/OrderItem/getOrderItem", Handler: "OrderItemRoutes.GetOrderItems", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/order/getAllOrder", Handler: "OrderRoutes.GetOrders", Middlewares: 4},
		{Method: "POST", Path: "/api/v1/otp/createAndSendOtp", Handler: "OtpRoutes.CreateOTP", Middlewares: 4},
		{Method: "POST", Path: "/api/v1/otp/verifyOtp", Handler: "OtpRoutes.VerifyOtp", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/payment/getAllPayment", Handler: "PaymentRoutes.GetPayment", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/restarunt/getAllRestarunt", Handler: "RestaruntRoutes.GetAllRestaurants", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/restarunt/getRestarunt", Handler: "RestaruntRoutes.GetRestaurant", Middlewares: 4},
		{Method: "PUT", Path: "/api/v1/restarunt/updateRestarunt", Handler: "RestaruntRoutes.UpdateRestaurant", Middlewares: 5},
		{Method: "DELETE", Path: "/api/v1/restarunt/deleteRestarunt", Handler: "RestaruntRoutes.DeleteRestaurant", Middlewares: 5},
		{Method: "GET", Path: "/api/v1/tables/GetAllRider", Handler: "RiderRoutes.GetRider", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/tables/GetAllTables", Handler: "TableRoutes.GetTables", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/transaction/GetAllTransaction", Handler: "TransactionRoutes.GetAllTransaction", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/transaction/getTransaction", Handler: "TransactionRoutes.GetByIdTransaction", Middlewares: 4},
		{Method: "POST", Path: "/api/v1/users/signup", Handler: "UserRoutes.Signup", Middlewares: 4},
		{Method: "POST", Path: "/api/v1/users/login", Handler: "UserRoutes.Login", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/users/getAllUsers", Handler: "UserRoutes.GetAllUser", Middlewares: 4},
		{Method: "GET", Path: "/api/v1/users/getUserDetails", Handler: "UserRoutes.GetUserDetails", Middlewares: 5},
		{Method: "PUT", Path: "/api/v1/users/updateUserDetails", Handler: "UserRoutes.UpdateUserDetails", Middlewares: 5},
		{Method: "POST", Path: "/api/v1/users/resetPassword", Handler: "UserRoutes.ResetPassword", Middlewares: 5},
		{Method: "POST", Path: "/api/v1/users/forgetPassword", Handler: "UserRoutes.ForgetPassword", Middlewares: 4},
		{Method: "DELETE", Path: "/api/v1/users/deleteAccount", Handler: "UserRoutes.DeleteAccount", Middlewares: 5},
	}

	byMethod := make(map[string]int)
	byGroup := make(map[string][]string)

	for _, endpoint := range endpoints {
		byMethod[endpoint.Method]++
		
		// Extract group from path (e.g., "/api/v1/users/signup" -> "users")
		if len(endpoint.Path) > 8 && endpoint.Path[:8] == "/api/v1/" {
			parts := strings.Split(endpoint.Path[8:], "/")
			if len(parts) > 0 {
				group := parts[0]
				if byGroup[group] == nil {
					byGroup[group] = make([]string, 0)
				}
				byGroup[group] = append(byGroup[group], endpoint.Path)
			}
		}
	}

	return EndpointInventory{
		Total:     len(endpoints),
		Endpoints: endpoints,
		ByMethod:  byMethod,
		ByGroup:   byGroup,
	}
}

// TrackRequest tracks a request for metrics
func (sm *StatusManager) TrackRequest(latency time.Duration, isError bool) {
	if sm.requestMetrics == nil {
		return
	}

	sm.requestMetrics.mutex.Lock()
	defer sm.requestMetrics.mutex.Unlock()

	sm.requestMetrics.totalRequests++
	sm.requestMetrics.totalLatency += latency

	if isError {
		sm.requestMetrics.errorCount++
	}

	// Update per-minute tracking
	now := time.Now()
	if now.Sub(sm.requestMetrics.lastMinute) >= time.Minute {
		sm.requestMetrics.minuteIndex = (sm.requestMetrics.minuteIndex + 1) % 60
		sm.requestMetrics.requestsPerMin[sm.requestMetrics.minuteIndex] = 0
		sm.requestMetrics.lastMinute = now
	}
	sm.requestMetrics.requestsPerMin[sm.requestMetrics.minuteIndex]++
}

// formatBytes formats bytes to human readable format
func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}