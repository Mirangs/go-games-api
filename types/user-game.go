package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Created time.Time

type UserGame struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	PointsGained int                `bson:"points_gained,omitempty" json:"points_gained,omitempty"`
	WinStatus    int8               `bson:"win_status,omitempty" json:"win_status,omitempty"`
	GameType     int8               `bson:"game_type,omitempty" json:"game_type,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Created      time.Time          `bson:"created,omitempty" json:"created,omitempty"`
}
