package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var transactionCollection *mongo.Collection = database.OpenCollection(database.Client, "transaction")

// GetInvoice returns a Gin handler function
func GetTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
