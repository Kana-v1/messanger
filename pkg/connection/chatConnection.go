package connection

import (
	"fmt"
	"messanger/pkg/chat"
	"strconv"
	"strings"

	"messanger/internal/logs"
	crypto "messanger/pkg/cryptography/symmetricCrypto"

	"github.com/gorilla/websocket"
)

const separateString = "SenderId"

var (
	ChatId   int64
	PeerId   int64
	UserId   int64
	Sessions []ChatSession
)

type ChatSession struct {
	Id              int64
	Peers           map[*User]*Peer
	PrivateKey      *crypto.CryptoKeys
	MessageReceived chan string
}

func init() {
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
			}

			for user, peer := range cs.Peers {
				if senderId != user.Id {
					msg, err := crypto.DecryptMessage([]byte(senderIdAndMessage[1]), cs.PrivateKey)
					if err != nil {
						logs.ErrorLog("chatError.log", fmt.Sprintf("Peer id: %v, Ssession id: %v, user id: %v; err:", peer.Id, cs.Id, user.Id), err)
					}
					peer.WriteMessage(websocket.BinaryMessage, msg)

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
