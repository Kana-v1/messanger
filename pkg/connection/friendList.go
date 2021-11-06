package connection

import "messanger/pkg/enums"

func (u *User) InBlackList(userId int64) bool {
	if collocutor, ok := u.UsersList[userId]; ok && collocutor == enums.InBlackList {
		return true
	}
	return false
}