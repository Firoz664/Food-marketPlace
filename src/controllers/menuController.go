package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

// GetInvoice returns a Gin handler function
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation here
	}
}
