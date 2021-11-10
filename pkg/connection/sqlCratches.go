package connection

import (
	"database/sql/driver"
	"encoding/json"
	"messanger/internal/logs"
	sql "messanger/pkg/repository/Sql"
)

func (ics *InactiveChatSessions) Scan(src interface{}) error {
	return json.Unmarshal([]byte(src.(string)), &ics)
}

func (ics *SessionId) Value() (driver.Value, error) {
	val, err := json.Marshal(ics.SessionId)
	return string(val), err
}

func SaveChatSessions() {
	sessions := make([]ChatSession, 0)
	for _, session := range Sessions {
		sessions = append(sessions, *session)
	}

	err := sql.SqlContext.AddValues("chat_session", sessions)
	if err != nil {
		logs.ErrorLog("sqlError.log", "Can not save chat sessions; err:", err)
	}
}

func GetChatSessions() map[int64]*ChatSession {
	sessions, err := sql.SqlContext.GetValues("chat_session", make([]ChatSession, 0))
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
	}
	return res
}

func SaveUsers() {
	users := make([]User, 0)
	for _, user := range Users {
		users = append(users, *user)
	}

	err := sql.SqlContext.AddValues("user", users)
	if err != nil {
		logs.ErrorLog("sqlError.log", "Can not save users; err:", err)
	}
}

func GetUsers() map[int64]*User {
	users, err := sql.SqlContext.GetValues("user", make([]User, 0))
	res := make(map[int64]*User)
	if err != nil {
		logs.ErrorLog("sqlError.log", "Can not get chat sessions; err:", err)
		return nil
	}
	for _, user := range users {
		u, ok := user.(User)
		if !ok {
			logs.ErrorLog("sqlError.log", "Can not extract chat session from interface{}", nil)
			return nil
		}
		res[u.Id] = &u
	}
	return res
}


