# Application Status Implementation

## Overview

This implementation provides comprehensive application status monitoring and health checks for the Food Marketplace API. The system tracks all major components and provides detailed insights into the application's health and performance.

## Features Implemented

### ✅ 1. **Configuration Management**
- **Environment-based configs** (development/staging/production)
- **YAML configuration files** with environment variable overrides
- **Viper-based configuration** loading and validation
- **JWT secrets** from environment variables
- **Database connection pooling** configuration

### ✅ 2. **Database Connectivity**
- **MongoDB connection** with proper pooling (50 max, 5 min connections)
- **Lazy initialization** pattern to prevent race conditions
- **Health checks** with timeout handling
- **Connection monitoring** and status reporting

### ✅ 3. **API Endpoints Inventory**
- **33 total endpoints** across 15 functional groups
- **REST API structure** with proper HTTP methods
- **Middleware tracking** (4-5 middlewares per endpoint)
- **Route organization** by functional domains

### ✅ 4. **Health Monitoring**
- **Component-level health checks** (Database, JWT, Config, Logger, Server, Routes)
- **System information** (Go version, memory usage, CPU count, goroutines)
- **Request metrics** (total requests, latency, error rate)
- **Uptime tracking** since server start

### ✅ 5. **Structured Logging**
- **Zap logger** with environment-based configuration
- **Request tracking** with correlation IDs
- **Error monitoring** with detailed stack traces
- **Performance metrics** logging

## API Endpoints

### Status & Health Checks
```bash
# Comprehensive application status
GET /status
GET /api/v1/status/

# Simple health check
GET /health
GET /healthz
GET /api/v1/status/health

# Kubernetes-style checks
GET /readyz    # Readiness check
GET /livez     # Liveness check

# System information
GET /system
GET /api/v1/status/system

# API endpoints inventory
GET /endpoints
GET /api/v1/status/endpoints

# Monitoring endpoints
GET /api/v1/monitoring/status
GET /api/v1/monitoring/metrics
GET /api/v1/monitoring/health
```

## Sample Responses

### Full Status Response
```json
{
  "status": "healthy",
  "environment": "development",
  "version": "1.0.0",
  "timestamp": "2025-08-22T21:23:14+05:30",
  "uptime": "3m12.155185139s",
  "components": {
    "database": {
      "status": "healthy",
      "lastChecked": "2025-08-22T21:23:14+05:30",
      "details": {
        "database": "food-marketplace-dev",
        "maxPoolSize": 50,
        "minPoolSize": 5,
        "timeout": "10s",
        "uri": "mongodb://localhost:27017"
      }
    },
    "configuration": {
      "status": "healthy",
      "details": {
        "environment": "development",
        "port": "9000",
        "corsOrigins": 3
      }
    },
    "jwt": {
      "status": "healthy",
      "details": {
        "secretConfigured": true,
        "accessTokenDuration": "24h0m0s",
        "refreshTokenDuration": "168h0m0s"
      }
    }
  },
  "system": {
    "goVersion": "go1.22.2",
    "goroutines": 10,
    "memoryUsage": "4.6 MB",
    "cpuCount": 16
  },
  "metrics": {
    "requestsTotal": 4,
    "requestsPerMinute": 0.067,
    "averageLatency": "75.713µs",
    "errorRate": 100
  },
  "endpoints": {
    "total": 33,
    "byMethod": {
      "GET": 18,
      "POST": 10,
      "PUT": 2,
      "DELETE": 3
    }
  }
}
```

### Health Check Response
```json
{
  "status": "healthy",
  "environment": "development",
  "version": "1.0.0",
  "timestamp": "2025-08-22T21:23:38+05:30",
  "uptime": "3m36.356111495s"
}
```

## Component Status Levels

- **healthy**: Component is functioning normally
- **degraded**: Component has issues but is still operational
- **unhealthy**: Component has critical issues

## Enterprise Features

### 🔧 **Configuration Management**
- Environment-specific configuration files
- Secure secrets management via environment variables
- Validation and default value handling
- Hot-reloadable configurations

### 📊 **Monitoring & Observability**
- Real-time component health monitoring
- Request metrics and performance tracking
- System resource utilization
- Structured logging with correlation IDs

### 🏥 **Health Checks**
- Kubernetes-compatible readiness/liveness probes
- Component-level health validation
- Graceful degradation handling
- Detailed error reporting

### 🛡️ **Security & Reliability**
- No hardcoded credentials
- Environment-specific CORS policies
- Database connection pooling
- Timeout handling and circuit breaker patterns

## Usage Examples

### Basic Health Check
```bash
curl http://localhost:9000/health
```

### Detailed Status
```bash
curl http://localhost:9000/status | jq .
```

### System Metrics
```bash
curl http://localhost:9000/system | jq '.metrics'
```

### API Inventory
```bash
curl http://localhost:9000/endpoints | jq '.total'
```

## Development vs Production

### Development
- Debug logging enabled
- Console log format
- Local MongoDB connection
- CORS allows localhost origins

### Production
- Info level logging
- JSON log format
- Secure MongoDB URI from environment
- Restricted CORS origins
- Enhanced security configurations

## Startup Sequence

1. **Configuration Loading**: YAML files + environment variables
2. **Database Initialization**: MongoDB connection with pooling
3. **Logger Setup**: Environment-based configuration
4. **Status Manager**: Component monitoring initialization
5. **Routes Registration**: All API endpoints including status routes
6. **Server Start**: HTTP server with configured timeouts

The application successfully demonstrates enterprise-level status monitoring with comprehensive health checks, metrics collection, and component tracking across all major subsystems.