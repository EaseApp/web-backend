package helper

import (
	"encoding/json"
	// "github.com/EaseApp/web-backend/src/app/models"
	"io"
	"log"
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

// RequireAPIToken requires token
func RequireAPIToken() {

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
