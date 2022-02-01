package models

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/pkg/scheduler"
)

type Invite struct {
	Code   string `json:"code"`
	ChatID int    `json:"chat_id,omitempty"`
	Aliens int    `json:"aliens"`
	Exp    int64  `json:"exp"`
}

type CreateMessage struct {
	ReplyTo *int
	UserID  *int
	RoomID  int
	Body    string
	Type    model.MessageType
	// CreatedAt int64 migrate to postgres
}

type MessageInfo struct {
	ID      int               `json:"id"`
	ReplyTo int               `json:"reply_to"`
	Author  int               `json:"author"`
	Body    string            `json:"body"`
	Type    model.MessageType `json:"type"`
	Time    int               `json:"time"`
}

/* // todo
msg_type:
	- system message
		- sender
		- event
		- body
	- user message
	- formatted message
msg_fields:
	- text
	- photo
	- file
	- vote
	- music
	- video
*/

type RoleReference struct {
	ID    *int
	Name  *string
	Color *model.HexColor
}

type RefreshSession struct {
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	ExpAt        int64
}

type UserInfo struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type AllowHolder struct {
	RoleID   *int
	Char     *model.CharType
	UserID   int
	MemberID int
}

type Unit struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Type   res.UnitType
}

type Chat struct {
	Unit    Unit
	Private bool
}

type DemoMember struct {
	UserID   int
	MemberID int
	IsOwner  bool
	Char     *model.CharType
	Muted    bool
}
type DefMember struct {
	UserID int
	ChatID int
}

type ScheduleInvite struct {
	ChatID int
	Code   string
	Exp    *int64
	Task   *scheduler.Task
}

type ScheduleRegisterSession struct {
	Email string
	Exp   int64
	Task  *scheduler.Task
}

type ScheduleRefreshSession struct {
	ID   int
	Exp  int64
	Task *scheduler.Task
}

type SubUser struct {
	MemberID *int
	ChatID   *int
}

type LoginRequisites struct {
	Email        string
	HashedPasswd string
}

type RegisterData struct {
	Domain       string
	Name         string
	Email        string
	HashPassword string
}
