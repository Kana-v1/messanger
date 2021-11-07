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

func (c *MySqlContext) AddValue(table string, value interface{}) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	res := c.DB.Table(table).Save(value)
	if res.Error != nil {
		return errors.Wrap(res.Error, "Can not add value")
	}
	return nil
}

func (c *MySqlContext) RemoveValue(table string, value interface{}) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	err := c.DB.Table(table).Delete(value).Error
	return errors.Wrap(err, fmt.Sprintf("Can not delete value %v from table %s;", value, table))
}

func (c *MySqlContext) GetValue(table string, value interface{}) (res interface{}, err error) {
	res, _, err = c.Exist(table, value)
	return res, err
}

func (c *MySqlContext) Exist(table string, value interface{}) (interface{}, bool, error) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	rec := c.DB.Table(table).Find(value)
	if rec.Error != nil {
		return nil, false, errors.Wrap(rec.Error, "Can not get value")
	}

	if rec != nil {
		return nil, false, nil
	}

	return rec, true, nil
}

func (c *MySqlContext) AccountExist(hashedLog []byte, hashedPass []byte, acc interface{}) (bool, error){
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
