package authorization

import (
	"bytes"
	sql "messanger/pkg/database/Sql"
)

type Account struct {
	Id       int64
	Log      []byte
	Password []byte
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

func AccountExist(hashedLog []byte, hashedPassword []byte) (int64, error) {
	accs := make([]Account, 0)
	if db, err := sql.SqlContext.GetAccounts(); err != nil {
		return -1, err
	} else if db == nil {
		return -1, nil
	} else {
		db.Find(&accs)
		for _, account := range accs {
			if bytes.Equal(account.Log, hashedLog) {
				return account.Id, nil
			}
		}
	}
	return -1, nil
}
