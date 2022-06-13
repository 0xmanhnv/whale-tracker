package services

import (
	"context"
	"fmt"
	"time"
	"whale-tracker/src/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateWhale(data interface{}) (*mongo.InsertOneResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var whaleCollection *mongo.Collection = database.GetCollection(database.DB, "whales")

	result, err := whaleCollection.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindOnenWhale(tokenAddress string, address string) bson.M {

	var whale bson.M

	var whaleCollection *mongo.Collection = database.GetCollection(database.DB, "whales")

	if err := whaleCollection.FindOne(context.TODO(), bson.M{
		"address":       address,
		"token_address": tokenAddress,
	}).Decode(&whale); err != nil {
		fmt.Println(address + "Chua co trong csdl")
		return nil
	}

	return whale
}

// func UpdatePointBlockNumber(number int) bson.M {
// 	var point bson.M
// 	var pointCollection *mongo.Collection = database.GetCollection(database.DB, "points")

// 	if err := pointCollection.FindOne(context.TODO(), bson.M{}).Decode(&point); err != nil {
// 		return nil
// 	}

// 	if
// }
