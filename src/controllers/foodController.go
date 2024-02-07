package controllers

import (
	"fmt"

	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

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
