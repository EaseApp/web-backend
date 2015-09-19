package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EaseApp/web-backend/src/db"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting server...")

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	// TODO: Use command line flag credentials.
	client, err := db.NewClient("localhost:28015")
	if err != nil {
		log.Fatal("Couldn't initialize database: ", err.Error())
	}
	defer client.Close()

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3001")
}
