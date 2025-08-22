package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	JWT      JWTConfig      `yaml:"jwt" mapstructure:"jwt"`
	AWS      AWSConfig      `yaml:"aws" mapstructure:"aws"`
	Email    EmailConfig    `yaml:"email" mapstructure:"email"`
	SMS      SMSConfig      `yaml:"sms" mapstructure:"sms"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"redis"`
	Logging  LoggingConfig  `yaml:"logging" mapstructure:"logging"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host         string        `yaml:"host" mapstructure:"host"`
	Port         string        `yaml:"port" mapstructure:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout" mapstructure:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout" mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `yaml:"idleTimeout" mapstructure:"idleTimeout"`
	Environment  string        `yaml:"environment" mapstructure:"environment"`
	CORS         CORSConfig    `yaml:"cors" mapstructure:"cors"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins     []string      `yaml:"allowedOrigins" mapstructure:"allowedOrigins"`
	AllowedMethods     []string      `yaml:"allowedMethods" mapstructure:"allowedMethods"`
	AllowedHeaders     []string      `yaml:"allowedHeaders" mapstructure:"allowedHeaders"`
	ExposedHeaders     []string      `yaml:"exposedHeaders" mapstructure:"exposedHeaders"`
	AllowCredentials   bool          `yaml:"allowCredentials" mapstructure:"allowCredentials"`
	MaxAge             time.Duration `yaml:"maxAge" mapstructure:"maxAge"`
	OptionsPassthrough bool          `yaml:"optionsPassthrough" mapstructure:"optionsPassthrough"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URI                    string        `yaml:"uri" mapstructure:"uri"`
	Name                   string        `yaml:"name" mapstructure:"name"`
	MaxPoolSize            uint64        `yaml:"maxPoolSize" mapstructure:"maxPoolSize"`
	MinPoolSize            uint64        `yaml:"minPoolSize" mapstructure:"minPoolSize"`
	MaxConnIdleTime        time.Duration `yaml:"maxConnIdleTime" mapstructure:"maxConnIdleTime"`
	ConnectionTimeout      time.Duration `yaml:"connectionTimeout" mapstructure:"connectionTimeout"`
	ServerSelectionTimeout time.Duration `yaml:"serverSelectionTimeout" mapstructure:"serverSelectionTimeout"`
	SocketTimeout          time.Duration `yaml:"socketTimeout" mapstructure:"socketTimeout"`
	HeartbeatInterval      time.Duration `yaml:"heartbeatInterval" mapstructure:"heartbeatInterval"`
	LocalThreshold         time.Duration `yaml:"localThreshold" mapstructure:"localThreshold"`
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	SecretKey            string        `yaml:"secretKey" mapstructure:"secretKey"`
	AccessTokenDuration  time.Duration `yaml:"accessTokenDuration" mapstructure:"accessTokenDuration"`
	RefreshTokenDuration time.Duration `yaml:"refreshTokenDuration" mapstructure:"refreshTokenDuration"`
	Issuer               string        `yaml:"issuer" mapstructure:"issuer"`
	Audience             string        `yaml:"audience" mapstructure:"audience"`
}

// AWSConfig holds AWS-related configuration
type AWSConfig struct {
	Region          string    `yaml:"region" mapstructure:"region"`
	AccessKeyID     string    `yaml:"accessKeyId" mapstructure:"accessKeyId"`
	SecretAccessKey string    `yaml:"secretAccessKey" mapstructure:"secretAccessKey"`
	S3              S3Config  `yaml:"s3" mapstructure:"s3"`
	SES             SESConfig `yaml:"ses" mapstructure:"ses"`
}

// S3Config holds S3-specific configuration
type S3Config struct {
	Bucket         string `yaml:"bucket" mapstructure:"bucket"`
	Region         string `yaml:"region" mapstructure:"region"`
	PresignedURL   bool   `yaml:"presignedUrl" mapstructure:"presignedUrl"`
	URLExpiryHours int    `yaml:"urlExpiryHours" mapstructure:"urlExpiryHours"`
}

// SESConfig holds SES-specific configuration
type SESConfig struct {
	FromEmail    string `yaml:"fromEmail" mapstructure:"fromEmail"`
	FromName     string `yaml:"fromName" mapstructure:"fromName"`
	ReplyToEmail string `yaml:"replyToEmail" mapstructure:"replyToEmail"`
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	Provider string         `yaml:"provider" mapstructure:"provider"` // sendgrid, ses, smtp
	SendGrid SendGridConfig `yaml:"sendgrid" mapstructure:"sendgrid"`
	SMTP     SMTPConfig     `yaml:"smtp" mapstructure:"smtp"`
}

// SendGridConfig holds SendGrid configuration
type SendGridConfig struct {
	APIKey    string `yaml:"apiKey" mapstructure:"apiKey"`
	FromEmail string `yaml:"fromEmail" mapstructure:"fromEmail"`
	FromName  string `yaml:"fromName" mapstructure:"fromName"`
}

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	TLS      bool   `yaml:"tls" mapstructure:"tls"`
}

// SMSConfig holds SMS service configuration
type SMSConfig struct {
	Provider string        `yaml:"provider" mapstructure:"provider"` // twilio, aws-sns
	Twilio   TwilioConfig  `yaml:"twilio" mapstructure:"twilio"`
	AWSSNS   AWSSNSConfig  `yaml:"awsSns" mapstructure:"awsSns"`
}

// TwilioConfig holds Twilio configuration
type TwilioConfig struct {
	AccountSID string `yaml:"accountSid" mapstructure:"accountSid"`
	AuthToken  string `yaml:"authToken" mapstructure:"authToken"`
	FromNumber string `yaml:"fromNumber" mapstructure:"fromNumber"`
}

// AWSSNSConfig holds AWS SNS configuration
type AWSSNSConfig struct {
	Region string `yaml:"region" mapstructure:"region"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string        `yaml:"host" mapstructure:"host"`
	Port         string        `yaml:"port" mapstructure:"port"`
	Password     string        `yaml:"password" mapstructure:"password"`
	DB           int           `yaml:"db" mapstructure:"db"`
	PoolSize     int           `yaml:"poolSize" mapstructure:"poolSize"`
	MinIdleConns int           `yaml:"minIdleConns" mapstructure:"minIdleConns"`
	DialTimeout  time.Duration `yaml:"dialTimeout" mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout" mapstructure:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout" mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `yaml:"idleTimeout" mapstructure:"idleTimeout"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level" mapstructure:"level"`           // debug, info, warn, error
	Format     string `yaml:"format" mapstructure:"format"`         // json, console
	OutputPath string `yaml:"outputPath" mapstructure:"outputPath"` // stdout, stderr, file path
	MaxSize    int    `yaml:"maxSize" mapstructure:"maxSize"`       // MB
	MaxBackups int    `yaml:"maxBackups" mapstructure:"maxBackups"`
	MaxAge     int    `yaml:"maxAge" mapstructure:"maxAge"` // days
	Compress   bool   `yaml:"compress" mapstructure:"compress"`
}

var AppConfig *Config

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// Set default configuration path if not provided
	if configPath == "" {
		configPath = "configs"
	}

	// Get environment
	env := getEnv("APP_ENV", "development")

	// Configure Viper
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")

	// Set environment variable prefix
	viper.SetEnvPrefix("FOOD_APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Warning: Config file not found, using environment variables and defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Set defaults
	setDefaults()

	// Unmarshal configuration
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Set global config
	AppConfig = config

	return config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "9000")
	viper.SetDefault("server.readTimeout", "30s")
	viper.SetDefault("server.writeTimeout", "30s")
	viper.SetDefault("server.idleTimeout", "120s")
	viper.SetDefault("server.environment", "development")

	// CORS defaults
	viper.SetDefault("server.cors.allowedOrigins", []string{"http://localhost:3000"})
	viper.SetDefault("server.cors.allowedMethods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	viper.SetDefault("server.cors.allowedHeaders", []string{"Origin", "Content-Length", "Content-Type", "Authorization"})
	viper.SetDefault("server.cors.exposedHeaders", []string{"Content-Length"})
	viper.SetDefault("server.cors.allowCredentials", true)
	viper.SetDefault("server.cors.maxAge", "12h")

	// Database defaults
	viper.SetDefault("database.uri", "mongodb://localhost:27017")
	viper.SetDefault("database.name", "food-marketplace")
	viper.SetDefault("database.maxPoolSize", 100)
	viper.SetDefault("database.minPoolSize", 5)
	viper.SetDefault("database.maxConnIdleTime", "30m")
	viper.SetDefault("database.connectionTimeout", "10s")
	viper.SetDefault("database.serverSelectionTimeout", "5s")
	viper.SetDefault("database.socketTimeout", "10s")
	viper.SetDefault("database.heartbeatInterval", "10s")
	viper.SetDefault("database.localThreshold", "15ms")

	// JWT defaults
	viper.SetDefault("jwt.accessTokenDuration", "24h")
	viper.SetDefault("jwt.refreshTokenDuration", "168h") // 7 days
	viper.SetDefault("jwt.issuer", "food-marketplace")
	viper.SetDefault("jwt.audience", "food-marketplace-users")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.poolSize", 10)
	viper.SetDefault("redis.minIdleConns", 5)
	viper.SetDefault("redis.dialTimeout", "5s")
	viper.SetDefault("redis.readTimeout", "3s")
	viper.SetDefault("redis.writeTimeout", "3s")
	viper.SetDefault("redis.idleTimeout", "5m")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.outputPath", "stdout")
	viper.SetDefault("logging.maxSize", 100)
	viper.SetDefault("logging.maxBackups", 3)
	viper.SetDefault("logging.maxAge", 28)
	viper.SetDefault("logging.compress", true)

	// Email defaults
	viper.SetDefault("email.provider", "sendgrid")

	// SMS defaults
	viper.SetDefault("sms.provider", "twilio")

	// AWS defaults
	viper.SetDefault("aws.region", "us-west-2")
	viper.SetDefault("aws.s3.region", "us-west-2")
	viper.SetDefault("aws.s3.presignedUrl", true)
	viper.SetDefault("aws.s3.urlExpiryHours", 24)
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate JWT secret key
	if config.JWT.SecretKey == "" {
		return fmt.Errorf("JWT secret key is required")
	}

	// Validate database URI
	if config.Database.URI == "" {
		return fmt.Errorf("database URI is required")
	}

	// Validate server port
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	// Validate environment
	validEnvs := []string{"development", "staging", "production"}
	if !contains(validEnvs, config.Server.Environment) {
		return fmt.Errorf("invalid environment: %s. Valid options: %v", config.Server.Environment, validEnvs)
	}

	return nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetConfigPath returns the absolute path to config directory
func GetConfigPath() string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "configs")
}

// IsProduction returns true if running in production environment
func IsProduction() bool {
	return AppConfig != nil && AppConfig.Server.Environment == "production"
}

// IsDevelopment returns true if running in development environment
func IsDevelopment() bool {
	return AppConfig != nil && AppConfig.Server.Environment == "development"
}

// GetDatabaseURI returns the complete database connection URI
func GetDatabaseURI() string {
	if AppConfig == nil {
		return "mongodb://localhost:27017"
	}
	return AppConfig.Database.URI
}

// GetJWTSecret returns the JWT secret key
func GetJWTSecret() string {
	if AppConfig == nil {
		return ""
	}
	return AppConfig.JWT.SecretKey
}