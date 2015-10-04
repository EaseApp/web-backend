package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/EaseApp/web-backend/src/db"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/EaseApp/web-backend/src/app/models/user"
	// r "github.com/dancannon/gorethink"
	// "strconv"
)

func FetchAllHandler(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)

	// username := vars["user"]
	table := vars["db"]
	// n, _ := strconv.Atoi(db)

	fmt.Fprintf(w, "Table: (%v). All records: (%v)", table, user.FetchAll(table))
}


func DBCountHandler(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	db := vars["db"]
	// n, _ := strconv.Atoi(db)

	// Move RecordCount out of user DAO and into generic DAO
	fmt.Fprintf(w, "Db (%v) has (%v) objects.", db, user.RecordCount(db))
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

	user.Init()

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	router.HandleFunc("/count/{db}", DBCountHandler);
	router.HandleFunc("/{db}", FetchAllHandler)
	// router.HandleFunc("/{user}/{db}", FetchAllHandler)



	defer client.Close()

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
