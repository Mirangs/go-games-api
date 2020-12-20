package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Routes struct {
	Client *mongo.Client
	Router *gin.Engine
}
