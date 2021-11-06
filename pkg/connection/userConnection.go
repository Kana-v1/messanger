package connection

import (
	"crypto/rsa"
	"fmt"
	"messanger/internal/logs"
	"messanger/pkg/chat"
	crypto "messanger/pkg/cryptography/symmetricCrypto"
	"messanger/pkg/enums"

	"github.com/gorilla/websocket"
)

var newUser chan string

//TODO use last user id from bd + 1
var usersId int64

const LastChatMessage = "You are the last user in chat, so chat and its history will be deleted"

type Peer struct {
	Id int64
	*websocket.Conn
	IsClosed bool
}

type User struct {
	Id         int64
	Name       string
	Sessions   []int64
	Peers      []Peer
	PublicKeys map[int64]*rsa.PublicKey //for each chat session use own public key and session's private key
	BlackList  []int64
}

func NewUser() *User {
	usersId++
	return &User{
		Id:         usersId,
		Name:       fmt.Sprintf("%v_user", usersId),
		Sessions:   make([]int64, 0),
		Peers:      make([]Peer, 0),
		PublicKeys: make(map[int64]*rsa.PublicKey),
		BlackList:  make([]int64, 0),
	}
}

func (u *User) disconnect() {
	for i, chatSession := range Sessions {
		for user := range chatSession.Peers {
			if u.Id == user.Id {
				chat.RemoveUser(chatSession.GetChannel(), u.Name)
				for u, p := range chatSession.Peers {
					if u.Id == user.Id {
						p.Close()
						delete(chatSession.Peers, u)
					}
				}
				if len(chatSession.Peers) == 1 {
					Sessions[i].deleteChat()
				}
				break
			}
		}
	}
}

func (u *User) Start(peer *Peer) {
	newUser = make(chan string, 2) //dont block method until somebody read from channel
	var chatSession *ChatSession
	var sessionId int64
	cryptoKeys := crypto.GenerateKeys()

	for i := range Sessions {
		for _, p := range Sessions[i].Peers {
			if p.Id == peer.Id {
				sessionId = Sessions[i].Id
				u.PublicKeys[sessionId] = &Sessions[i].PrivateKey.PublicKey

				chatSession = &Sessions[i]
				chatSession.Peers[u] = peer
				u.Sessions = append(u.Sessions, Sessions[i].Id)

				newUser <- u.Name

				break
			}
		}
	}

	if chatSession == nil {

		InactiveSessions.Mutex.Lock()
		if InactiveSessions.List.Len() > 0 {
			session := InactiveSessions.List.Front()
			InactiveSessions.List.Remove(session)

			cs, ok := session.Value.(ChatSession)
			if !ok {
				logs.ErrorLog("InactiveSessions.log", "Can not get inactive session from list", nil)
			} else {
				chatSession = &cs
				chatSession.State = enums.ChatActive
			}
			InactiveSessions.Mutex.Unlock()
		} else {
			InactiveSessions.Mutex.Unlock()

			sessionId = ChatId
			ChatId++

			chatSession = &ChatSession{
				Id:              sessionId,
				PrivateKey:      cryptoKeys,
				MessageReceived: make(chan string),
				Messages:        make([]Message, 0),
			}
		}
		userPeer := make(map[*User]*Peer)
		userPeer[u] = peer
		if chatSession.Peers == nil {
			chatSession.Peers = make(map[*User]*Peer)
		}
		chatSession.Peers = userPeer
		u.PublicKeys[sessionId] = &chatSession.PrivateKey.PublicKey
		chatSession.StartSubscriber()

		Sessions = append(Sessions, *chatSession)
	}

	go func() {
		for {
			_, msg, err := peer.ReadMessage()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					for user := range chatSession.Peers {
						if user.Id == u.Id {
							user.disconnect()
							break
						}
					}
				} else if !peer.IsClosed {
					logs.ErrorLog("websocketErrors.log", "error while starting chat session, err:", err)
				}
				return
			}
			encrypedMessage, err := crypto.EncryptMessage(msg, u.PublicKeys[sessionId])
			if err != nil {
				logs.ErrorLog("chatErrors.log", fmt.Sprintf("Session id:%v, user id: %v", sessionId, u.Id), err)
				return
			}
			message := fmt.Sprintf(`%v%v%s`, u.Id, separateString, string(encrypedMessage))
			channel := fmt.Sprint(sessionId) + "-channel"

			chat.SendToChannel(message, channel)
		}
	}()
}
