package helpers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"

	"net/http"
	// "golang.org/x/net/websocket"

	"github.com/EaseApp/web-backend/src/app/models"
	"github.com/gorilla/mux"
)

<<<<<<< HEAD
var querier *models.Querier
=======
var querier *models.ModelQuerier
>>>>>>> master

type errorResponse struct {
	ErrCode    int    `json:"error_code"`
	ErrMessage string `json:"error"`
}

// Init sets up the helpers global ModelQuerier.
func Init(q *models.ModelQuerier) {
	querier = q
}

// SendError sends and logs the given error.
func SendError(errorCode int, err error, w http.ResponseWriter) {
	w.WriteHeader(errorCode)
	log.Printf("Error: Returning status code %d with error message %s.\n", errorCode, err)
	resp := errorResponse{ErrCode: errorCode, ErrMessage: err.Error()}
	json.NewEncoder(w).Encode(resp)
}

func SendSocketError(err error, conn *websocket.Conn) {
	resp := errorResponse{ErrCode: 500, ErrMessage: err.Error()}
	byteArray, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error: Cannot marshal JSON.")
		return
	}
	err = conn.WriteMessage(1, byteArray)
	if err != nil {
		log.Println("Error: Can't send socket message")
	}
}

// RequireAPIToken requires that the given route has a valid APIToken
// and passes the user it represents to the handler.
func RequireAPIToken(
	handler func(http.ResponseWriter, *http.Request, *models.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if auth == "" {
			friendlyErr := errors.New("No Authorization token provided.")
			SendError(http.StatusUnauthorized, friendlyErr, w)
			return
		}

		user := querier.FindUserByAPIToken(auth)
		if user == nil {
			friendlyErr := errors.New("Authorization token does not match.")
			SendError(http.StatusUnauthorized, friendlyErr, w)
			return
		}
		handler(w, req, user)
	}
}

// RequireAppToken requires that the given route has a valid AppToken.
// It requires that the route contains `username` and `app_name`.
func RequireAppToken(
	handler func(http.ResponseWriter, *http.Request, *models.Application)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		appName := vars["app_name"]
		username := vars["username"]
		appToken := req.Header.Get("Authorization")
		app, err := querier.AuthenticateApplication(username, appName, appToken)
		if err != nil {
			SendError(http.StatusUnauthorized, err, w)
		} else {
			handler(w, req, app)
		}
	}
}
