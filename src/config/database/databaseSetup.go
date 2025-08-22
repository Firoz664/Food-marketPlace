package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/foodmngtapp/food-management-apps/src/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

var (
	Client *mongo.Client
	once   sync.Once
	logger *zap.Logger
)

// DatabaseManager manages database connections and operations
type DatabaseManager struct {
	Client *mongo.Client
	Config *config.DatabaseConfig
	Logger *zap.Logger
}

// InitializeDatabase initializes the database connection with configuration
func InitializeDatabase(cfg *config.DatabaseConfig) (*mongo.Client, error) {
	var err error
	
	once.Do(func() {
		// Initialize logger
		if cfg != nil && config.IsProduction() {
			logger, _ = zap.NewProduction()
		} else {
			logger, _ = zap.NewDevelopment()
		}

		// Use config if provided, otherwise use defaults
		dbURI := "mongodb://localhost:27017"
		if cfg != nil && cfg.URI != "" {
			dbURI = cfg.URI
		}

		// Create MongoDB client options with enhanced configuration
		clientOptions := options.Client().ApplyURI(dbURI)

		if cfg != nil {
			// Set connection pool options
			clientOptions.SetMaxPoolSize(cfg.MaxPoolSize)
			clientOptions.SetMinPoolSize(cfg.MinPoolSize)
			clientOptions.SetMaxConnIdleTime(cfg.MaxConnIdleTime)
			
			// Set timeout options
			clientOptions.SetConnectTimeout(cfg.ConnectionTimeout)
			clientOptions.SetServerSelectionTimeout(cfg.ServerSelectionTimeout)
			clientOptions.SetSocketTimeout(cfg.SocketTimeout)
			
			// Set heartbeat and monitoring options
			clientOptions.SetHeartbeatInterval(cfg.HeartbeatInterval)
			clientOptions.SetLocalThreshold(cfg.LocalThreshold)
		} else {
			// Set default values
			clientOptions.SetMaxPoolSize(100)
			clientOptions.SetMinPoolSize(5)
			clientOptions.SetMaxConnIdleTime(30 * time.Minute)
			clientOptions.SetConnectTimeout(10 * time.Second)
			clientOptions.SetServerSelectionTimeout(5 * time.Second)
			clientOptions.SetSocketTimeout(10 * time.Second)
		}

		// Create context for connection
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Create MongoDB client
		Client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			logger.Fatal("Failed to create MongoDB client", zap.Error(err))
			return
		}

		// Ping the database to verify connection
		err = Client.Ping(ctx, readpref.Primary())
		if err != nil {
			logger.Fatal("Failed to ping MongoDB", zap.Error(err))
			return
		}

		dbName := "food-marketplace"
		if cfg != nil && cfg.Name != "" {
			dbName = cfg.Name
		}

		maxPool := uint64(100)
		minPool := uint64(5)
		if cfg != nil {
			maxPool = cfg.MaxPoolSize
			minPool = cfg.MinPoolSize
		}

		logger.Info("Successfully connected to MongoDB",
			zap.String("uri", obfuscateURI(dbURI)),
			zap.String("database", dbName),
			zap.Uint64("maxPoolSize", maxPool),
			zap.Uint64("minPoolSize", minPool),
		)

		// Create indexes for optimal performance (run in background)
		go func() {
			if err := CreateAllIndexes(Client, dbName); err != nil {
				logger.Error("Failed to create database indexes", zap.Error(err))
			} else {
				logger.Info("Database indexes created successfully")
			}
		}()
	})

	return Client, err
}

// DBInstance returns the singleton MongoDB client instance
func DBInstance() *mongo.Client {
	if Client == nil {
		// Fallback to simple initialization if not already initialized
		cfg := &config.DatabaseConfig{
			URI:                    "mongodb://localhost:27017",
			Name:                   "food-marketplace",
			MaxPoolSize:            100,
			MinPoolSize:            5,
			MaxConnIdleTime:        30 * time.Minute,
			ConnectionTimeout:      10 * time.Second,
			ServerSelectionTimeout: 5 * time.Second,
			SocketTimeout:          10 * time.Second,
			HeartbeatInterval:      10 * time.Second,
			LocalThreshold:         15 * time.Millisecond,
		}
		
		client, err := InitializeDatabase(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		return client
	}
	return Client
}

// OpenCollection opens a collection with the given name
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := "food-marketplace"
	
	// Use config database name if available
	if config.AppConfig != nil && config.AppConfig.Database.Name != "" {
		dbName = config.AppConfig.Database.Name
	}
	
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}

// GetDatabase returns the database instance
func GetDatabase(client *mongo.Client) *mongo.Database {
	dbName := "food-marketplace"
	
	// Use config database name if available
	if config.AppConfig != nil && config.AppConfig.Database.Name != "" {
		dbName = config.AppConfig.Database.Name
	}
	
	return client.Database(dbName)
}

// HealthCheck performs a health check on the database connection
func HealthCheck(ctx context.Context) error {
	if Client == nil {
		return fmt.Errorf("database client is not initialized")
	}
	
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	return Client.Ping(ctx, readpref.Primary())
}

// Close gracefully closes the database connection
func Close(ctx context.Context) error {
	if Client == nil {
		return nil
	}
	
	if logger != nil {
		logger.Info("Closing MongoDB connection")
	}
	
	return Client.Disconnect(ctx)
}

// obfuscateURI removes sensitive information from URI for logging
func obfuscateURI(uri string) string {
	// Simple obfuscation - in production, use more sophisticated method
	if len(uri) > 20 {
		return uri[:10] + "***" + uri[len(uri)-7:]
	}
	return "***"
}

// NewDatabaseManager creates a new database manager instance
func NewDatabaseManager(cfg *config.DatabaseConfig) (*DatabaseManager, error) {
	client, err := InitializeDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	
	return &DatabaseManager{
		Client: client,
		Config: cfg,
		Logger: logger,
	}, nil
}
