package user

import (
	"crypto/rand"
	"errors"
	r "github.com/dancannon/gorethink"
	"github.com/easeapp/web-backend/config/db"
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

func NewUser(username, password string) (*User, error) {
	user := new(User)
	user.CreatedAt = time.Now()
	user.Username = username
	randToken := make([]byte, 30)
	_, err := rand.Read(randToken)
	if err != nil {
		log.Println("Error: Couldn't generate random API token.")
		log.Println(err)
		return nil, err
	}
	user.ApiToken = string(randToken)
	byteHash, err :=
		bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.PasswordHash = string(byteHash)
	if err != nil {
		log.Println("Error: Couldn't hash password.")
		log.Println(err)
		return nil, err
	}
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
