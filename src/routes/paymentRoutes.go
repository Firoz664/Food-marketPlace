package routes

import (
	paymentController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(router *gin.Engine) {
	paymentGroup := router.Group("/api/v1/payment")
	{
		paymentGroup.GET("/getAllPayment", paymentController.GetPayment())
	}
}
