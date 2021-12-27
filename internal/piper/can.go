package piper

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/pkg/kit"
)

const all = false

// как вариант сделать 3 уровня доступа
// 0 - только owner
// 1 - owner + admin
// 2 - owner + admin + moder

// Can
// very influential and monolithic(to Pipeline)!!
type Can struct {
	pl *Pipeline
}

func (n *Node) owner(userId, chatId int) (such bool) {
	if !n.repos.Chats.UserIsChatOwner(userId, chatId) {
		n.Err = resp.Error(resp.ErrBadRequest, "пользователь не является владельцем чата")
		return
	}
	return true
}

func (n *Node) admin(userId, chatId int) (such bool) {
	if !n.repos.Chats.UserIs(userId, chatId, rules.Admin) {
		n.Err = resp.Error(resp.ErrBadRequest, "пользователь не является администратором чата")
		return
	}
	return true
}

func (n *Node) moder(userId, chatId int) (such bool) {
	if !n.repos.Chats.UserIs(userId, chatId, rules.Moder) {
		n.Err = resp.Error(resp.ErrBadRequest, "пользователь не является модером чата")
		return
	}
	return true
}

func (n *Node) CanCreateInvite(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		n.owner(uid, cid),
		n.admin(uid, cid),
	)
}

func (n *Node) CanBan(userID, targetID, chatID int) (fail bool) {
	demoMembers := n.repos.Chats.DemoMembers(chatID, 0, userID, targetID)
	if demoMembers[0] == nil || demoMembers[1] == nil {
		n.Err = resp.Error(resp.ErrBadRequest, "не удалось найти мембрса")
		return true
	}
	if !(demoMembers[0].IsOwner ||
		*demoMembers[0].Char == model.CharTypeAdmin) {
		n.Err = resp.Error(resp.ErrBadRequest, "недостаточно прав")
		return true
	}
	if getCharLevel(demoMembers[0].Char) < getCharLevel(demoMembers[1].Char) {
		n.Err = resp.Error(resp.ErrBadRequest, "недостаточно прав")
		return true
	}
	return false
}

func (n *Node) CanCreateRole(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		n.owner(uid, cid),
		n.admin(uid, cid),
	)
}

func (n *Node) CanCreateRoom(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		n.owner(uid, cid),
		n.admin(uid, cid),
	)
}

func (n *Node) CanGiveRole(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		n.owner(uid, cid),
		n.admin(uid, cid),
	)
}

func (n *Node) CanLeaveFromChat(uid, cid int) (fail bool) {
	return !n.owner(uid, cid)
}

func (n *Node) CanObserveInvites(uid, cid int) (fail bool) {
	return all
}

func (n *Node) CanObserveCountMembers(uid, cid int) (fail bool) {
	return all
}

func (n *Node) CanObserveRoles(uid, cid int) (fail bool) {
	return all
}

func (n *Node) CanObserveBanlist(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		n.owner(uid, cid),
		n.admin(uid, cid),
	)
}

func (n *Node) CanObserveMembers(uid, cid int) (fail bool) {
	return all
}

func (n *Node) CanObserveOwner(uid, cid int) (fail bool) {
	return all
}

func (n *Node) CanObserveRooms(uid, cid int) (fail bool) {
	return all
}

func (n *Node) CanUpdateRoom(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		n.owner(uid, cid),
		n.admin(uid, cid),
	)
}

func (n *Node) CanUpdateChat(uid, cid int) (fail bool) {
	return !n.owner(uid, cid)
}

func (n *Node) CanTakeRole(MemberID, targetMemberID, chatID int) (fail bool) {
	demoMembers := n.repos.Chats.DemoMembers(0, 1, MemberID, targetMemberID)
	if !demoMembers[0].IsOwner {
		n.Err = resp.Error(resp.ErrBadRequest, "недостаточно прав")
		return true
	}
	if demoMembers[0] == nil ||
		!(*demoMembers[0].Char == model.CharTypeAdmin ||
			*demoMembers[0].Char == model.CharTypeModer) {
		n.Err = resp.Error(resp.ErrBadRequest, "недостаточно прав")
		return true
	}
	if getCharLevel(demoMembers[0].Char) < getCharLevel(demoMembers[1].Char) {
		n.Err = resp.Error(resp.ErrBadRequest, "недостаточно прав")
		return true
	}
	return false
}

func getCharLevel(char *model.CharType) int {
	level := 0
	for i, _char := range rules.CharLevels {
		if *char == _char {
			level = i + 1
		}
	}
	return level
}
