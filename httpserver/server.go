package httpserver

import (
	"github.com/infinit-lab/yolanda/config"
	l "github.com/infinit-lab/yolanda/logutils"
	"container/list"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"sync"
)

type handler struct {
	handler    http.HandlerFunc
	checkToken bool
}

type httpHandler struct {
	httpRoutes    map[string]*httpHandler
	getHandler    handler
	putHandler    handler
	postHandler   handler
	deleteHandler handler
}

type websocketHandler struct {
	socketRoutes map[string]*websocketHandler
	handler      IWebsocketHandler
	checkToken   bool
}

type httpServer struct {
	server       http.Server
	fileHandler  http.Handler
	isServe      bool
	httpRoutes   map[string]*httpHandler
	socketRoutes map[string]*websocketHandler
	upgrader     websocket.Upgrader
	mutex        sync.Mutex
	index        int
	sockets      map[int]*websocket.Conn
	filters      *list.List
	tokenChecker ITokenChecker
}

var server httpServer

func init() {
	port := config.GetInt("server.port")
	l.Trace("Get port ", port)
	if port == 0 {
		port = 8088
		l.Info("Port reset to ", port)
	}
	dir := config.GetString("server.dir")
	l.Trace("Get dir ", dir)
	if len(dir) == 0 {
		dir = "www"
		l.Info("Dir reset to ", dir)
	}
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	server.server = http.Server{
		Addr:    addr,
		Handler: &server,
	}
	server.fileHandler = http.FileServer(http.Dir(dir))
	server.httpRoutes = make(map[string]*httpHandler)
	server.isServe = false
	server.socketRoutes = make(map[string]*websocketHandler)
	server.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	server.index = 0
	server.sockets = make(map[int]*websocket.Conn)
	server.filters = list.New()
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.TraceF("Request method is %s, url is %s", r.Method, r.URL.Path)

	httpFunction, checkToken := s.httpFunc(r.Method, r.URL.Path)
	if httpFunction != nil {
		if checkToken && s.tokenChecker != nil {
			if err := s.tokenChecker.CheckToken(r); err != nil {
				l.Error("Failed to check token. error: ", err)
				ResponseError(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}

		for f := s.filters.Front(); f != nil; f = f.Next() {
			if filter, ok := f.Value.(IFilter); ok {
				if err := filter.Filter(r); err != nil {
					l.Error("Failed to filter. error: ", err)
					ResponseError(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		}
		httpFunction(w, r)
		return
	}

	websocketHandler, checkToken := s.websocketHandler(r.URL.Path)
	if websocketHandler != nil {
		if checkToken && s.tokenChecker != nil {
			if err := s.tokenChecker.CheckToken(r); err != nil {
				l.Error("Failed to check token. error: ", err)
				ResponseError(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}
		ws, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			l.Error("Failed to upgrade websocket. error: ", err)
			ResponseError(w, "升级Websocket失败", http.StatusBadRequest)
			return
		}
		nodeId := s.insertSocket(ws)
		defer func() {
			_ = ws.Close()
			s.removeSocket(nodeId)
			websocketHandler.Disconnected(nodeId)
		}()

		websocketHandler.NewConnection(nodeId, r)
		s.loop(nodeId, ws, websocketHandler)
		return
	}
	s.fileHandler.ServeHTTP(w, r)
}

func (s *httpServer) registerHttpHandlerFunc(method string, url string, handler http.HandlerFunc, checkToken bool) {
	if s.isServe {
		return
	}
	path := strings.Split(url, "/")
	routes := s.httpRoutes
	for key, value := range path {
		v, ok := routes[value]
		if !ok {
			child := new(httpHandler)
			child.httpRoutes = make(map[string]*httpHandler)
			routes[value] = child
			v = child
		}
		if (key + 1) == len(path) {
			switch method {
			case http.MethodGet:
				v.getHandler.handler = handler
				v.getHandler.checkToken = checkToken
				break
			case http.MethodPut:
				v.putHandler.handler = handler
				v.putHandler.checkToken = checkToken
				break
			case http.MethodPost:
				v.postHandler.handler = handler
				v.postHandler.checkToken = checkToken
				break
			case http.MethodDelete:
				v.deleteHandler.handler = handler
				v.deleteHandler.checkToken = checkToken
				break
			default:
				break
			}
		} else {
			routes = v.httpRoutes
		}
	}
}

func (s *httpServer) httpFunc(method string, url string) (http.HandlerFunc, bool) {
	path := strings.Split(url, "/")
	f, checkToken := s.fullyMatch(method, path)
	if f == nil {
		f, checkToken = s.partlyMatch(method, path)
	}
	return f, checkToken
}

func (s *httpServer) fullyMatch(method string, path []string) (http.HandlerFunc, bool) {
	routes := s.httpRoutes
	isLast := false
	for key, value := range path {
		v, ok := routes[value]
		if !ok {
			v, ok = routes["#"]
			if !ok {
				return nil, false
			} else {
				isLast = true
			}
		}
		if (key+1) == len(path) || isLast {
			return funcByMethod(method, *v)
		} else {
			routes = v.httpRoutes
		}
	}
	return nil, false
}

func (s *httpServer) partlyMatch(method string, path []string) (http.HandlerFunc, bool) {
	routes := s.httpRoutes
	isLast := false
	for key, value := range path {
		v, ok := routes["+"]
		if !ok {
			v, ok = routes[value]
			if !ok {
				v, ok = routes["#"]
				if !ok {
					return nil, false
				} else {
					isLast = true
				}
			}
		}
		if (key+1) == len(path) || isLast {
			return funcByMethod(method, *v)
		} else {
			routes = v.httpRoutes
		}
	}
	return nil, false
}

func funcByMethod(method string, h httpHandler) (http.HandlerFunc, bool) {
	var temp handler
	switch method {
	case http.MethodGet:
		temp = h.getHandler
		break
	case http.MethodPut:
		temp = h.putHandler
		break
	case http.MethodPost:
		temp = h.postHandler
		break
	case http.MethodDelete:
		temp = h.deleteHandler
		break
	default:
		break
	}
	return temp.handler, temp.checkToken
}

func (s *httpServer) registerWebsocketHandler(url string, handler IWebsocketHandler, checkToken bool) {
	if s.isServe {
		return
	}
	path := strings.Split(url, "/")
	routes := s.socketRoutes
	for key, value := range path {
		v, ok := routes[value]
		if !ok {
			child := new(websocketHandler)
			child.socketRoutes = make(map[string]*websocketHandler)
			routes[value] = child
			v = child
		}
		if (key + 1) == len(path) {
			v.handler = handler
			v.checkToken = checkToken
		} else {
			routes = v.socketRoutes
		}
	}
}

func (s *httpServer) websocketHandler(url string) (IWebsocketHandler, bool) {
	path := strings.Split(url, "/")
	h, checkToken := s.fullyMatchWebsocketHandler(path)
	if h == nil {
		h, checkToken = s.partlyMatchWebsocketHandler(path)
	}
	return h, checkToken
}

func (s *httpServer) fullyMatchWebsocketHandler(path []string) (IWebsocketHandler, bool) {
	routes := s.socketRoutes
	isLast := false
	for key, value := range path {
		v, ok := routes[value]
		if !ok {
			v, ok = routes["#"]
			if !ok {
				return nil, false
			} else {
				isLast = true
			}
		}
		if (key+1) == len(path) || isLast {
			return v.handler, v.checkToken
		} else {
			routes = v.socketRoutes
		}
	}
	return nil, false
}

func (s *httpServer) partlyMatchWebsocketHandler(path []string) (IWebsocketHandler, bool) {
	routes := s.socketRoutes
	isLast := false
	for key, value := range path {
		v, ok := routes["+"]
		if !ok {
			v, ok = routes[value]
			if !ok {
				v, ok = routes["#"]
				if !ok {
					return nil, false
				} else {
					isLast = true
				}
			}
		}
		if (key+1) == len(path) || isLast {
			return v.handler, v.checkToken
		} else {
			routes = v.socketRoutes
		}
	}
	return nil, false
}

func (s *httpServer) insertSocket(socket *websocket.Conn) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.index++
	nodeId := s.index
	s.sockets[nodeId] = socket
	return nodeId
}

func (s *httpServer) removeSocket(nodeId int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sockets, nodeId)
}

func (s *httpServer) loop(nodeId int, ws *websocket.Conn, handler IWebsocketHandler) {
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			l.Error("Socket disconnected. error: ", err)
			return
		}
		switch messageType {
		case websocket.TextMessage:
			handler.ReadMessage(nodeId, message)
			break
		case websocket.BinaryMessage:
			handler.ReadBytes(nodeId, message)
			break
		default:
			break
		}
	}
}

func (s *httpServer) socket(nodeId int) *websocket.Conn {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.sockets[nodeId]
}

func (s *httpServer) writeMessage(nodeId int, message []byte) error {
	ws := s.socket(nodeId)
	if ws == nil {
		return errors.New("获取Websocket失败")
	}
	return ws.WriteMessage(websocket.TextMessage, message)
}

func (s *httpServer) writeBytes(nodeId int, bytes []byte) error {
	ws := s.socket(nodeId)
	if ws == nil {
		return errors.New("获取Websocket失败")
	}
	return ws.WriteMessage(websocket.BinaryMessage, bytes)
}

func (s *httpServer) closeSocket(nodeId int) error {
	ws := s.socket(nodeId)
	if ws == nil {
		return errors.New("获取Websocket失败")
	}
	return ws.Close()
}

func (s *httpServer) registerFilter(filter IFilter) {
	s.filters.PushBack(filter)
}

func (s *httpServer) registerTokenChecker(checker ITokenChecker) {
	s.tokenChecker = checker
}
