package models

import "time"

// User struct holds the data for a user
type User struct {
	ID                  string `gorethink:"id,omitempty"`
	Username            string
	Email               string
	PasswordHash        string
	APIToken            string
	LoginToken          string
	LoginTokenUpdatedAt time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
