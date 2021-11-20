package connection

import (
	sql "messanger/pkg/database/Sql"

	"gorm.io/gorm"
)

func (cs *ChatSession) AfterFind(tx *gorm.DB) (err error) {
	cs.GetChatSessionPeers(true)
	return
}

func (cs *ChatSession) BeforeCreate(tx *gorm.DB) (err error) {
	sql.SqlContext.DB.Exec("DELETE FROM inactive_chat_sessions")
	cs.SaveChatSessionPeers()
	return
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.GetUserPeers()
	u.GetChatPublicKey()
	u.GetUserFriends()
	return
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.SaveChatPublicKey()
	u.SaveFriendFriends()
	return
}

func (cs SessionId) BeforeCreate(tx *gorm.DB) (err error) {
	sql.SqlContext.DB.Exec("DELETE FROM session_ids")
	return
}
