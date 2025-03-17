package routes

import (
	"rest_api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("/", controllers.GetUser)
		userGroup.POST("/", controllers.CreateUser)
		userGroup.GET("/:id", controllers.GetUserByID)
	}
}
