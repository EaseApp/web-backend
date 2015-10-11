package user

import (
	// "encoding/json"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	// "errors"
	"net/http"
	"time"

	. "github.com/EaseApp/web-backend/src/app/dao"
	. "github.com/EaseApp/web-backend/src/app/models"

	"github.com/EaseApp/web-backend/src/app/helpers"

	r "github.com/dancannon/gorethink"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	// "net/url"
)

var session *r.Session

// SignInHandler takes username and password and checks whether the hashes match
func SignInHandler(w http.ResponseWriter, req *http.Request) {
	obj, err := helper.DecodeIOStreamToJSON(req.Body)
	if err == io.EOF {
		fmt.Fprintf(w, "No creds provided.")
	} else if err != nil {
		log.Println(err)
		fmt.Fprintf(w, PrintObj(err))
	} else {
		username := obj["username"]
		sUsername, usernameOk := username.(string)

		password := obj["password"]
		sPassword, passwordOk := password.(string)
		if usernameOk && passwordOk {
			result := AttemptLogin(sUsername, sPassword)
			if result != nil {
				// fmt.Fprintf(w, "%v", PrintObj(helper.SuccessfulRequest(result.LoginToken)))
				fmt.Fprintf(w, "%v", result.LoginToken)
			} else {
				// fmt.Fprintf(w, PrintObj(helper.FailedRequest("Ayyy lmao")))
				fmt.Fprintf(w, "Ayy lmao")
			}
		} else {
			fmt.Fprintf(w, "Creds are not strings")
		}
	}
}

// SignUpHandler takes username and password in URL, makes a new user, and returns a token
func SignUpHandler(w http.ResponseWriter, req *http.Request) {
	obj, err := helper.DecodeIOStreamToJSON(req.Body)
	if err == io.EOF {
		fmt.Fprintf(w, "No creds provided.")
	} else if err != nil {
		log.Println(err)
		fmt.Fprintf(w, PrintObj(err))
	} else {
		username := obj["username"]
		sUsername, usernameOk := username.(string)

		password := obj["password"]
		sPassword, passwordOk := password.(string)

		if usernameOk && passwordOk {
			user, err := NewUser(sUsername, sPassword)
			if err != nil {
				log.Println(w, "Error: %v", err)
			}
			err = Save(user)
			if err != nil {
				log.Println(w, "Error2: %v", err)
				fmt.Fprintf(w, "%v", err)
			} else {
				fmt.Fprintf(w, user.LoginToken)
			}
		} else {
			fmt.Fprintf(w, "Creds are not strings")
		}
	}
}

// FetchAllHandler returns everything in a table
func FetchAllHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	db := vars["db"]

	fmt.Fprint(w, FetchAll(db))
}

// AttemptLogin attempts to find a user then returns the user
func AttemptLogin(username, password string) *User {
	user := Find(username)
	if user == nil {
		return nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil
	}
	return user
}

// NewUser initializes a blank user with a username and password
func NewUser(username, password string) (*User, error) {
	user := new(User)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LoginTokenUpdatedAt = time.Now()
	user.Username = username
	randToken := make([]byte, 30)
	_, err := rand.Read(randToken)
	if err != nil {
		log.Println("Error: Couldn't generate random API token.")
		log.Println(err)
		return nil, err
	}
	user.APIToken = string(randToken)

	randToken = make([]byte, 30)
	_, err = rand.Read(randToken)
	if err != nil {
		log.Println("Error: Couldn't generate random API token for Login.")
		log.Println(err)
		return nil, err
	}
	user.LoginToken = string(randToken)

	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.PasswordHash = string(byteHash)
	if err != nil {
		log.Println("Error: Couldn't hash password.")
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func NewStaticUserHandler(w http.ResponseWriter, req *http.Request) {
	_ = InsertStaticUser()
	http.Redirect(w, req, "/users", http.StatusFound)
}
