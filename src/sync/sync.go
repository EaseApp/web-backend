package sync

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "github.com/codegangsta/negroni"

	"github.com/gorilla/websocket"

	"github.com/gorilla/mux"
)

var applications map[string][]Connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SyncServer struct {
	r *mux.Router
}

// ServeHTTP serves requests from the EaseServer's mux while allowing
// cross origin access.
func (s *SyncServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}

	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func createRouting() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/sub", subHandler)
	router.HandleFunc("/pub/{application}", pubHandler)
	return router
}

func NewSyncServer() *SyncServer {
	applications = make(map[string][]Connection)
	return &SyncServer{r: createRouting()}
}

// Connection holds connection data
type Connection struct {
	Conn *websocket.Conn
}

func subHandler(w http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// msg := make([]byte, 512)

	_, p, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
	}

	name := string(p)

	applications[name] = append(applications[name], Connection{ws})
	log.Println(applications)

	if err = ws.WriteMessage(1, p); err != nil {
		log.Println(err)
		return
	}

	if err != nil {
		log.Println(err)
	}
}

func publish(application string, data []byte) {
	log.Println("Sync is publishing to: " + application)
	if len(applications[application]) > 0 {
		for _, element := range applications[application] {
			err := element.Conn.WriteMessage(1, data)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func decodeData(req *http.Request) ([]byte, error) {
	bodyBytes, err := ioutil.ReadAll(req.Body)

	// var data JsonData
	// err := json.NewDecoder(req.Body).Decode(&data)
	return bodyBytes, err
}

// pubHandler triggers a publishing event
func pubHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	application := vars["application"]

	data, err := decodeData(req)
	if err != nil {
		log.Println(err)
	}

	publish(application, data)
	fmt.Fprintf(w, "You just published!")

}
