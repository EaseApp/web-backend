package controller

import(
	"fmt"
	"log"
	r "github.com/dancannon/gorethink"
	"encoding/json"
	"strconv"
	"crypto/rand"
	"time"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"github.com/gorilla/mux"
	"errors"
	// "net/url"
	// "github.com/EaseApp/web-backend/src/app/models/user"
)

type User struct{
	Id string `gorethink:"id,omitempty"`
	Username string
	Email	string
	PasswordHash string
	ApiToken string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var session *r.Session

func (user *User) Save() error {
	// Check that a user with the given username doesn't already exist.
	otherUser := Find(user.Username)
	log.Println("USERNAME ", Find(user.Username))
	if otherUser != nil && user.Id != otherUser.Id {
		return errors.New("Error: A user with that name already exists.")
	}

	_, err := r.Table("users").Insert(user).RunWrite(session)
	if err != nil {
		friendlyErr := errors.New("Error: Couldn't save user.")
		log.Println(friendlyErr)
		log.Println(err)
		return friendlyErr
	}
	return nil
}

func SignInHandler(w http.ResponseWriter, req *http.Request){
	username := req.URL.Query().Get("u")
	password := req.URL.Query().Get("p")

	result := AttemptLogin(username, password)
	if result != nil{
			fmt.Fprintf(w, "Successful login")
	} else {
			fmt.Fprintf(w, "Login failed")
	}

}

func SignUpHandler(w http.ResponseWriter, req *http.Request){
	username := req.URL.Query().Get("u")
	password := req.URL.Query().Get("p")

	user, err := NewUser(username, password)
	if err != nil {
			log.Println(w, "Error: %v", err)
	}
	err = user.Save()
	if err != nil {
			log.Println(w, "Error2: %v", err)
			fmt.Fprintf(w, "Problem saving")
	} else {
		fmt.Fprintf(w, "%v", "Signed Up")
	}
}

func FetchAllHandler(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	table := vars["db"]

	fmt.Fprintf(w, "Table: (%v). All records: (%v)", table, FetchAll(table))
}

// Initialize connection and set global session variable
func Init() {
    var err error
    session, err = r.Connect(r.ConnectOpts{
        Address:  "localhost:28015",
        Database: "ease",
    })
    if err != nil {
        log.Println(err)
        return
    }
}


// Custom method to find user by Rethink ID
func FindUser(Id string) (string){
	result, err := r.DB("ease").Table("users").Get(Id).Run(session)
	if err != nil{
		log.Println(err)
	}
	return printObj(result)
}

// Custom method to add seed or extra data
func InsertStaticUser() string{
	var data = map[string]interface{}{
		"Username": fmt.Sprintf("User-%v", time.Now()),
		"Email": "email@domain.com",
		"PasswordHash": "static_passwordhash",
		"ApiToken": "Idk what this is yet",
		"CreatedAt": time.Now(),
		"UpdatedAt": time.Now(),
  }
	result, err := r.DB("ease").Table("users").Insert(data).RunWrite(session)
	if err != nil{
		log.Println(err)
		return ""
	}
	log.Println(printObj(result))
	return result.GeneratedKeys[0]
}

// Custom method to return db record count
func RecordCount(db string) string {
    cursor, err := r.DB("ease").Table(db).Count().Run(session)
    if err != nil {
        log.Println(err)
        return strconv.Itoa(-1)
    }
    var cnt int
    cursor.One(&cnt)
    cursor.Close()
    return printObj(cnt)
}

// Custom method to get last n objects
func GetNth(n int) (string){
	result, err := r.DB("ease").Table("users").Limit(n).Run(session)
	if err != nil{
		log.Println(err)
	}
	return printObj(result)
}

// Custom method to make large string of all db records
func FetchAll(table string) string{
    rows, err := r.Table(table).Run(session)
    if err != nil {
        log.Println(err)
        return "Table "+table+" doesn't exist"
    }
    // Read records into persons slice
    var records []User
    err2 := rows.All(&records)
    if err2 != nil {
        log.Println(err2)
        return "error caught2"
    }
		result := ""
    for _, p := range records {
        result += printObj(p)
    }
		return result
}

// Provided method by gorethink example for printing in memory object
func printObj(v interface{}) (string) {
    vBytes, _ := json.Marshal(v)
    return (string(vBytes))
}

// Provided method
func Find(username string) *User {
	res, err := r.Table("users").Filter(map[string]string{
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

// Provided method
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

// Provided method
func NewUser(username, password string) (*User, error) {
	user := new(User)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
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


/*

import (

	"log"

)



*/
