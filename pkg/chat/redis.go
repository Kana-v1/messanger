package chat

import (
	"fmt"
	"messanger/internal/logs"
	nosql "messanger/pkg/repository/noSql"
	"os"
	"sync"

	"github.com/go-redis/redis"
)

var (
	Client        *redis.Client
	redisHost     string
	redisPassword string
	redisContext  *nosql.RedisContext
)

func init() {

	var exit bool
	redisHost, exit = os.LookupEnv("REDIS_HOST")
	if !exit {
		logs.FatalLog("", "missing REDIS_HOST env var", nil)
	}

	redisPassword, exit = os.LookupEnv("REDIS_PASSWORD")
	if !exit {
		logs.FatalLog("", "missing REDIS_PASSWORD env var", nil)
	}
	if redisPassword == "" {
		logs.InfoLog("", "REDIS_PASSWORD is empty", nil)
	}
	Client = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	fmt.Println("Redis client started")
	redisContext = &nosql.RedisContext{
		Client: Client,
		Mutex:  new(sync.Mutex),
	}
}

func RemoveUser(session string, user string) {
	redisContext.RemoveValue(session, user)
}

func CreateUser(session string, user string) {
	redisContext.AddValue(session, user)
}

func SendToChannel(msg string, channel string) {
	err := Client.Publish(channel, msg).Err()
	if err != nil {
		fmt.Printf("can not publish message '%v' to channel '%v', err: %v", msg, channel, err)
	}
}
