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
	})
	if err != nil {
		log.Println("Error connecting to RethinkdB:")
		log.Println(err.Error())
		return nil, err
	}

	log.Println("Successfully connected to RethinkDB.")

	return &Client{Session: session}, nil
}

// SetUpDatabase creates the needed tables.
func (c *Client) SetUpDatabase() error {
	err := createEaseDb(c)
	if err != nil {
		// log.Println(err)
	}
	err = createUserTable(c)
	if err != nil {
		// log.Println(err)
	}
	err = createApplicationTable(c)
	if err != nil {
		// log.Println(err)
	}

	// Probably not the best, but everytime we start the server, we expect these errors
	return nil
}

func createEaseDb(c *Client) error {
	_, err := r.DBCreate("ease").RunWrite(c.Session)
	return err
}

func createUserTable(c *Client) error {
	_, err := r.DB("ease").TableCreate("users").RunWrite(c.Session)
	return err
}

func createApplicationTable(c *Client) error {
	_, err := r.DB("ease").TableCreate("applications").RunWrite(c.Session)
	return err
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
