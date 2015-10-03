package db

import (
	"log"

	r "github.com/dancannon/gorethink"
)

// Client holds a session connection to a RethinkDB database.
type Client struct {
	Session *r.Session
}

// NewClient connects to the database with the given address.
func NewClient(addr string) (*Client, error) {

	log.Println("Connecting to RethinkDB...")

	session, err := r.Connect(r.ConnectOpts{
		Address: addr,
		MaxIdle: 10,
		MaxOpen: 10,
		Database: "ease",
	})
	if err != nil {
		log.Println("Error connecting to RethinkdB:")
		log.Println(err.Error())
		return nil, err
	}

	log.Println("Successfully connected to RethinkDB.")

	return &Client{Session: session}, nil
}

func CreateEaseDb(c *Client){
	_, err := r.DBCreate("ease").RunWrite(c.Session)
	if err != nil{
		log.Println(err)
	}
}

func CreateUserDb(c *Client){
	_, err := r.DB("ease").TableCreate("users").RunWrite(c.Session)
	if err != nil {
		log.Println(err)
	}
}

// Close closes the connection to the database.
func (c *Client) Close() error {
	log.Println("Closing connection to RethinkDB...")
	err := c.Session.Close()
	if err != nil {
		log.Println("Error closing connection to RethinkDB:")
		log.Println(err.Error())
		return err
	}
	log.Println("Successfully closed connection to RethinkDB.")
	return nil
}
