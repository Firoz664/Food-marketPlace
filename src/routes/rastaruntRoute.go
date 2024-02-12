package routes

import (
	AuthController "github.com/foodmngtapp/food-management-apps/src/common/middleware"
	restaruntController "github.com/foodmngtapp/food-management-apps/src/controllers"

	"github.com/gin-gonic/gin"
)

func RestaruntRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/v1/restarunt")

	{

		userGroup.POST("/addRestarunt", restaruntController.CreateRestaurant())
		userGroup.POST("/login", restaruntController.Login())

		userGroup.GET("/getAllRestarunt", restaruntController.GetAllRestarunt())
		userGroup.GET("/getRestarunt", AuthController.Authentication(), restaruntController.GetRestaurant())
		userGroup.PUT("/updateRestarunt", AuthController.Authentication(), restaruntController.UpdateRestaurant())
		userGroup.DELETE("/deleteRestarunt", AuthController.Authentication(), restaruntController.DeleteRestaurant())

	}
}
