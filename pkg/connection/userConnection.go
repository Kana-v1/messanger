package connection

import (
	"fmt"
	"messanger/internal/collections"
	"messanger/internal/logs"
	"messanger/pkg/chat"

	"github.com/gorilla/websocket"
)

type Peer struct {
	Id int64
	*websocket.Conn
}

type User struct {
	Id         int64
	Name       string
	PrivateKey string
	PublicKey  string
	Peers      []Peer
}

func (uc *User) disconnect() {
	for _, chatSession := range Sessions {
		if user := collections.Contains(chatSession.Users); user != nil {

			if len(chatSession.Users) == 2 { //TODO reuse such empty sessions instead of creating new one
				chatSession.Peer.Close()
				for _, user := range chatSession.Users {
					chat.RemoveUser(chatSession.GetChannel(), user.Name)
				}
			} else {
				chat.RemoveUser(chatSession.GetChannel(), uc.Name)

				for i, user := range chatSession.Users {
					if user.Id == uc.Id {
						usersInChat := make([]User, 0)
						usersInChat = append(chatSession.Users[:i], chatSession.Users[i+1:]...)
						chatSession.Users = usersInChat
						return
					}
				}
			}

		}
	}
}

func (u *User) Start(peer *Peer) {
	var chatSession *ChatSession
	for _, session := range Sessions {
		if session.Peer.Id == peer.Id {
			chatSession = &session
			chatSession.Users = append(chatSession.Users, *u)
			break
		}
	}

	if chatSession == nil {
		ChatId++
		chatSession = &ChatSession{
			Id:    ChatId,
			Users: []User{*u},
			Peer:  peer,
		}
		chatSession.StartSubscriber()

		Sessions = append(Sessions, *chatSession)
	}

	go func() {
		for {
			_, msg, err := peer.ReadMessage()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					for _, user := range chatSession.Users {
						user.disconnect()
					}
				}
				logs.ErrorLog("websocketErrors.log", "error while starting chat session, err:", err)
				return
			}
			sender := fmt.Sprintf(`{%v}`, u.Id)
			channel := string(chatSession.Id) + "-channel"

			chat.SendToChannel(fmt.Sprintf(sender, msg), channel)
		}
	}()
}
