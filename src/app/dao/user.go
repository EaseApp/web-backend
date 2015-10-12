package dao

import (
	"errors"
	"fmt"
	"log"
	"time"

	. "github.com/EaseApp/web-backend/src/app/models"
	r "github.com/dancannon/gorethink"
)

// Find attempts to find user by username
func Find(username string) *User {
	res, err := r.DB("ease").Table("users").Filter(map[string]string{
		"Username": username,
	}).Run(session)
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

// FindUser attempts to find a user by RethinkDB ID and print the User
func FindUser(ID string) string {
	result, err := r.DB("ease").Table("users").Get(ID).Run(session)
	if err != nil {
		log.Println(err)
	}
	return PrintObj(result)
}

// InsertStaticUser creates a basic user. This was just a test method to make sure the insert was working. May be moved to a utility file or removed.
func InsertStaticUser() string {
	var data = map[string]interface{}{
		"Username":     fmt.Sprintf("User-%v", time.Now()),
		"Email":        "email@domain.com",
		"PasswordHash": "static_passwordhash",
		"ApiToken":     "Idk what this is yet",
		"CreatedAt":    time.Now(),
		"UpdatedAt":    time.Now(),
	}
	result, err := r.DB("ease").Table("users").Insert(data).RunWrite(session)
	if err != nil {
		log.Println(err)
		return ""
	}
	log.Println(PrintObj(result))
	return result.GeneratedKeys[0]
}

func Save(user *User) error {
	// Check that a user with the given username doesn't already exist.
	otherUser := Find(user.Username)
	log.Println("USERNAME ", Find(user.Username))
	if otherUser != nil && user.ID != otherUser.ID {
		return errors.New("Error: A user with that name already exists.")
	}

	_, err := r.DB("ease").Table("users").Insert(user).RunWrite(session)
	if err != nil {
		friendlyErr := errors.New("Error: Couldn't save user.")
		log.Println(friendlyErr)
		log.Println(err)
		return friendlyErr
	}
	return nil
}
