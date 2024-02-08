package routes

import (
	"github.com/gin-gonic/gin"
)

// Register all routes
func RegisterRoutes(router *gin.Engine) {
	FileUploadRoutes(router)
	FoodRoutes(router)
	InvoiceRoutes(router)
	MenuRoutes(router)
	NotesRoutes(router)
	NotificationsRoutes(router)
	OrderItemRoutes(router)
	OrderRoutes(router)
	OtpRoutes(router)
	PaymentRoutes(router)
	TableRoutes(router)
	TransactionRoutes(router)
	UserRoutes(router)

}
