package sql

import (
	"database/sql"
	"fmt"
	"messanger/internal/logs"
	"reflect"
	"sync"
)

type MySqlContext struct {
	*sql.DB
	mutex *sync.Mutex
}

func (c *MySqlContext) AddValue(table string, values ...interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	tran, err := c.Begin()
	if err != nil {
		logs.ErrorLog("", "Error while creating transaction", err)
		return
	}
	for _, value := range values {
		_, err := tran.Exec("INSERT ? INTO ?", value, table)
		if err != nil {
			logs.ErrorLog("", fmt.Sprintf("Error while adding new value(%v) to the table %v", value, table), err)
			err = tran.Rollback()
			logs.ErrorLog("", "Can not rollback tran", err)
			return
		}
	}
	tran.Commit()
}

//every value have to has key
func (c *MySqlContext) RemoveValue(table string, values ...interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	tran, err := c.Begin()
	if err != nil {
		logs.ErrorLog("", "Error while creating transaction", err)
		return
	}

	for _, value := range values {
		id, err := getValueId(value)
		if err != nil {
			logs.ErrorLog("", "", err)
			return
		}
		if _, ok := id.(int64); !ok {
			if _, ok = id.(int); !ok {
				logs.ErrorLog("", fmt.Sprintf("Incorrect id type in the table %s, row - %v", table, value), nil)
				continue
			}
		}

		tran.Exec("DELETE FROM ? WHERE Id = ?", table, id)
		if err != nil {
			logs.ErrorLog("", fmt.Sprintf("Error while deleting value(%v) from the table %v", value, table), err)
			err = tran.Rollback()
			logs.ErrorLog("", "Can not rollback tran", err)
			return
		}
	}
	tran.Commit()
}

func (c *MySqlContext) GetValue(table string, id int64) (res interface{}) {
	rows, err := c.Query("SELECT * FROM ? WHERE Id = ?", table, id)
	if err != nil {
		logs.ErrorLog("", fmt.Sprintf("Can not get rows from table %s", table), err)
		return nil
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		logs.ErrorLog("", "", err)
		return nil
	}

	res = []interface{}{
		new(int64),//id
		new(string),//name
		new(string),//open key
		new(string),//login
		new(string),//password
	}

	err = rows.Scan(cols)
	if err != nil {
		logs.ErrorLog("", "", err)
		return nil
	}

	return res
}

func getValueId(value interface{}) (interface{}, error) {
	val := reflect.ValueOf(value).Elem()

	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Name == "Id" {
			return val.Field(i).Interface(), nil
		}
	}

	return nil, fmt.Errorf("can not get id from value %v", value)
}
