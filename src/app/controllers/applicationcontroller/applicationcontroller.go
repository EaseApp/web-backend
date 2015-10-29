package applicationcontroller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EaseApp/web-backend/src/app/controllers/helpers"
	"github.com/EaseApp/web-backend/src/app/models"
)

var querier *models.ApplicationQuerier

// Init sets the hacky global UserQuerier to the given querier.
// This is to simplify the code because for this school project, we don't need
// to have perfect dependency injection practices.
func Init(applicationQuerier *models.ApplicationQuerier) {
	querier = applicationQuerier
}

type userParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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
