package connection

import (
	"crypto/rsa"
	"fmt"
	"messanger/internal/logs"
	"messanger/pkg/chat"
	crypto "messanger/pkg/cryptography/symmetricCrypto"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/sha3"
)

var usersUpdated chan bool
var id int64

type Peer struct {
	Id int64
	*websocket.Conn
}

type User struct {
	Id        int64
	Name      string
	Sessions  []int64
	Peers     []Peer
	PublicKeys map[int64]*rsa.PublicKey //for each chat session use own private key and session's public key
}

func NewUser() *User {
	id++
	return &User{
		Id:        id,
		Name:      fmt.Sprintf("%v_user", id),
		Sessions:  make([]int64, 0),
		Peers:     make([]Peer, 0),
		PublicKeys: make(map[int64]*rsa.PublicKey),
	}
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

				usersUpdated <- true

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
			Id:        sessionId,
			Peers:     userPeer,
			PrivateKey: cryptoKeys,
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
						user.disconnect()
					}
				} else {
					logs.ErrorLog("websocketErrors.log", "error while starting chat session, err:", err)
				}
				return
			}
			encrypedMessage, err := crypto.EncryptMessage(msg, u.PublicKeys[sessionId])
			if err != nil {
				logs.ErrorLog("chatErrors.log", fmt.Sprintf("Session id:%v, user id: %v", sessionId, u.Id), err)
				return
			}
			//symbol which separate sender id & message doesn't have to have collissions with symbols inside the message
			message := fmt.Sprintf(`%v%v%s`, u.Id, sha3.New224().Sum([]byte(separateString)), string(encrypedMessage))
			channel := fmt.Sprint(sessionId) + "-channel"

			chat.SendToChannel(message, channel)
		}
	}()
}
