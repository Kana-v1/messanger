package connection

import (
	"fmt"
	"messanger/internal/logs"
	"messanger/pkg/chat"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/sha3"
)

var usersUpdated chan bool

type Peer struct {
	Id int64
	*websocket.Conn
}

type User struct {
	Id         int64
	Name       string
	PrivateKey string
	PublicKey  string
	Sessions   []int64
	Peers      []Peer
}

func (u *User) disconnect() {
	for _, chatSession := range Sessions {
		for user := range chatSession.Peers {
			if u.Id == user.Id {
				if len(chatSession.Peers) == 2 { //TODO reuse such empty sessions instead of creating new one
					for u, p := range chatSession.Peers {
						p.Close()
						chat.RemoveUser(chatSession.GetChannel(), u.Name)
						
					}
				} else {
					chat.RemoveUser(chatSession.GetChannel(), u.Name)
					delete(chatSession.Peers, u)
				}
				break
			}
		}
	}
}

func (u *User) Start(peer *Peer) {
	usersUpdated = make(chan bool, 2) //dont want to block method until somebody read from channel
	var chatSession *ChatSession
	for i := range Sessions {
		for _, p := range Sessions[i].Peers {
			if p.Id == peer.Id {
				chatSession = &Sessions[i]
				chatSession.Peers[u] = peer
				u.Sessions = append(u.Sessions, Sessions[i].Id)
				usersUpdated <- true
				break
			}
		}
	}

	if chatSession == nil {
		ChatId++
		userPeer := make(map[*User]*Peer)
		userPeer[u] = peer
		chatSession = &ChatSession{
			Id:    ChatId,
			Peers: userPeer,
		}
		chatSession.StartSubscriber()

		Sessions = append(Sessions, *chatSession)
	}

	go func() {
		for {
			_, msg, err := peer.ReadMessage()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					for user := range chatSession.Peers {
						user.disconnect()
					}
				} else {
					logs.ErrorLog("websocketErrors.log", "error while starting chat session, err:", err)
				}
				return
			}
			//symbol which separate sender id & message dont have to has collissions with symbols inside the message
			message := fmt.Sprintf(`%v%v%s`, u.Id, sha3.New224().Sum([]byte(separateString)), string(msg))
			channel := fmt.Sprint(chatSession.Id) + "-channel"

			chat.SendToChannel(message, channel)
		}
	}()
}
