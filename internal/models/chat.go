package models

// general chat model
type Chat struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	Domain  string `json:"domain"`
	Name    string `json:"name"`
}

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

// RESPONSE models
type ChatInfo struct {
	ID           int    `json:"id"`
	OwnerID      int    `json:"owner_id"`
	Domain       string `json:"domain"`
	Name         string `json:"name"`
	CountMembers int    `json:"count_members"`
}
type ChatData struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	Domain  string `json:"domain"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
}
type ChatMembersCount struct {
	Count int `json:"count"`
}
type ListChatInfo struct {
	Chats []ChatInfo `json:"chats"`
}
type ListChatMembers struct {
	Members []UserInfo `json:"members"`
}
type ChatID struct {
	ID int `json:"id"`
}
type InviteLink struct {
	Code   string `json:"code"`
	ChatID int    `json:"chat_id,omitempty"`
	Aliens int    `json:"aliens"`
	Exp    int64  `json:"exp"`
}
type InviteLinkInput struct {
	Aliens   int   `json:"aliens"`
	LifeTime int64 `json:"lifetime"`
}
type CreateInviteLink struct {
	ChatID int   `json:"chat_id"`
	Aliens int   `json:"aliens"`
	Exp    int64 `json:"exp"`
}
type InviteLinks struct {
	Links []InviteLink
}
type MemberInfo struct {
	RoleID   int `json:"role_id"`
	JoinedAt int `json:"joined_at"`
}

// todo: user is chat admin
