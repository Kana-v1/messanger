package jwt

import (
	"encoding/json"
	"fmt"
	"messanger/internal/logs"
	"messanger/pkg/authorization"
	"messanger/pkg/chat"
	"messanger/pkg/cryptography/hash"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type claims struct {
	Log string `json:"log"`
	jwt.StandardClaims
}

var jwtKey = []byte("secrey_key")

const TOKEN_LIFE_TIME = 5 * time.Minute

func SignIn(c echo.Context) error {
	logData := new(authorization.LogData)
	err := json.NewDecoder(c.Request().Body).Decode(logData)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Can not serialize request body; err: %v", err.Error()))
	}

	hashedLog := hash.Hash([]byte(logData.Log))
	hashedPassword := hash.Hash([]byte(logData.Password))
	accId, err := authorization.AccountExist(hashedLog, hashedPassword)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if accId == -1 {
		return c.String(http.StatusBadRequest, "Incorrect log or password")
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &claims{
		Log: string(hashedLog),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not create token")
	}

	redisKey := redisAccKey(strconv.FormatInt(accId, 10))
	setToken(tokenString, redisKey)

	go deleteToken(redisKey, TOKEN_LIFE_TIME) //not sure if this goroutine won't interfere gc to collect function's garbage

	return c.String(http.StatusAccepted, strconv.FormatInt(accId, 10))
}

func RefreshToken(c echo.Context) error {
	accId := c.Param("accId")
	claims, err := IsAuthorized(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if claims == nil {
		return c.String(http.StatusEarlyHints, "Too early to refresh token")
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not create token")
	}
	key := redisAccKey(accId)
	err = setToken(tokenString, key)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusAccepted, "TokenRefreshed")
}

func IsAuthorized(c echo.Context) (*claims, error) {
	unAuthorized := errors.New(strconv.Itoa(http.StatusUnauthorized))
	accId := c.Param("accId")
	if accId == "" {
		return nil, errors.New(strconv.Itoa(http.StatusInternalServerError))
	}

	redisKey := redisAccKey(accId)
	authData, err := chat.RedisContext.GetValueThreadUnsafe(redisKey)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Can not get data by %s key; err:", redisKey))
	}

	if len(authData) == 0 {
		return nil, nil
	}

	tokenStr := authData[0]

	claims := new(claims)

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, unAuthorized
		}
		return nil, err
	}

	if !tkn.Valid {
		return nil, unAuthorized
	}

	return claims, nil
}

func deleteToken(key string, expirationTime time.Duration) {
	time.Sleep(expirationTime)
	err := chat.RedisContext.Clear(key)
	if err != nil {
		logs.ErrorLog("redisError.log", "", err)
	}
}

func setToken(token string, key string) error {
	err := chat.RedisContext.Clear(key)
	if err != nil {
		logs.ErrorLog("redisError.log", "", err)
		return errors.New("redis error")
	}
	chat.RedisContext.AddValue(key, token)
	return nil
}

func redisAccKey(id string) string {
	return fmt.Sprintf("%v-token", accId)
}
