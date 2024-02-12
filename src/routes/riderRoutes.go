package routes

import (
	riderController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func RiderRoutes(router *gin.Engine) {
	tablesGroup := router.Group("api/v1/tables")
	{
		tablesGroup.GET("/GetAllRider", riderController.GetRider())
	}
}
