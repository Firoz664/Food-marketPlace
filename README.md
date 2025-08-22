# Food Marketplace API

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/foodmngtapp/food-management-apps)
[![API Documentation](https://img.shields.io/badge/docs-available-brightgreen.svg)](#api-documentation)

A comprehensive, enterprise-grade REST API for food marketplace management built with Go, featuring user management, restaurant operations, order processing, and real-time status monitoring.

## 📋 Table of Contents

- [Features](#-features)
- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [API Documentation](#-api-documentation)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Monitoring](#-monitoring)
- [Security](#-security)
- [Contributing](#-contributing)
- [License](#-license)

## 🚀 Features

### Core Features
- **User Management**: Registration, authentication, profile management
- **Restaurant Operations**: Menu management, restaurant profiles, verification
- **Order Processing**: Order creation, tracking, payment integration
- **File Management**: S3 integration for file uploads
- **Real-time Notifications**: OTP verification, email notifications
- **Transaction Management**: Payment processing and tracking

### Enterprise Features
- **Environment-based Configuration**: Development, staging, production configs
- **Health Monitoring**: Comprehensive component health checks
- **Request Metrics**: Performance tracking and monitoring
- **Structured Logging**: JSON logging with correlation IDs
- **Database Connection Pooling**: Optimized MongoDB connections
- **JWT Authentication**: Secure token-based authentication
- **CORS Security**: Environment-specific origin policies
- **API Rate Limiting**: Request throttling and abuse prevention

### Technical Features
- **RESTful API Design**: Following REST principles and best practices
- **Middleware Architecture**: Modular request processing pipeline
- **Error Handling**: Comprehensive error responses with proper HTTP codes
- **Input Validation**: Request validation using Go validator
- **Status Monitoring**: Real-time application and component status
- **Docker Support**: Containerized deployment ready

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- MongoDB 4.4 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/foodmngtapp/food-management-apps.git
   cd food-management-apps
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start MongoDB**
   ```bash
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   
   # Or use your existing MongoDB instance
   ```

5. **Run the application**
   ```bash
   source .env
   go run main.go
   ```

6. **Verify installation**
   ```bash
   curl http://localhost:9000/health
   ```

### Docker Quick Start

```bash
# Build and run with Docker Compose
docker-compose up -d

# Check status
curl http://localhost:9000/health
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file based on `.env.example`:

```env
# Application Environment
APP_ENV=development

# Server Configuration
FOOD_APP_SERVER_PORT=9000
FOOD_APP_FRONTEND_URL=http://localhost:3000
FOOD_APP_ADMIN_URL=http://localhost:3001

# Database Configuration
FOOD_APP_DATABASE_URI=mongodb://localhost:27017
FOOD_APP_DATABASE_NAME=food-marketplace

# JWT Configuration
FOOD_APP_JWT_SECRET_KEY=your-super-secret-jwt-key-change-in-production

# AWS Configuration
FOOD_APP_AWS_REGION=us-west-2
FOOD_APP_AWS_ACCESS_KEY_ID=your-aws-access-key-id
FOOD_APP_AWS_SECRET_ACCESS_KEY=your-aws-secret-access-key
FOOD_APP_AWS_S3_BUCKET=food-marketplace-uploads

# Email Configuration
FOOD_APP_EMAIL_PROVIDER=sendgrid
FOOD_APP_SENDGRID_API_KEY=your-sendgrid-api-key
FOOD_APP_FROM_EMAIL=noreply@foodmarketplace.com

# SMS Configuration
FOOD_APP_SMS_PROVIDER=twilio
FOOD_APP_TWILIO_ACCOUNT_SID=your-twilio-account-sid
FOOD_APP_TWILIO_AUTH_TOKEN=your-twilio-auth-token

# Redis Configuration (Optional)
FOOD_APP_REDIS_HOST=localhost
FOOD_APP_REDIS_PORT=6379
```

### Configuration Files

The application uses YAML configuration files for different environments:

- `configs/config.development.yaml` - Development settings
- `configs/config.staging.yaml` - Staging settings  
- `configs/config.production.yaml` - Production settings

### Database Configuration

```yaml
database:
  uri: "${FOOD_APP_DATABASE_URI}"
  name: "${FOOD_APP_DATABASE_NAME}"
  maxPoolSize: 100
  minPoolSize: 5
  connectionTimeout: "10s"
  serverSelectionTimeout: "5s"
```

## 📚 API Documentation

### Base URL
- **Development**: `http://localhost:9000`
- **Production**: `https://api.foodmarketplace.com`

### Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```bash
Authorization: Bearer <your-jwt-token>
```

### API Endpoints Overview

| Group | Endpoints | Description |
|-------|-----------|-------------|
| **Health** | 6 endpoints | Health checks and status monitoring |
| **Users** | 8 endpoints | User registration, authentication, profile management |
| **Restaurants** | 4 endpoints | Restaurant management and verification |
| **Food** | 2 endpoints | Menu items and food management |
| **Orders** | 2 endpoints | Order processing and tracking |
| **Payments** | 1 endpoint | Payment processing |
| **Files** | 3 endpoints | File upload and management |
| **OTP** | 2 endpoints | OTP generation and verification |
| **Admin** | 1 endpoint | Administrative functions |

### Core API Endpoints

#### Health & Status
```http
GET /health              # Simple health check
GET /status              # Comprehensive application status
GET /system              # System information and metrics
GET /endpoints           # API endpoints inventory
GET /readyz              # Kubernetes readiness probe
GET /livez               # Kubernetes liveness probe
```

#### Authentication
```http
POST /api/v1/users/signup     # User registration
POST /api/v1/users/login      # User login
POST /api/v1/users/resetPassword   # Password reset
POST /api/v1/users/forgetPassword  # Forgot password
```

#### User Management
```http
GET  /api/v1/users/getAllUsers        # Get all users (paginated)
GET  /api/v1/users/getUserDetails     # Get user profile
PUT  /api/v1/users/updateUserDetails  # Update user profile
DELETE /api/v1/users/deleteAccount    # Delete user account
```

#### Restaurant Management
```http
GET    /api/v1/restarunt/getAllRestarunt  # Get all restaurants
GET    /api/v1/restarunt/getRestarunt     # Get restaurant details
PUT    /api/v1/restarunt/updateRestarunt  # Update restaurant
DELETE /api/v1/restarunt/deleteRestarunt  # Delete restaurant
POST   /api/v1/admin/verifyRestaurant     # Verify restaurant (admin)
```

#### Food & Menu
```http
POST /api/v1/food/createFood    # Create food item
GET  /api/v1/food/getAllFood    # Get all food items
GET  /api/v1/menu/getAllMenu    # Get menu items
```

#### File Management
```http
POST   /api/v1/files/singleFileUpload     # Upload single file
POST   /api/v1/files/multipleUploadFiles  # Upload multiple files
DELETE /api/v1/files/deleteFile           # Delete file
```

### Request/Response Examples

#### User Registration
```bash
curl -X POST http://localhost:9000/api/v1/users/signup \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "password": "securepassword123",
    "mobileNumber": "+1234567890",
    "profileType": "USER"
  }'
```

Response:
```json
{
  "data": {
    "status": 200,
    "message": "User signup successful",
    "result": {
      "id": "64f1234567890abcdef123456",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "profileType": "USER",
      "token": "eyJhbGciOiJIUzI1NiIs...",
      "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
    }
  }
}
```

#### Health Check
```bash
curl http://localhost:9000/health
```

Response:
```json
{
  "status": "healthy",
  "environment": "development",
  "version": "1.0.0",
  "timestamp": "2025-08-22T21:23:38+05:30",
  "uptime": "3m36.356111495s"
}
```

### Error Responses

All endpoints return consistent error responses:

```json
{
  "status": false,
  "statusCode": 400,
  "error": "Validation failed",
  "details": "Email is required"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error
- `503` - Service Unavailable

## 🛠 Development

### Project Structure

```
food-marketplace/
├── configs/                 # Configuration files
│   ├── config.development.yaml
│   ├── config.staging.yaml
│   └── config.production.yaml
├── src/
│   ├── common/             # Shared utilities
│   │   ├── functions/      # Helper functions
│   │   ├── helpers/        # Utility helpers
│   │   ├── logger/         # Logging utilities
│   │   ├── middleware/     # HTTP middleware
│   │   ├── rateLimiter/    # Rate limiting
│   │   └── status/         # Status monitoring
│   ├── config/             # Configuration management
│   │   └── database/       # Database setup
│   ├── controllers/        # HTTP handlers
│   ├── models/            # Data models
│   └── routes/            # Route definitions
├── template/              # HTML templates
├── .env.example          # Environment variables template
├── docker-compose.yml    # Docker composition
├── Dockerfile           # Container definition
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── main.go             # Application entry point
└── README.md           # This file
```

### Development Setup

1. **Install development dependencies**
   ```bash
   # Install air for hot reloading (optional)
   go install github.com/cosmtrek/air@latest
   
   # Install golangci-lint for code quality
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. **Run with hot reloading**
   ```bash
   air
   ```

3. **Code formatting and linting**
   ```bash
   # Format code
   gofmt -s -w .
   
   # Run linter
   golangci-lint run
   ```

### Coding Standards

- **Go Code Style**: Follow standard Go conventions
- **Error Handling**: Always handle errors explicitly
- **Logging**: Use structured logging with zap
- **Comments**: Document public functions and complex logic
- **Testing**: Write unit tests for business logic
- **Validation**: Validate all input data

### Git Workflow

1. Create feature branch: `git checkout -b feature/your-feature`
2. Make changes and commit: `git commit -m "feat: add user authentication"`
3. Push branch: `git push origin feature/your-feature`
4. Create pull request
5. Code review and merge

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestUserSignup ./src/controllers
```

### Test Structure

```bash
# Unit tests
src/controllers/user_test.go
src/models/user_test.go

# Integration tests
tests/integration/
├── user_api_test.go
├── restaurant_api_test.go
└── health_check_test.go

# End-to-end tests
tests/e2e/
└── api_flow_test.go
```

### Test Coverage

Maintain minimum 80% test coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🚀 Deployment

### Docker Deployment

1. **Build Docker image**
   ```bash
   docker build -t food-marketplace-api .
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Environment-specific deployment**
   ```bash
   # Staging
   docker-compose -f docker-compose.staging.yml up -d
   
   # Production
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Kubernetes Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: food-marketplace-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: food-marketplace-api
  template:
    metadata:
      labels:
        app: food-marketplace-api
    spec:
      containers:
      - name: api
        image: food-marketplace-api:latest
        ports:
        - containerPort: 9000
        env:
        - name: APP_ENV
          value: "production"
        livenessProbe:
          httpGet:
            path: /livez
            port: 9000
          initialDelaySeconds: 30
        readinessProbe:
          httpGet:
            path: /readyz
            port: 9000
          initialDelaySeconds: 5
```

### Production Checklist

- [ ] Environment variables configured
- [ ] Database connection secured
- [ ] JWT secret keys rotated
- [ ] HTTPS/TLS enabled
- [ ] CORS origins restricted
- [ ] Rate limiting configured
- [ ] Monitoring and alerting setup
- [ ] Backup strategy implemented
- [ ] Load balancer configured
- [ ] Auto-scaling policies set

## 📊 Monitoring

### Health Checks

The application provides multiple health check endpoints:

```bash
# Basic health check
curl http://localhost:9000/health

# Comprehensive status
curl http://localhost:9000/status

# Kubernetes probes
curl http://localhost:9000/readyz  # Readiness
curl http://localhost:9000/livez   # Liveness
```

### Metrics and Monitoring

#### Application Metrics
- Request count and latency
- Error rates and types
- Component health status
- System resource usage

#### Component Monitoring
- **Database**: Connection pool status, query performance
- **JWT**: Token validation success/failure rates
- **File Storage**: Upload success rates, storage usage
- **External APIs**: Response times, error rates

### Logging

Structured JSON logging with correlation IDs:

```json
{
  "level": "info",
  "ts": "2025-08-22T21:23:14.123Z",
  "caller": "controllers/user.go:45",
  "msg": "User created successfully",
  "correlationId": "req-123e4567-e89b-12d3-a456-426614174000",
  "userId": "64f1234567890abcdef123456",
  "email": "user@example.com",
  "duration": "45ms"
}
```

### Alerting

Set up alerts for:
- Application health status changes
- High error rates (>5%)
- High response times (>1s p95)
- Database connection issues
- Memory/CPU usage spikes

## 🔒 Security

### Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt with salt rounds
- **Input Validation**: Comprehensive request validation
- **CORS Protection**: Environment-specific origin policies
- **Rate Limiting**: Request throttling per IP/user
- **Secure Headers**: Standard security headers
- **Environment Isolation**: Separate configs per environment

### Security Best Practices

1. **Never commit secrets** to version control
2. **Use environment variables** for all sensitive data
3. **Implement proper error handling** without exposing internals
4. **Validate all inputs** on the server side
5. **Use HTTPS** in production
6. **Implement proper logging** without sensitive data
7. **Regular security updates** for dependencies

### Security Checklist

- [ ] JWT secrets are strong and rotated
- [ ] Database credentials are secured
- [ ] API keys are not hardcoded
- [ ] CORS is properly configured
- [ ] Rate limiting is enabled
- [ ] Input validation is comprehensive
- [ ] Error messages don't expose sensitive info
- [ ] Dependencies are up to date

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### Getting Started

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

### Development Guidelines

- Follow Go coding standards
- Write clear commit messages
- Add documentation for new features
- Maintain test coverage >80%
- Update README for significant changes

### Pull Request Process

1. Update documentation if needed
2. Add tests for new features
3. Ensure CI/CD pipeline passes
4. Request review from maintainers
5. Address review feedback
6. Merge after approval

### Code Review Criteria

- Code quality and style
- Test coverage and quality
- Documentation completeness
- Performance implications
- Security considerations

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Go community for excellent tooling
- MongoDB team for robust database solutions
- Gin framework for HTTP routing
- All contributors and maintainers

## 📞 Support

### Documentation
- [API Documentation](#api-documentation)
- [Configuration Guide](#configuration)
- [Deployment Guide](#deployment)

### Community
- GitHub Issues: [Report bugs or request features](https://github.com/foodmngtapp/food-management-apps/issues)
- Discussions: [Join the community discussion](https://github.com/foodmngtapp/food-management-apps/discussions)

### Professional Support
For enterprise support and consulting, please contact: support@foodmarketplace.com

---

**Made with ❤️ by the Food Marketplace Team**

For the latest updates and releases, please visit our [GitHub repository](https://github.com/foodmngtapp/food-management-apps).