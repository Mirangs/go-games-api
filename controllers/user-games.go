package controllers

import (
	"encoding/json"
	"github.com/Mirangs/bm-go-test-task/services"
	. "github.com/Mirangs/bm-go-test-task/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

type UserGamesController struct {
	Client *mongo.Client
}

func (h UserGamesController) GetUserGames(c *gin.Context) {
	filters := c.DefaultQuery("filters", "{}")
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"error": "Invalid limit",
		})
		return
	}
	skip, err := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"error": "Invalid skip",
		})
		return
	}
	userGames := services.UserGamesServices{Client: h.Client}.GetUserGames(limit, skip, filters)
	c.JSON(200, gin.H{
		"user-games": userGames,
	})
}
func (h UserGamesController) GetUserGameById(c *gin.Context) {
	userGameId := c.Param("id")
	userGame := services.UserGamesServices{Client: h.Client}.GetUserGameById(userGameId)
	c.JSON(200, gin.H{
		"user-game": userGame,
	})
}
func (h UserGamesController) CreateUserGame(c *gin.Context) {
	var body UserGame
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	res := services.UserGamesServices{Client: h.Client}.CreateUserGame(body)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(201, gin.H{"message": string(jsonRes)})
}
func (h UserGamesController) UpdateUserGame(c *gin.Context) {
	var body User
	userId := c.Param("id")
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res := services.UsersServices{Client: h.Client}.UpdateUser(userId, body)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(201, gin.H{"message": string(jsonRes)})
}
func (h UserGamesController) DeleteUserGame(c *gin.Context) {
	userId := c.Param("id")
	res := services.UserGamesServices{Client: h.Client}.DeleteUserGame(userId)
	if res == "" {
		c.JSON(400, gin.H{"error": "Invalid id"})
	}

	c.JSON(200, gin.H{"message": res})
}
