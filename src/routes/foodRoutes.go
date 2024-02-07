package routes

import (
	foodController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(router *gin.Engine) {
	foodGroup := router.Group("/api/v1/food")
	{
		// Directly use GetInvoice as a handler without invoking it.
		foodGroup.GET("/getAllFood", foodController.GetFoods())
	}
}
