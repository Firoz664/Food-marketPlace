package routes

import (
	orderItemController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(router *gin.Engine) {
	orderItemGroup := router.Group("/api/v1/OrderItem")
	{
		orderItemGroup.GET("/getOrderItem", orderItemController.GetOrderItems())
	}
}
