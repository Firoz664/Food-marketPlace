package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var noteCollection *mongo.Collection = database.OpenCollection(database.Client, "note")

// GetInvoice returns a Gin handler function
func GetNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
