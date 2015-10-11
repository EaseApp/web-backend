package config

import (
	"github.com/EaseApp/web-backend/src/app/controllers/application"
	"github.com/EaseApp/web-backend/src/app/controllers/home"
	"github.com/EaseApp/web-backend/src/app/controllers/user"
	"github.com/gorilla/mux"
	// "net/http"
)

// CreateRouting establishes the routing for the server
func CreateRouting() *mux.Router {
	router := mux.NewRouter()
	// router.Host("{listen}.domain.com").Path("/").HandlerFunc(PubSubHandler).Name("root")

	router.HandleFunc("/", home.IndexHandler)

	router.HandleFunc("/static/user/new", user.NewStaticUserHandler)

	// These should be POST, but it's easier to test with GET
	router.HandleFunc("/users/sign_in", user.SignInHandler)
	router.HandleFunc("/users/sign_up", user.SignUpHandler)

	router.HandleFunc("/count/{db}", home.DBCountHandler)
	router.HandleFunc("/{db}", user.FetchAllHandler)

	router.HandleFunc("/{client}/{application}", application.QueryApplicationHandler)
	router.HandleFunc("/{client}/{application}/new", application.CreateApplicationHandler).Methods("POST")
	// router.HandleFunc("/{client}/{application}/pubsub", websocket.Handler(EchoServer))
	router.HandleFunc("/{client}/{application}/{id}", application.UpdateApplicationHandler).Methods("PUT")
	router.HandleFunc("/{client}/{application}/{id}", application.DeleteApplicationHandler).Methods("DELETE")

	return router
}
