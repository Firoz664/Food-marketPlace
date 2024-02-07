package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var paymentCollection *mongo.Collection = database.OpenCollection(database.Client, "payment")

// GetInvoice returns a Gin handler function
func GetPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
