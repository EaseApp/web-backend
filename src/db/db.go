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
// It also sets up the database tables needed for Ease.
func NewClient(addr string) (*Client, error) {

	log.Println("Connecting to RethinkDB...")

	// TODO: Make this function take in a default database and set it here.
	session, err := r.Connect(r.ConnectOpts{
		Address:  addr,
		MaxIdle:  10,
		MaxOpen:  10,
		Database: "test",
	})
	r.SetVerbose(true)
	if err != nil {
		log.Println("Error connecting to RethinkdB:")
		log.Println(err.Error())
		return nil, err
	}

	setupTables(session)

	log.Println("Successfully connected to RethinkDB.")

	return &Client{Session: session}, nil
}

// setupTables sets up the needed tables for Ease.
func setupTables(session *r.Session) {
	// TODO: Make this use the default database param.
	r.DB("test").TableCreate("users").Exec(session)
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
