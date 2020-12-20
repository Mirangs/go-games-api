package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Email     string             `bson:"email" json:"email" binding:"required" validate:"email"`
	LastName  string             `bson:"last_name" json:"last_name" binding:"required"`
	Country   string             `bson:"country" json:"country" binding:"required"`
	City      string             `bson:"city" json:"city" binding:"required"`
	Gender    string             `bson:"gender" json:"gender" binding:"required"`
	BirthDate time.Time          `bson:"birth_date" json:"birth_date" binding:"required"`
}

type UserRatingRes struct {
	User       User  `json:"user"`
	CountGames int64 `json:"count_games"`
}
