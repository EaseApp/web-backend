package main

import (
	"ease-backend/config/db"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	err := db.Init()
	if err != nil {
		log.Fatalln("Couldn't connect to DB. Quitting...")
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
