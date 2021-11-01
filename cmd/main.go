package main

import (
	"context"
	"fmt"
	"messanger/pkg/chat"
	"net/http"
	"os"
	"os/signal"
	"time"

	"messanger/internal/logs"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}
var port string
var id int

func main() {
	if err := godotenv.Load(); err != nil {
		logs.FatalLog("", "No .env file found", nil)
	}
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
	stopServer(server)
}

func webSocketHandler(rw http.ResponseWriter, req *http.Request) {
	body := make([]byte, 0)
	_, err := req.Body.Read(body)
	if err != nil {
		logs.ErrorLog("", fmt.Sprintf("Invalid request body: %v", string(body)), err)
	}

	user := new(pkg.User)
	//json.Unmarshal(body, user)
	id++
	user.Id = int64(id)

	peer, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		logs.FatalLog("", "websocket connection failed", err)
	}
	chatSession := chat.ChatSession{user}
	chatSession.Start(1)
}

func stopServer(s http.Server) {
	stop := make(chan os.Signal)
	signal.Notify(stop)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Shutdown(ctx)

}
