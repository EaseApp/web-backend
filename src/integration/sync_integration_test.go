package integration

import (
	"github.com/EaseApp/web-backend/src/sync"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSocketConnection(t *testing.T) {
	server := setUpSyncServer(t)

	testcases := []struct {
		subscribeTo  string
		publishTo    string
		publishData  string
		expectedData string
		conn         *websocket.Conn
	}{
		{
			// Subscribt to an application, get the relavent data
			subscribeTo:  "test",
			publishTo:    "test",
			publishData:  `{"data":"test"}`,
			expectedData: `{"data":"test"}`,
		},
		{
			// Subscribe to one application, make sure you don't get another app's data
			subscribeTo:  "anApp",
			publishTo:    "differentApp",
			publishData:  `{"data":"test"}`,
			expectedData: "",
		},
	}

	for _, testcase := range testcases {
		port := strings.Split(server.URL, ":")[2]
		testcase.conn = openConnection("ws://localhost:" + port + "/sub")
		defer testcase.conn.Close()

		sendSocketData(testcase.conn, testcase.subscribeTo)
		assert.Equal(t, testcase.subscribeTo, grabSocketData(testcase.conn))

		resp := sendJSON(testcase.publishData, "", server.URL, "/pub/"+testcase.publishTo, "POST", t)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		actual := grabSocketData(testcase.conn)
		assert.Equal(t, testcase.expectedData, actual)
	}
	defer server.Close()
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
	mux := sync.NewSyncServer()
	log.Println("Sync server running...")
	return httptest.NewServer(mux)
}
