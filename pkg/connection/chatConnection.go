package connection

import (
	"fmt"
	"messanger/pkg/chat"
	"messanger/pkg/enums"
	"strconv"
	"strings"
	"sync"
	"time"

	"messanger/internal/logs"
	crypto "messanger/pkg/cryptography/symmetricCrypto"

	"github.com/gorilla/websocket"
)

const separateString = "SenderId"

var (
	ChatId           int64
	PeerId           int64
	UserId           int64
	Sessions         map[int64]*ChatSession //TODO get from db, write to db. read all sessions to slice and write/read to db
	InactiveSessions *InactiveChatSessions
	newUser          chan string

	Users map[int64]*User //TODO get from db, write to db. read all users to slice and write/read to db

	//TODO use last user id from bd + 1
	usersId int64
)

type InactiveChatSessions struct {
	ChatSessionsId []int64
	Mutex          *sync.Mutex `gorm:"-"`
}

type Message struct {
	Message []byte
	Sender  int64
	Time    time.Time
}

type ChatSession struct {
	Id              int64
	Peers           map[int64]*Peer //int64 - userId
	PrivateKey      []byte          //[]byte encode to  privateKey
	MessageReceived chan string
	Messages        []Message
	State           enums.ChatSessionState
}

func init() {
	InactiveSessions = &InactiveChatSessions{
		ChatSessionsId: make([]int64, 0),
		Mutex:          new(sync.Mutex),
	}
	Sessions = make(map[int64]*ChatSession)
	Users = make(map[int64]*User)
}

func (cs *ChatSession) StartSubscriber() {
	go func() {
		channel := cs.GetChannel()
		sub := chat.Client.Subscribe(channel)
		for userId := range cs.Peers {
			for _, user := range Users {
				if user.Id == userId {
					chat.CreateUser(cs.GetChannel(), user.Name)
				}
			}
		}
		messages := sub.Channel()
		messageMutex := new(sync.Mutex)

		for message := range messages {
			//message structure
			//sender id	separateString    	  encrypted message
			//↓			 	↓	   		   	  ↓
			//1          	SeparateString   [//xx]asx]a]sx]as[sad[]d[a]]
			senderIdAndMessage := strings.Split(message.Payload, separateString)
			senderId, err := strconv.ParseInt(senderIdAndMessage[0], 10, 64)
			if err != nil {
				logs.ErrorLog("messagesErrors", fmt.Sprintf("Can not get sender id, file: %s", "chatConnection.go; err:"), err)
			}

			go func() {
				for {
					select {
					case user := <-newUser:
						chat.CreateUser(cs.GetChannel(), user)

						for i := range Sessions {
							if Sessions[i].Id == cs.Id {
								cs.Peers = Sessions[i].Peers
								break
							}
						}
					default:
						continue
					}
					continue
				}
			}()

			for receiverId, peer := range cs.Peers {
				var receiver *User
				for i := range Users {
					if Users[i].Id == receiverId {
						receiver = Users[i]
						break
					}
				}
				if receiver == nil {
					logs.ErrorLog("getMessageError.log", fmt.Sprintf("Can not find message receiver with id %v", receiverId), nil)
					continue
				}

				if senderId != receiver.Id && !receiver.InBlackList(senderId) {
					privateKey, err := crypto.EncodePrivateKey(cs.PrivateKey)
					if err != nil {
						logs.ErrorLog("cryptoKeys.log", "", err)
						return
					}
					msg, err := crypto.DecryptMessage([]byte(senderIdAndMessage[1]), privateKey)
					if err != nil {
						logs.ErrorLog("chatError.log", fmt.Sprintf("Peer id: %v, Ssession id: %v, user id: %v; err:", peer.Id, cs.Id, receiver.Id), err)
					}
					peer.WriteMessage(websocket.BinaryMessage, msg)

					//maybe it can bee too many gorutines if there are 100 users in chat that writes at the same time, but if u have 100 users in chat u probably have 100 chats that run async
					go func() {
						messageMutex.Lock()
						cs.Messages = append(cs.Messages, Message{
							Message: []byte(senderIdAndMessage[1]),
							Sender:  senderId,
							Time:    time.Now(),
						})
						messageMutex.Unlock()
					}()

					if string(msg) == LastChatMessage {
						cs.MessageReceived <- string(msg)
					}
				}
			}
		}
	}()
}

func (cs *ChatSession) GetChannel() string {
	return fmt.Sprint(cs.Id) + "-channel"
}

func (cs *ChatSession) deleteChat() {
	for userId, p := range cs.Peers {
		var u *User
		for i := range Users {
			if Users[i].Id == userId {
				u = Users[i]
				break
			}
		}
		if u == nil {
			logs.ErrorLog("deleteChatError.log", fmt.Sprintf("Can not find message user with id %v", userId), nil)
			continue
		}
		p.IsClosed = true

		publicKey, err := crypto.EncodePublicKey(u.PublicKeys[cs.Id])
		if err != nil {
			logs.ErrorLog("cryptoKeys.log", "", err)
			return
		}
		msg, err := crypto.EncryptMessage([]byte(LastChatMessage), publicKey)
		if err != nil {
			logs.ErrorLog("", "can not encrypt message while deleting chat. Err:", err)
		}

		message := fmt.Sprintf(`%v%v%s`, -1, separateString, string(msg))
		channel := fmt.Sprint(cs.Id) + "-channel"
		chat.SendToChannel(message, channel)

		<-cs.MessageReceived
		p.Close()

		cs.State = enums.ChatClosed
		cs.Messages = make([]Message, 0)
		cs.Peers = make(map[int64]*Peer)

		InactiveSessions.ChatSessionsId = append(InactiveSessions.ChatSessionsId, cs.Id)
	}
}
