package integration

import (
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
			// Subscribe to an application, get the relavent data
			subscribeTo:                   testUser.ID + "_" + "test",
			firstSocketMessage:            `{"username": "` + testUser.ID + `", "application": "test", "authorization": "` + testApplication.AppToken + `"}`,
			firstSocketResponse:           `{"status": "success"}`,
			publishTo:                     "test",
			publishData:                   `{"data": "test"}`,
			expectedData:                  `{"data": "test"}`,
			expectedApplicationStatusCode: http.StatusOK,
		},
		{
			// Subscribe to one application, shouldn't get data with bad Publish call.
			subscribeTo:                   "ronswanson_test",
			firstSocketMessage:            `{"username": "ronswanson", "application": "test", "authorization": "123"}`,
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
		resp := sendJSON(testcase.publishData, testApplication.AppToken, syncServer.URL, path, "POST", t) // Publish to an app
		assert.Equal(t, testcase.expectedApplicationStatusCode, resp.StatusCode)
		actual := grabSocketData(conn)
		assert.Equal(t, testcase.expectedData, actual)

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
