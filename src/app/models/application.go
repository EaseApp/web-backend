package models

import (
	"errors"
	"fmt"
	"log"

	r "github.com/dancannon/gorethink"
)

// Application holds attributes for an Ease user's applications.
type Application struct {
	Name      string `gorethink:"name" json:"name"`
	AppToken  string `gorethink:"app_token" json:"app_token"`
	TableName string `gorethink:"table_name" json:"-"`
}

// newApplication creates a new application with a token and the given name.
func newApplication(user *User, appName string) (*Application, error) {
	appToken, err := generateRandomString(30)
	if err != nil {
		log.Println("Error: Couldn't generate random API token.")
		log.Println(err)
		return nil, err
	}
	tableName := fmt.Sprintf("%v_%v", user.Username, appName)
	return &Application{
		Name:      appName,
		AppToken:  appToken,
		TableName: tableName,
	}, nil
}

// CreateApplication creates a new application on the given user.
func (querier *ModelQuerier) CreateApplication(user *User, appName string) (*Application, error) {
	app, err := newApplication(user, appName)
	if err != nil {
		return nil, err
	}

	// Create a table for the new application.
	_, err = r.DB("test").TableCreate(app.TableName).RunWrite(querier.session)
	if err != nil {
		return nil, err
	}

	user.Applications = append(user.Applications, *app)
	user, err = querier.Save(user)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// DeleteApplication handles deleting an application and dropping its table.
func (querier *ModelQuerier) DeleteApplication(user *User, appName string) (*User, error) {

	// Search for the app to delete.
	var appToDelete Application
	var newApps []Application
	for i, app := range user.Applications {
		if app.Name == appName {
			newApps = append(user.Applications[:i], user.Applications[i+1:]...)
			appToDelete = app
		}
	}

	// If an app with that name does not exist.
	if newApps == nil {
		return nil, errors.New("Could not find application with that name")
	}

	// Drop the app's table.
	_, err := r.DB("test").TableDrop(appToDelete.TableName).RunWrite(querier.session)
	if err != nil {
		return nil, err
	}

	// Resave the user with the updated application list.
	user.Applications = newApps
	return querier.Save(user)
}

// AuthenticateApplication checks that the given username, app name, and
// app token are valid, and if so returns the given application.
func (querier *ModelQuerier) AuthenticateApplication(
	username, appName, appToken string) (*Application, error) {
	user := querier.FindUser(username)
	if user == nil {
		return nil, errors.New("Couldn't find user with that name")
	}

	for _, app := range user.Applications {
		if app.Name == appName && app.AppToken == appToken {
			return &app, nil
		}
	}
	return nil, errors.New("Invalid application token")
}
