package controllers

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var riderCollection *mongo.Collection

// getRiderCollection returns the rider collection with lazy initialization
func getRiderCollection() *mongo.Collection {
	if riderCollection == nil {
		client := database.DBInstance()
		riderCollection = database.OpenCollection(client, "rider")
	}
	return riderCollection
}
var vehiclesCollection *mongo.Collection

// getVehiclesCollection returns the vehicle collection with lazy initialization
func getVehiclesCollection() *mongo.Collection {
	if vehiclesCollection == nil {
		client := database.DBInstance()
		vehiclesCollection = database.OpenCollection(client, "vehicle") // Fixed typo: "vechile" -> "vehicle"
	}
	return vehiclesCollection
}

// GetInvoice returns a Gin handler function
func GetRider() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
