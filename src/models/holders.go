package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Holder struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Address      string             `json:"address,omitempty" validate:"required,unique"`
	TokenAddress string             `json:"token_address,omitempty" validate:"required"`
}
