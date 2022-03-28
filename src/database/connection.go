package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// database
var DB *mongo.Client

func CreateDBInstance() {
	fmt.Println(os.Getenv("MONGO_URI"))
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	DB = client

	fmt.Println("Connected to MongoDB")
}

//getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(os.Getenv("DB_NAME")).Collection(collectionName)
	return collection
}
