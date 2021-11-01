package chat

import (
	"fmt"
	"log"
	nosql "messanger/pkg/repository/noSql"
	"os"
	"sync"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

var (
	Client        *redis.Client
	redisHost     string
	redisPassword string
	redisContext  *nosql.RedisContext
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	redisHost = os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatal("missing REDIS_HOST env var")
	}

	redisPassword = os.Getenv("REDIS_PASSWORD")
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
	//startSubscriber()
}

func RemoveUser(users string, user string) {
	redisContext.RemoveValue(users, user)
}

func CreateUser(users string, user string) {
	redisContext.AddValue(users, user)
}

func SendToChannel(msg string, channel string) {
	err := Client.Publish(channel, msg).Err()
	if err != nil {
		fmt.Printf("can not publish message '%v' to channel '%v', err: %v", msg, channel, err)
	}
}
