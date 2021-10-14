package models

// general chat model
type Chat struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	Domain  string `json:"domain"`
	Name    string `json:"name"`
}

//  INTERNAL models

// RECEIVED models
type CreateChatInput struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type CreateRoomInput struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type UpdateChatInput struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

// RESPONSE models
type ChatInfo struct {
	OwnerDomain int    `json:"owner_domain"`
	Domain      string `json:"domain"`
	Name        string `json:"name"`
}

type ChatMembers struct {
	Users []UserInfo
}

type ChatRooms struct {
	Rooms []RoomInfo
}
