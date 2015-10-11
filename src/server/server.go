package server

import (
	"github.com/EaseApp/web-backend/config"
	"github.com/EaseApp/web-backend/src/app/dao"
	"github.com/EaseApp/web-backend/src/db"
	"github.com/codegangsta/negroni"
	"log"
)

// StartServer starts the server
func StartServer() {
	log.Println("Starting server...")

	// TODO: Use command line flag credentials.
	client, err := db.NewClient("localhost:28015")
	if err != nil {
		log.Fatal("Couldn't initialize database: ", err.Error())
	}
	client.SetUpDatabase()
	// CreatePubSubServer()

	err = client.SetUpDatabase()
	if err != nil {
		log.Println(err)
	}
	dao.InitGeneric(client.Session)
	defer client.Close()

	n := negroni.Classic()
	n.UseHandler(config.CreateRouting())
	n.Run(":3000")
}
