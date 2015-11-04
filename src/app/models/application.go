package models

// Application holds attributes for an Ease user's applications.
type Application struct {
	ID       string `gorethink:"id,omitempty" json:"id"`
	Name     string `gorethink:"name" json:"name"`
	AppToken string `gorethink:"app_token" json:"app_token"`
}
