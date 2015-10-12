package server

import (
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

	router := mux.Router()

	//	router.HandleFunc().Methods("POST")
	return router
}
