package server

import (
	"messanger/internal/logs"
	"messanger/pkg/connection"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

func WebSocketHandler(rw http.ResponseWriter, req *http.Request) {
	// body := make([]byte, 0)
	// _, err := req.Body.Read(body)
	// if err != nil {
	// 	logs.ErrorLog("", fmt.Sprintf("Invalid request body: %v, fileName: %s, method:%s", string(body), main.go, webSocketHandler), err)
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	webSockerConn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logs.FatalLog("", "websocket connection failed", err)
		return
	}

	user := connection.NewUser()
	peer := &connection.Peer{
		Conn:     webSockerConn,
		Id:       3,
		IsClosed: false,
	}

	user.Peers = append(user.Peers, *peer)
	user.Start(peer)
}
