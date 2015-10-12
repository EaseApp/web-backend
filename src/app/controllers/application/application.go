package application

import (
	"fmt"
	"log"
	"net/http"
	// "strconv"

	"io"

	"github.com/EaseApp/web-backend/src/app/dao"
	"github.com/EaseApp/web-backend/src/app/helper"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
	// "github.com/EaseApp/web-backend/src/app/utils"
)

var session *r.Session

// Init connection and set global session variable
func Init(s *r.Session) {
	if s == nil {
		log.Fatal("Generic DAO initialize failure")
	}
	session = s
}

// QueryApplicationHandler  Search application given a JSON object in the body
func QueryApplicationHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	client := vars["client"]
	application := vars["application"]

	obj, err := helper.DecodeIOStreamToJSON(req.Body)
	if err == io.EOF {
		emptyMap := make(map[string]interface{})
		fmt.Fprintf(w, "%v", dao.QueryApplication(client, application, emptyMap))
	} else if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Decoding error.")
	} else {
		log.Println(obj)
		result := dao.QueryApplication(client, application, obj)
		fmt.Fprintf(w, "%v", result)
	}
}

// CreateApplicationHandler Creates new application
func CreateApplicationHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	client := vars["client"]
	application := vars["application"]
	result := dao.CreateApplication(client, application)

	fmt.Fprintf(w, "%v", result)
}

// UpdateApplicationHandler Update endpoint
func UpdateApplicationHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	client := vars["client"]
	application := vars["application"]
	id := vars["id"]
	log.Println(id)

	obj, err := helper.DecodeIOStreamToJSON(req.Body)
	if err == io.EOF {
		fmt.Fprintf(w, "No update or invalid object provided")
	}
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Decoding error.")
	} else {
		log.Println(obj)
		result := dao.UpdateApplication(client, application, id, obj)
		fmt.Fprintf(w, "%v", result)
	}
	// Add pub-sub
}

// DeleteApplicationHandler Destroy endpoint
func DeleteApplicationHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	client := vars["client"]
	application := vars["application"]
	id := vars["id"]

	result := dao.DeleteApplication(client, application, id)
	fmt.Fprintf(w, "%v", result)
}

// PubSubApplicationHandler sets up pub sub for a connection
func PubSubApplicationHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Sup")
}
