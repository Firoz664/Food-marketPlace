package routes

import (
	noteficationController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func NotificationsRoutes(router *gin.Engine) {
	notificationsGroup := router.Group("/api/v1/notification")
	{
		notificationsGroup.GET("/getAllNotification", noteficationController.GetNotification())
	}
}
