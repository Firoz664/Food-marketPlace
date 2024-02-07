package controllers

import (
	"net/http"

	logger "github.com/foodmngtapp/food-management-apps/src/common/logger"
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

// Dummy user struct for demonstration purposes
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user
		dummyUser := User{
			ID:   "1",
			Name: "John Doe",
		}

		// Log fetching user operation
		logger.Log.Info("Fetched user successfully", zap.String("userID", dummyUser.ID))
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User fetched successfully",
			"user":    dummyUser,
		})
	}
}

func GetUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user
		dummyUser := User{
			ID:   "1",
			Name: "John Doe",
		}

		// Log fetching user operation
		logger.Log.Info("Fetched user successfully", zap.String("userID", dummyUser.ID))

		// Return the dummy user object along with status code and message
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User fetched successfully",
			"user":    dummyUser,
		})
	}
}

func UpdateUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dummy user before update
		dummyUser := User{
			ID:   "1",
			Name: "John Doe",
		}

		// Attempt to bind the incoming JSON payload to the updatedUser struct
		var updatedUser User
		if err := c.ShouldBindJSON(&updatedUser); err != nil {
			// Use logger.Error to log the error
			logger.Log.Error("Error binding JSON to User struct", zap.Error(err))

			// Return a 400 Bad Request response with the error message
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Log the successful update operation
		logger.Log.Info("Updated user details successfully", zap.String("userID", dummyUser.ID))

		// Simulate updating the user's name and return the updated user object
		dummyUser.Name = updatedUser.Name

		// Return the updated user object along with status code and message
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User details updated successfully",
			"user":    dummyUser,
		})
	}
}
