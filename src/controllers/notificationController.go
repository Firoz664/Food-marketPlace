package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var notificationCollection *mongo.Collection = database.OpenCollection(database.Client, "notification")

// GetInvoice returns a Gin handler function
func GetNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
