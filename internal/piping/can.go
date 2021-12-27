package piping

import (
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

func (c *Can) owner(userId, chatId int) (such bool) {
	if !c.pl.repos.Chats.UserIsChatOwner(userId, chatId) {
		c.pl.Err = resp.Error(resp.ErrBadRequest, "пользователь не является владельцем чата")
		return
	}
	return true
}

func (c *Can) admin(userId, chatId int) (such bool) {
	if !c.pl.repos.Chats.UserIs(userId, chatId, rules.Admin) {
		c.pl.Err = resp.Error(resp.ErrBadRequest, "пользователь не является администратором чата")
		return
	}
	return true
}

func (c *Can) moder(userId, chatId int) (such bool) {
	if !c.pl.repos.Chats.UserIs(userId, chatId, rules.Moder) {
		c.pl.Err = resp.Error(resp.ErrBadRequest, "пользователь не является участником чата")
		return
	}
	return true
}

func (c *Can) CreateInvite(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		c.owner(uid, cid),
		c.admin(uid, cid),
	)
}

func (c *Can) Ban(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		c.owner(uid, cid),
		c.admin(uid, cid),
	)
}

func (c *Can) CreateRole(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		c.owner(uid, cid),
		c.admin(uid, cid),
	)
}

func (c *Can) CreateRoom(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		c.owner(uid, cid),
		c.admin(uid, cid),
	)
}

func (c *Can) GiveRole(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		c.owner(uid, cid),
		c.admin(uid, cid),
	)
}

func (c *Can) LeaveFromChat(uid, cid int) (fail bool) {
	return !c.owner(uid, cid)
}

func (c *Can) ObserveInvites(uid, cid int) (fail bool) {
	return all
}

func (c *Can) ObserveCountMembers(uid, cid int) (fail bool) {
	return all
}

func (c *Can) ObserveRoles(uid, cid int) (fail bool) {
	return all
}

func (c *Can) ObserveBanlist(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		c.owner(uid, cid),
		c.admin(uid, cid),
	)
}

func (c *Can) ObserveMembers(uid, cid int) (fail bool) {
	return all
}

func (c *Can) ObserveOwner(uid, cid int) (fail bool) {
	return all
}

func (c *Can) ObserveRooms(uid, cid int) (fail bool) {
	return all
}

func (c *Can) UpdateRoom(uid, cid int) (fail bool) {
	return !kit.LeastOne(
		c.owner(uid, cid),
		c.admin(uid, cid),
	)
}
