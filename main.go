package main

import (
	"log"

	"github.com/EaseApp/web-backend/src/db"
	"github.com/EaseApp/web-backend/src/server"

	"github.com/codegangsta/negroni"
)

func main() {
	log.Println("Starting server...")

	// TODO: Use command line flag credentials.
	client, err := db.NewClient("localhost:28015")
	if err != nil {
		log.Fatal("Couldn't initialize database: ", err.Error())
	}
	defer client.Close()

	router := server.CreateRoutingMux(client)

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3001")
}
