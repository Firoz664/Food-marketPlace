package routes

import (
	orderController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	orderGroup := router.Group("/api/v1/order")
	{
		orderGroup.GET("/getAllOrder", orderController.GetOrders())
	}
}
