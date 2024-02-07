package routes

import (
	notesController "github.com/foodmngtapp/food-management-apps/src/controllers"
	"github.com/gin-gonic/gin"
)

func NotesRoutes(router *gin.Engine) {
	noteGroup := router.Group("/api/v1/notes")
	{
		noteGroup.GET("/getAllNotes", notesController.GetNotes())
	}
}
