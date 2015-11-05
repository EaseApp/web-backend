package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/EaseApp/web-backend/src/lib"
	r "github.com/dancannon/gorethink"
)

// Application holds attributes for an Ease user's applications.
type Application struct {
	Name      string `gorethink:"name" json:"name"`
	AppToken  string `gorethink:"app_token" json:"app_token"`
	TableName string `gorethink:"table_name" json:"-"`
}

// appDoc is the type of each document in an application table.
type appDoc struct {
	ID   string      `gorethink:"id"`
	Name string      `gorethink:"name"`
	Data interface{} `gorethink:"data"`
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

// SaveApplicationData saves the given data to the application's table at the given path.
func (querier *ModelQuerier) SaveApplicationData(
	app *Application, path lib.Path, data interface{}) error {
	if path.IsRoot() {
		return errors.New("Cannot save data to application root")
	}
	res, err := r.Table(app.TableName).Filter(map[string]string{"name": path.TopLevelDocName}).Run(querier.session)
	if err != nil {
		return err
	}

	// Find the ID of the top-level doc.
	var docID string

	// If the top-level doc for this query doesn't exist yet, it needs to be created.
	if res.IsNil() {
		insertRes, err := r.Table(app.TableName).Insert(
			map[string]interface{}{"name": path.TopLevelDocName, "data": nil}).RunWrite(querier.session)
		if err != nil {
			return err
		}
		docID = insertRes.GeneratedKeys[0]
	} else {
		var doc appDoc
		err = res.One(&doc)
		if err != nil {
			return err
		}
		docID = doc.ID
	}

	// Generate the nested data query.
	query := path.ToNestedQuery(data)

	// Upsert the given data at the nested path.
	_, err = r.Table(app.TableName).Get(docID).Update(query).RunWrite(querier.session)

	return err
}

// ReadApplicationData reads the application's data at the given path and returns it.
func (querier *ModelQuerier) ReadApplicationData(
	app *Application, path lib.Path) (interface{}, error) {
	// Send back all the documents if root.
	if path.IsRoot() {
		res, err := r.Table(app.TableName).Filter(map[string]string{}).Run(querier.session)
		if err != nil {
			return nil, err
		}

		var docs []appDoc
		err = res.All(&docs)
		if err != nil {
			return nil, err
		}

		// Convert the documents to the pure user data.
		docsData := make(map[string]interface{})
		for _, doc := range docs {
			docsData[doc.Name] = doc.Data
		}
		return docsData, nil
	}

	res, err := r.Table(app.TableName).Filter(map[string]string{"name": path.TopLevelDocName}).Run(querier.session)
	if err != nil {
		return nil, err
	}

	// If the top-level doc for this query doesn't exist, return nil.
	if res.IsNil() {
		return nil, nil
	}

	var doc map[string]interface{}
	err = res.One(&doc)
	if err != nil {
		return nil, err
	}

	// If nested data isn't requested, return all the doc's data.
	if len(path.RemainingSegments) == 0 {
		return doc["data"], nil
	}

	nextMapLevel, ok := doc["data"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	// Dive into the nested maps.
	for idx, segment := range path.RemainingSegments {
		// Try to get the next nested level for each remaining segment.
		_, ok = nextMapLevel[segment]
		if !ok {
			return nil, nil
		}

		// Return the final data if this is the last segment
		if idx == len(path.RemainingSegments)-1 {
			return nextMapLevel[segment], nil
		}

		nextMapLevel, ok = nextMapLevel[segment].(map[string]interface{})
		// The nest doesn't go any further, so return nil.
		if !ok {
			return nil, nil
		}
	}

	// This should never be reached.
	log.Println("ERROR: This should never be reached.")
	return nil, nil
}
