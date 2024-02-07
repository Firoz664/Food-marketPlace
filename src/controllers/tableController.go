package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

// GetInvoice returns a Gin handler function
func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
