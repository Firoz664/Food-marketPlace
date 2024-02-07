package routes

import (
	invoiceController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(router *gin.Engine) {
	invoiceGroup := router.Group("/api/v1/invoice")
	{
		// Directly use GetInvoice as a handler without invoking it.
		invoiceGroup.GET("/getAllInvoice", invoiceController.GetInvoice())
	}
}
