package models

import r "github.com/dancannon/gorethink"

// UserQuerier queries the user table and logs users in.
type ApplicationQuerier struct {
	session *r.Session
}

// Application holds attributes for an Ease user's applications.
type Application struct {
	ID       string `gorethink:"id,omitempty" json:"id"`
	Name     string `gorethink:"name" json:"name"`
	AppToken string `gorethink:"app_token" json:"app_token"`
}
