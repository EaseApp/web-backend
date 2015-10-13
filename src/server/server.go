package server

import (
	"fmt"
	"net/http"

	"github.com/EaseApp/web-backend/src/app/controllers"
	"github.com/EaseApp/web-backend/src/app/models"
	"github.com/EaseApp/web-backend/src/db"
	"github.com/gorilla/mux"
)

// CreateRoutingMux sets up the routing for the server.
func CreateRoutingMux(client *db.Client) *mux.Router {

	// Set up the queriers and controllers.
	userQuerier := models.NewUserQuerier(client.Session)

	controllers.InitUserController(userQuerier)

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	// User routes.
	router.HandleFunc("/users/sign_up", controllers.SignUpHandler).Methods("POST")
	router.HandleFunc("/users/sign_in", controllers.SignInHandler).Methods("POST")

	return router
}
