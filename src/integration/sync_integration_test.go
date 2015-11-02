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
)

func TestSocketConnection(t *testing.T) {
	server := setUpSyncServer(t)

	testcases := []struct {
		connectionDb string
		publishData  string
	}{
		{
			connectionDb: "test",
			publishData:  `{"data":"test"}`,
		},
	}

	for _, testcase := range testcases {
		port := strings.Split(server.URL, ":")[2]
		conn := openConnection("ws://localhost:" + port + "/sub")
		defer conn.Close()

		sendSocketData(conn, testcase.connectionDb)
		assert.Equal(t, testcase.connectionDb, grabSocketData(conn))
	}
	defer server.Close()
}

func grabSocketData(conn *websocket.Conn) string {
	_, p, _ := conn.ReadMessage()
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
	conn, resp, err := DefaultDialer.Dial(url, header)
	if err != nil {
		log.Println(err)
	}
	log.Println("Resp: ", resp)
	return conn
}

func setUpSyncServer(t *testing.T) *httptest.Server {
	mux := sync.NewSyncServer()
	log.Println("Sync server running")
	return httptest.NewServer(mux)
}
