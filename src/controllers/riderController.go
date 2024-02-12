package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var riderCollection *mongo.Collection = database.OpenCollection(database.Client, "rider")
var vehiclesCollection *mongo.Collection = database.OpenCollection(database.Client, "vechile")

// GetInvoice returns a Gin handler function
func GetRider() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
