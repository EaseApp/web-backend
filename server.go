package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/easeapp/web-backend/config/db"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

var session *r.Session

func init() {
    var err error
    session, err = r.Connect(r.ConnectOpts{
        Address:  "localhost:28015",
        Database: "test",
    })
    if err != nil {
        fmt.Println(err)
        return
    }
}

// Struct tags are used to map struct fields to fields in the database
type User struct {
	Id           string `gorethink:"id,omitempty"`
	Username     string
	PasswordHash string
	ApiToken     string
	CreatedAt    time.Time
}

func HomePage(w http.ResponseWriter, req *http.Request) {

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

  rows, _ := r.Table("users").Run(session)
  fmt.Println(rows.Next(&rows))
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
