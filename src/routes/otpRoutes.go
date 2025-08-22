package routes

import (
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	otpController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func OtpRoutes(router *gin.Engine) {
	otpGroup := router.Group("/api/v1/otp")
	{
		otpGroup.POST("/createAndSendOtp", otpController.CreateOTP(getOtpCollection()))
		otpGroup.POST("/verifyOtp", otpController.VerifyOtp())

	}
}

// getOtpCollection returns the otp collection with lazy initialization
func getOtpCollection() *mongo.Collection {
	client := database.DBInstance()
	return database.OpenCollection(client, "otp")
}
