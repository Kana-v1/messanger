package server

import (
	"messanger/internal/logs"
	"messanger/pkg/connection"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

func WebSocketHandler(c echo.Context) error {
	// body := make([]byte, 0)
	// _, err := req.Body.Read(body)
	// if err != nil {
	// 	logs.ErrorLog("", fmt.Sprintf("Invalid request body: %v, fileName: %s, method:%s", string(body), main.go, webSocketHandler), err)
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	webSockerConn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logs.FatalLog("", "websocket connection failed", err)
		return c.String(http.StatusInternalServerError, "")
	}

	user := connection.NewUser()
	peer := &connection.Peer{
		Conn:     webSockerConn,
		Id:       3,
		IsClosed: false,
	}

	user.Peers = append(user.Peers, *peer)
	user.Start(peer)
	return c.String(http.StatusOK, "")
}
