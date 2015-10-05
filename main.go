package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/EaseApp/web-backend/src/db"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	user "github.com/EaseApp/web-backend/src/app/controllers/user"

)



// func SignUpHandler(w http.ResponseWriter, req *http.Request){
// 	vars := mux.Vars(req)
// 	table := vars["db"]
//
// 	fmt.Fprintf(w, "Table: (%v). All records: (%v)", table, user.FetchAll(table))
// }


func NewStaticUserHandler(w http.ResponseWriter, req *http.Request){
	_ = user.InsertStaticUser()
	http.Redirect(w, req, "/users", http.StatusFound)
}

func DBCountHandler(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	db := vars["db"]

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
	router.HandleFunc("/{db}", user.FetchAllHandler)
	router.HandleFunc("/count/{db}", DBCountHandler)
	router.HandleFunc("/static/user/new", NewStaticUserHandler)

	router.HandleFunc("/users/sign_in", user.SignInHandler)
	router.HandleFunc("/users/sign_up", user.SignUpHandler)
	// router.HandleFunc("/{user}/{db}", FetchAllHandler)

	defer client.Close()

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
