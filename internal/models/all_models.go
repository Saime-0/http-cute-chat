package models

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

type Invite struct {
	Code   string `json:"code"`
	ChatID int    `json:"chat_id,omitempty"`
	Aliens int    `json:"aliens"`
	Exp    int64  `json:"exp"`
}

type CreateMessage struct {
	ReplyTo *int
	Author  *int
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
	Lifetime     int64
}

type UserInput struct { // pass & email
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type AllowsDB struct { // db table be like
	Action rules.AllowActionType
	Group  rules.AllowGroupType
	Value  string
}

type AllowHolder struct {
	RoleID   *int
	Char     rules.CharType
	UserID   int
	MemberID int
}

type Unit struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Type   rules.UnitType
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
