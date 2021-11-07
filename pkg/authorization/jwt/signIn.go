package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"messanger/pkg/authorization"
	"messanger/pkg/cryptography/hash"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type claims struct {
	Log string `json:"log"`
	jwt.StandardClaims
}

var jwtKey = []byte("secrey_key")

func SignIn(c echo.Context) error {
	logData := new(authorization.LogData)
	err := json.NewDecoder(c.Request().Body).Decode(logData)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Can not serialize request body; err: %v", err.Error()))
	}

	hashedLog := hash.Hash([]byte(logData.Log))
	hashedPassword := hash.Hash([]byte(logData.Password))
	exist, err := authorization.AccountExist(hashedLog, hashedPassword)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !exist {
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

	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	return c.String(http.StatusAccepted, "Successfully loged in")
}

func RefreshToken(c echo.Context) error {
	_, claims := IsAuthorized(c)

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return c.String(http.StatusEarlyHints, "Too early to refresh token")
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not create token")
	}

	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	return c.String(http.StatusAccepted, "TokenRefreshed")
}

func IsAuthorized(c echo.Context) (error, *claims) {
	unAuthorized := errors.New(strconv.Itoa(http.StatusUnauthorized))
	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return unAuthorized, nil
		}

		return err, nil
	}

	tokenStr := cookie.Value
	claims := new(claims)

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return unAuthorized, nil
		}
		return err, nil
	}

	if !tkn.Valid {
		return unAuthorized, nil
	}

	return nil, claims
}
