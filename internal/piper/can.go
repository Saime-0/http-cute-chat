package piper

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
)

const all = false

// Can
// very influential and monolithic(to Pipeline)!!
type Can struct {
	pl *Pipeline
}

func (n *Node) CanCreateInvite(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}

func (n *Node) CanBan(userID, targetID, chatID int) (fail bool) {
	demoMembers := n.repos.Chats.DemoMembers(chatID, 0, userID, targetID)
	return n.diffLevelCheck(
		false,
		false,
		model.CharTypeAdmin,
		demoMembers[0],
		demoMembers[1],
	)
}

func (n *Node) CanUnban(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}

func (n *Node) CanCreateRole(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}

func (n *Node) CanCreateRoom(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}

func (n *Node) CanCreateAllow(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}

func (n *Node) CanGiveRole(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}
func (n *Node) CanGiveChar(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}
func (n *Node) CanMuteMember(uid, cid int) (fail bool) {
	return n.levelCheck(
		moder,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}

func (n *Node) CanLeaveFromChat(uid, cid int) (fail bool) {
	if n.repos.Chats.UserIsChatOwner(uid, cid) {
		n.SetError(resp.ErrBadRequest, "невозможно выйти из чата")
		return
	}
	return true
}

func (n *Node) CanObserveInvites(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}

func (n *Node) CanObserveBanlist(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}

func (n *Node) CanUpdateRoom(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}
func (n *Node) CanUpdateRole(uid, cid int) (fail bool) {
	return n.levelCheck(
		admin,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}

func (n *Node) CanUpdateChat(uid, cid int) (fail bool) {
	return n.levelCheck(
		owner,
		n.repos.Chats.DemoMembers(cid, 0, uid)[0],
	)
}

func (n *Node) CanTakeRole(MemberID, targetMemberID int) (fail bool) {
	demoMembers := n.repos.Chats.DemoMembers(0, 1, MemberID, targetMemberID)
	return n.diffLevelCheck(
		false,
		true,
		model.CharTypeAdmin,
		demoMembers[0],
		demoMembers[1],
	)
}

func (n *Node) CanTakeChar(MemberID, targetMemberID int) (fail bool) {
	demoMembers := n.repos.Chats.DemoMembers(0, 1, MemberID, targetMemberID)
	return n.diffLevelCheck(
		true,
		false,
		model.CharTypeAdmin,
		demoMembers[0],
		demoMembers[1],
	)
}
