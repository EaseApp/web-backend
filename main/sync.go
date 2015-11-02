package main

import (

	"log"
	"github.com/codegangsta/negroni"
	"github.com/EaseApp/web-backend/src/sync"

)

// This function runs the main webserver for the sync service.
// go run main/sync.go
func main() {
	// Make websocket
	log.Println("Starting sync server")

	router := sync.NewSyncServer()

	// Make web server
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8000")
}
