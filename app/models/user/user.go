package user

import (
	"crypto/rand"
	r "github.com/dancannon/gorethink"
	"github.com/easeapp/web-backend/config/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
	"errors"
)

type User struct {
	Id           string `gorethink:"id,omitempty"`
	Username     string
	PasswordHash string
	ApiToken     string
	CreatedAt    time.Time
}

func NewUser(username, password string) *User, error {
	user := new(User)
	user.CreatedAt = time.Now()
	user.Username = username
	randToken := make([]byte, 30)
	rand.Read(randToken)
	if err != nil {
		friendlyErr := errors.New("Error: Couldn't generate random API token.")
		log.Println(friendlyErr)
		log.Println(err)
		return nil, friendlyErr
	}
	user.ApiToken = string(randToken)
	user.PasswordHash = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return user, nil
}

func (user *User) Save() error {
	_, err := r.Table("users").Insert(user).RunWrite(db.Session)
	if err != nil {
		friendlyErr := errors.New("Error: Couldn't save user.")
		log.Println(friendlyErr)
		log.Println(err)
		return friendlyErr
	}
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
