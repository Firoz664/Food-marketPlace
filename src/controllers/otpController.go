package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var otpCollection *mongo.Collection = database.OpenCollection(database.Client, "otp")

// GetInvoice returns a Gin handler function
func GetOtps() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
