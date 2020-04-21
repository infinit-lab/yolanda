package httpserver

import (
	l "github.com/infinit-lab/yolanda/logutils"
	"flag"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

func handleHttp(w http.ResponseWriter, r *http.Request) {
	l.Trace("Get url ", r.URL.Path)
	_, err := w.Write([]byte(r.URL.Path))
	if err != nil {
		l.Error("Failed to write response")
	}
}

func testHandleHttp(t *testing.T, url string) {
	resp, err := http.Get("http://127.0.0.1:8088" + url)
	if err != nil {
		t.Errorf("Failed to get %s", url)
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("Failed to get %s, code is %d", url, resp.StatusCode)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if string(body) != url {
		t.Errorf("The body is wrong, expect is %s, actual is %s", url, string(body))
		return
	}
	l.TraceF("Success to test %s", url)
}

type TestWebsocketHandler struct {
	nodeId int
}

func (h *TestWebsocketHandler) NewConnection(nodeId int, r *http.Request) {
	l.Trace("New connection. nodeId: ", nodeId)
	h.nodeId = nodeId
}

func (h *TestWebsocketHandler) Disconnected(nodeId int) {
	l.Trace("Disconnected. nodeId: ", nodeId)
}

func (h *TestWebsocketHandler) ReadBytes(nodeId int, bytes []byte) {
	l.Trace("Read bytes: ", string(bytes))
	if err := WebsocketWriteBytes(nodeId, bytes); err != nil {
		l.Error("Write bytes error: ", err)
	} else {
		l.Trace("Write bytes success")
	}
}

func (h *TestWebsocketHandler) ReadMessage(nodeId int, message []byte) {
	l.Trace("Read message: ", string(message))
	if err := WebsocketWriteMessage(nodeId, message); err != nil {
		l.Error("Write message error: ", err)
	} else {
		l.Trace("Write message success")
	}
}

func TestRegisterHttpHandlerFunc(t *testing.T) {
	RegisterHttpHandlerFunc(http.MethodGet, "/api/1/2", handleHttp, false)
	RegisterHttpHandlerFunc(http.MethodGet, "/api/1/+/3", handleHttp, false)
	RegisterHttpHandlerFunc(http.MethodGet, "/api/1/4", handleHttp, false)
	RegisterHttpHandlerFunc(http.MethodGet, "/api/1/4/3/#", handleHttp, false)
}

func TestRegisterWebsocketHandler(t *testing.T) {
	h := new(TestWebsocketHandler)
	RegisterWebsocketHandler("/api/2/+", h, false)
}

func TestListenAndServe(t *testing.T) {
	l.Trace("Start test ListenAndServe")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		l.Trace("Start test request")
		testHandleHttp(t, "/api/1/2")
		testHandleHttp(t, "/api/1/1/3")
		testHandleHttp(t, "/api/1/2/3")
		testHandleHttp(t, "/api/1/4")
		testHandleHttp(t, "/api/1/4/3/1")
		testHandleHttp(t, "/api/1/4/3/2")
	}()

	wg.Add(1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		defer wg.Done()
		addr := flag.String("Websocket Client", "127.0.0.1:8088", "Websocket Server")
		u := url.URL {
			Scheme: "ws",
			Host: *addr,
			Path: "/api/2/1",
		}
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			t.Fatal("Failed to connect to server. error: ", err)
		}
		_ = conn.WriteMessage(websocket.TextMessage, []byte("123456"))
		messageType, message, err := conn.ReadMessage()
		if messageType != websocket.TextMessage {
			t.Error("Read message text type error")
		} else {
			if string(message) != "123456" {
				t.Errorf("Read message expect %s, actual %s", "123456", string(message))
			} else {
				l.Trace("Read message: ", string(message))
			}
		}
		_ = conn.WriteMessage(websocket.BinaryMessage, []byte("123456"))
		messageType, message, err = conn.ReadMessage()
		if messageType != websocket.BinaryMessage {
			t.Error("Read message binary type error")
		} else {
			if string(message) != "123456" {
				t.Errorf("Read message expect %s, actual %s", "123456", string(message))
			} else {
				l.Trace("Read binary: ", string(message))
			}
		}
		_ = conn.Close()
		time.Sleep(100 * time.Millisecond)
	}()

	go func() {
		wg.Wait()
		_ = Shutdown()
	}()

	l.Trace("ListenAndServe")
	_ = ListenAndServe()
}

