package authorization

import sql "messanger/pkg/repository/Sql"

type Account struct {
	Id       int64
	Log      []byte `gorm:schema:"-"`
	Password []byte `gorm:schema:"-"`
}

type LogData struct {
	Log      string `json:"log" xml:"log"`
	Password string `json:"password" xml:"password"`
}

type Tabler interface {
	TableName() string
}

func (Account) TableName() string {
	return "Accounts"
}

func AccountExist(hashedLog []byte, hashedPassword []byte) (bool, error) {
	logData := new(LogData)
	return sql.SqlContext.AccountExist(hashedLog, hashedPassword, *logData)
}
