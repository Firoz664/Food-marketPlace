package routes

import (
	AuthController "github.com/foodmngtapp/food-management-apps/src/common/middleware"
	restaruntController "github.com/foodmngtapp/food-management-apps/src/controllers"

	"github.com/gin-gonic/gin"
)

func RestaruntRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/v1/restarunt")

	{

		userGroup.GET("/getAllRestarunt", restaruntController.GetAllRestaurants())
		userGroup.GET("/getRestarunt", restaruntController.GetRestaurant())

		userGroup.PUT("/updateRestarunt", AuthController.Authentication(), restaruntController.UpdateRestaurant())
		userGroup.DELETE("/deleteRestarunt", AuthController.Authentication(), restaruntController.DeleteRestaurant())

	}
}
