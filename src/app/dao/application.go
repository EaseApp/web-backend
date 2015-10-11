package dao

import (
	// "encoding/json"
	// "fmt"
	r "github.com/dancannon/gorethink"
	"log"
	// "strconv"
)

// QueryApplication looks for matching objects in an application
func QueryApplication(client, application string, matchObject map[string]interface{}) string {
	tableName := GetTableName(client, application)
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

// CreateApplication creates an application given a user
func CreateApplication(client, application string) string {
	tableName := GetTableName(client, application)
	_, err := r.DB("ease").TableCreate(tableName).RunWrite(session)
	if err != nil {
		log.Println(err)
		return "-1"
	}
	return tableName
}

// DeleteApplication deletes a record from an application
func DeleteApplication(client, application, id string) string {
	tableName := GetTableName(client, application)
	result, err := r.DB("ease").Table(tableName).Get(id).Delete().RunWrite(session)
	if err != nil {
		log.Println(err)
		return "-1"
	}
	return PrintObj(result)
}

// UpdateApplication takes a record ID and a new object and replaces the record with the new object
func UpdateApplication(client, application, id string, object map[string]interface{}) string {
	tableName := GetTableName(client, application)
	result, err := r.DB("ease").Table(tableName).Get(id).Update(object).RunWrite(session)
	if err != nil {
		log.Println(err)
		return "-1"
	}
	return PrintObj(result)
}
