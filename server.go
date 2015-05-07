package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/easeapp/web-backend/config/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"html/template"
	"path"
	

)


func HomePage(w http.ResponseWriter, r *http.Request){

	// Hash example
	var hash map[string]string
	hash = make(map[string]string)
	hash["hello"] = "world"

	fp := path.Join("app", "templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, hash); err != nil {
  	http.Error(w, err.Error(), http.StatusInternalServerError)
  }

}

func main() {
	// Connection to RethinkDB
	err := db.Init()
	if err != nil {
		log.Fatalln("Couldn't connect to DB. Quitting...")
	}
	defer db.Close()

	// Routing configurations
	router := mux.NewRouter()

	router.HandleFunc("/", HomePage)

	router.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Current path: %s\n\n", req.URL)
	})

	// Start web server
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
