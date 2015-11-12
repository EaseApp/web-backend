package sync

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	// "github.com/codegangsta/negroni"

	"github.com/EaseApp/web-backend/src/app/controllers/helpers"
	"github.com/EaseApp/web-backend/src/app/models"
	"github.com/EaseApp/web-backend/src/db"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

func createRouting(client *db.Client) *mux.Router {
	querier := models.NewModelQuerier(client.Session)

	helpers.Init(querier)

	router := mux.NewRouter()
	router.HandleFunc("/sub", subHandler)
	router.HandleFunc("/pub/{username}/{app_name}", helpers.RequireAppToken(pubHandler)).Methods("POST")
	return router
}

func NewSyncServer(client *db.Client) *SyncServer {
	applications = make(map[string][]Connection)
	return &SyncServer{r: createRouting(client)}
}

type applicationParams struct {
	Username      string
	Application   string
	Authorization string
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

	// TODO: Investigate. This ReadMessage method might block. Meaning, if you dont get a message immediately you're holding the server. Need to investigate.
	_, p, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
		friendlyErr := errors.New("Reading application error.")
		helpers.SendSocketError(friendlyErr, ws)
		return
	}

	var params applicationParams
	err = json.NewDecoder(bytes.NewReader(p)).Decode(&params)
	if err != nil {
		log.Println(err)
		return
	}

	if helpers.IsValidAppToken(params.Username, params.Application, params.Authorization) {
		appName := helpers.GetAppName(params.Username, params.Application)
		applications[appName] = append(applications[appName], Connection{ws})
		log.Println(applications)
		success := `{"status": "success"}`
		if err = ws.WriteMessage(1, []byte(success)); err != nil {
			log.Println(err)
			return
		}
	} else {
		failed := `{"status": "failed"}`
		if err = ws.WriteMessage(1, []byte(failed)); err != nil {
			log.Println(err)
			return
		}
	}

	if err != nil {
		log.Println(err)
		return
	}
}

func publish(application string, data []byte) {
	log.Println("Sync is publishing to: " + application)
	for _, element := range applications[application] {
		// err := element.Conn.WriteMessage(1, data)

		w, err := element.Conn.NextWriter(1)
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
			log.Println(err)
			return
		}
		if err := w.Close(); err != nil {
			log.Println(err)
			return
		}
	}
}

func decodeData(req *http.Request) ([]byte, error) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	return bodyBytes, err
}

// pubHandler triggers a publishing event
func pubHandler(w http.ResponseWriter, req *http.Request, app *models.Application) {
	data, err := decodeData(req)
	if err != nil {
		log.Println(err)
	}

	go publish(app.TableName, data)
	fmt.Fprintf(w, "You just published!")

}
