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

// todo: user is chat admin
