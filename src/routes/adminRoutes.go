package routes

import (
	adminController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	tablesGroup := router.Group("api/v1/admin")
	{
		tablesGroup.POST("/verifyRestaurant", adminController.VerifyRestaurant())
	}
}
