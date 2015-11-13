package integration

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EaseApp/web-backend/src/sync"
	"github.com/gorilla/websocket"
)

/*
func TestSocketConnection(t *testing.T) {
	syncServer := setUpSyncServer(t)
	webServer, client := setUpServer(t)
	apiToken := createTestUser(webServer.URL, t)
	appToken := createTestApp(webServer.URL, apiToken, "test", t)

	testcases := []struct {
		subscribeTo                   string
		publishTo                     string
		publishData                   string
		expectedData                  string
		expectedApplicationStatusCode int
	}{
		{
			// Subscribe to an application, get the relavent data
			subscribeTo:                   "ronswanson_test",
			publishTo:                     "test",
			publishData:                   `{"data":"test"}`,
			expectedData:                  `{"data":"test"}`,
			expectedApplicationStatusCode: http.StatusOK,
		},
		{
			// Subscribe to one application, shouldn't get data with bad Publish call.
			subscribeTo:                   "wrongAppKey",
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

		sendSocketData(conn, testcase.subscribeTo)
		assert.Equal(t, testcase.subscribeTo, grabSocketData(conn))

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
*/

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
