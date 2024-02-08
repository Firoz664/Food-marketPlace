package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// Dummy user struct for demonstration purposes

func S3FileUploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user

		// Log fetching user operation

		// Return the dummy user object along with status code and message
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User fetched successfully",
		})
	}
}
func MultipleUploadFiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user

		// Log fetching user operation
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User fetched successfully",
		})
	}
}
