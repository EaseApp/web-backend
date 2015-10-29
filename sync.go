package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"

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
	for _, element := range applications[application] {
		log.Println("Writing to " + application)
		err := element.Conn.WriteMessage(1, data)
		if err != nil {
			log.Println(err)
		}
	}
}

// pubHandler triggers a publishing event
func pubHandler(w http.ResponseWriter, req *http.Request) {
	publish("test", []byte("Some data!"))
	fmt.Fprintf(w, "You just published!")
}

func main() {
	// Make websocket
	log.Println("Starting sync server")
	applications = make(map[string][]Connection)
	// err := http.ListenAndServe(":8000", nil)
	// if err != nil {
	// 	panic("ListenAndServe: " + err.Error())
	// }

	router := mux.NewRouter()
	router.HandleFunc("/pub", pubHandler)
	router.HandleFunc("/sub", subHandler)

	// Make web server
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8000")

}
