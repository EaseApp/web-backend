package dao

import (
	"encoding/json"
	"fmt"
	r "github.com/dancannon/gorethink"
	"log"
	"strconv"
)

var session *r.Session

func Init(s *r.Session) {
	session = s
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

func getTableName(client, application string) string {
	return fmt.Sprintf("%v_%v", client, application)
}

func QueryApplication(client, application string, matchObject map[string]interface{}) string {
	tableName := getTableName(client, application)
	res, err := r.DB("ease").Table(tableName).Filter(matchObject).Run(session)
	if err != nil {
		log.Println(err)
		return "Error"
	}
	var records []map[string]interface{}
	err = res.All(&records)
	if err != nil {
		log.Println(err)
		return "error caught"
	}
	result := ""
	for _, p := range records {
		result += PrintObj(p)
	}
	return result
}

func CreateApplication(client, application string) string {
	tableName := getTableName(client, application)
	_, err := r.DB("ease").TableCreate(tableName).RunWrite(session)
	if err != nil {
		log.Println(err)
		return "-1"
	}
	return tableName
}

func DeleteApplication(client, application, id string) string {
	tableName := getTableName(client, application)
	result, err := r.DB("ease").Table(tableName).Get(id).Delete().RunWrite(session)
	if err != nil {
		log.Println(err)
		return "-1"
	}
	return PrintObj(result)
}

func UpdateApplication(client, application, id string, object map[string]interface{}) string {
	tableName := getTableName(client, application)
	result, err := r.DB("ease").Table(tableName).Get(id).Update(object).RunWrite(session)
	if err != nil {
		log.Println(err)
		return "-1"
	}
	return PrintObj(result)
}
