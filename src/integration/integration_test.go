package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EaseApp/web-backend/src/app/models"
	"github.com/EaseApp/web-backend/src/db"
	"github.com/EaseApp/web-backend/src/server"
	r "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var localDBAddr = "localhost:28015"

type errorResp struct {
	Err string `json:"error"`
}

func TestSignUp(t *testing.T) {
	server := setUpServer(t)
	defer server.Close()

	testcases := []struct {
		input            string
		expectedCode     int
		expectedError    string
		expectedUsername string
	}{
		{
			input:            `{"username": "user", "password": "pass"}`,
			expectedCode:     http.StatusOK,
			expectedError:    "",
			expectedUsername: "user",
		},
		{
			input:            `{"username": "user", "password": "pass"}`,
			expectedCode:     http.StatusBadRequest,
			expectedError:    "A user with that name already exists",
			expectedUsername: "",
		},
		{
			input:            `{"password": "pass"}`,
			expectedCode:     http.StatusBadRequest,
			expectedError:    "Username or password cannot be blank",
			expectedUsername: "",
		},
	}

	for _, testcase := range testcases {
		resp := sendJSON(testcase.input, server.URL, "/users/sign_up", "POST", t)

		assert.Equal(t, testcase.expectedCode, resp.StatusCode)

		// No error expected.
		if testcase.expectedError == "" {
			var errStruct errorResp
			json.NewDecoder(resp.Body).Decode(&errStruct)
			assert.Equal(t, testcase.expectedError, errStruct.Err)
		} else { // Error expected.
			var userStruct models.User
			json.NewDecoder(resp.Body).Decode(&userStruct)
			assert.Equal(t, testcase.expectedUsername, userStruct.Username)
		}
	}
}

func TestSignIn(t *testing.T) {
	server := setUpServer(t)
	defer server.Close()

	// Create a test user for signing in with.
	resp := sendJSON(`{"username":"ronswanson","password":"meat"}`,
		server.URL, "/users/sign_up", "POST", t)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var testUser models.User
	err := json.NewDecoder(resp.Body).Decode(&testUser)
	require.NoError(t, err)
	userAPIToken := testUser.APIToken

	testcases := []struct {
		input            string
		expectedCode     int
		expectedError    string
		expectedUsername string
		expectedAPIToken string
	}{
		{
			input:            `{"username": "ronswanson", "password": "meat"}`,
			expectedCode:     http.StatusOK,
			expectedError:    "",
			expectedUsername: "user",
			expectedAPIToken: userAPIToken,
		},
		{
			input:            `{"username": "anneperkins", "password": "pass"}`,
			expectedCode:     http.StatusUnauthorized,
			expectedError:    "Couldn't find user with that username",
			expectedUsername: "",
			expectedAPIToken: "",
		},
		{
			input:            `{"username": "anneperkins", "password": "pass"}`,
			expectedCode:     http.StatusUnauthorized,
			expectedError:    "Password was invalid",
			expectedUsername: "",
			expectedAPIToken: "",
		},
	}

	for _, testcase := range testcases {
		resp := sendJSON(testcase.input, server.URL, "/users/sign_in", "POST", t)

		assert.Equal(t, testcase.expectedCode, resp.StatusCode)

		// No error expected.
		if testcase.expectedError == "" {
			var errStruct errorResp
			json.NewDecoder(resp.Body).Decode(&errStruct)
			assert.Equal(t, testcase.expectedError, errStruct.Err)
		} else { // Error expected.
			var userStruct models.User
			json.NewDecoder(resp.Body).Decode(&userStruct)
			assert.Equal(t, testcase.expectedUsername, userStruct.Username)
			assert.Equal(t, testcase.expectedAPIToken, userStruct.APIToken)
		}
	}

}

func sendJSON(jsonInput, url, path, method string, t *testing.T) *http.Response {
	var jsonStr = []byte(jsonInput)
	req, err := http.NewRequest(method, url+path, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Could not make server request: " + err.Error())
	}

	return resp
}

func setUpServer(t *testing.T) *httptest.Server {
	client := getDBClient(t)
	mux := server.CreateRoutingMux(client)
	return httptest.NewServer(mux)
}

func getDBClient(t *testing.T) *db.Client {
	client, err := db.NewClient(localDBAddr)
	require.NoError(t, err)

	// Clear the user table for the tests.
	r.DB("test").TableDrop("users").Exec(client.Session)
	r.Wait(r.WaitOpts{WaitFor: "ready_for_writes", Timeout: 300})
	r.DB("test").TableCreate("users").Exec(client.Session)
	r.Wait(r.WaitOpts{WaitFor: "ready_for_writes", Timeout: 300})
	return client
}
