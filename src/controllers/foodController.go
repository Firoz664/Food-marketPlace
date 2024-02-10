package controllers

import (
	"fmt"
	"net/http"

	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/foodmngtapp/food-management-apps/src/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve user ID from context
		userID, _ := c.Get("userID")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		// Retrieve user details from MongoDB
		var user models.User
		err := userCollection.FindOne(c, bson.M{"_id": userID, "isAdmin": true}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not an admin or not found"})
			return
		}

		// Only admins are allowed to create food items
		// Implement the logic to create a food item here

		// Simulate successful creation of a food item
		foodName := "New Food Item"
		fmt.Printf("Admin '%s' created food: %s\n", user.Email, foodName)

		// Return success response
		c.JSON(http.StatusOK, gin.H{
			"message":     "Food created successfully",
			"foodName":    foodName,
			"createdBy":   user.Email,
			"createdByID": user.ID,
		})
	}
}

// GetFoods returns a Gin handler function
func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulating a list of foods fetched from somewhere
		foods := []string{"Apple", "Banana", "Carrot"}
		fmt.Println("Food------>>>", foods)
		// Returning the list of foods in JSON format
		c.JSON(200, gin.H{
			"message": "Foods fetched successfully",
			"foods":   foods,
		})
	}
}
