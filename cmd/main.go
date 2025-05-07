package main

import (
	"fmt"

	"github.com/neifen/htmx-login/api/server"
)

func main() {
	store, err := server.NewPostGresStore()
	if err != nil {
		fmt.Print("didnt work to make new postgres, delete log tho")
		// log.Fatal(err)
	}

	api := server.NewAPIHandler(":1323", store)
	api.Run()
}
