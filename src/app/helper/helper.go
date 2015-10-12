package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/EaseApp/web-backend/src/app/dao"
	"github.com/gorilla/mux"
)

func DecodeIOStreamToJSON(body io.Reader) (map[string]interface{}, error) {
	decoder := json.NewDecoder(body)
	var m map[string]interface{}
	err := decoder.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func hasValidPermissions(token string, req *http.Request) bool {
	vars := mux.Vars(req)
	c := vars["client"]
	client := dao.Find(c)
	if client.APIToken == token {
		return true
	}
	return false
}

// RequireAPIToken requires token
func RequireAPIToken(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")

		// log.Println("Token:" + token)
		if token != "" {
			log.Println(token)
			if hasValidPermissions(token, req) {
				handler(w, req)
			} else {
				fmt.Fprintf(w, "Token invalid")
			}
		} else {
			fmt.Fprintf(w, "Token not provided.")
		}
	}
}

// RequireAPIToken requires token
func RequireLoginToken() {
}

// func SuccessfulRequest(data string) models.Response {
// 	var res Response
// 	res.Status = 200
// 	res.Message = "Successul"
// 	res.Data = ""
// 	return res
// }
//
// func FailedRequest(data interface{}) models.Response {
// 	var res Response
// 	res.Status = 500
// 	res.Message = "Failed"
// 	res.Data = ""
// 	return res
// }

// RequestSoRiduclousThatItPurposefullyCrashesTheServer ... just in case
func RequestSoRiduclousThatItPurposefullyCrashesTheServer(funny, evenFunnier int) {
	if funny == 24 && evenFunnier == 25 {
		log.Fatal("DIVE DIVE DIVE!")
	}
}
