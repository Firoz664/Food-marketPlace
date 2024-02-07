package routes

import (
	tablesController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func TableRoutes(router *gin.Engine) {
	tablesGroup := router.Group("api/v1/tables")
	{
		tablesGroup.GET("/GetAllTables", tablesController.GetTables())
	}
}
