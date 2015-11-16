package applicationcontroller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/EaseApp/web-backend/src/app/controllers/helpers"
	"github.com/EaseApp/web-backend/src/app/models"
	"github.com/EaseApp/web-backend/src/lib"
	"github.com/gorilla/mux"
)

var querier *models.ModelQuerier

// Init sets the hacky global ModelQuerier to the given querier.
// This is to simplify the code because for this school project, we don't need
// to have perfect dependency injection practices.
func Init(querierX *models.ModelQuerier) {
	querier = querierX
}

var testingOnlySyncServerURL string

// TestingOnlySetSyncServerURL sets the sync server URL for use in the tests.
// This is because we can't get the URL ahead of time with httptest.
func TestingOnlySetSyncServerURL(syncServerURL string) {
	testingOnlySyncServerURL = syncServerURL
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

// appDataReqParams holds the params needed for the data handlers below.
type appDataReqParams struct {
	PathStr  string      `json:"path"`
	Data     interface{} `json:"data"`
	Path     lib.Path    `json:"-"`
	Username string      `json:"username"`
}

var successResponse = struct {
	Success bool `json:"success"`
}{true}

// SaveApplicationDataHandler handles saving app data.
func SaveApplicationDataHandler(w http.ResponseWriter, req *http.Request, app *models.Application) {
	params, err := parseAppDataParams(w, req)
	if err != nil {
		return
	}

	err = querier.SaveApplicationData(app, params.Path, params.Data)
	if err != nil {
		friendlyErr := errors.New("Failed to save application data")
		log.Println(friendlyErr, ": ", err)
		helpers.SendError(http.StatusInternalServerError, friendlyErr, w)
		return
	}

	go sendPublishEvent(app, "SAVE", params)

	json.NewEncoder(w).Encode(successResponse)
}

// ReadApplicationDataHandler handles reading app data.
func ReadApplicationDataHandler(w http.ResponseWriter, req *http.Request, app *models.Application) {
	path, err := lib.ParsePath(req.URL.Query().Get("path"))
	if err != nil {
		helpers.SendError(http.StatusBadRequest, err, w)
		return
	}

	data, err := querier.ReadApplicationData(app, path)
	if err != nil {
		friendlyErr := errors.New("Failed to read application data")
		log.Println(friendlyErr, ": ", err)
		helpers.SendError(http.StatusInternalServerError, friendlyErr, w)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// DeleteApplicationDataHandler handles deleting app data.
func DeleteApplicationDataHandler(w http.ResponseWriter, req *http.Request, app *models.Application) {
	params, err := parseAppDataParams(w, req)
	if err != nil {
		return
	}

	err = querier.DeleteApplicationData(app, params.Path)
	if err != nil {
		friendlyErr := errors.New("Failed to delete application data")
		log.Println(friendlyErr, ": ", err)
		helpers.SendError(http.StatusInternalServerError, friendlyErr, w)
		return
	}

	go sendPublishEvent(app, "DELETE", params)

	json.NewEncoder(w).Encode(successResponse)
}

// sendPublishEvent sends a publish event to the sync server.
func sendPublishEvent(app *models.Application, action string, params appDataReqParams) *http.Response {
	url := "http://localhost:8000"
	if testingOnlySyncServerURL != "" {
		url = testingOnlySyncServerURL
	}

	buff := bytes.NewBuffer(nil)
	json.NewEncoder(buff).Encode(map[string]interface{}{
		"path":   params.PathStr,
		"data":   params.Data,
		"action": action,
	})

	log.Println("Sending Publish Event:", params.Username)
	resp := sendJSON(buff, app.AppToken, url, "/pub/"+params.Username+"/"+app.TableName, "POST")
	log.Println("Publish response: ", resp)
	return resp
}

func sendJSON(data io.Reader, token, url, path, method string) *http.Response {
	req, err := http.NewRequest(method, url+path, data)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	return resp
}

// parseAppDataParams parses the given app data params and sends an error if they're invalid.
func parseAppDataParams(w http.ResponseWriter, req *http.Request) (appDataReqParams, error) {
	var params appDataReqParams

	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		friendlyErr := errors.New("Invalid JSON: " + err.Error())
		helpers.SendError(http.StatusBadRequest, friendlyErr, w)
		return params, err
	}

	params.Path, err = lib.ParsePath(params.PathStr)
	if err != nil {
		friendlyErr := errors.New("Invalid Path: " + err.Error())
		helpers.SendError(http.StatusBadRequest, friendlyErr, w)
		return params, err
	}

	vars := mux.Vars(req)
	username := vars["username"]
	params.Username = username

	return params, nil
}
