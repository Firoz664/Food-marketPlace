package routes

import (
	"github.com/gin-gonic/gin"
)

// Register all routes
func RegisterRoutes(router *gin.Engine) {
	AdminRoutes(router)
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
	RestaruntRoutes(router)
	RiderRoutes(router)
	TableRoutes(router)
	TransactionRoutes(router)
	UserRoutes(router)

}
