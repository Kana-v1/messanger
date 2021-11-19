package chat

import (
	"fmt"
	"messanger/internal/logs"
	nosql "messanger/pkg/database/noSql"
	"os"
	"sync"

	"github.com/go-redis/redis"
)

var (
	Client        *redis.Client
	redisHost     string
	redisPassword string
	RedisContext  *nosql.RedisContext
)

func init() {

	var exist bool
	redisHost, exist = os.LookupEnv("REDIS_HOST")
	if !exist {
		logs.FatalLog("", "missing REDIS_HOST env var", nil)
	}

	redisPassword, exist = os.LookupEnv("REDIS_PASSWORD")
	if !exist {
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
	RedisContext = &nosql.RedisContext{
		Client: Client,
		Mutex:  new(sync.Mutex),
	}
}

func RemoveUser(session string, user string) {
	RedisContext.RemoveValue(session, user)
}

func CreateUser(session string, user string) {
	RedisContext.AddValue(session, user)
}

func SendToChannel(msg string, channel string) {
	err := Client.Publish(channel, msg).Err()
	if err != nil {
		fmt.Printf("can not publish message '%v' to channel '%v', err: %v", msg, channel, err)
	}
}
