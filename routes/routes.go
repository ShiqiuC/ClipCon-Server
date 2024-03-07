package routes

import (
	"ClipCon-Server/controllers"
	"ClipCon-Server/database"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	database.Init()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/clipboard", controllers.CreateClipboardItem)
		v1.GET("/clipboard", controllers.GetClipboardItems)
	}

	return router
}
