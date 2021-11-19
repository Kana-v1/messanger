package nosql

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

type RedisContext struct {
	*redis.Client
	Mutex *sync.Mutex
}

func (c *RedisContext) AddValue(key string, values ...interface{}) {
	c.Mutex.Lock()
	for _, val := range values {
		res := c.SAdd(key, val)
		fmt.Println(res.Err())
	}
	c.Mutex.Unlock()
}

func (c *RedisContext) RemoveValue(key string, values ...interface{}) {
	c.Mutex.Lock()
	for _, val := range values {
		c.SRem(key, val)
	}
	c.Mutex.Unlock()
}

func (c *RedisContext) GetValue(key string) (string, error) {
	return c.Get(key).Result() //тут хуйня
}
