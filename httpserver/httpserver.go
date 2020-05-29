package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	l "github.com/infinit-lab/yolanda/logutils"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ResponseBody struct {
	Result bool   `json:"result"`
	Error  string `json:"error"`
}

type IWebsocketHandler interface {
	NewConnection(nodeId int, r *http.Request)
	Disconnected(nodeId int)
	ReadBytes(nodeId int, bytes []byte)
	ReadMessage(nodeId int, message []byte)
}

type IFilter interface {
	Filter(r *http.Request, checkToken bool) error
}

type ITokenChecker interface {
	CheckToken(r *http.Request) error
}

func ListenAndServe() error {
	server.isServe = true
	return server.server.ListenAndServe()
}

func Shutdown() error {
	server.isServe = false
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	return server.server.Shutdown(ctx)
}

func RegisterHttpHandlerFunc(method string, url string, function http.HandlerFunc, checkToken bool) {
	server.registerHttpHandlerFunc(method, url, function, checkToken)
}

func RegisterWebsocketHandler(url string, handler IWebsocketHandler, checkToken bool) {
	server.registerWebsocketHandler(url, handler, checkToken)
}

func WebsocketWriteBytes(nodeId int, bytes []byte) error {
	return server.writeBytes(nodeId, bytes)
}

func WebsocketWriteMessage(nodeId int, message []byte) error {
	return server.writeMessage(nodeId, message)
}

func WebsocketClose(nodeId int) error {
	return server.closeSocket(nodeId)
}

func RegisterFilter(filter IFilter) {
	server.registerFilter(filter)
}

func RegisterTokenChecker(checker ITokenChecker) {
	server.registerTokenChecker(checker)
}

func GetRequestBody(r *http.Request, body interface{}) error {
	defer func() {
		_ = r.Body.Close()
	}()
	str, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(str, body); err != nil {
		return err
	}
	return nil
}

func RemoteIp(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("XRealIP"); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get("XForwardedFor"); ip != "" {
		remoteAddr = ip
	} else {
		var err error
		remoteAddr, _, err = net.SplitHostPort(remoteAddr)
		if err != nil {
			return ""
		}
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func ResponseError(w http.ResponseWriter, err string, code int) {
	var response ResponseBody
	response.Result = false
	response.Error = err
	data, _ := json.Marshal(&response)
	http.Error(w, string(data), code)
}

func Response(w http.ResponseWriter, body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		l.Error("Failed to write response. error: ", err)
	}
}

func GetId(url string, key string) string {
	path := strings.Split(url, "/")
	for i, v := range path {
		if v == key && (i+1) < len(path) {
			return path[i+1]
		}
	}
	return ""
}

func GetIdInt(url string, key string) (int, error) {
	id := GetId(url, key)
	if id == "" {
		return 0, errors.New("无效ID")
	}
	tempId, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.New("无效ID")
	}
	return tempId, nil
}
