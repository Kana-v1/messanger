package connection

import (
	"crypto/rsa"
	"fmt"
	"messanger/internal/logs"
	"messanger/pkg/chat"
	crypto "messanger/pkg/cryptography/symmetricCrypto"

	"github.com/gorilla/websocket"
)

var newUser chan string
var id int64

const LastChatMessage = "You are the last user in chat, so chat and its history deleted"

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
	PublicKeys map[int64]*rsa.PublicKey //for each chat session use own private key and session's public key
}

func NewUser() *User {
	id++
	return &User{
		Id:         id,
		Name:       fmt.Sprintf("%v_user", id),
		Sessions:   make([]int64, 0),
		Peers:      make([]Peer, 0),
		PublicKeys: make(map[int64]*rsa.PublicKey),
	}
}

func (u *User) disconnect() {
	for _, chatSession := range Sessions {
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
					chatSession.deleteChat()
				}
				break
			}
		}
	}
}

func (cs *ChatSession) deleteChat() {
	for u, p := range cs.Peers {
		p.IsClosed = true
		msg, err := crypto.EncryptMessage([]byte(LastChatMessage), u.PublicKeys[cs.Id])
		if err != nil {
			logs.ErrorLog("", "can not encrypt message while deleting chat. Err:", err)
		}
		message := fmt.Sprintf(`%v%v%s`, -1, separateString, string(msg))
		channel := fmt.Sprint(cs.Id) + "-channel"
		chat.SendToChannel(message, channel)
		message = <- cs.MessageReceived
		if  message == LastChatMessage {
			p.Close()
		}
	}
}

func (u *User) Start(peer *Peer) {
	newUser = make(chan string, 2) //dont want to block method until somebody read from channel
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
		sessionId = ChatId
		ChatId++
		userPeer := make(map[*User]*Peer)
		userPeer[u] = peer

		chatSession = &ChatSession{
			Id:              sessionId,
			Peers:           userPeer,
			PrivateKey:      cryptoKeys,
			MessageReceived: make(chan string),
		}
		u.PublicKeys[sessionId] = &cryptoKeys.PublicKey
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
