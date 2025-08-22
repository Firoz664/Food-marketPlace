package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var notificationCollection *mongo.Collection

// getNotificationCollection returns the notification collection with lazy initialization
func getNotificationCollection() *mongo.Collection {
	if notificationCollection == nil {
		client := database.DBInstance()
		notificationCollection = database.OpenCollection(client, "notification")
	}
	return notificationCollection
}

// GetInvoice returns a Gin handler function
func GetNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
