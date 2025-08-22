package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection

// getTableCollection returns the table collection with lazy initialization
func getTableCollection() *mongo.Collection {
	if tableCollection == nil {
		client := database.DBInstance()
		tableCollection = database.OpenCollection(client, "table")
	}
	return tableCollection
}

// GetInvoice returns a Gin handler function
func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
