package routes

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	otpController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var otpCollection *mongo.Collection = database.OpenCollection(database.Client, "otp")

func OtpRoutes(router *gin.Engine) {
	otpGroup := router.Group("/api/v1/otp")
	{
		otpGroup.POST("/createAndSendOtp", otpController.CreateOTP(otpCollection))
		otpGroup.POST("/verifyOtp", otpController.VerifyOtp())

	}
}
