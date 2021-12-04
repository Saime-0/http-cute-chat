package models

import "github.com/saime-0/http-cute-chat/internal/api/rules"

// INTERNAL models
// RECEIVED models
type CreateChat struct {
	Domain  string `json:"domain"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

type UpdateChatData struct {
	Domain  string `json:"domain"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

type ChatMembersCount struct {
	Count int `json:"count"`
}

type Members struct {
	Members []Member `json:"members"`
}
type ChatID struct {
	ID int `json:"id"`
}
type Invite struct {
	Code   string `json:"code"`
	ChatID int    `json:"chat_id,omitempty"`
	Aliens int    `json:"aliens"`
	Exp    int64  `json:"exp"`
}
type InviteInput struct {
	Aliens   int   `json:"aliens"`
	LifeTime int64 `json:"lifetime"`
}
type CreateInvite struct {
	ChatID int   `json:"chat_id"`
	Aliens int   `json:"aliens"`
	Exp    int64 `json:"exp"`
}
type Invites struct {
	Invites []Invite
}
type MemberInfo struct {
	RoleID   int `json:"role_id"`
	JoinedAt int `json:"joined_at"`
}

// news
type Member struct {
	User     User
	RoleID   *int // deprecated
	Char     rules.CharType
	JoinedAt int64
	Muted    bool
	Frozen   bool
}

// general room model
type Dialog struct {
	ID          int      `json:"id"`
	MessagesIDs []int    `json:"messages_ids"`
	Companion   UserInfo `json:"companion"`
}

// INTERNAL models
// RECEIVED models

// RESPONSE models
type DialogInfo struct {
	Companion UserInfo `json:"companion"`
}
type ListDialogs struct {
	Dialogs []DialogInfo `json:"dialogs"`
}

// RECEIVED models
type CreateMessage struct {
	ReplyTo int    `json:"reply_to"`
	Author  int    `json:"author"`
	Body    string `json:"body"`
	// UnitType    rules.MessageType `json:"type"`
}

// RESPONSE models
type MessageInfo struct {
	ID      int               `json:"id"`
	ReplyTo int               `json:"reply_to"`
	Author  int               `json:"author"`
	Body    string            `json:"body"`
	Type    rules.MessageType `json:"type"`
	Time    int               `json:"time"`
}
type MessagesList struct {
	Messages []MessageInfo `json:"messages"`
}
type MessageID struct {
	ID int `json:"id"`
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

type Response struct {
	ID      int    `json:"id"`
	ReplyTo int    `json:"reply_to"`
	Author  int    `json:"author"`
	Body    string `json:"body"`
	Type    int    `json:"type"`
}

type CreateRole struct {
	Name  string `json:"role_name"`
	Color string `json:"color"`
}
type UpdateRole struct {
	RoleName      string `json:"role_name"`
	Color         string `json:"color"`
	Visible       bool   `json:"visible"`
	ManageRooms   bool   `json:"manage_rooms"`
	RoomID        int    `json:"room_id"`
	ManageChat    bool   `json:"manage_chat"`
	ManageRoles   bool   `json:"manage_roles"`
	ManageMembers bool   `json:"manage_members"`
}
type Role struct {
	ID    int
	Name  string
	Color string
}

type RoleID struct {
	ID int `json:"id"`
}

type CreateRoom struct {
	Name      string      `json:"name"`
	ParentID  int         `json:"parent_id"`
	Note      string      `json:"note"`
	MsgFormat FormPattern `json:"msg_format"`
	Restricts Allows      `json:"restricts"`
}
type UpdateRoomData struct {
	Name    string `json:"name"`
	Note    string `json:"note"`
	Private bool   `json:"private"`
}

// RESPONSE models
type RoomInfo struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
	Note     string `json:"note"`
	Private  bool   `json:"private"`
}
type ListRoomInfo struct {
	Rooms []RoomInfo `json:"rooms"`
}
type RoomID struct {
	ID int `json:"id"`
}

type RefreshSession struct {
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	Exp          int64  `json:"exp"`
	CreatedAt    int64  `json:"created_at"`
}

type CustomClaims struct {
	UserID    int   `json:"user-id"`
	ExpiresAt int64 `json:"exp"`
}

type CreateUser struct {
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserInput struct { // pass & email
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateUserData struct {
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateUserSettings struct {
	AppSettings string `json:"app_settings"`
}
type TokenForRefreshPair struct {
	RefreshToken string `json:"refresh_token"`
}

// RESPONSE models
type UserInfo struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}
type MeData struct {
	Email    string `json:"email"`
	Password string
}
type UserSettings struct {
	AppSettings string `json:"app_settings"`
}
type FreshTokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type ListUserInfo struct {
	Users []UserInfo `json:"users"`
}
type AllowsDB struct { // db table be like
	Action rules.AllowActionType
	Group  rules.AllowGroupType
	Value  string
}
type Allows struct { // gql scheme be like
	Read  AllowHolders
	Write AllowHolders
}
type AllowHolders struct {
	Roles []int
	Chars []rules.CharType
	Users []int
}

type Unit struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Type   rules.UnitType
}

type Me struct {
	Unit Unit
	User MeData
}
type Chat struct {
	Unit    Unit
	Private bool
}
type User struct {
	Unit Unit
}
type Users struct {
	Users []User
}
type InviteInfo struct {
	Unit         *Unit `json:"unit"`
	Private      bool  `json:"private"`
	CountMembers int   `json:"count_members"`
}
type AllowV2 struct {
	Action rules.AllowActionType
	Group  rules.AllowGroupType
	Value  string
}
