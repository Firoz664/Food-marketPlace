package controllers

import (
	"net/http"

	"github.com/foodmngtapp/food-management-apps/src/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func VerifyRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve user ID from context
		userID := c.Query("userId")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		// Retrieve user details from MongoDB
		var user models.User
		err := userCollection.FindOne(c, bson.M{"userId": userID, "profileType": "ADMIN"}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not an admin or not found"})
			return
		}
		// Check if payload contains status field
		var payload struct {
			Status       string `json:"status" binding:"required"`
			RestaurantID string `json:"restaurantId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Retrieve restaurant ID from payload
		restaurantID, err := primitive.ObjectIDFromHex(payload.RestaurantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}

		// Retrieve restaurant details from MongoDB
		var restaurant models.Restaurant
		err = restaurantCollection.FindOne(c, bson.M{"_id": restaurantID}).Decode(&restaurant)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}

		// Update restaurant status based on payload
		_, err = restaurantCollection.UpdateOne(c, bson.M{"_id": restaurantID}, bson.M{"$set": bson.M{"status": payload.Status}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restaurant status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Restaurant status updated successfully", "restaurantID": restaurantID, "status": payload.Status})
	}
}

// GetFoods returns a Gin handler function
// func GetFoods() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Simulating a list of foods fetched from somewhere
// 		foods := []string{"Apple", "Banana", "Carrot"}
// 		fmt.Println("Food------>>>", foods)
// 		// Returning the list of foods in JSON format
// 		c.JSON(200, gin.H{
// 			"message": "Foods fetched successfully",
// 			"foods":   foods,
// 		})
// 	}
// }
