package main

import (
	"log"

	"github.com/EaseApp/web-backend/src/db"
	"github.com/EaseApp/web-backend/src/sync"
	"github.com/codegangsta/negroni"
)

// This function runs the main webserver for the sync service.
// go run main/sync.go
func main() {
	// Make websocket
	log.Println("Starting sync server")

	// TODO: Use command line flag credentials.
	client, err := db.NewClient("localhost:28015")
	if err != nil {
		log.Fatal("Couldn't initialize database: ", err.Error())
	}
	defer client.Close()

	router := sync.NewServer(client)

	// Make web server
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8000")
}
