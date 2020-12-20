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
	"time"
)

type UsersController struct {
	Client *mongo.Client
}

func getDateErrorMessage(field string) string {
	return "Invalid " + field + " format, please use dd-mm-yyyy format"
}

func (h UsersController) GetUsers(c *gin.Context) {
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
	users := services.UsersServices{Client: h.Client}.GetUsers(limit, skip, filters)
	c.JSON(200, gin.H{
		"users": users,
	})
}
func (h UsersController) GetUserById(c *gin.Context) {
	userId := c.Param("id")
	user := services.UsersServices{Client: h.Client}.GetUserById(userId)
	code := 200
	c.JSON(code, gin.H{
		"user": user,
	})
}
func (h UsersController) CreateUser(c *gin.Context) {
	var body User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	res := services.UsersServices{Client: h.Client}.CreateUser(body)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(201, gin.H{"message": string(jsonRes)})
}
func (h UsersController) UpdateUser(c *gin.Context) {
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
func (h UsersController) DeleteUser(c *gin.Context) {
	userId := c.Param("id")
	res := services.UsersServices{Client: h.Client}.DeleteUser(userId)
	if res == "" {
		c.JSON(400, gin.H{"error": "Invalid id"})
	}

	c.JSON(200, gin.H{"message": res})
}
func (h UsersController) GetUsersRating(c *gin.Context) {
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
	usersRating := services.UsersServices{Client: h.Client}.GetUsersRating(limit, skip)
	c.JSON(200, gin.H{
		"rating": usersRating,
	})
}
func (h UsersController) GetGamesStatistics(c *gin.Context) {
	userId := c.Param("id")
	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")
	if startDate == "" || endDate == "" {
		c.JSON(400, gin.H{"error": "Please specify startDate and endDate query param"})
		return
	}
	dateLayout := "2-1-2006"

	parsedStartDate, err := time.Parse(dateLayout, startDate)
	if err != nil {
		c.JSON(400, gin.H{"error": getDateErrorMessage("startDate")})
		return
	}

	parsedEndDate, err := time.Parse(dateLayout, endDate)
	if err != nil {
		c.JSON(400, gin.H{"error": getDateErrorMessage("endDate")})
		return
	}

	if parsedStartDate.After(parsedEndDate) {
		c.JSON(400, gin.H{"error": "startDate should not be after endDate"})
		return
	}

	res, err := services.UsersServices{Client: h.Client}.GetGamesStatistics(userId, parsedStartDate, parsedEndDate)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	c.JSON(200, gin.H{"user_id": userId, "statistics": res})
}
