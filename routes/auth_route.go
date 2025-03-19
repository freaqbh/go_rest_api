package routes

import (
	"rest_api/controllers"
	"rest_api/middlewares"

	"github.com/gin-gonic/gin"
)

// AuthRoutes - Rute autentikasi
func AuthRoutes(router *gin.Engine) {
	router.POST("/login", controllers.LoginHandler)

	protected := router.Group("/protected")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			username, _ := c.Get("username")
			c.JSON(200, gin.H{"message": "Welcome!", "username": username})
		})
	}
}
