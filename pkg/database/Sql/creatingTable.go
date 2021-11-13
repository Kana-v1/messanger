package sql

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func (c *MySqlContext) CreateTables(engine string, tableTemplate ...interface{}) error {
	for _, template := range tableTemplate {
		err := c.createTable(template, engine)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Can not create table %v", template))
		}
	}
	return nil
}
func (c *MySqlContext) createTable(tableTemplate interface{}, engine string) error {
	if engine == "" {
		engine = "InnoDB"
	}

	if !c.DB.Migrator().HasTable(tableTemplate) {
		c.DB.Set("gorm:table_options", fmt.Sprintf("ENGINE=%s", engine)).Migrator().CreateTable(tableTemplate)
	} else {
		return fmt.Errorf("table %s exist", reflect.TypeOf(tableTemplate))
	}
	return nil
}
