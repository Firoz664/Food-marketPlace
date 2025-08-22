package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection

// getMenuCollection returns the menu collection with lazy initialization
func getMenuCollection() *mongo.Collection {
	if menuCollection == nil {
		client := database.DBInstance()
		menuCollection = database.OpenCollection(client, "menu")
	}
	return menuCollection
}

// GetInvoice returns a Gin handler function
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
