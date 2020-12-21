package main

import (
	"context"
	"encoding/json"
	"github.com/Mirangs/bm-go-test-task/collections"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/Mirangs/bm-go-test-task/helpers"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserJSON struct {
	Email     string `json:"email" binding:"required" validate:"email"`
	LastName  string `json:"last_name" binding:"required"`
	Country   string `json:"country" binding:"required"`
	City      string `json:"city" binding:"required"`
	Gender    string `json:"gender" binding:"required"`
	BirthDate string `json:"birth_date" binding:"required"`
}

type UserGameJSON struct {
	PointsGained int    `json:"points_gained,string"`
	WinStatus    int8   `json:"win_status,string"`
	GameType     int8   `json:"game_type,string"`
	Created      string `json:"created"`
}

type usersJSONRes struct {
	Objects []UserJSON `json:"objects"`
}

type userGameJSONRes struct {
	Objects []UserGameJSON `json:"objects"`
}

func parseUsersJSON() (users []UserJSON) {
	jsonFile, err := os.Open("./data/users_go.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var fileRes usersJSONRes
	err = json.Unmarshal(byteValue, &fileRes)
	if err != nil {
		log.Error(err)
	}
	return fileRes.Objects
}
func parseUserGamesJSON() (usersGames []UserGameJSON) {
	jsonFile, err := os.Open("./data/games.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var fileRes userGameJSONRes
	json.Unmarshal(byteValue, &fileRes)
	return fileRes.Objects
}
func createIndexes(client *mongo.Client) (err error) {
	_, err = collections.Collections{Client: client}.Users().Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"email": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = collections.Collections{Client: client}.UserGames().Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"user_id": 1,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	return nil
}
func insertUsers(client *mongo.Client) (insertedIds []interface{}) {
	users := parseUsersJSON()
	userCollection := client.Database("Leads").Collection("users")
	log.Debug("Started inserting users...")
	usersMongo := make([]interface{}, 0)
	for _, user := range users {
		birthDate, err := time.Parse("Monday, January 2, 2006 3:04 PM", user.BirthDate)
		if err != nil {
			log.Error(err)
		}
		validYear := rand.Intn(2015-1980) + 1980
		validBirthDate := time.Date(validYear, birthDate.Month(), birthDate.Day(), birthDate.Hour(), birthDate.Minute(), birthDate.Second(), birthDate.Nanosecond(), birthDate.Location())
		newUser := bson.M{
			"email":      user.Email,
			"last_name":  user.LastName,
			"country":    user.Country,
			"city":       user.City,
			"gender":     user.Gender,
			"birth_date": primitive.NewDateTimeFromTime(validBirthDate),
		}
		usersMongo = append(usersMongo, newUser)
	}
	insertRes, _ := userCollection.InsertMany(context.TODO(), usersMongo)
	log.WithFields(log.Fields{"insertedUsers": len(insertRes.InsertedIDs)}).Debug("Finished inserting users")
	return insertRes.InsertedIDs
}
func insertUserGames(client *mongo.Client, userIds []interface{}) {
	userGames := parseUserGamesJSON()
	userGamesCollection := client.Database("Leads").Collection("user_games")
	redisClient := helpers.ConnectToRedis()

	log.Debug("Started inserting user games...")
	var foundUserIds = make([]primitive.ObjectID, 0, len(userIds))
	for _, userID := range userIds {
		idToObjectID, ok := userID.(primitive.ObjectID)
		if !ok {
			log.Warn("Cannot cast userId to ObjectID")
			continue
		}
		foundUserIds = append(foundUserIds, idToObjectID)
	}

	for _, foundUserID := range foundUserIds {
		var randGames = make([]interface{}, 0)

		rand.Seed(time.Now().Unix())
		for i := 0; i < 5000; i++ {
			randGame := userGames[rand.Intn(len(userGames))]
			created, err := time.Parse("1/2/2006 3:04 PM", randGame.Created)
			if err != nil {
				log.Error(err)
				continue
			}
			var newUserGame = bson.M{
				"_id":           primitive.NewObjectID(),
				"points_gained": randGame.PointsGained,
				"win_status":    randGame.WinStatus,
				"game_type":     randGame.GameType,
				"user_id":       foundUserID,
				"created":       created,
			}
			randGames = append(randGames, newUserGame)
		}

		_, err := userGamesCollection.InsertMany(context.TODO(), randGames)
		if err != nil {
			log.Error(err)
		}
		redisClient.ZAdd(context.TODO(), "z:USER_COUNT_GAMES", &redis.Z{Score: 5000, Member: foundUserID.Hex()})
	}
	log.Debug("Inserted user games")
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.Info("Started DB initialization...")
	client, err := helpers.ConnectToMongo()
	if err != nil {
		log.Fatal("Error connecting to Mongo " + err.Error())
		return
	}
	err = createIndexes(client)
	if err != nil {
		log.Fatal(err)
		return
	}

	insertedIds := insertUsers(client)
	insertUserGames(client, insertedIds)
}
