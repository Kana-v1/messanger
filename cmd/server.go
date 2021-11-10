package main

import (
	"context"
	"fmt"
	_ "messanger/configs"
	"messanger/internal/logs"
	"messanger/pkg/connection"
	mySql "messanger/pkg/repository/Sql"
	"messanger/pkg/server"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

var port string

func main() {

	port = os.Getenv("PORT")
	if port == "" {
		logs.FatalLog("", "missing PORT env var", nil)
	}
	startMySqlServer()

	e := echo.New()
	s := http.Server{
		Addr:    ":" + port,
		Handler: e,
	}

	e.GET("/*", server.WebSocketHandler)
	e.POST("/SignUp", server.SignUp)
	e.POST("/SignIn", server.SignIn)
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

func startMySqlServer() {
	port, exist := os.LookupEnv("MYSQL_PORT")
	if !exist {
		logs.FatalLog("mysql.log", "Port does not defined in .env file", nil)
	}
	user, exist := os.LookupEnv("MYSQL_USER")
	if !exist {
		logs.FatalLog("mysql.log", "User does not defined in .env file", nil)
	}
	ip, exist := os.LookupEnv("MYSQL_IP")
	if !exist {
		logs.FatalLog("mysql.log", "Ip does not defined in .env file", nil)
	}
	password, exist := os.LookupEnv("MYSQL_PASSWORD")
	if !exist {
		logs.FatalLog("mysql.log", "Password does not defined in .env file", nil)
	}
	title, exist := os.LookupEnv("MYSQL_DB_NAME")
	if !exist {
		logs.FatalLog("mysql.log", "Password does not defined in .env file", nil)
	}
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, password, ip, port, title)), &gorm.Config{})
	if err != nil {
		logs.FatalLog("mysql.log", "Can not create mySql server; err: ", err)
	}

	mySql.SqlContext = &mySql.MySqlContext{
		Mutex: new(sync.RWMutex),
		DB:    db,
	}

	mySql.SqlContext.CreateTables("",
		// authorization.Account{
		// 	Id:       0,
		// 	Log:      make([]byte, 0),
		// 	Password: make([]byte, 0),
		// },
		connection.InactiveChatSessions{
			ChatSessionsId: make([]int64, 0),
		},
		// connection.Message{
		// 	Message: make([]byte, 0),
		// 	Sender:  -1,
		// 	Time:    time.Now(),
		// },
		// connection.ChatSession{
		// 	Id:         -1,
		// 	Peers:      make(map[int64]connection.Peer),
		// 	PrivateKey: make([]byte, 0),
		// 	Messages:   make([]connection.Message, 0),
		// 	State:      1,
		// },
		// connection.Peer{
		// 	Id:       -1,
		// 	IsClosed: true,
		// },
		// connection.User{
		// 	Id:       -1,
		// 	Name:     "SomeName",
		// 	Sessions: make([]connection.SessionId, 0),
		// },
		// connection.SessionId{
		// 	UserId:    -1,
		// 	SessionId: -1,
		// },
	)

}
