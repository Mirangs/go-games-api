package collections

import "go.mongodb.org/mongo-driver/mongo"

type Collections struct {
	Client *mongo.Client
}

func (c Collections) Users() *mongo.Collection {
	return c.Client.Database("Leads").Collection("users")
}

func (c Collections) UserGames() *mongo.Collection {
	return c.Client.Database("Leads").Collection("user_games")
}
