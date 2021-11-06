package main

import (
	"context"
	"fmt"
	_ "messanger/configs"
	"messanger/pkg/connection"
	"net/http"
	"os"
	"os/signal"
	"time"

	"messanger/internal/logs"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}
var port string
var peer *connection.Peer

func main() {
	a := os.Getenv("REDIS_HOST")
	if a == "" {
		a = "123"
	}
	fmt.Println(a)

	port = os.Getenv("PORT")
	if port == "" {
		logs.FatalLog("", "missing PORT env var", nil)
	}

	http.Handle("/chat/", http.HandlerFunc(webSocketHandler))
	server := http.Server{Addr: ":" + port, Handler: nil}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logs.FatalLog("", "failed to start server", err)
		}
	}()
	dontStop := make(chan int)
	<-dontStop
	stopServer(server)
}

func webSocketHandler(rw http.ResponseWriter, req *http.Request) {
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
	peer = &connection.Peer{
		Conn:     webSockerConn,
		Id:       3,
		IsClosed: false,
	}

	user.Peers = append(user.Peers, *peer)
	user.Start(peer)
}

func stopServer(s http.Server) {
	stop := make(chan os.Signal)
	signal.Notify(stop)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Shutdown(ctx)

}
