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
	Id                int64
	ChatSessionPeerId int64
	*websocket.Conn   `gorm:"-"` //TODO recreate websocket conn when start db server and read users/peers from db
	IsClosed          bool
}

type SessionId struct {
	UserId    int64
	SessionId int64
}

type User struct {
	Id         int64
	Name       string
	Sessions   []SessionId
	Peers      []Peer                   `gorm:"-"`
	PublicKeys map[int64][]byte         `gorm:"-"` //[]byte encodes to public key and vice versa, for each chat session use own public key and session's private key. Saving key to the db obviously is not the best practise, but its ok for now
	UsersList  map[int64]enums.UserType `gorm:"-"`
}

func NewUser() *User {
	UserId++
	user := &User{
		Id:         UserId,
		Name:       fmt.Sprintf("%v_user", UserId),
		Sessions:   make([]SessionId, 0),
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
				for i, userPeer := range u.Peers {
					if userPeer.Id == p.Id {
						peers := u.Peers[:i]
						peers = append(peers, u.Peers[i+1:]...)
						u.Peers = peers
					}
				}
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
				newUser <- u.Name
				break
			}
		}
	}

	if chatSession == nil {
		inactiveChatSessionMutex.Lock()
		if len(InactiveSessions) > 0 {
			sessionId = InactiveSessions[0].ChatSessionId
			InactiveSessions = InactiveSessions[1:]
			chatSession = Sessions[sessionId]
			chatSession.State = enums.ChatActive
			inactiveChatSessionMutex.Unlock()
		} else {
			inactiveChatSessionMutex.Unlock()

			sessionId = ChatId
			ChatId++
			chatSession = &ChatSession{
				Id:              sessionId,
				PrivateKey:      crypto.DecodePrivateKey(privateKey),
				MessageReceived: make(chan string),
				Messages:        make([]Message, 0),
				State:           enums.ChatActive,
			}
		}
		chatSession.StartSubscriber()
	}

	userPeer := make(map[int64]Peer)
	userPeer[u.Id] = *peer
	if chatSession.Peers == nil {
		chatSession.Peers = make(map[int64]Peer)
	}
	chatSession.Peers[u.Id] = *peer
	u.PublicKeys[sessionId] = crypto.GetPublicKeyFromPrivateKey(chatSession.PrivateKey)
	sId := &SessionId{
		SessionId: chatSession.Id,
	}
	u.Sessions = append(u.Sessions, *sId)

	Sessions[chatSession.Id] = chatSession

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
