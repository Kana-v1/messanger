package connection

import (
	sql "messanger/pkg/database/Sql"
)

func LoadChatSession() ([]ChatSession, error) {
	sessions := make([]ChatSession, 0)
	err := sql.SqlContext.DB.Preload("Messages").Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func LoadUsers() ([]User, error) {
	users := make([]User, 0)
	err := sql.SqlContext.DB.Preload("Sessions")/*.Preload("Peers")*/.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func LoadInactiveSessions() ([]InactiveChatSession, error) {
	inactiveSessions := make([]InactiveChatSession, 0)
	err := sql.SqlContext.DB.Find(&inactiveSessions).Error
	if err != nil {
		return nil, err
	}
	return inactiveSessions, err
}
