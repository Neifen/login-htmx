package main

import (
	"errors"
	"log"

	"github.com/joho/godotenv"
	"github.com/neifen/htmx-login/api/server"
)

func main() {
	// 1. load env variable
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	store, err := server.NewPostGresStore()
	if err != nil {
		// need to be able to set up db, otherwise fail
		log.Fatal(err)
	}

	api := server.NewAPIHandler(":1323", store)
	api.Run()
}

func loadEnv() error {
	err := godotenv.Load(".env")
	if err == nil {
		return nil
	}

	err = godotenv.Load("/run/secrets/dot-env")
	if err != nil {
		return errors.New("error loading secret .env file")
	}

	return nil
}
