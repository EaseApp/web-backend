package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/EaseApp/web-backend/src/db"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/EaseApp/web-backend/src/app/models/user"
	r "github.com/dancannon/gorethink"
)

func DbsHandler(w http.ResponseWriter, req *http.Request, session *r.Session){
	vars := mux.Vars(req)

	username := vars["user"]
	db := vars["db"]



	fmt.Fprintf(w, "Welcome to db. User: (%v). Db: (%v). Finding user: (%v)", username, db, user.FindUser(username, session))

}

func main() {
	log.Println("Starting server...")

	// TODO: Use command line flag credentials.
	client, err := db.NewClient("localhost:28015")
	if err != nil {
		log.Fatal("Couldn't initialize database: ", err.Error())
	}
	db.CreateEaseDb(client)
	db.CreateUserDb(client)

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	router.HandleFunc("/{user}/{db}", func(w http.ResponseWriter, req *http.Request) {
		DbsHandler(w, req, client.Session)
	})



	defer client.Close()

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
