package models

// Querier.go holds the ModelQuerier type.
// All other models use it to query the database.

import (
	r "github.com/dancannon/gorethink"
)

// ModelQuerier queries the database for each model.
type ModelQuerier struct {
	session *r.Session
}

// NewModelQuerier returns a new ModelQuerier.
func NewModelQuerier(session *r.Session) *ModelQuerier {
	return &ModelQuerier{session: session}
}
