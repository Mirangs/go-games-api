package main

import (
	"github.com/Mirangs/bm-go-test-task/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/Mirangs/bm-go-test-task/helpers"
)

func main() {
	err := godotenv.Load()
	gin.SetMode(os.Getenv("GIN_MODE"))
	if err != nil {
		log.Fatal("Error loading .env file " + err.Error())
		return
	}
	mongoClient, err := helpers.ConnectToMongo()
	if err != nil {
		log.Fatal(err)
		return
	}
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

	router.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
}
