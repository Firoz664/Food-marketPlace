package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection

// getOrderCollection returns the order collection with lazy initialization
func getOrderCollection() *mongo.Collection {
	if orderCollection == nil {
		client := database.DBInstance()
		orderCollection = database.OpenCollection(client, "order")
	}
	return orderCollection
}

// GetInvoice returns a Gin handler function
func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
