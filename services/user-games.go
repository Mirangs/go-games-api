package services

import (
	"context"
	"encoding/json"
	"github.com/Mirangs/bm-go-test-task/helpers"
	. "github.com/Mirangs/bm-go-test-task/types"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserGamesServices struct {
	Client *mongo.Client
}

func (s UserGamesServices) GetUserGames(limit, skip int64, filters string) []UserGame {
	var mongoFilters bson.M
	err := json.Unmarshal([]byte(filters), &mongoFilters)
	if err != nil {
		log.Error(err)
	}
	users := make([]UserGame, 0)
	if limit == 0 {
		limit = 10
	}

	cur, err := s.Client.Database("Leads").Collection("user_games").Find(context.TODO(), mongoFilters, options.Find().SetLimit(limit).SetSkip(skip))
	defer cur.Close(context.TODO())
	if err != nil {
		log.Fatal(err)
		return users
	}

	for cur.Next(context.TODO()) {
		user := UserGame{}
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
			continue
		}

		users = append(users, user)
	}

	return users
}
func (s UserGamesServices) GetUserGameById(id string) UserGame {
	var user UserGame
	userIdMongo, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return user
	}
	err = s.Client.Database("Leads").Collection("user_games").FindOne(context.TODO(), bson.M{"_id": userIdMongo}).Decode(&user)
	if err != nil {
		log.Error(err)
	}
	return user
}
func (s UserGamesServices) CreateUserGame(userGameData UserGame) interface{} {
	userGameData.ID = primitive.NewObjectID()
	userGameData.Created = time.Now()
	res, err := s.Client.Database("Leads").Collection("user_games").InsertOne(context.TODO(), userGameData)
	if err != nil {
		log.Error(err)
		return nil
	}
	redisClient := helpers.ConnectToRedis()
	redisClient.ZIncrBy(context.TODO(), "z:USER_COUNT_GAMES", 1, userGameData.UserID.Hex())
	return res.InsertedID
}
func (s UserGamesServices) UpdateUserGame(id string, userGameData UserGame) interface{} {
	userGameId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return ""
	}
	var updatedUserGame UserGame
	err = s.Client.Database("Leads").Collection("user_games").FindOneAndUpdate(context.TODO(), bson.M{"_id": userGameId}, bson.M{"$set": userGameData}).Decode(&updatedUserGame)
	if err != nil {
		log.Error(err)
		return ""
	}
	return updatedUserGame.ID
}

// TODO: add redis cache decrementing
func (s UserGamesServices) DeleteUserGame(id string) interface{} {
	userGameId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return ""
	}
	var deletedUserGame UserGame
	err = s.Client.Database("Leads").Collection("user_games").FindOneAndDelete(context.TODO(), bson.M{"_id": userGameId}).Decode(&deletedUserGame)
	if err != nil {
		log.Error(err)
		return ""
	}
	return deletedUserGame.ID
}
