package integration

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/EaseApp/web-backend/src/sync"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSocketConnection(t *testing.T) {
	syncServer := setUpSyncServer(t)
	webServer, client := setUpServer(t)
	apiToken := createTestUser(webServer.URL, t)
	appToken := createTestApp(webServer.URL, apiToken, "test", t)

	testcases := []struct {
		subscribeTo                   string
		firstSocketResponse           string
		publishTo                     string
		publishData                   string
		expectedData                  string
		expectedApplicationStatusCode int
	}{
		{
			// Subscribe to an application, get the relavent data
			subscribeTo:                   "ronswanson_test",
			publishTo:                     "test",
			firstSocketResponse:           `{"status": "success"}`,
			publishData:                   `{"data": "test"}`,
			expectedData:                  `{"data": "test"}`,
			expectedApplicationStatusCode: http.StatusOK,
		},
	}

	for _, testcase := range testcases {
		port := strings.Split(syncServer.URL, ":")[2]
		conn := openConnection("ws://localhost:" + port + "/sub") // Connect to sync annonymously
		defer conn.Close()

		sendSocketData(conn, `{"username": "ronswanson", "application": "test", "authorization": "`+appToken+`"}`)
		assert.Equal(t, testcase.firstSocketResponse, grabSocketData(conn), "Socket response failed")

		path := "/pub/ronswanson/" + testcase.publishTo
		// log.Println(syncServer.URL, path)
		resp := sendJSON(testcase.publishData, appToken, syncServer.URL, path, "POST", t) // Publish to an app
		assert.Equal(t, testcase.expectedApplicationStatusCode, resp.StatusCode)
		actual := grabSocketData(conn)
		assert.Equal(t, testcase.expectedData, actual)

	}
	r.DB("test").TableDrop("ronswanson_" + "test").RunWrite(client.Session)
	defer syncServer.Close()
}

// func createApplication(username, appName string, client *db.Client) string {
// 	user, err := models.NewUser(username, "pass")
// 	querier := models.NewModelQuerier(client.Session)
// 	querier.Save(user)
// 	application, err := querier.CreateApplication(user, appName)
// 	return application.AppToken
// }

func TestSocketAuthFailure(t *testing.T) {
	syncServer := setUpSyncServer(t)
	webServer, client := setUpServer(t)
	apiToken := createTestUser(webServer.URL, t)
	appToken := createTestApp(webServer.URL, apiToken, "test", t)

	testcases := []struct {
		subscribeTo                   string
		firstSocketMessage            string
		firstSocketResponse           string
		publishTo                     string
		publishData                   string
		expectedData                  string
		expectedApplicationStatusCode int
	}{
		{
			// Subscribe to one application, shouldn't get data with bad Publish call.
			subscribeTo:                   "wrongAppKey",
			firstSocketMessage:            `{"username": "user", "application": "wrongAppKey", "authorization": "123"}`,
			firstSocketResponse:           `{"status": "failed"}`,
			publishTo:                     "user_differentApp",
			publishData:                   `{"data":"test"}`,
			expectedData:                  "",
			expectedApplicationStatusCode: http.StatusUnauthorized,
		},
	}

	for _, testcase := range testcases {
		port := strings.Split(syncServer.URL, ":")[2]
		conn := openConnection("ws://localhost:" + port + "/sub") // Connect to sync annonymously
		defer conn.Close()

		sendSocketData(conn, testcase.firstSocketMessage)
		assert.Equal(t, testcase.firstSocketResponse, grabSocketData(conn), "Socket response failed")

		path := "/pub/ronswanson/" + testcase.publishTo
		// log.Println(syncServer.URL, path)
		resp := sendJSON(testcase.publishData, appToken, syncServer.URL, path, "POST", t) // Publish to an app
		assert.Equal(t, testcase.expectedApplicationStatusCode, resp.StatusCode, "Bad response from endpoint")
		actual := grabSocketData(conn)
		assert.Equal(t, testcase.expectedData, actual, "Publish data does not match this data")

	}
	r.DB("test").TableDrop("ronswanson_" + "test").RunWrite(client.Session)
	defer syncServer.Close()
}

func grabSocketData(conn *websocket.Conn) string {
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, p, err := conn.ReadMessage() // This methods blocks. Make sure to set Deadline
	if err != nil {
		log.Println(err)
	}
	return string(p)

}

func sendSocketData(conn *websocket.Conn, data string) error {
	return conn.WriteMessage(1, []byte(data))
}

func openConnection(url string) *websocket.Conn {
	log.Println("URL: ", url)
	var DefaultDialer = &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
	}

	header := make(http.Header)
	conn, _, err := DefaultDialer.Dial(url, header)
	if err != nil {
		log.Println(err)
	}
	return conn
}

func setUpSyncServer(t *testing.T) *httptest.Server {
	client := getDBClient(t)
	mux := sync.NewSyncServer(client)
	log.Println("Sync server running...")
	return httptest.NewServer(mux)
}
