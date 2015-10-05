package controller

import(
  "fmt"
	"log"
  "net/http"
	// "strconv"
  "github.com/gorilla/mux"
	r "github.com/dancannon/gorethink"
  dao "github.com/EaseApp/web-backend/src/app/dao"
	// "encoding/json"
)

var session *r.Session

// Initialize connection and set global session variable
func Init(s *r.Session) {
	if s == nil{
		log.Fatal("Generic DAO initialize failure")
	}
	session = s
}

func QueryApplicationHandler(w http.ResponseWriter, req *http.Request){
  vars := mux.Vars(req)
	client := vars["client"]
  application := vars["application"]
  result := dao.QueryApplication(client, application, "{}")

  fmt.Fprintf(w, "%v", result)
}

func CreateApplicationHandler(w http.ResponseWriter, req *http.Request){
  vars := mux.Vars(req)
	client := vars["client"]
  application := vars["application"]
  result := dao.CreateApplication(client, application)

  fmt.Fprintf(w, "%v", result)
}
