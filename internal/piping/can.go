package piping

import (
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/pkg/kit"
)

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
	}
	return
}

func (c *Can) admin(userId, chatId int) (such bool) {
	if !c.pl.repos.Chats.UserIs(userId, chatId, rules.Admin) {
		c.pl.Err = resp.Error(resp.ErrBadRequest, "пользователь не является администратором чата")
	}
	return
}

func (c *Can) moder(userId, chatId int) (such bool) {
	if !c.pl.repos.Chats.UserIs(userId, chatId, rules.Moder) {
		c.pl.Err = resp.Error(resp.ErrBadRequest, "пользователь не является участником чата")
	}
	return
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
