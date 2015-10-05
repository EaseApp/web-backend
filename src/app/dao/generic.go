package dao

import(
  	"strconv"
  	r "github.com/dancannon/gorethink"
  	"log"
  	"encoding/json"
    "fmt"
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
func PrintObj(v interface{}) (string) {
    vBytes, _ := json.Marshal(v)
    return (string(vBytes))
}

func QueryApplication(client, application, matchObject string)(string){
  tableName := fmt.Sprintf("%v_%v", client, application)
  cursor, err := r.DB("ease").Table(tableName).Count().Run(session)
  if err != nil{
    log.Println(err)
    return strconv.Itoa(-1)
  }
  var cnt int
  cursor.One(&cnt)
  cursor.Close()
  return PrintObj(cnt)
}

func CreateApplication(client, application string)(string){
  tableName := fmt.Sprintf("%v_%v", client, application)
  _, err := r.DB("ease").TableCreate(tableName).RunWrite(session)
	if err != nil {
		log.Println(err)
    return "-1"
	}
  return tableName
}
