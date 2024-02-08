package routes

import (
	AuthController "github.com/foodmngtapp/food-management-apps/src/common/middleware"
	userController "github.com/foodmngtapp/food-management-apps/src/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/v1/users")
	{

		userGroup.POST("/signup", userController.Signup())
		userGroup.POST("/login", userController.Login())

		userGroup.GET("/getAllUsers", userController.GetAllUser())
		userGroup.GET("/getUserDetails", AuthController.Authentication(), userController.GetUserDetails())
		userGroup.PUT("/updateUserDetails", AuthController.Authentication(), userController.UpdateUserDetails())

		userGroup.POST("/resetPassword", AuthController.Authentication(), userController.ResetPassword())
		userGroup.POST("/forgetPassword", userController.ForgetPassword())

		userGroup.POST("/deleteAccount", AuthController.Authentication(), userController.DeleteAccount())

	}
}
