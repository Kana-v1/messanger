package sql

import (
	"fmt"
	"sync"

	"gorm.io/gorm"

	"github.com/pkg/errors"
)

var SqlContext *MySqlContext

type MySqlContext struct {
	DB    *gorm.DB
	Mutex *sync.RWMutex
}

func (c *MySqlContext) AddValues(table string, values ...interface{}) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for _, value := range values {
		res := c.DB.Table(table).Save(value)
		if res.Error != nil {
			return errors.Wrap(res.Error, "Can not add value")
		}
	}
	return nil
}

func (c *MySqlContext) RemoveValues(table string, value ...interface{}) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	err := c.DB.Table(table).Delete(value).Error
	return errors.Wrap(err, fmt.Sprintf("Can not delete value %v from table %s;", value, table))
}

func (c *MySqlContext) GetValues(table string, values ...interface{}) (res []interface{}, err error) {
	res, _, err = c.Exist(table, values)
	return res, err
}

func (c *MySqlContext) Exist(table string, values ...interface{}) ([]interface{}, bool, error) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	res := make([]interface{}, 0)
	for _, value := range values {
		rec := c.DB.Table(table).Find(&value)
		if rec.Error != nil {
			return nil, false, errors.Wrap(rec.Error, fmt.Sprintf("Can not get value %v", value))
		}
		if rec == nil {
			return nil, false, nil
		}
		res = append(res, rec)
	}
	return res, true, nil
}

func (c *MySqlContext) AccountExist(hashedLog []byte, hashedPass []byte, acc interface{}) (bool, error) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	err := c.DB.Table("Accounts").Where(map[string]interface{}{"log": hashedLog, "password": hashedPass}).Find(acc).Error
	if err != nil {
		return false, errors.Wrap(err, "Can not check if account exists")
	}

	if acc != nil {
		return true, nil
	}

	return false, nil
}
