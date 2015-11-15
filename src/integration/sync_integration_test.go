package integration

import (
	"github.com/EaseApp/web-backend/src/app/controllers/applicationcontroller"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/EaseApp/web-backend/src/sync"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSocketConnection(t *testing.T) {
	syncServer := setUpSyncServer(t)
	webServer := setUpServer(t)
	testUser := createTestUser(webServer.URL, t)
	testApplication := createTestApp(webServer.URL, testUser.APIToken, "test", t)
	applicationcontroller.TestingOnlySetSyncServerURL(syncServer.URL)

	testcases := []struct {
		firstSocketMessage            string
		firstSocketResponse           string
		publishTo                     string
		publishData                   string
		expectedData                  string
		expectedApplicationStatusCode int
	}{
		{
			// Subscribe to an application, get the relavent data
			firstSocketMessage:            `{"username": "` + testUser.Username + `", "table_name": "` + testApplication.TableName + `", "authorization": "` + testApplication.AppToken + `"}`,
			firstSocketResponse:           `{"status": "success"}`,
			publishTo:                     testApplication.TableName,
			publishData:                   `{"path":"/hello","data": "world"}`,
			expectedData:                  `{"action":"SAVE","data":"world","path":{"OriginalString":"/hello","TopLevelDocName":"hello","RemainingSegments":[]}}`,
			expectedApplicationStatusCode: http.StatusOK,
		},
	}

	for _, testcase := range testcases {
		port := strings.Split(syncServer.URL, ":")[2]
		conn := openConnection("ws://localhost:" + port + "/sub") // Connect to sync annonymously
		defer conn.Close()

		sendSocketData(conn, testcase.firstSocketMessage)
		assert.Equal(t, testcase.firstSocketResponse, grabSocketData(conn), "Socket response failed")

		resp := sendJSON(testcase.publishData,
			testApplication.AppToken, webServer.URL, "/data/"+testUser.Username+"/"+testApplication.TableName, "POST", t)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		actual := grabSocketData(conn)
		assert.Equal(t, testcase.expectedData, strings.Trim(actual, "\n"))

		resp = sendJSON(testcase.publishData,
			testApplication.AppToken, webServer.URL, "/data/"+testUser.Username+"/"+testApplication.TableName, "DELETE", t)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		actual = grabSocketData(conn)
		expectedDeleteResponse := `{"action":"DELETE","data":"world","path":{"OriginalString":"/hello","TopLevelDocName":"hello","RemainingSegments":[]}}`
		assert.Equal(t, expectedDeleteResponse, strings.Trim(actual, "\n"))
	}
	defer syncServer.Close()
}

func TestSocketAuth(t *testing.T) {
	syncServer := setUpSyncServer(t)
	webServer := setUpServer(t)
	testUser := createTestUser(webServer.URL, t)
	testApplication := createTestApp(webServer.URL, testUser.APIToken, "test", t)
	applicationcontroller.TestingOnlySetSyncServerURL(syncServer.URL)

	testcases := []struct {
		firstSocketMessage            string
		firstSocketResponse           string
		publishTo                     string
		publishData                   string
		expectedData                  string
		expectedApplicationStatusCode int
	}{
		{
			// Subscribe to one application, shouldn't get data with bad Publish call.
			firstSocketMessage:            `{"username": "ronswanson", "table_name": "test", "authorization": "123"}`,
			firstSocketResponse:           `{"status": "failed"}`,
			publishTo:                     "user_differentApp",
			publishData:                   `{"data":"test"}`,
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
		resp := sendJSON(testcase.publishData, testApplication.AppToken, syncServer.URL, path, "POST", t) // Publish to an app
		assert.Equal(t, testcase.expectedApplicationStatusCode, resp.StatusCode)
	}
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
	// log.Println("URL: ", url)
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
	mux := sync.NewServer(client)
	log.Println("Sync server running...")
	return httptest.NewServer(mux)
}
