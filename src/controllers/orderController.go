package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

// GetInvoice returns a Gin handler function
func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
