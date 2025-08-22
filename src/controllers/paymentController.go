package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var paymentCollection *mongo.Collection

// getPaymentCollection returns the payment collection with lazy initialization
func getPaymentCollection() *mongo.Collection {
	if paymentCollection == nil {
		client := database.DBInstance()
		paymentCollection = database.OpenCollection(client, "payment")
	}
	return paymentCollection
}

// GetInvoice returns a Gin handler function
func GetPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
