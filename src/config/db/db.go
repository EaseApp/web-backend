package db

import (
	r "github.com/dancannon/gorethink"
	"log"
)

var Session *r.Session

func Init() error {

	log.Println("Connecting to RethinkDB...")

	// TODO: Set up actual production configuration.
	Session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "test",
		MaxIdle:  10,
		MaxOpen:  10,
	})
	if err != nil {
		log.Println("Error connecting to RethinkdB:")
		log.Println(err.Error())
		return err
	}

	Session.SetMaxOpenConns(5)

	log.Println("Successfully connected to RethinkDB.")

	setupTables()

	return nil
}

func Close() error {
	log.Println("Closing connection to RethinkDB...")
	err := Session.Close()
	if err != nil {
		log.Println("Error closing connection to RethinkDB:")
		log.Println(err.Error())
		return err
	}
	log.Println("Successfully closed connection to RethinkDB.")
	return nil
}

func setupTables() {

	log.Println("Setting up tables...")

	// Set up the initial user table.
	// TODO Maybe handle the table already existing?
	//r.DB("test").TableCreate("users").RunWrite(Session)
	//r.Table("users").IndexCreate("Username").RunWrite(Session)

	log.Println("Done setting up tables.")
}
