package helper

import (
	"encoding/json"
	"io"
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
