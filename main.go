package main

import (
	"github.com/Mirangs/bm-go-test-task/routes"
	"github.com/gin-gonic/gin"

	"github.com/Mirangs/bm-go-test-task/helpers"
)

func main() {
	mongoClient := helpers.ConnectToMongo()
	router := gin.Default()

	routes := routes.Routes{
		Client: mongoClient,
		Router: router,
	}
	routes.Users()
	routes.UserGames()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found"})
	})
	router.Run("localhost:3000")
}
