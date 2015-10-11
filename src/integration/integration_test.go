package integration

import (
	"bytes"
	"github.com/EaseApp/web-backend/src/server"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func sendJSON(jsonInput, path string, t *testing.T) string {
	var jsonStr = []byte(jsonInput)
	req, err := http.NewRequest("POST", "http://localhost:3000"+path, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return (string(body))
}

func SetUpServer() {
	go server.StartServer()
}

func TestSetUpServer(t *testing.T) {
	SetUpServer()
	time.Sleep(time.Second)
	body := sendJSON(`{"username":"user", "password":"password"}`, "/users/sign_up", t)
	assert.NotEqual(t, body, "")
}
