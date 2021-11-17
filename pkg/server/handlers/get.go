package handlers

import (
	"encoding/json"
	"messanger/internal/logs"
	"messanger/pkg/connection"
	crypto "messanger/pkg/cryptography/symmetricCrypto"
	"messanger/pkg/enums"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	users := make([]connection.User, 0)
	for _, user := range connection.Users {
		users = append(users, *user)
	}
	jsonUsers, err := json.Marshal(users)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, string(jsonUsers))
}

type ChatForFront struct {
	Id         int64
	Users      []connection.User
	Messages   []connection.Message
	PrivateKey []byte `json:"-"`
}

func GetChats(c echo.Context) error {
	usersInChats := make(map[int64][]connection.User) // key - chat id; value - users in that chat
	for i, user := range connection.Users {
		for _, chats := range user.Sessions {
			if _, ok := usersInChats[chats.SessionId]; !ok {
				usersInChats[chats.SessionId] = []connection.User{*connection.Users[i]}
			}
			usersInChats[chats.SessionId] = append(usersInChats[chats.SessionId], *connection.Users[i])
		}
	}

	sessions := make([]ChatForFront, 0)
	for _, session := range connection.Sessions {
		if session.State != enums.ChatClosed {
			messages  := make([]connection.Message, 0)
			messages = append(messages, session.Messages...)
			
			sessions = append(sessions, ChatForFront{
				Id:         session.Id,
				Users:      usersInChats[session.Id],
				Messages:   messages,
				PrivateKey: session.PrivateKey,
			})
		}
	}

	for i := range sessions {
		privateKey, err := crypto.EncodePrivateKey(sessions[i].PrivateKey)
		if err != nil {
			logs.ErrorLog("cryptoKeys.log", "", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		for j := range sessions[i].Messages {
			msg, err := crypto.DecryptMessage(sessions[i].Messages[j].Message, privateKey)
			if err != nil {
				logs.ErrorLog("httpError.log", "Can not decrypt message", err)
				return c.String(http.StatusInternalServerError, err.Error())
			}
			sessions[i].Messages[j].Message = msg
		}
	}

	jsonChats, err := json.Marshal(sessions)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, string(jsonChats))
}
