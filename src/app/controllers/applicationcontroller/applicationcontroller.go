package applicationcontroller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/EaseApp/web-backend/src/app/controllers/helpers"
	"github.com/EaseApp/web-backend/src/app/models"
	"github.com/gorilla/mux"
)

var querier *models.ModelQuerier

// Init sets the hacky global ModelQuerier to the given querier.
// This is to simplify the code because for this school project, we don't need
// to have perfect dependency injection practices.
func Init(querierX *models.ModelQuerier) {
	querier = querierX
}

// CreateApplicationHandler handles creating applications for the authenticated user.
func CreateApplicationHandler(w http.ResponseWriter, req *http.Request, user *models.User) {
	vars := mux.Vars(req)
	appName := vars["application"]
	newApp, err := querier.CreateApplication(user, appName)

	// TODO: Check that the application doesn't already exist.
	if err != nil {
		friendlyErr := errors.New("Could not create application")
		log.Println(friendlyErr.Error() + ".  Error: " + err.Error())
		helpers.SendError(http.StatusInternalServerError, friendlyErr, w)
		return
	}
	json.NewEncoder(w).Encode(newApp)
}

// ListApplicationsHandler handles listing the applications for the authenticated user.
func ListApplicationsHandler(w http.ResponseWriter, req *http.Request, user *models.User) {
	json.NewEncoder(w).Encode(user.Applications)
}

// DeleteApplicationHandler handles deleting the authenticated user's application.
func DeleteApplicationHandler(w http.ResponseWriter, req *http.Request, user *models.User) {
	vars := mux.Vars(req)
	appName := vars["application"]

	user, err := querier.DeleteApplication(user, appName)
	if err != nil {
		friendlyErr := errors.New("Failed to delete application")
		log.Println(friendlyErr.Error() + ".  Error: " + err.Error())
		helpers.SendError(http.StatusInternalServerError, friendlyErr, w)
		return
	}

	json.NewEncoder(w).Encode(user.Applications)
}
