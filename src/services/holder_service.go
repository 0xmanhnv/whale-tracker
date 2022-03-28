package services

import (
	"context"
	"fmt"
	"time"
	"whale-tracker/src/database"
	"whale-tracker/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateHolder(holder models.Holder) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var oldHolder models.Holder

	var holderCollection *mongo.Collection = database.GetCollection(database.DB, "holders")

	err := holderCollection.FindOne(ctx, bson.M{"address": "0x58c34146316a9a60BFA5dA1d7F451e46BDd51215"}).Decode(&holder)
	if err != nil {
		fmt.Println("Error")
	}

	if oldHolder.Id == primitive.NilObjectID {
		result, _ := holderCollection.InsertOne(ctx, holder)
		fmt.Println(result)
	}
}
