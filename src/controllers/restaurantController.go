package controllers

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/foodmngtapp/food-management-apps/src/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var restaurantCollection *mongo.Collection = database.OpenCollection(database.Client, "restaurant")
var restaurantValidate = validator.New()

func CreateRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("userId")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
			return
		}

		var user models.User
		err := userCollection.FindOne(context.TODO(), bson.M{"userId": userID, "profileType": "RESTAURANT"}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found!"})
			return
		}

		var restaurant models.Restaurant
		if err := c.ShouldBindJSON(&restaurant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		validationErr := restaurantValidate.Struct(restaurant)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if restaurant.RestaurantName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Restaurant name is required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		restaurant.UserID = user.UserID
		restaurant.Status = "Pending"
		result, err := restaurantCollection.InsertOne(ctx, restaurant)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"status":  true,
			"message": "Restaurant created successfully",
			"result":  result,
		})
	}
}

func GetRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the restaurant ID from the request parameters
		restaurantID := c.Query("restaurantId")
		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
			return
		}

		// Convert the restaurant ID string to an ObjectID
		objID, err := primitive.ObjectIDFromHex(restaurantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}

		// Create a context with a timeout for MongoDB operations
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Query the database to find the restaurant by its ID
		var restaurant models.Restaurant
		err = restaurantCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&restaurant)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restaurant"})
			return
		}

		c.JSON(http.StatusOK, restaurant)
	}
}

func GetAllRestaurants() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Pagination parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		skip := (page - 1) * limit

		// Search parameters
		search := c.Query("search")

		// MongoDB aggregation pipeline for pagination and search
		pipeline := make([]bson.M, 0)
		matchStage := bson.M{"$match": bson.M{"status": "Approve"}}
		pipeline = append(pipeline, matchStage)

		if search != "" {
			searchQuery := bson.M{
				"$or": []bson.M{
					{"restaurantName": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
					{"description": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
				},
			}
			pipeline = append(pipeline, bson.M{"$match": searchQuery})
		}

		cursor, err := restaurantCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restaurants"})
			return
		}
		defer cursor.Close(ctx)

		var restaurants []models.Restaurant
		if err := cursor.All(ctx, &restaurants); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode restaurants"})
			return
		}

		// Calculate total count
		totalCount := len(restaurants)
		// Apply pagination
		startIndex := skip
		endIndex := minHelper(skip+limit, totalCount)
		paginatedRestaurants := restaurants[startIndex:endIndex]

		// Calculate total pages
		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"status":      true,
				"statusCode":  http.StatusOK,
				"message":     "Fetch All Restaurants",
				"totalPages":  totalPages,
				"currentPage": page,
				"totalCount":  totalCount,
				"result":      paginatedRestaurants,
			},
		})

	}
}
func minHelper(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func UpdateRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the restaurant ID from the request parameters
		restaurantID := c.Query("restaurantId")
		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
			return
		}

		// Convert the restaurant ID string to an ObjectID
		objID, err := primitive.ObjectIDFromHex(restaurantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}

		// Parse the updated restaurant data from the request body
		var updatedRestaurant models.Restaurant
		if err := c.ShouldBindJSON(&updatedRestaurant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		// Create a context with a timeout for MongoDB operations
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Update the restaurant profile in the database
		_, err = restaurantCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updatedRestaurant})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restaurant"})
			return
		}

		// Respond with success message
		c.JSON(http.StatusOK, gin.H{"message": "Restaurant profile updated successfully"})
	}
}

func DeleteRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID := c.Query("restaurantId")
		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
			return
		}
		objID, err := primitive.ObjectIDFromHex(restaurantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err = restaurantCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete restaurant"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Restaurant profile deleted successfully"})
	}
}
