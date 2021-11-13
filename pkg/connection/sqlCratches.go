package connection

import (
	"fmt"
	"messanger/internal/logs"
	"messanger/pkg/enums"
	sql "messanger/pkg/repository/Sql"
)

var sessionPeers []ChatSessionPeer
var sessionPeersSavedCount = 0

var userPublicKeys []UserPublicKey
var userFriendList []UserFriendList

func SaveChatSessions() {
	sessions := make([]ChatSession, 0)
	for _, session := range Sessions {
		sessions = append(sessions, *session)
	}

	err := sql.SqlContext.AddValuesThreadSafe("chat_sessions", sessions)
	if err != nil {
		logs.ErrorLog("sqlError.log", "Can not save chat sessions; err:", err)
	}
}

func GetChatSessions() map[int64]*ChatSession {
	sessions, err := sql.SqlContext.GetValuesThreadSafe("chat_sessions", make([]ChatSession, 0))
	res := make(map[int64]*ChatSession)
	if err != nil {
		logs.ErrorLog("sqlError.log", "Can not get chat sessions; err:", err)
		return nil
	}
	for _, session := range sessions {
		s, ok := session.(ChatSession)
		if !ok {
			logs.ErrorLog("sqlError.log", "Can not extract chat session from interface{}", nil)
			return nil
		}
		res[s.Id] = &s
		if ChatId < s.Id {
			ChatId = s.Id
		}
	}
	ChatId++
	return res
}

func SaveUsers() {
	users := make([]User, 0)
	for _, user := range Users {
		users = append(users, *user)
	}

	err := sql.SqlContext.AddValuesThreadSafe("users", users)
	if err != nil {
		logs.ErrorLog("sqlError.log", "Can not save users; err:", err)
	}
}

func GetUsers() map[int64]*User {
	users, err := sql.SqlContext.GetValuesThreadSafe("users", make([]User, 0))
	res := make(map[int64]*User)
	if err != nil {
		logs.ErrorLog("sqlError.log", "Can not get chat users; err:", err)
		return nil
	}
	for _, user := range users {
		u, ok := user.(User)
		if !ok {
			logs.ErrorLog("sqlError.log", "Can not extract chat session from interface{}", nil)
			return nil
		}
		res[u.Id] = &u
		if UserId < u.Id {
			UserId = u.Id
		}
	}
	UserId++
	return res
}

func SaveInactiveSessions() {
	var sessionsToUpdate []ChatSession
	for _, session := range InactiveSessions {
		err := sql.SqlContext.AddValuesThreadSafe("inactive_chat_sessions", session)
		if err != nil {
			logs.ErrorLog("sqlError.log", "Can not save inactive chat session; err:", err)
			break
		}
		sessionsToUpdate = append(sessionsToUpdate, *Sessions[session.ChatSessionId])
	}
	err := sql.SqlContext.UpdateValues(sessionsToUpdate)
	logs.ErrorLog("sqlError.log", "Can not update session's columns; err: ", err)
}

func GetInactiveSession() []InactiveChatSession {
	sessions, err := sql.SqlContext.GetValuesThreadSafe("inactive_chat_sessions", make([]InactiveChatSession, 0))
	if err != nil {
		logs.ErrorLog("sqlError.log", "Can not get inactive chat sessions; err:", err)
		return nil
	}
	res := make([]InactiveChatSession, 0)
	for _, session := range sessions {
		s, ok := session.(InactiveChatSession)
		if !ok {
			logs.ErrorLog("sqlError.log", "Can not extract inactive chat session from interface{}", nil)
			return nil
		}
		res = append(res, s)
	}
	return res
}

type ChatSessionPeer struct {
	Id        int64
	SessionId int64
	Peer      Peer
	UserId    int64
}

func (cs *ChatSession) SaveChatSessionPeers() {
	sesPeers := make([]ChatSessionPeer, 0)
	for userId, peer := range cs.Peers {

		isAlreadyExist := false
		for _, sp := range sessionPeers {
			if sp.Peer.Id == peer.Id && sp.SessionId == cs.Id {
				isAlreadyExist = true
			}
		}

		if !isAlreadyExist {
			chatSessionPeer := &ChatSessionPeer{
				SessionId: cs.Id,
				Peer:      peer,
				UserId:    userId,
			}

			sesPeers = append(sesPeers, *chatSessionPeer)
			sessionPeers = append(sessionPeers, *chatSessionPeer)
		}
	}

	err := sql.SqlContext.AddValues("chat_session_peers", sesPeers)
	if err != nil {
		logs.ErrorLog("sqlError.log", fmt.Sprintf("Can not save chatSession(Id: %v) peers; err:", cs.Id), err)
	}

	sessionPeersSavedCount++
	if sessionPeersSavedCount == len(Sessions) {
		sessionPeers = make([]ChatSessionPeer, 0)
	}
}

func (cs *ChatSession) GetChatSessionPeers(safe bool) {
	var sessionsPeers []interface{}
	var err error

	if sessionPeers == nil {
		if safe {
			sessionsPeers, err = sql.SqlContext.GetValuesThreadSafe("chat_session_peers", make([]ChatSessionPeer, 0))
		} else {
			sessionsPeers, err = sql.SqlContext.GetValues("chat_session_peers", make([]ChatSessionPeer, 0))
		}

		if err != nil {
			logs.ErrorLog("sqlError.log", "Can not get chatSessionPeers; err:", err)
			return
		}
		for _, sessionPeer := range sessionsPeers {
			sp, ok := sessionPeer.(ChatSessionPeer)
			if !ok {
				logs.ErrorLog("sqlError.log", "Can not extract chatSessionPeer from interface{}", nil)
				return
			}
			sessionPeers = append(sessionPeers, sp)
		}
	}
	for _, sessionPeer := range sessionPeers {
		if sessionPeer.SessionId == cs.Id {
			cs.Peers[sessionPeer.UserId] = sessionPeer.Peer
		}

	}

}

func (u *User) GetUserPeers() {
	var userPeers []interface{}
	var err error

	if sessionPeers == nil {
		userPeers, err = sql.SqlContext.GetValuesThreadSafe("chat_session_peers", make([]ChatSessionPeer, 0))
		if err != nil {
			logs.ErrorLog("sqlError.log", "Can not get chatSessionPeers; err:", err)
			return
		}
	}
	for _, userPeer := range userPeers {
		up, ok := userPeer.(ChatSessionPeer)
		if !ok {
			logs.ErrorLog("sqlError.log", "Can not extract chatSessionPeer from interface{}", nil)
			return
		}
		if up.UserId == u.Id {
			u.Peers = append(u.Peers, up.Peer)
		}
	}
}

//each session save peers and each peer contains user id, so user dont have so save its own peers

type UserPublicKey struct {
	UserId    int64
	ChatId    int64
	PublicKey []byte
}

func (u *User) GetChatPublicKey() {
	if userPublicKeys == nil {
		upk, err := sql.SqlContext.GetValuesThreadSafe("user_public_keys", make([]UserPublicKey, 0))
		if err != nil {
			logs.ErrorLog("sqlError.log", "Can not get userPublicKey; err:", err)
			return
		}

		for _, userPubKey := range upk {
			userPublicKey, ok := userPubKey.(UserPublicKey)
			if !ok {
				logs.ErrorLog("sqlError.log", "Can not extract userPublicKey from interface{}", nil)
				return
			}
			userPublicKeys = append(userPublicKeys, userPublicKey)
		}
	}

	for _, userPublicKey := range userPublicKeys {
		if userPublicKey.UserId == u.Id {
			u.PublicKeys[userPublicKey.ChatId] = userPublicKey.PublicKey
		}
	}
}

func (u *User) SaveChatPublicKey() {
	userPublicKeys := make([]UserPublicKey, 0)
	for chatId, publicKey := range u.PublicKeys {
		userPublicKey := &UserPublicKey{
			UserId:    u.Id,
			ChatId:    chatId,
			PublicKey: publicKey,
		}
		userPublicKeys = append(userPublicKeys, *userPublicKey)
	}

	err := sql.SqlContext.AddValues("user_public_keys", userPublicKeys)
	if err != nil {
		logs.ErrorLog("sqlError.log", fmt.Sprintf("Can not save user's (Id: %v) public keys; err:", u.Id), err)
	}
}

type UserFriendList struct {
	UserId     int64
	FriendId   int64
	FriendType enums.UserType
}

func (u *User) GetUserFriends() {
	if userFriendList == nil {
		ufl, err := sql.SqlContext.GetValuesThreadSafe("user_friend_lists", make([]UserFriendList, 0))
		if err != nil {
			logs.ErrorLog("sqlError.log", "Can not get userFriendList; err:", err)
			return
		}

		for _, userFriend := range ufl {
			uf, ok := userFriend.(UserFriendList)
			if !ok {
				logs.ErrorLog("sqlError.log", "Can not extract userFriendList from interface{}", nil)
				return
			}
			userFriendList = append(userFriendList, uf)
		}
	}

	for _, userFriend := range userFriendList {
		if userFriend.UserId == u.Id {
			u.UsersList[userFriend.FriendId] = userFriend.FriendType
		}
	}
}
func (u *User) SaveFriendFriends() {
	friendList := make([]UserFriendList, 0)
	for friendId, friendType := range u.UsersList {
		userFriend := &UserFriendList{
			UserId:     u.Id,
			FriendId:   friendId,
			FriendType: friendType,
		}
		friendList = append(friendList, *userFriend)
	}
	err := sql.SqlContext.AddValues("user_friend_lists", friendList)
	if err != nil {
		logs.ErrorLog("sqlError.log", fmt.Sprintf("Can not save user's(Id: %v) friends", u.Id), err)
	}
}
