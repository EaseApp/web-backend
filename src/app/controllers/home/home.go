package home

import (
	"fmt"
	"github.com/EaseApp/web-backend/src/app/dao"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

// DeleteApplicationHandler Destroy endpoint
func IndexHandler(w http.ResponseWriter, req *http.Request) {
	render(w, "index.html")
}

func render(w http.ResponseWriter, tmpl string) {
	tmpl = fmt.Sprintf("src/app/views/%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, "hello")
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func DBCountHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	db := vars["db"]
	// Move RecordCount out of user DAO and into generic DAO
	fmt.Fprintf(w, "Db (%v) has (%v) objects.", db, dao.RecordCount(db))
}
