package main

import (
	"context"
	"fmt"
	_ "messanger/configs"
	"messanger/internal/logs"
	"messanger/pkg/server"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
)

var port string

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

	e := echo.New()
	s := http.Server{
		Addr:    ":" + port,
		Handler: e,
	}

	e.GET("/*", server.WebSocketHandler)
	logs.FatalLog("server.log", "Can not start server", s.ListenAndServe())

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logs.FatalLog("", "failed to start server", err)
		}
	}()
	dontStop := make(chan int)
	<-dontStop
	stopServer(s)
}

func stopServer(s http.Server) {
	stop := make(chan os.Signal)
	signal.Notify(stop)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Shutdown(ctx)

}
