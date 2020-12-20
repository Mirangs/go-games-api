package services

import (
	"context"
	"encoding/json"
	"github.com/Mirangs/bm-go-test-task/collections"
	"github.com/Mirangs/bm-go-test-task/helpers"
	. "github.com/Mirangs/bm-go-test-task/types"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UsersServices struct {
	Client *mongo.Client
}

func (s UsersServices) GetUsers(limit, skip int64, filters string) []User {
	var mongoFilters bson.M
	err := json.Unmarshal([]byte(filters), &mongoFilters)
	if err != nil {
		log.Error(err)
	}
	users := make([]User, 0)
	if limit == 0 {
		limit = 10
	}

	cur, err := s.Client.Database("Leads").Collection("users").Find(context.TODO(), mongoFilters, options.Find().SetLimit(limit).SetSkip(skip))
	defer cur.Close(context.TODO())
	if err != nil {
		log.Fatal(err)
		return users
	}

	for cur.Next(context.TODO()) {
		user := User{}
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
			continue
		}

		users = append(users, user)
	}

	return users
}
func (s UsersServices) GetUserById(id string) User {
	var user User
	userIdMongo, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return user
	}
	err = s.Client.Database("Leads").Collection("users").FindOne(context.TODO(), bson.M{"_id": userIdMongo}).Decode(&user)
	if err != nil {
		log.Error(err)
	}
	return user
}
func (s UsersServices) CreateUser(userData User) interface{} {
	userData.ID = primitive.NewObjectID()
	res, err := s.Client.Database("Leads").Collection("users").InsertOne(context.TODO(), userData)
	if err != nil {
		log.Error(err)
		return nil
	}
	return res.InsertedID
}
func (s UsersServices) UpdateUser(id string, userData User) interface{} {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return ""
	}
	var updatedUser User
	s.Client.Database("Leads").Collection("users").FindOneAndUpdate(context.TODO(), bson.M{"_id": userId}, bson.M{"$set": userData}).Decode(&updatedUser)
	return updatedUser.ID
}

// TODO: add redis cache cleaning
func (s UsersServices) DeleteUser(id string) interface{} {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return ""
	}
	var deletedUser User
	s.Client.Database("Leads").Collection("users").FindOneAndDelete(context.TODO(), bson.M{"_id": userId}).Decode(&deletedUser)
	return deletedUser.ID
}
func (s UsersServices) GetUsersRating(limit, skip int64) []UserRatingRes {
	redisClient := helpers.ConnectToRedis()
	usersRatingIds, err := redisClient.ZRevRangeByScoreWithScores(context.TODO(), "z:USER_COUNT_GAMES", &redis.ZRangeBy{
		Count:  limit,
		Min:    "-inf",
		Max:    "+inf",
		Offset: skip,
	}).Result()
	if err != nil {
		log.Error(err)
		return nil
	}

	mongoUsers := make([]User, 0)
	mongoUserRatingIds := make([]primitive.ObjectID, 0)
	for _, userRatingId := range usersRatingIds {
		res, err := primitive.ObjectIDFromHex(userRatingId.Member.(string))
		if err != nil {
			log.Error(err)
			return nil
		}
		mongoUserRatingIds = append(mongoUserRatingIds, res)
	}

	cur, err := s.Client.Database("Leads").Collection("users").Find(context.TODO(), bson.M{"_id": bson.M{"$in": mongoUserRatingIds}})
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var user User
		if err := cur.Decode(&user); err != nil {
			log.Error(err)
			continue
		}
		mongoUsers = append(mongoUsers, user)
	}

	mongoUsersSorted := make([]UserRatingRes, 0)
	for _, userRatingId := range usersRatingIds {
		for _, mongoUser := range mongoUsers {
			if mongoUser.ID.Hex() == userRatingId.Member.(string) {
				mongoUsersSorted = append(mongoUsersSorted, UserRatingRes{User: mongoUser, CountGames: int64(userRatingId.Score)})
			}
		}
	}

	return mongoUsersSorted
}
func (s UsersServices) GetGamesStatistics(userId string, startDate, endDate time.Time) (data []bson.M, err error) {
	userIdParsed, err := primitive.ObjectIDFromHex(userId)
	log.Info(userIdParsed)
	if err != nil {
		return
	}
	dateProjectStage := bson.M{
		"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$created"},
	}
	groupedByDayStages := []bson.D{
		{{"$project", bson.M{
			"date":    dateProjectStage,
			"created": true,
		}}},
		{{"$group", bson.M{
			"_id":          "$date",
			"games_played": bson.M{"$sum": 1},
		}}},
		{{"$project", bson.M{
			"date":         "$_id",
			"games_played": true,
			"_id":          false,
		}}},
	}
	withGameTypeStages := []bson.D{
		{{"$project", bson.M{
			"date":      dateProjectStage,
			"created":   true,
			"game_type": true,
		}}},
		{{"$group", bson.M{
			"_id": bson.M{
				"date":      "$date",
				"game_type": "$game_type",
			},
			"games_played": bson.M{"$sum": 1},
		}}},
		{{"$project", bson.M{
			"date":         "$_id.date",
			"game_type":    "$_id.game_type",
			"games_played": true,
			"_id":          false,
		}}},
	}
	pipeline := []bson.D{
		{{"$match", bson.M{
			"user_id": userIdParsed,
			"created": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		{{"$facet", bson.M{
			"group_by_day":   groupedByDayStages,
			"with_game_type": withGameTypeStages,
		}}},
	}

	cur, err := collections.Collections{Client: s.Client}.UserGames().Aggregate(context.TODO(), pipeline)
	if err != nil {
		return
	}
	var stats []bson.M
	if err = cur.All(context.TODO(), &stats); err != nil {
		return
	}

	return stats, nil
}
