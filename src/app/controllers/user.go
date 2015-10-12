package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EaseApp/web-backend/src/app/models"
)

var querier *models.UserQuerier

// InitUserController sets the hacky global UserQuerier to the given querier.
// This is to simplify the code because for this school project, we don't need
// to have perfect dependency injection practices.
func InitUserController(userQuerier *models.UserQuerier) {
	querier = userQuerier
}

type userParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignInHandler handles user sign ins.
func SignInHandler(w http.ResponseWriter, req *http.Request) {
}

// SignUpHandler handles user sign up.
func SignUpHandler(w http.ResponseWriter, req *http.Request) {
}

// parseUserParams parses user params and returns an error to the user
// if they are invalid.
// Returns the params if successful.
// Returns error if they were invalid.
func parseUserParams(w http.ResponseWriter, req *http.Request) (userParams, error) {
	var params userParams
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		friendlyErr := errors.New("Invalid User Params: " + err.Error())
		sendError(http.StatusBadRequest, friendlyErr, w)
		return params, friendlyErr
	}
	if params.Password == "" || params.Username == "" {
		friendlyErr := errors.New("Username or password cannot be blank")
		sendError(http.StatusBadRequest, friendlyErr, w)
		return params, friendlyErr
	}
	return params, nil
}
