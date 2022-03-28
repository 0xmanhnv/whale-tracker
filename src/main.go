package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"whale-tracker/src/database"
	"whale-tracker/src/models"
	"whale-tracker/src/services"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func LoadEnv() {
	pwd, _ := os.Getwd()

	// load .env file
	err := godotenv.Load(path.Join(pwd, ".env"))
	fmt.Println(path.Join(pwd, ".env"))

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func bootstrap() {
	LoadEnv() // init load env
	database.CreateDBInstance()
}

func main() {
	bootstrap()

	// fmt.Print(os.Getenv("APP"))

	// database.CreateDBInstance()

	// logs := handles.LoadLogs(97, "0xE3233fdb23F1c27aB37Bd66A19a1f1762fCf5f3F", 16446905, 16446905+5000)

	// fmt.Println(logs)

	// handles.LogHandle(logs)

	holder := models.Holder{
		Id:           primitive.NewObjectID(),
		Address:      "0x58c34146316a9a60BFA5dA1d7F451e46BDd51215",
		TokenAddress: "0xE3233fdb23F1c27aB37Bd66A19a1f1762fCf5f3F",
	}

	services.CreateHolder(holder)
}
