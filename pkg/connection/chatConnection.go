package connection

import (
	"fmt"
	"messanger/pkg/chat"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	ChatId   int64
	PeerId   int64
	UserId   int64
	Sessions []ChatSession
)

type ChatSession struct {
	Id      int64
	Users   []User
	Peer    *Peer
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
			from := strings.Split(message.Payload, ":")[0]

			for _, user := range cs.Users {
				if from != user.Name {
					cs.Peer.WriteMessage(websocket.TextMessage, []byte(message.Payload))
				}
			}
		}
	}()
}

func (cs *ChatSession) GetChannel() string{
	return fmt.Sprint(cs.Id) + "-channel"
}
