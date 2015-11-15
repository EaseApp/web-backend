package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EaseApp/web-backend/src/app/controllers/applicationcontroller"
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
			input:            `{"username": "theuser", "password": "pass"}`,
			expectedCode:     http.StatusOK,
			expectedError:    "",
			expectedUsername: "theuser",
		},
		{
			input:            `{"username": "theuser", "password": "pass"}`,
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

func TestDeleteApplication(t *testing.T) {
	server := setUpServer(t)
	defer server.Close()

	testUser := createTestUser(server.URL, t)

	// Create two applications.
	resp := sendJSON("", testUser.APIToken, server.URL, "/users/applications/bestappevar", "POST", t)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp = sendJSON("", testUser.APIToken, server.URL, "/users/applications/lol", "POST", t)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	testcases := []struct {
		token         string
		appToDelete   string
		appNames      []string
		expectedCode  int
		expectedError string
	}{
		// Invalid token.
		{
			token:         "badtoken",
			appToDelete:   "lol",
			appNames:      nil,
			expectedCode:  http.StatusUnauthorized,
			expectedError: "Authorization token does not match.",
		},
		// Invalid app name.
		{
			token:         testUser.APIToken,
			appToDelete:   "idontexist",
			appNames:      nil,
			expectedCode:  http.StatusInternalServerError,
			expectedError: "Failed to delete application",
		},
		// Valid token and the app deleted.
		{
			token:         testUser.APIToken,
			appToDelete:   "bestappevar",
			appNames:      []string{"lol"},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		// Valid token and both apps deleted.
		{
			token:         testUser.APIToken,
			appToDelete:   "lol",
			appNames:      []string{},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
	}

	for _, testcase := range testcases {
		resp := sendJSON("", testcase.token, server.URL, "/users/applications/"+testcase.appToDelete, "DELETE", t)

		assert.Equal(t, testcase.expectedCode, resp.StatusCode)

		// Error expected.
		if testcase.expectedError != "" {
			var errStruct errorResp
			json.NewDecoder(resp.Body).Decode(&errStruct)
			assert.Equal(t, testcase.expectedError, errStruct.Err)
		} else { // No error expected.
			var apps []models.Application
			json.NewDecoder(resp.Body).Decode(&apps)
			for i, app := range apps {
				assert.Equal(t, testcase.appNames[i], app.Name)
				assert.NotEmpty(t, app.AppToken)
			}
		}
	}
}

func TestListApplications(t *testing.T) {
	server := setUpServer(t)
	defer server.Close()

	testUser := createTestUser(server.URL, t)

	// Create two applications.
	resp := sendJSON("", testUser.APIToken, server.URL, "/users/applications/bestappevar", "POST", t)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp = sendJSON("", testUser.APIToken, server.URL, "/users/applications/lol", "POST", t)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	testcases := []struct {
		token         string
		appNames      []string
		expectedCode  int
		expectedError string
	}{
		// Invalid token.
		{
			token:         "badtoken",
			appNames:      nil,
			expectedCode:  http.StatusUnauthorized,
			expectedError: "Authorization token does not match.",
		},
		// Valid token and two apps returned.
		{
			token:         testUser.APIToken,
			appNames:      []string{"bestappevar", "lol"},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
	}

	for _, testcase := range testcases {
		resp := sendJSON("", testcase.token, server.URL, "/users/applications", "GET", t)

		assert.Equal(t, testcase.expectedCode, resp.StatusCode)

		// Error expected.
		if testcase.expectedError != "" {
			var errStruct errorResp
			json.NewDecoder(resp.Body).Decode(&errStruct)
			assert.Equal(t, testcase.expectedError, errStruct.Err)
		} else { // No error expected.
			var apps []models.Application
			json.NewDecoder(resp.Body).Decode(&apps)
			for i, app := range apps {
				assert.Equal(t, testcase.appNames[i], app.Name)
				assert.NotEmpty(t, app.AppToken)
			}
		}
	}
}

func TestCreateApplication(t *testing.T) {
	server := setUpServer(t)
	defer server.Close()

	testUser := createTestUser(server.URL, t)

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
			expectedError: "Authorization token does not match.",
		},
		// Valid token and created app.
		{
			token:         testUser.APIToken,
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
}

func TestSignIn(t *testing.T) {
	server := setUpServer(t)
	defer server.Close()

	testUser := createTestUser(server.URL, t)

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
			expectedAPIToken: testUser.APIToken,
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

func TestSaveReadAndDeleteAppDataEndpoints(t *testing.T) {
	// This test isn't as extensive as some of the other ones because these components are already
	// tested well in models/application_test.

	server := setUpServer(t)
	syncServer := setUpSyncServer(t)

	defer server.Close()
	defer syncServer.Close()

	applicationcontroller.TestingOnlySetSyncServerURL(syncServer.URL)

	appToken := createTestApplication(server.URL, t)

	resp := sendJSON(`{"path":"/hello", "data":{"nested":"objects", "yes": 1}}`,
		appToken, server.URL, "/data/ronswanson/bestappevar", "POST", t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp = sendJSON("", appToken, server.URL, "/data/ronswanson/bestappevar?path=/hello", "GET", t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data1 interface{}
	err := json.NewDecoder(resp.Body).Decode(&data1)
	assert.NoError(t, err)
	assert.Equal(t, interface{}(map[string]interface{}{"nested": "objects", "yes": float64(1)}), data1)

	resp = sendJSON(`{"path":"/hello/yes"}`,
		appToken, server.URL, "/data/ronswanson/bestappevar", "DELETE", t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp = sendJSON("", appToken, server.URL, "/data/ronswanson/bestappevar?path=/hello", "GET", t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data2 interface{}
	err = json.NewDecoder(resp.Body).Decode(&data2)
	assert.NoError(t, err)
	assert.Equal(t, interface{}(map[string]interface{}{"nested": "objects"}), data2)

	// Verify data can't be accessed with a bad token.
	resp = sendJSON("", "iamtoken", server.URL, "/data/ronswanson/bestappevar?path=/hello", "GET", t)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var data3 interface{}
	err = json.NewDecoder(resp.Body).Decode(&data3)
	assert.NoError(t, err)
	assert.Equal(t, interface{}(map[string]interface{}{"error_code": float64(401), "error": "Invalid application token"}), data3)
}

var testUserUsername = "ronswanson"
var testUserPass = "meat"
var testAppName = "bestappevar"

// createTestApplication creates a test application (under a new user) and returns its app token.
func createTestApplication(url string, t *testing.T) string {
	testUser := createTestUser(url, t)

	resp := sendJSON("", testUser.APIToken, url, "/users/applications/bestappevar", "POST", t)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var app models.Application
	err := json.NewDecoder(resp.Body).Decode(&app)
	require.NoError(t, err)
	return app.AppToken
}

// createTestUser creates a test user and returns its api token.
func createTestUser(url string, t *testing.T) models.User {
	resp := sendJSON(`{"username":"ronswanson","password":"meat"}`,
		"", url, "/users/sign_up", "POST", t)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var user models.User
	err := json.NewDecoder(resp.Body).Decode(&user)
	require.NoError(t, err)
	return user
}

func createTestApp(url, apiToken, appName string, t *testing.T) models.Application {
	resp := sendJSON("", apiToken, url, "/users/applications/"+appName, "POST", t)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var application models.Application
	err := json.NewDecoder(resp.Body).Decode(&application)
	require.NoError(t, err)
	// log.Println("APPLICATION:", application)
	return application
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

func setUpServer(t *testing.T) *httptest.Server {
	client := getDBClient(t)
	mux := server.NewEaseServer(client)
	return httptest.NewServer(mux)
}

func getDBClient(t *testing.T) *db.Client {
	client, err := db.NewClient(localDBAddr)
	require.NoError(t, err)

	// Wait for the db to be ready.  Needed for Travis.
	r.Wait().Exec(client.Session)

	// Clear the user table for the tests.
	r.DB("test").Table("users").Delete().RunWrite(client.Session)

	return client
}
