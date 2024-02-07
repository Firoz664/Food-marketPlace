package routes

import (
	transactionController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func TransactionRoutes(router *gin.Engine) {
	transactionGroup := router.Group("/api/v1/transaction")
	{
		transactionGroup.GET("/GetAllTransaction", transactionController.GetTransaction())
	}
}
