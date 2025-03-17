package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"rest_api/database"
	"rest_api/routes"
)

func main() {
	database.ConnectDB()

	r := gin.Default()

	routes.UserRoutes(r)

	fmt.Println("server running on port 8000")
	r.Run(":8000")
}
