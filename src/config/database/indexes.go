package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexDefinition represents a database index
type IndexDefinition struct {
	Collection string
	Index      mongo.IndexModel
	Name       string
}

// CreateAllIndexes creates all necessary indexes for optimal performance
func CreateAllIndexes(client *mongo.Client, databaseName string) error {
	db := client.Database(databaseName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	indexes := getAllIndexes()
	
	log.Printf("🔧 Creating %d database indexes for optimal performance...", len(indexes))
	
	for _, idx := range indexes {
		collection := db.Collection(idx.Collection)
		
		// Check if index already exists
		cursor, err := collection.Indexes().List(ctx)
		if err != nil {
			log.Printf("❌ Error listing indexes for %s: %v", idx.Collection, err)
			continue
		}
		
		indexExists := false
		var existingIndexes []bson.M
		if err = cursor.All(ctx, &existingIndexes); err == nil {
			for _, existingIndex := range existingIndexes {
				if name, ok := existingIndex["name"].(string); ok && name == idx.Name {
					indexExists = true
					break
				}
			}
		}
		cursor.Close(ctx)
		
		if indexExists {
			log.Printf("✅ Index %s already exists on %s", idx.Name, idx.Collection)
			continue
		}
		
		// Create the index
		_, err = collection.Indexes().CreateOne(ctx, idx.Index)
		if err != nil {
			log.Printf("❌ Failed to create index %s on %s: %v", idx.Name, idx.Collection, err)
			continue
		}
		
		log.Printf("✅ Created index %s on %s", idx.Name, idx.Collection)
	}
	
	log.Println("🎉 Database indexing completed successfully!")
	return nil
}

// getAllIndexes returns all index definitions for the application
func getAllIndexes() []IndexDefinition {
	return []IndexDefinition{
		// ============ USERS COLLECTION INDEXES ============
		{
			Collection: "users",
			Name:       "idx_users_email_unique",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "email", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("idx_users_email_unique").
					SetBackground(true),
			},
		},
		{
			Collection: "users",
			Name:       "idx_users_mobile_unique", 
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "mobileNumber", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("idx_users_mobile_unique").
					SetBackground(true).
					SetSparse(true), // Allow null values
			},
		},
		{
			Collection: "users",
			Name:       "idx_users_profile_type",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "profileType", Value: 1}},
				Options: options.Index().
					SetName("idx_users_profile_type").
					SetBackground(true),
			},
		},
		{
			Collection: "users",
			Name:       "idx_users_verification_status",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "isVerify", Value: 1},
					{Key: "isComplete", Value: 1},
				},
				Options: options.Index().
					SetName("idx_users_verification_status").
					SetBackground(true),
			},
		},
		{
			Collection: "users",
			Name:       "idx_users_created_at",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "createdAt", Value: -1}}, // Descending for recent first
				Options: options.Index().
					SetName("idx_users_created_at").
					SetBackground(true),
			},
		},
		{
			Collection: "users",
			Name:       "idx_users_last_active",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "lastActive", Value: -1}},
				Options: options.Index().
					SetName("idx_users_last_active").
					SetBackground(true),
			},
		},
		{
			Collection: "users",
			Name:       "idx_users_location_city_state",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "location.city", Value: 1},
					{Key: "location.state", Value: 1},
				},
				Options: options.Index().
					SetName("idx_users_location_city_state").
					SetBackground(true),
			},
		},
		{
			Collection: "users",
			Name:       "idx_users_geospatial",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "location.coordinates", Value: "2dsphere"},
				},
				Options: options.Index().
					SetName("idx_users_geospatial").
					SetBackground(true).
					SetSparse(true),
			},
		},
		
		// ============ RESTAURANTS COLLECTION INDEXES ============
		{
			Collection: "restaurants",
			Name:       "idx_restaurants_user_id",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "userId", Value: 1}},
				Options: options.Index().
					SetName("idx_restaurants_user_id").
					SetBackground(true),
			},
		},
		{
			Collection: "restaurants",
			Name:       "idx_restaurants_status_delivering",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "status", Value: 1},
					{Key: "isDelivering", Value: 1},
				},
				Options: options.Index().
					SetName("idx_restaurants_status_delivering").
					SetBackground(true),
			},
		},
		{
			Collection: "restaurants",
			Name:       "idx_restaurants_location_city_state",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "location.city", Value: 1},
					{Key: "location.state", Value: 1},
				},
				Options: options.Index().
					SetName("idx_restaurants_location_city_state").
					SetBackground(true),
			},
		},
		{
			Collection: "restaurants",
			Name:       "idx_restaurants_geospatial",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "location.coordinates", Value: "2dsphere"},
				},
				Options: options.Index().
					SetName("idx_restaurants_geospatial").
					SetBackground(true).
					SetSparse(true),
			},
		},
		{
			Collection: "restaurants",
			Name:       "idx_restaurants_cuisines",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "cuisines", Value: 1}},
				Options: options.Index().
					SetName("idx_restaurants_cuisines").
					SetBackground(true),
			},
		},
		{
			Collection: "restaurants",
			Name:       "idx_restaurants_name_text",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "restaurantName", Value: "text"},
					{Key: "description", Value: "text"},
				},
				Options: options.Index().
					SetName("idx_restaurants_name_text").
					SetBackground(true).
					SetDefaultLanguage("english"),
			},
		},
		
		// ============ FOOD COLLECTION INDEXES ============ 
		{
			Collection: "foods",
			Name:       "idx_foods_menu_id",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "menuId", Value: 1}},
				Options: options.Index().
					SetName("idx_foods_menu_id").
					SetBackground(true),
			},
		},
		{
			Collection: "foods",
			Name:       "idx_foods_name_text",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "foodName", Value: "text"}},
				Options: options.Index().
					SetName("idx_foods_name_text").
					SetBackground(true).
					SetDefaultLanguage("english"),
			},
		},
		{
			Collection: "foods",
			Name:       "idx_foods_price",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "price", Value: 1}},
				Options: options.Index().
					SetName("idx_foods_price").
					SetBackground(true),
			},
		},
		{
			Collection: "foods",
			Name:       "idx_foods_created_at",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "createdAt", Value: -1}},
				Options: options.Index().
					SetName("idx_foods_created_at").
					SetBackground(true),
			},
		},
		
		// ============ MENU COLLECTION INDEXES ============
		{
			Collection: "menus",
			Name:       "idx_menus_category",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "category", Value: 1}},
				Options: options.Index().
					SetName("idx_menus_category").
					SetBackground(true),
			},
		},
		{
			Collection: "menus",
			Name:       "idx_menus_name_text",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "name", Value: "text"}},
				Options: options.Index().
					SetName("idx_menus_name_text").
					SetBackground(true).
					SetDefaultLanguage("english"),
			},
		},
		{
			Collection: "menus",
			Name:       "idx_menus_date_range",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "startDate", Value: 1},
					{Key: "endDate", Value: 1},
				},
				Options: options.Index().
					SetName("idx_menus_date_range").
					SetBackground(true).
					SetSparse(true),
			},
		},
		
		// ============ OTP COLLECTION INDEXES ============
		{
			Collection: "otps",
			Name:       "idx_otps_email",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "email", Value: 1}},
				Options: options.Index().
					SetName("idx_otps_email").
					SetBackground(true).
					SetSparse(true),
			},
		},
		{
			Collection: "otps",
			Name:       "idx_otps_mobile",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "mobileNumber", Value: 1}},
				Options: options.Index().
					SetName("idx_otps_mobile").
					SetBackground(true).
					SetSparse(true),
			},
		},
		{
			Collection: "otps",
			Name:       "idx_otps_expiry_ttl",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "expiredAt", Value: 1}},
				Options: options.Index().
					SetName("idx_otps_expiry_ttl").
					SetBackground(true).
					SetExpireAfterSeconds(0), // TTL index for automatic cleanup
			},
		},
		{
			Collection: "otps",
			Name:       "idx_otps_verification_status",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "isVerify", Value: 1},
					{Key: "createdAt", Value: -1},
				},
				Options: options.Index().
					SetName("idx_otps_verification_status").
					SetBackground(true),
			},
		},
		
		// ============ PAYMENTS COLLECTION INDEXES ============
		{
			Collection: "payments",
			Name:       "idx_payments_user_id",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "userId", Value: 1}},
				Options: options.Index().
					SetName("idx_payments_user_id").
					SetBackground(true),
			},
		},
		{
			Collection: "payments",
			Name:       "idx_payments_payment_id_unique",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "paymentId", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("idx_payments_payment_id_unique").
					SetBackground(true).
					SetSparse(true),
			},
		},
		{
			Collection: "payments",
			Name:       "idx_payments_status_created",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "paymentStatus", Value: 1},
					{Key: "createdAt", Value: -1},
				},
				Options: options.Index().
					SetName("idx_payments_status_created").
					SetBackground(true),
			},
		},
		
		// ============ INVOICES COLLECTION INDEXES ============
		{
			Collection: "invoices",
			Name:       "idx_invoices_invoice_id_unique",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "invoiceId", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("idx_invoices_invoice_id_unique").
					SetBackground(true),
			},
		},
		{
			Collection: "invoices",
			Name:       "idx_invoices_order_id",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "orderId", Value: 1}},
				Options: options.Index().
					SetName("idx_invoices_order_id").
					SetBackground(true),
			},
		},
		{
			Collection: "invoices",
			Name:       "idx_invoices_payment_status",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "paymentStatus", Value: 1},
					{Key: "paymentDueDate", Value: 1},
				},
				Options: options.Index().
					SetName("idx_invoices_payment_status").
					SetBackground(true),
			},
		},
		
		// ============ DRIVERS COLLECTION INDEXES ============
		{
			Collection: "drivers",
			Name:       "idx_drivers_email_unique",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "email", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("idx_drivers_email_unique").
					SetBackground(true),
			},
		},
		{
			Collection: "drivers",
			Name:       "idx_drivers_phone_unique",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "phoneNumber", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("idx_drivers_phone_unique").
					SetBackground(true),
			},
		},
		{
			Collection: "drivers",
			Name:       "idx_drivers_license_unique",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "license", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("idx_drivers_license_unique").
					SetBackground(true),
			},
		},
		{
			Collection: "drivers", 
			Name:       "idx_drivers_status",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "status", Value: 1}},
				Options: options.Index().
					SetName("idx_drivers_status").
					SetBackground(true),
			},
		},
		
		// ============ VEHICLES COLLECTION INDEXES ============
		{
			Collection: "vehicles",
			Name:       "idx_vehicles_plate_unique",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "plateNumber", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("idx_vehicles_plate_unique").
					SetBackground(true),
			},
		},
		{
			Collection: "vehicles",
			Name:       "idx_vehicles_owner_id",
			Index: mongo.IndexModel{
				Keys: bson.D{{Key: "ownerId", Value: 1}},
				Options: options.Index().
					SetName("idx_vehicles_owner_id").
					SetBackground(true),
			},
		},
		{
			Collection: "vehicles",
			Name:       "idx_vehicles_make_model",
			Index: mongo.IndexModel{
				Keys: bson.D{
					{Key: "make", Value: 1},
					{Key: "model", Value: 1},
					{Key: "year", Value: -1},
				},
				Options: options.Index().
					SetName("idx_vehicles_make_model").
					SetBackground(true),
			},
		},
	}
}

// DropAllIndexes drops all custom indexes (useful for development/testing)
func DropAllIndexes(client *mongo.Client, databaseName string) error {
	db := client.Database(databaseName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collections := []string{
		"users", "restaurants", "foods", "menus", 
		"otps", "payments", "invoices", "drivers", "vehicles",
	}
	
	log.Println("🗑️  Dropping all custom indexes...")
	
	for _, collectionName := range collections {
		collection := db.Collection(collectionName)
		
		// Get all indexes
		cursor, err := collection.Indexes().List(ctx)
		if err != nil {
			log.Printf("❌ Error listing indexes for %s: %v", collectionName, err)
			continue
		}
		
		var indexes []bson.M
		if err = cursor.All(ctx, &indexes); err != nil {
			log.Printf("❌ Error reading indexes for %s: %v", collectionName, err)
			continue
		}
		cursor.Close(ctx)
		
		for _, index := range indexes {
			if name, ok := index["name"].(string); ok && name != "_id_" {
				_, err := collection.Indexes().DropOne(ctx, name)
				if err != nil {
					log.Printf("❌ Failed to drop index %s on %s: %v", name, collectionName, err)
				} else {
					log.Printf("✅ Dropped index %s on %s", name, collectionName)
				}
			}
		}
	}
	
	log.Println("🎉 Index dropping completed!")
	return nil
}

// GetIndexInfo returns information about all indexes
func GetIndexInfo(client *mongo.Client, databaseName string) error {
	db := client.Database(databaseName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collections := []string{
		"users", "restaurants", "foods", "menus", 
		"otps", "payments", "invoices", "drivers", "vehicles",
	}
	
	fmt.Println("\n📊 DATABASE INDEX INFORMATION")
	fmt.Println(strings.Repeat("=", 50))
	
	for _, collectionName := range collections {
		collection := db.Collection(collectionName)
		
		cursor, err := collection.Indexes().List(ctx)
		if err != nil {
			log.Printf("❌ Error listing indexes for %s: %v", collectionName, err)
			continue
		}
		
		var indexes []bson.M
		if err = cursor.All(ctx, &indexes); err != nil {
			log.Printf("❌ Error reading indexes for %s: %v", collectionName, err)
			continue
		}
		cursor.Close(ctx)
		
		fmt.Printf("\n🗂️  Collection: %s (%d indexes)\n", collectionName, len(indexes))
		for _, index := range indexes {
			if name, ok := index["name"].(string); ok {
				keys := index["key"]
				unique := ""
				if val, exists := index["unique"]; exists && val.(bool) {
					unique = " [UNIQUE]"
				}
				sparse := ""
				if val, exists := index["sparse"]; exists && val.(bool) {
					sparse = " [SPARSE]"
				}
				ttl := ""
				if val, exists := index["expireAfterSeconds"]; exists {
					ttl = fmt.Sprintf(" [TTL: %v]", val)
				}
				fmt.Printf("  • %s: %v%s%s%s\n", name, keys, unique, sparse, ttl)
			}
		}
	}
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	return nil
}