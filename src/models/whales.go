package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Whale struct {
	Id           primitive.ObjectID `bson:"_id"`
	Address      string             `bson:"address" validate:"required"`
	TokenAddress string             `bson:"token_address" validate:"required"`
}

type Whales []Whale
