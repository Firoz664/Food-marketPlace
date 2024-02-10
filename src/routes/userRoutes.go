package routes

import (
	AuthController "github.com/foodmngtapp/food-management-apps/src/common/middleware"
	userController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/foodmngtapp/food-management-apps/src/config/database"

	"github.com/gin-gonic/gin"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/v1/users")

	{

		userGroup.POST("/signup", userController.Signup())
		userGroup.POST("/login", userController.Login())

		userGroup.GET("/getAllUsers", userController.GetAllUser())
		userGroup.GET("/getUserDetails", AuthController.Authentication(), userController.GetUserDetails(userCollection))
		userGroup.PUT("/updateUserDetails", AuthController.Authentication(), userController.UpdateUserDetails(userCollection))

		userGroup.POST("/resetPassword", AuthController.Authentication(), userController.ResetPassword(userCollection))
		userGroup.POST("/forgetPassword", userController.ForgetPassword())

		userGroup.POST("/deleteAccount", AuthController.Authentication(), userController.DeleteAccount())

	}
}
