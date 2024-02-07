package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

// GetInvoice returns a Gin handler function
func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
