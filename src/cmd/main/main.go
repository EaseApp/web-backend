package main

import (
	"github.com/EaseApp/web-backend/src/config/db"

	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {

	err := db.Init()
	if err != nil {
		log.Fatalln("Couldn't connect to database. Quitting...")
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3001")
}
