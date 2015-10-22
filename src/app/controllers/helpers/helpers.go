package helpers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/EaseApp/web-backend/src/app/models"
)

var querier *models.UserQuerier

type errorResponse struct {
	ErrCode    int    `json:"error_code"`
	ErrMessage string `json:"error"`
}

// Init sets up the helpers global UserQuerier.
func Init(q *models.UserQuerier) {
	querier = q
}

// SendError sends and logs the given error.
func SendError(errorCode int, err error, w http.ResponseWriter) {
	w.WriteHeader(errorCode)
	log.Printf("Error: Returning status code %d with error message %s.\n", errorCode, err)
	resp := errorResponse{ErrCode: errorCode, ErrMessage: err.Error()}
	json.NewEncoder(w).Encode(resp)
}

// RequireAPIToken requires that the given route has a valid APIToken
// and passes the user it represents to the handler.
func RequireAPIToken(
	handler func(http.ResponseWriter, *http.Request, *models.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// TODO: implement this.
		// Make a new UserQuerier function to find a user by API token.
		// If one isn't found, do:
		// SendError(http.StatusUnAuthorized, errors.New("Invalid API token.", w)
		// If one is found, call handler like the below code:
		var user *models.User
		auth := req.Header.Get("Authorization")
		if auth == "" {
			friendlyErr := errors.New("No Authorization token provided.")
			SendError(http.StatusUnauthorized, friendlyErr, w)
			return
		}

		user = querier.FindUserByAPIToken(auth)
		if user == nil {
			friendlyErr := errors.New("Authorization token does not match.")
			SendError(http.StatusUnauthorized, friendlyErr, w)
			return
		}
		handler(w, req, user)
	}
}
