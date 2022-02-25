package piper

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
)

func charPtr(m model.CharType) *model.CharType {
	return &m
}

var charLevels = map[*model.CharType]int{
	nil:                          0,
	charPtr(model.CharTypeModer): 1,
	charPtr(model.CharTypeAdmin): 2,
}

type mlevel = model.CharType

const (
	owner mlevel = model.CharTypeAdmin
	admin mlevel = model.CharTypeAdmin
	moder mlevel = model.CharTypeModer
)

func (n *Node) diffLevelCheck(applyToOwner, applyToSelfChar bool, minCharLevel mlevel, master, slave *models.DemoMember) (bad bool) {
	if master == nil || slave == nil {
		n.SetError(resp.ErrBadRequest, "invalid memberID value")
		return true
	}
	if !applyToOwner && slave.IsOwner {
		n.SetError(resp.ErrBadRequest, "невозможно применить на этого участника")
		return true
	}
	if master.IsOwner {
		return false
	}
	if !applyToSelfChar && (charLevels[master.Char] == charLevels[slave.Char]) {
		n.SetError(resp.ErrBadRequest, "недостаточно прав")
		return true
	}
	if charLevels[master.Char] < charLevels[charPtr(minCharLevel)] && charLevels[master.Char] < charLevels[slave.Char] {
		n.SetError(resp.ErrBadRequest, "недостаточно прав")
		return true
	}
	return false
}

func (n *Node) levelCheck(minCharLevel mlevel, demo *models.DemoMember) (bad bool) {
	if demo == nil {
		n.SetError(resp.ErrBadRequest, "не является участником чата")
		return true
	}
	if !(demo.IsOwner || charLevels[demo.Char] >= charLevels[charPtr(minCharLevel)]) {
		n.SetError(resp.ErrBadRequest, "недостаточно прав")
		return true
	}
	return false
}
