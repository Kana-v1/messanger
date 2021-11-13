package connection

import "gorm.io/gorm"

func (cs *ChatSession) AfterFind(tx *gorm.DB) (err error) {
	cs.GetChatSessionPeers(true)
	return
}

func (cs *ChatSession) BeforeCreate(tx *gorm.DB) (err error) {
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
