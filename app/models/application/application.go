package models

import (
	r "github.com/dancannon/gorethink"
	"time"
)

type Application struct {
	Id        string `gorethink:"id,omitempty"`
	CreatedAt time.Time
}


func NewApplication() *Application {
	app := new(Application)
	app.CreatedAt = time.Now()
}
