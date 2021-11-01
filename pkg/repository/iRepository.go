package repository

type IRepository interface {
	AddValue(table string, values ...interface{})
	RemoveValue(table string, values ...interface{})
	GetValue(table string, id int64) (res interface{})
}