package controllers

import (
	"net/http"

	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var transactionCollection *mongo.Collection = database.OpenCollection(database.Client, "transaction")

func GetByIdTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user

		// Log fetching user operation
		// logger.Log.Info("Fetched user successfully", zap.String("userID"))

		// Return the dummy user object along with status code and message
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Fatch All trnasaction fetched successfully",
		})
	}
}
func GetAllTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user

		// Log fetching user operation

		// Return the dummy user object along with status code and message
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Fatch All trnasaction fetched successfully",
		})
	}
}
