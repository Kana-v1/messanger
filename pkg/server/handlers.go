package server

import (
	"messanger/internal/logs"
	"messanger/pkg/authorization/jwt"
	"messanger/pkg/connection"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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

func SignIn(c echo.Context) error {
	return jwt.SignIn(c)
}

func SignUp(c echo.Context) error {
	return jwt.SignUp(c)
}

func RefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwt.RefreshToken(c)
		if err := next(c); err != nil {
			return err
		}
		return nil
	}
}

func IsAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		if err, _ := jwt.IsAuthorized(c); err != nil {
			return c.String(http.StatusUnauthorized, "You are unauthorized")
		}

		if err := next(c); err != nil {
			return err
		}
		return nil
	}
}
