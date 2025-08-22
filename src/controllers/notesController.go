package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var noteCollection *mongo.Collection

// getNoteCollection returns the note collection with lazy initialization
func getNoteCollection() *mongo.Collection {
	if noteCollection == nil {
		client := database.DBInstance()
		noteCollection = database.OpenCollection(client, "note")
	}
	return noteCollection
}

// GetInvoice returns a Gin handler function
func GetNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
