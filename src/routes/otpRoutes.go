package routes

import (
	otpController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func OtpRoutes(router *gin.Engine) {
	otpGroup := router.Group("/api/v1/otp")
	{
		otpGroup.GET("/getOtp", otpController.GetOtps())
	}
}
