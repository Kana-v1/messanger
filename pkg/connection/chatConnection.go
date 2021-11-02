package connection

import (
	"fmt"
	"messanger/pkg/chat"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/sha3"
	"messanger/internal/logs"
)

const separateString = "SenderId"

var (
	ChatId   int64
	PeerId   int64
	UserId   int64
	Sessions []ChatSession
)

type ChatSession struct {
	Id    int64
	Users []User
	Peer  *Peer
}

func init() {
	Sessions = make([]ChatSession, 0)
}

func (cs *ChatSession) StartSubscriber() {
	go func() {
		channel := cs.GetChannel()
		sub := chat.Client.Subscribe(channel)
		messages := sub.Channel()
		for message := range messages {
			//message structure
			//sender id; hashed separateString message itself
			//↓			 ↓					 ↓↓				↓
			//1asdasdasdasdasdasdasdasdasdasdsSomeRealMessage
			senderIdAndMessage := strings.Split(message.Payload, fmt.Sprint(sha3.New224().Sum([]byte(separateString))))
			senderId, err := strconv.ParseInt(senderIdAndMessage[0], 10, 64)
			if err != nil {
				logs.ErrorLog("messagesErrors", fmt.Sprintf("Can not get sender id, file: %s", "chatConnection.go"), err)
			}
			msg := senderIdAndMessage[1]

			for _, user := range cs.Users {
				if senderId != user.Id {
					cs.Peer.WriteMessage(websocket.TextMessage, []byte(msg))
				}
			}
		}
	}()
}

func (cs *ChatSession) GetChannel() string {
	return fmt.Sprint(cs.Id) + "-channel"
}
