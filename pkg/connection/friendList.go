package connection

import "messanger/pkg/enums"

func (u *User) AddToBlackList(userId int64) {
	u.UsersList[userId] = enums.InBlackList
}
func (u *User) AddToFriendList(userId int64) {
	u.UsersList[userId] = enums.Friend
}

func (u *User) InBlackList(userId int64) bool {
	if collocutor, ok := u.UsersList[userId]; ok && collocutor == enums.InBlackList {
		return true
	}
	return false
}
