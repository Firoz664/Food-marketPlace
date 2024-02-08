package routes

import (
	fileUploadController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func FileUploadRoutes(router *gin.Engine) {
	fileUploadGroup := router.Group("/api/v1/files/")
	{
		// Directly use GetInvoice as a handler without invoking it.
		fileUploadGroup.POST("/singleFileUpload", fileUploadController.S3FileUploadFile())
		fileUploadGroup.POST("/multipleUploadFiles", fileUploadController.MultipleUploadFiles())

	}
}
