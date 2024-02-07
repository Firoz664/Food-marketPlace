package routes

import (
	menuController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	menuGroup := router.Group("/api/v1/menu")
	{
		menuGroup.GET("/getAllMenu", menuController.GetMenu())
	}
}
