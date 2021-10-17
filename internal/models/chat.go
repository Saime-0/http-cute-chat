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
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type UpdateChatData struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}
type ChatName struct {
	Name string `json:"name"`
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
type ListChatRooms struct {
	Rooms []RoomInfo `json:"rooms"`
}
