package nosql

import (
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
		c.SAdd(key, val)
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

func (c *RedisContext) GetValue(key string, id int64) (result interface{}) {
	return c.Get(key)
}
