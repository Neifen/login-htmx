package main

import (
	"log"

	"github.com/neifen/htmx-login/api/server"
)

func main() {
	store, err := server.NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	api := server.NewAPIHandler(":1323", store)
	api.Run()
}
