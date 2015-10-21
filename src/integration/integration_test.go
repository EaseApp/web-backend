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
	server, _ := setUpServer(t)
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
		resp := sendJSON(testcase.input, "", server.URL, "/users/sign_up", "POST", t)

		assert.Equal(t, testcase.expectedCode, resp.StatusCode)

		// Error expected.
		if testcase.expectedError != "" {
			var errStruct errorResp
			json.NewDecoder(resp.Body).Decode(&errStruct)
			assert.Equal(t, testcase.expectedError, errStruct.Err)
		} else { // No error expected.
			var userStruct models.User
			json.NewDecoder(resp.Body).Decode(&userStruct)
			assert.Equal(t, testcase.expectedUsername, userStruct.Username)
		}
	}
}

func TestCreateApplication(t *testing.T) {
	server, client := setUpServer(t)
	defer server.Close()

	apiToken := createTestUser(server.URL, t)

	testcases := []struct {
		token         string
		appName       string
		expectedCode  int
		expectedError string
	}{
		// Invalid token.
		{
			token:         "badtoken",
			appName:       "lol",
			expectedCode:  http.StatusUnauthorized,
			expectedError: "Not authorized",
		},
		// Valid token and created app.
		{
			token:         apiToken,
			appName:       "bestappevar",
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
	}

	for _, testcase := range testcases {
		resp := sendJSON("", testcase.token, server.URL, "/users/applications/"+testcase.appName, "POST", t)

		assert.Equal(t, testcase.expectedCode, resp.StatusCode)

		// Error expected.
		if testcase.expectedError != "" {
			var errStruct errorResp
			json.NewDecoder(resp.Body).Decode(&errStruct)
			assert.Equal(t, testcase.expectedError, errStruct.Err)
		} else { // No error expected.
			var app models.Application
			json.NewDecoder(resp.Body).Decode(&app)
			assert.Equal(t, testcase.appName, app.Name)
			assert.NotEmpty(t, app.AppToken)
		}
	}

	// Delete the created application table.
	r.DB("test").TableDrop("ronswanson_bestappevar").RunWrite(client.Session)
}

func TestSignIn(t *testing.T) {
	server, _ := setUpServer(t)
	defer server.Close()

	userAPIToken := createTestUser(server.URL, t)

	testcases := []struct {
		input            string
		expectedCode     int
		expectedError    string
		expectedUsername string
		expectedAPIToken string
	}{
		// Valid login.
		{
			input:            `{"username": "ronswanson", "password": "meat"}`,
			expectedCode:     http.StatusOK,
			expectedError:    "",
			expectedUsername: "ronswanson",
			expectedAPIToken: userAPIToken,
		},
		// Invalid username.
		{
			input:            `{"username": "anneperkins", "password": "pass"}`,
			expectedCode:     http.StatusUnauthorized,
			expectedError:    "Couldn't find user with that username",
			expectedUsername: "",
			expectedAPIToken: "",
		},
		// Invalid password.
		{
			input:            `{"username": "ronswanson", "password": "pass"}`,
			expectedCode:     http.StatusUnauthorized,
			expectedError:    "Password was invalid",
			expectedUsername: "",
			expectedAPIToken: "",
		},
	}

	for _, testcase := range testcases {
		resp := sendJSON(testcase.input, "", server.URL, "/users/sign_in", "POST", t)

		assert.Equal(t, testcase.expectedCode, resp.StatusCode)

		// Error expected.
		if testcase.expectedError != "" {
			var errStruct errorResp
			json.NewDecoder(resp.Body).Decode(&errStruct)
			assert.Equal(t, testcase.expectedError, errStruct.Err)
		} else { // No error expected.
			var userStruct models.User
			json.NewDecoder(resp.Body).Decode(&userStruct)
			assert.Equal(t, testcase.expectedUsername, userStruct.Username)
			assert.Equal(t, testcase.expectedAPIToken, userStruct.APIToken)
		}
	}

}

var testUserUsername = "ronswanson"
var testUserPass = "meat"

// createTestUser creates a test user and returns its api token.
func createTestUser(url string, t *testing.T) string {
	resp := sendJSON(`{"username":"ronswanson","password":"meat"}`,
		"", url, "/users/sign_up", "POST", t)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var user models.User
	err := json.NewDecoder(resp.Body).Decode(&user)
	require.NoError(t, err)
	return user.APIToken
}

func sendJSON(jsonInput, token, url, path, method string, t *testing.T) *http.Response {
	var jsonStr = []byte(jsonInput)
	req, err := http.NewRequest(method, url+path, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Could not make server request: " + err.Error())
	}

	return resp
}

func setUpServer(t *testing.T) (*httptest.Server, *db.Client) {
	client := getDBClient(t)
	mux := server.NewEaseServer(client)
	return httptest.NewServer(mux), client
}

func getDBClient(t *testing.T) *db.Client {
	client, err := db.NewClient(localDBAddr)
	require.NoError(t, err)

	// Wait for the db to be ready.  Needed for Travis.
	r.Wait().Exec(client.Session)

	// Clear the user table for the tests.
	r.DB("test").Table("users").Delete().Run(client.Session)

	return client
}
