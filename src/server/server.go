package server

import (
	"fmt"
	"net/http"

	"github.com/EaseApp/web-backend/src/app/controllers/helpers"
	"github.com/EaseApp/web-backend/src/app/controllers/usercontroller"
	"github.com/EaseApp/web-backend/src/app/models"
	"github.com/EaseApp/web-backend/src/db"
	"github.com/gorilla/mux"
)

// EaseServer serves while allowing cross origin access.
type EaseServer struct {
	r *mux.Router
}

// NewEaseServer creates a new handler for Ease.
func NewEaseServer(client *db.Client) *EaseServer {
	return &EaseServer{r: createRoutingMux(client)}
}

// ServeHTTP serves requests from the EaseServer's mux while allowing
// cross origin access.
func (s *EaseServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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

// createRoutingMux sets up the routing for the server.
func createRoutingMux(client *db.Client) *mux.Router {

	// Set up the queriers and controllers.
	querier := models.NewModelQuerier(client.Session)

	usercontroller.Init(querier)
	helpers.Init(querier)

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	// User routes.
	router.HandleFunc("/users/sign_up", usercontroller.SignUpHandler).Methods("POST")
	router.HandleFunc("/users/sign_in", usercontroller.SignInHandler).Methods("POST")
	router.HandleFunc("/users/applications/{application}",
		helpers.RequireAPIToken(usercontroller.CreateApplicationHandler)).Methods("POST")
	router.HandleFunc("/users/applications",
		helpers.RequireAPIToken(usercontroller.ListApplicationsHandler)).Methods("GET")
	router.HandleFunc("/users/applications/{application}",
		helpers.RequireAPIToken(usercontroller.DeleteApplicationHandler)).Methods("DELETE")

	return router
}
