package routes

import (
	AuthController "github.com/foodmngtapp/food-management-apps/src/common/middleware"
	userController "github.com/foodmngtapp/food-management-apps/src/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/v1/user")
	{
		userGroup.GET("/getAllUsers", userController.GetUser())
		userGroup.GET("/getUserDetails", AuthController.Authentication(), userController.GetUserDetails())
		userGroup.GET("/updateUserDetails", userController.UpdateUserDetails())

	}
}
