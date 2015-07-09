package user

import (
	"config/db"
	"crypto/rand"
	"errors"
	r "github.com/dancannon/gorethink"
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
	// Check that a user with the given username doesn't already exist.
	otherUser := Find(user.Username)
	if otherUser != nil && user.Id != otherUser.Id {
		return errors.New("Error: A user with that name already exists.")
	}

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
	res, err := r.Table("users").Filter(map[string]string{
		"username": username,
	}).Run(db.Session)
	if err != nil || res.IsNil() {
		return nil
	}
	var user *User
	err = res.One(&user)
	if err != nil {
		return nil
	}
	return user
}

func AttemptLogin(username, password string) *User {
	user := Find(username)
	if user == nil {
		return nil
	}
	err :=
		bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil
	}
	return user
}
