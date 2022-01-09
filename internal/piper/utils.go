package piper

import (
	"fmt"
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
	fmt.Printf("%#v\n%#v\n", master, slave) // debug
	if master == nil || slave == nil {
		n.Err = resp.Error(resp.ErrBadRequest, "мемберса не существует")
		return true
	}
	if !applyToOwner && slave.IsOwner {
		n.Err = resp.Error(resp.ErrBadRequest, "невозможно применить на этого участника")
		return true
	}
	if master.IsOwner {
		return false
	}
	if !applyToSelfChar && (charLevels[master.Char] == charLevels[slave.Char]) {
		n.Err = resp.Error(resp.ErrBadRequest, "недостотачный уровень")
		return true
	}
	if charLevels[master.Char] < charLevels[charPtr(minCharLevel)] && charLevels[master.Char] < charLevels[slave.Char] {
		n.Err = resp.Error(resp.ErrBadRequest, "недостотачный уровень")
		return true
	}
	return false
}

func (n *Node) levelCheck(minCharLevel mlevel, demo *models.DemoMember) (bad bool) {
	if demo == nil {
		n.Err = resp.Error(resp.ErrBadRequest, "не является участником чата")
		return true
	}
	if !(demo.IsOwner || charLevels[demo.Char] >= charLevels[charPtr(minCharLevel)]) {
		n.Err = resp.Error(resp.ErrBadRequest, "недостаточно прав")
		return true
	}
	return false
}
