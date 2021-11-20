package nosql

import (
	"fmt"
	"messanger/internal/logs"
	"sync"

	"github.com/pkg/errors"

	"github.com/go-redis/redis"
)

type RedisContext struct {
	*redis.Client
	Mutex *sync.RWMutex
}

func (c *RedisContext) AddValue(key string, values ...interface{}) {
	c.Mutex.Lock()
	for i := len(values) - 1; i >= 0; i-- { //set and get in the same order
		err := c.SAdd(key, values[i]).Err()
		if err != nil {
			logs.ErrorLog("redisError.log", fmt.Sprintf("Can not add value %v; err:", values[i]), err)
		}
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

func (c *RedisContext) Clear(key string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	members := c.SMembers(key)
	err := members.Err()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Can not delete data by key %s", key))
	}

	len := len(members.Args())
	err = c.SPopN(key, int64(len)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisContext) GetValueThreadUnsafe(key string) ([]string, error) {
	res := c.SMembers(key)
	return res.Val(), res.Err()
}
