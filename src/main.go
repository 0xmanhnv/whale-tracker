package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"whale-tracker/src/handles"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	pwd, _ := os.Getwd()

	// load .env file
	err := godotenv.Load(path.Join(pwd, ".env"))

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func bootstrap() {
	LoadEnv() // init load env
}

func main() {
	bootstrap()

	fmt.Print(os.Getenv("APP"))

	logs := handles.LoadLogs(97, 16446905)

	handles.LogHandle(logs, "0xE3233fdb23F1c27aB37Bd66A19a1f1762fCf5f3F")
}
