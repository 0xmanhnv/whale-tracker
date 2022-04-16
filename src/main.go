package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"
	"whale-tracker/src/blockchain"
	"whale-tracker/src/database"
	"whale-tracker/src/handles"
	"whale-tracker/src/models"
	"whale-tracker/src/responses"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	// crons.Init()
}

func InitBSCWhaleTracker() {
	// init client
	client := blockchain.GetClient(56)
	defer client.Close()

	/*
		load transfer logs of tokens
	*/
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	topics := [][]common.Hash{
		{logTransferSigHash},
	}

	// tokens
	addresses := []common.Address{
		common.HexToAddress("0x0eb3a705fc54725037cc9e008bdede697f62f335"),
		common.HexToAddress("0x3244b3b6030f374bafa5f8f80ec2f06aaf104b64"),
	}

	fromBlock := 16618405
	toBlock := fromBlock + 1000

	logs := handles.LoadLogs(client, addresses, int64(fromBlock), int64(toBlock), topics)
	// end load logs

	// fmt.Println(logs)

	handles.LogHandleToWhale(client, logs, addresses)

	// holder := models.Holder{
	// 	Id:           primitive.NewObjectID(),
	// 	Address:      "0x58c34146316a9a60BFA5dA1d7F451e46BDd51215",
	// 	TokenAddress: "0xE3233fdb23F1c27aB37Bd66A19a1f1762fCf5f3F",
	// }

	// services.CreateHolder(holder)

	// services.GetTokenHolderWithCovalenthqAPI("0xE3233fdb23F1c27aB37Bd66A19a1f1762fCf5f3F", "56")
}

func main() {
	bootstrap()
	// InitBSCWhaleTracker()

	router := mux.NewRouter()

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		json.NewEncoder(rw).Encode(map[string]string{"data": "Hello from Mux & mongoDB"})
	}).Methods("GET")

	router.HandleFunc("/whales", func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var whales models.Whales

		var whaleCollection *mongo.Collection = database.GetCollection(database.DB, "whales")
		results, err := whaleCollection.Find(ctx, bson.M{})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)

		for results.Next(ctx) {
			var singleWhale models.Whale
			if err = results.Decode(&singleWhale); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
			}
			whales = append(whales, singleWhale)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		response := responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": whales}}
		json.NewEncoder(rw).Encode(response)
	}).Methods("GET")

	log.Fatal(http.ListenAndServe(":6000", router))
}
