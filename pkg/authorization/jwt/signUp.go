package jwt

import (
	"bytes"
	"messanger/internal/logs"
	"messanger/pkg/authorization"
	"messanger/pkg/cryptography/hash"
	sql "messanger/pkg/database/Sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

var accId int64 = 1

func SignUp(c echo.Context) error {
	logData := new(authorization.LogData)
	if err := (&echo.DefaultBinder{}).BindBody(c, &logData); err != nil {
		return err
	}

	hashedPassword := hash.Hash([]byte(logData.Password))
	hashedLog := hash.Hash([]byte(logData.Log))
	//TODO use for users same id as for acc

	newAcc := &authorization.Account{
		Id:       accId,
		Log:      hashedLog,
		Password: hashedPassword,
	}
	accs := make([]authorization.Account, 0)

	if db, err := sql.SqlContext.GetAccounts(); err != nil {
		logs.ErrorLog("sqlError.log", "", err)
		return c.String(http.StatusInternalServerError, "Account has not been registered")
	} else if db == nil {
		return c.String(http.StatusConflict, "Account with same log already exist")
	} else {
		db.Find(&accs)
		for _, account := range accs {
			if bytes.Equal(account.Log, newAcc.Log) {
				return c.String(http.StatusConflict, "Account with same log already exist")
			}
		}
	}

	err := sql.SqlContext.AddValuesThreadSafe("Accounts", newAcc)
	if err != nil {
		logs.ErrorLog("sqlError.log", "", err)
		return c.String(http.StatusInternalServerError, "Account has not been registered")
	}

	return c.String(http.StatusOK, "Account has been registered")
}
