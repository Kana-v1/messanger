package main

import (
	"context"
	"fmt"
	_ "messanger/configs"
	"messanger/internal/logs"
	"messanger/pkg/connection"
	mySql "messanger/pkg/database/Sql"
	"messanger/pkg/server/handlers"
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

	e.GET("/*", handlers.WebSocketHandler)
	e.POST("/SignUp", handlers.SignUp)
	e.POST("/SignIn", handlers.SignIn)
	e.GET("/api/get/users", handlers.GetUsers)
	e.GET("api/get/chats", handlers.GetChats)
	e.GET("/checkAuthorize/:accId", handlers.IsAuthorized)

	go func() {
		stopServer(s)
	}()
	go func() {
		for {
			<-time.After(1 * time.Minute)
			SaveData()
		}
	}()
	logs.FatalLog("server.log", "Can not start server", s.ListenAndServe())
	dontStop := make(chan int)
	<-dontStop
}

func stopServer(s http.Server) {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	SaveData()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

	// mySql.SqlContext.CreateTables("", authorization.Account{
	// 	Id:       0,
	// 	Log:      make([]byte, 0),
	// 	Password: make([]byte, 0),
	// },
	// connection.InactiveChatSession{
	// 	ChatSessionId: -1,
	// },
	//  connection.Message{
	// 	Id:            -1,
	// 	ChatSessionId: -1,
	// 	Message:       make([]byte, 0),
	// 	Sender:        -1,
	// 	Time:          "",
	// },)
	// connection.ChatSessionPeer {
	// 	SessionId: -1,
	// 	Peer:  connection.Peer{Id: -1, IsClosed: false},
	// 	UserId: -1,
	// },
	// connection.ChatSession{
	// 	Id:         -1,
	// 	PrivateKey: make([]byte, 0),
	// 	Messages:   make([]connection.Message, 0),
	// 	State:      1,
	// },
	// connection.Peer{
	// 	Id:       -1,
	// 	IsClosed: true,
	// },
	// connection.UserPublicKey{
	// 	UserId:    -1,
	// 	ChatId:    -1,
	// 	PublicKey: make([]byte, 0),
	// },
	// connection.UserFriendList{
	// 	UserId:     -1,
	// 	FriendId:   -1,
	// 	FriendType: 0,
	// },
	// connection.User{
	// 	Id:         -1,
	// 	Name:       "SomeName",
	// 	Sessions:   make([]connection.SessionId, 0),
	// 	PublicKeys: make(map[int64][]byte),
	// 	UsersList:  make(map[int64]enums.UserType),
	// },
	// connection.SessionId{
	// 	UserId:    -1,
	// 	SessionId: -1,
	// })
	GetData()

}

func SaveData() {
	connection.SaveChatSessions()
	connection.SaveUsers()
	connection.SaveInactiveSessions()
}

func GetData() {
	connection.Sessions, connection.InactiveSessions = connection.GetChatSessions()
	connection.Users = connection.GetUsers()
	connection.InactiveSessions = connection.GetInactiveSession()
}
