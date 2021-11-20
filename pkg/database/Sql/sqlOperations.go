package sql

import (
	"fmt"
	"reflect"
	"sync"

	"gorm.io/gorm"

	"github.com/pkg/errors"
)

var SqlContext *MySqlContext

type MySqlContext struct {
	DB    *gorm.DB
	Mutex *sync.RWMutex
}

func (c *MySqlContext) AddValuesThreadSafe(table string, values ...interface{}) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for _, value := range values {
		err := c.DB.Table(table).Create(value).Error
		if err != nil {
			err := c.DB.Table(table).Save(value).Error
			if err != nil {
				return errors.Wrap(err, "Can not add or update value")
			}
		}
	}
	return nil
}

func (c *MySqlContext) AddValues(table string, values ...interface{}) error {
	for _, value := range values {
		err := c.DB.Table(table).Create(value).Error
		if err != nil {
			return errors.Wrap(err, "Can not add value")
		}
	}
	return nil
}

func (c *MySqlContext) RemoveValuesThreadSafe(table string, value ...interface{}) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	err := c.DB.Table(table).Delete(value).Error
	return errors.Wrap(err, fmt.Sprintf("Can not delete value %v from table %s;", value, table))
}

func (c *MySqlContext) RemoveValues(table string, value ...interface{}) error {
	err := c.DB.Table(table).Delete(value).Error
	return errors.Wrap(err, fmt.Sprintf("Can not delete value %v from table %s;", value, table))
}

func (c *MySqlContext) GetValuesThreadSafe(table string, value interface{}) (err error) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	val := reflect.ValueOf(value)
	err = c.DB.Table(table).Unscoped().Find(val.Pointer()).Error
	return err
}

func (c *MySqlContext) GetValues(table string, value interface{}) (err error) {
	return c.DB.Table(table).Unscoped().Find(value).Error
}

func (c *MySqlContext) ExistThreadSafe(table string, value interface{}) (bool, error) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	return c.Exist(table, &value)
}

func (c *MySqlContext) Exist(table string, value interface{}) (bool, error) {
	err := c.DB.Table(table).Find(&value).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, errors.Wrap(err, fmt.Sprintf("Can not get value %v", value))
	}
	return true, nil
}

func (c *MySqlContext) GetAccounts() (*gorm.DB, error) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	db := c.DB.Table("Accounts")
	if db.Error != nil {
		return nil, errors.Wrap(db.Error, "Can not check if account exists")
	}
	if db.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return db, nil
}

func (c *MySqlContext) UpdateValues(values ...interface{}) error {
	for _, value := range values {
		err := c.DB.Updates(value).Error
		if err != nil {
			return errors.Wrap(err, "Can not update value")
		}
	}
	return nil
}
