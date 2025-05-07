package main

import (
	"log"

	"github.com/neifen/htmx-login/api/server"
)

func main() {
	store, err := server.NewTestStore()
	if err != nil {
		log.Fatal(err)
	}

	api := server.NewAPIHandler(":1323", store)
	api.Run()
}
