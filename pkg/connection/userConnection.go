package connection

import (
	"fmt"
	"messanger/internal/logs"
	"messanger/pkg/chat"
	crypto "messanger/pkg/cryptography/symmetricCrypto"
	"messanger/pkg/enums"

	"github.com/gorilla/websocket"
)

const LastChatMessage = "You are the last user in chat, so chat and its history will be deleted"

type Peer struct {
	Id              int64
	*websocket.Conn //TODO recreate websocket conn foe when start db server and read users/peers from db
	IsClosed        bool
}

type User struct {
	Id         int64
	Name       string
	Sessions   []int64
	Peers      []Peer
	PublicKeys map[int64][]byte //[]byte encodes to public key and vice versa, for each chat session use own public key and session's private key
	UsersList  map[int64]enums.UserType
}

func NewUser() *User {
	usersId++
	user := &User{
		Id:         usersId,
		Name:       fmt.Sprintf("%v_user", usersId),
		Sessions:   make([]int64, 0),
		Peers:      make([]Peer, 0),
		PublicKeys: make(map[int64][]byte),
		UsersList:  make(map[int64]enums.UserType),
	}
	Users[user.Id] = user
	return user
}

func (u *User) disconnect() {
	for i, chatSession := range Sessions {
		for userId, p := range chatSession.Peers {
			if u.Id == userId {
				chat.RemoveUser(chatSession.GetChannel(), u.Name)
				p.Close()
				delete(chatSession.Peers, u.Id)
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
	privateKey := crypto.GenerateKeys()

	for _, session := range Sessions {
		for _, p := range session.Peers {
			if p.Id == peer.Id {
				sessionId = session.Id
				u.PublicKeys[sessionId] = crypto.GetPublicKeyFromPrivateKey(session.PrivateKey)

				chatSession = session
				chatSession.Peers[u.Id] = peer
				u.Sessions = append(u.Sessions, session.Id)

				newUser <- u.Name

				break
			}
		}
	}

	if chatSession == nil {
		InactiveSessions.Mutex.Lock()
		if len(InactiveSessions.ChatSessionsId) > 0 {
			sessionId := InactiveSessions.ChatSessionsId[0]
			InactiveSessions.ChatSessionsId = InactiveSessions.ChatSessionsId[1:]
			chatSession = Sessions[sessionId] //надо проверить что стейт сессии в массиве меняется когда меняется стейт на 2 строки ниже
			chatSession.State = enums.ChatActive
			InactiveSessions.Mutex.Unlock()
		} else {
			InactiveSessions.Mutex.Unlock()

			sessionId = ChatId
			ChatId++
			chatSession = &ChatSession{
				Id:              sessionId,
				PrivateKey:      crypto.DecodePrivateKey(privateKey),
				MessageReceived: make(chan string),
				Messages:        make([]Message, 0),
			}
		}
		userPeer := make(map[int64]*Peer)
		userPeer[u.Id] = peer
		if chatSession.Peers == nil {
			chatSession.Peers = make(map[int64]*Peer)
		}
		chatSession.Peers = userPeer
		u.PublicKeys[sessionId] = crypto.GetPublicKeyFromPrivateKey(chatSession.PrivateKey)
		chatSession.StartSubscriber()

		Sessions[chatSession.Id] = chatSession
	}

	go func() {
		for {
			_, msg, err := peer.ReadMessage()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					for userId := range chatSession.Peers {
						if userId == u.Id {
							Users[userId].disconnect()
							break
						}
					}
				} else if !peer.IsClosed {
					logs.ErrorLog("websocketErrors.log", "error while starting chat session, err:", err)
				}
				return
			}

			publicKey, err := crypto.EncodePublicKey(u.PublicKeys[sessionId])
			if err != nil {
				logs.ErrorLog("cryptoKeys.log", "", err)
				return
			}
			encrypedMessage, err := crypto.EncryptMessage(msg, publicKey)
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
