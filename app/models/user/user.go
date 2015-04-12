package user

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	Id           string `gorethink:"id,omitempty"`
	Username     string
	PasswordHash string
	ApiToken     string
	CreatedAt    time.Time
}

func NewUser(username, password string) *User {
	user := new(User)
	user.CreatedAt = time.Now()
	user.Username = username
	randToken := make([]byte, 30)
	rand.Read(randToken)
	if err != nil {
		log.Println("Error: Couldn't generate random API token.")
		return nil
	}
	user.ApiToken = string(randToken)
	user.PasswordHash = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return user
}

func (user *User) Save() error {
	// TODO save to database.
	return nil
}

func Find(username string) *User {
	// TODO find user.
	return nil
}

func AttemptLogin(username, password string) (*User, error) {
	// TODO find user with given user name, then try password.
	return nil, nil
}
