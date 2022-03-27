package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	pwd, _ := os.Getwd()

	// load .env file
	err := godotenv.Load(path.Join(pwd, "..", ".env"))

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
}
