package models

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"time"

	r "github.com/dancannon/gorethink"
	"golang.org/x/crypto/bcrypt"
)

// User holds attribute for an Ease user.
type User struct {
	ID           string        `gorethink:"id,omitempty" json:"id"`
	Username     string        `gorethink:"username" json:"username"`
	PasswordHash string        `gorethink:"password_hash" json:"-"`
	APIToken     string        `gorethink:"api_token" json:"api_token"`
	CreatedAt    time.Time     `gorethink:"created_at" json:"created_at"`
	Applications []Application `gorethink:"applications" json:"applications"`
}

// NewUser creates a new user with tokens and a hashed password.
func NewUser(username, password string) (*User, error) {
	user := &User{}
	user.CreatedAt = time.Now()
	user.Username = username
	apiToken, err := generateRandomString(30)
	if err != nil {
		log.Println("Error: Couldn't generate random API token.")
		log.Println(err)
		return nil, err
	}
	user.APIToken = apiToken

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

// FindUserByAPIToken finds a user by an API token.
func (querier *ModelQuerier) FindUserByAPIToken(token string) *User {
	res, err := r.Table("users").Filter(map[string]string{
		"api_token": token,
	}).Run(querier.session)
	defer res.Close()
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

// Save saves the given user and returns it.
// It verifies that the given username isn't already taken.
// Returns the updated user.
func (querier *ModelQuerier) Save(user *User) (*User, error) {
	// Check that a user with the given username doesn't already exist.
	otherUser := querier.FindUser(user.Username)
	if otherUser != nil && user.ID != otherUser.ID {
		return nil, errors.New("A user with that name already exists")
	}

	// Upsert the user.
	res, err := r.Table("users").Insert(
		user, r.InsertOpts{Conflict: "replace"},
	).RunWrite(querier.session)

	if err != nil {
		friendlyErr := errors.New("Couldn't save user")
		log.Println(friendlyErr)
		log.Println(err)
		return nil, friendlyErr
	}

	// Get the user's ID if one was generated.
	if user.ID == "" {
		user.ID = res.GeneratedKeys[0]
	}

	return user, nil
}

// FindUser finds the user with the given username.  Returns nil if none found.
func (querier *ModelQuerier) FindUser(username string) *User {
	res, err := r.Table("users").Filter(map[string]string{
		"username": username,
	}).Run(querier.session)
	defer res.Close()
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

// AttemptLogin attempts to login the user with the given username and password.
// Returns the user if successful, nil if failed.
func (querier *ModelQuerier) AttemptLogin(username, password string) (*User, error) {
	user := querier.FindUser(username)
	if user == nil {
		return nil, errors.New("Couldn't find user with that username")
	}
	err :=
		bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("Password was invalid")
	}
	return user, nil
}

// Possible token chars.
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generateRandomString generates a random string of length n.
func generateRandomString(n int) (string, error) {
	b := make([]rune, n)
	for i := range b {
		randInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[int(randInt.Int64())]
	}
	return string(b), nil
}
