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

	"container/list"

	"github.com/gorilla/websocket"
)

const separateString = "SenderId"

var (
	ChatId           int64
	PeerId           int64
	UserId           int64
	Sessions         []ChatSession
	InactiveSessions *InactiveSession
)

type InactiveSession struct {
	List  *list.List
	Mutex *sync.Mutex
}

type Message struct {
	Message []byte
	Sender  int64
	Time    time.Time
}

type ChatSession struct {
	Id              int64
	Peers           map[*User]*Peer
	PrivateKey      *crypto.CryptoKeys
	MessageReceived chan string
	Messages        []Message
	State           enums.ChatSessionState
}

func init() {
	InactiveSessions = &InactiveSession{
		List:  list.New(),
		Mutex: new(sync.Mutex),
	}
	Sessions = make([]ChatSession, 0)
}

func (cs *ChatSession) StartSubscriber() {
	go func() {
		channel := cs.GetChannel()
		sub := chat.Client.Subscribe(channel)
		for user := range cs.Peers {
			chat.CreateUser(cs.GetChannel(), user.Name)
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
							}
							break
						}
					default:
						continue
					}
					continue
				}
			}()

			for receiver, peer := range cs.Peers {
				if senderId != receiver.Id && !receiver.InBlackList(senderId){
					msg, err := crypto.DecryptMessage([]byte(senderIdAndMessage[1]), cs.PrivateKey)
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
	for u, p := range cs.Peers {
		p.IsClosed = true
		msg, err := crypto.EncryptMessage([]byte(LastChatMessage), u.PublicKeys[cs.Id])
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
		cs.Peers = make(map[*User]*Peer)

		InactiveSessions.List.PushBack(*cs)
	}
}
