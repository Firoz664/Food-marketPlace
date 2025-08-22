package routes

import (
	AuthController "github.com/foodmngtapp/food-management-apps/src/common/middleware"

	userController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/foodmngtapp/food-management-apps/src/config/database"

	"github.com/gin-gonic/gin"
)

// getUserCollection returns the user collection with lazy initialization
func getUserCollection() *mongo.Collection {
	client := database.DBInstance()
	return database.OpenCollection(client, "users") // Note: using "users" to match the controller
}

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/v1/users")

	{

		userGroup.POST("/signup", userController.Signup())
		userGroup.POST("/login", userController.Login())

		userGroup.GET("/getAllUsers", userController.GetAllUser())
		userGroup.GET("/getUserDetails", AuthController.Authentication(), userController.GetUserDetails(getUserCollection()))
		userGroup.PUT("/updateUserDetails", AuthController.Authentication(), userController.UpdateUserDetails(getUserCollection()))

		userGroup.POST("/resetPassword", AuthController.Authentication(), userController.ResetPassword(getUserCollection()))
		userGroup.POST("/forgetPassword", userController.ForgetPassword())

		userGroup.DELETE("/deleteAccount", AuthController.Authentication(), userController.DeleteAccount())

	}
}
