package usercontroller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/EaseApp/web-backend/src/app/controllers/helpers"
	"github.com/EaseApp/web-backend/src/app/models"
	"github.com/gorilla/mux"
)

var querier *models.Querier

// Init sets the hacky global Querier to the given querier.
// This is to simplify the code because for this school project, we don't need
// to have perfect dependency injection practices.
func Init(Querier *models.Querier) {
	querier = Querier
}

type userParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignInHandler handles user sign ins.
func SignInHandler(w http.ResponseWriter, req *http.Request) {
	params, err := parseUserParams(w, req)
	if err != nil {
		return
	}

	user, err := querier.AttemptLogin(params.Username, params.Password)
	if err != nil {
		helpers.SendError(http.StatusUnauthorized, err, w)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// SignUpHandler handles user sign up.
func SignUpHandler(w http.ResponseWriter, req *http.Request) {
	params, err := parseUserParams(w, req)
	if err != nil {
		return
	}

	user, err := models.NewUser(params.Username, params.Password)
	if err != nil {
		friendlyErr := errors.New("Could not create user")
		log.Println(friendlyErr.Error() + ".  Error: " + err.Error())
		helpers.SendError(http.StatusInternalServerError, friendlyErr, w)
		return
	}

	user, err = querier.Save(user)
	if err != nil {
		helpers.SendError(http.StatusBadRequest, err, w)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// CreateApplicationHandler handles creating applications for the authenticated user.
func CreateApplicationHandler(w http.ResponseWriter, req *http.Request, user *models.User) {
	vars := mux.Vars(req)
	application := vars["application"]
	newApp, err := querier.CreateApplication(user, application)

	if err != nil {
		friendlyErr := errors.New("Could not create application")
		log.Println(err)
		helpers.SendError(http.StatusInternalServerError, friendlyErr, w)
		return
	}
	json.NewEncoder(w).Encode(newApp)
}

// ListApplicationsHandler handles listing the applications for the authenticated user.
func ListApplicationsHandler(w http.ResponseWriter, req *http.Request, user *models.User) {
	json.NewEncoder(w).Encode(user.Applications)
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
		helpers.SendError(http.StatusBadRequest, friendlyErr, w)
		return params, friendlyErr
	}
	if params.Password == "" || params.Username == "" {
		friendlyErr := errors.New("Username or password cannot be blank")
		helpers.SendError(http.StatusBadRequest, friendlyErr, w)
		return params, friendlyErr
	}
	return params, nil
}
