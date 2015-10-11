package dao

import (
	"encoding/json"
	"fmt"
	. "github.com/EaseApp/web-backend/src/app/models"
	r "github.com/dancannon/gorethink"

	"log"
	"strconv"
)

var session *r.Session

// InitGeneric sets the rethinkDB session. This name  must be unique in the dao package. TODO: Use DRY pattern for this
func InitGeneric(s *r.Session) {
	session = s
}

// GetNth returns a string of the last n objects
func GetNth(n int) string {
	result, err := r.DB("ease").Table("users").Limit(n).Run(session)
	if err != nil {
		log.Println(err)
	}
	return PrintObj(result)
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
	return PrintObj(cnt)
}

// Provided method by gorethink example for printing in memory object
func PrintObj(v interface{}) string {
	vBytes, _ := json.Marshal(v)
	return (string(vBytes))
}

func GetTableName(client, application string) string {
	return fmt.Sprintf("%v_%v", client, application)
}

// FetchAll makes a large string of all db records in a given table
func FetchAll(table string) string {
	rows, err := r.DB("ease").Table(table).Run(session)
	if err != nil {
		log.Println(err)
		return "Table " + table + " doesn't exist"
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
		result += PrintObj(p)
	}
	return result
}
