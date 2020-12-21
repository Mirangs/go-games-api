package helpers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func ConnectToMongo() (client *mongo.Client, err error) {
	mongoURI := os.Getenv("MONGO_URI")
	client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
		return
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Debug("Connected to MongoDB: ", mongoURI)
	return client, nil
}
