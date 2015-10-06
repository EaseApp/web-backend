package controller

import(
  "fmt"
	"log"
  "net/http"
	// "strconv"
  "github.com/gorilla/mux"
	r "github.com/dancannon/gorethink"
  dao "github.com/EaseApp/web-backend/src/app/dao"
	"encoding/json"
  "io"
)

var session *r.Session

// Initialize connection and set global session variable
func Init(s *r.Session) {
	if s == nil{
		log.Fatal("Generic DAO initialize failure")
	}
	session = s
}

func decodeIOStreamToJSON(body io.Reader)(map[string]interface{}, error){
  decoder := json.NewDecoder(body)
  var m map[string]interface{}
  err := decoder.Decode(&m)

  if err != nil{
    return nil, err
  }
  return m, nil
}

// Search application given a JSON object in the body
func QueryApplicationHandler(w http.ResponseWriter, req *http.Request){
  vars := mux.Vars(req)
	client := vars["client"]
  application := vars["application"]

  obj, err := decodeIOStreamToJSON(req.Body)
  if err == io.EOF{
    emptyMap := make(map[string]interface{})
    fmt.Fprintf(w, "%v", dao.QueryApplication(client, application, emptyMap))
  } else if err != nil{
    log.Println(err)
    fmt.Fprintf(w, "Decoding error.")
  } else{
    log.Println(obj)
    result := dao.QueryApplication(client, application, obj)
    fmt.Fprintf(w, "%v", result)
  }
}

// Creates new application
func CreateApplicationHandler(w http.ResponseWriter, req *http.Request){
  vars := mux.Vars(req)
	client := vars["client"]
  application := vars["application"]
  result := dao.CreateApplication(client, application)

  fmt.Fprintf(w, "%v", result)
}

// Update endpoint
func UpdateApplicationHandler(w http.ResponseWriter, req *http.Request){
  vars := mux.Vars(req)
	client := vars["client"]
  application := vars["application"]
  id := vars["id"]
  log.Println(id)

  obj, err := decodeIOStreamToJSON(req.Body)
  if err == io.EOF{
    fmt.Fprintf(w, "No update or invalid object provided")
  }
  if err != nil{
    log.Println(err)
    fmt.Fprintf(w, "Decoding error.")
  } else {
    log.Println(obj)
    result := dao.UpdateApplication(client, application, id, obj)
    fmt.Fprintf(w, "%v", result)
  }
  // Add pub-sub
}

// Destroy endpoint
// func DestroyApplicationHandler(w http.ResponseWriter, req *http.Request){
// }


func PubSubApplicationHandler(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w,"Sup")
}
