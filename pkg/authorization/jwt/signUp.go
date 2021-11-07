package jwt

import (
	"messanger/internal/logs"
	"messanger/pkg/authorization"
	"messanger/pkg/cryptography/hash"
	sql "messanger/pkg/repository/Sql"
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

	if _, exist, err := sql.SqlContext.Exist("Accounts", newAcc); err != nil {
		logs.ErrorLog("sqlError.log", "", err)
		return c.String(http.StatusInternalServerError, "Account has not been registered")
	} else if exist {
		return c.String(http.StatusConflict, "Account with same log already exist")
	}

	err := sql.SqlContext.AddValue("Accounts", newAcc)
	if err != nil {
		logs.ErrorLog("sqlError.log", "", err)
		return c.String(http.StatusInternalServerError, "Account has not been registered")
	}

	return c.String(http.StatusOK, "Account has been registered")
}
