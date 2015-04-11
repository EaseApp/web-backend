package db

import (
	r "github.com/dancannon/gorethink"
	"log"
)

var session *r.Session

func Init() error {

	log.Println("Connecting to RethinkDB...")

	// TODO: Set up actual configuration.
	session, err := r.Connect(r.ConnectOpts{
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

	session.SetMaxOpenConns(5)

	log.Println("Successfully connected to RethinkDB.")
	return nil
}

func Close() error {
	log.Println("Closing connection to RethinkDB...")
	err := session.Close()
	if err != nil {
		log.Println("Error closing connection to RethinkDB:")
		log.Println(err.Error())
		return err
	}
	log.Println("Successfully closed connection to RethinkDB.")
	return nil
}
