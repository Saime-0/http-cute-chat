package models

// general room model
type Room struct {
	ID         int    `json:"id"`
	ChatID     int    `json:"chat_id"`
	ParentRoom int    `json:"parent_room"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
}

// INTERNAL models
// RECEIVED models
type CreateRoom struct {
	ChatID     int    `json:"chat_id"`
	Name       string `json:"name"`
	ParentRoom int    `json:"parent_room"`
	Desc       string `json:"desc"`
}
type UpdateRoomData struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// RESPONSE models
type RoomInfo struct {
	ID         int    `json:"id"`
	ParentRoom int    `json:"parent_room"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
}

// todo: change parent room and clear parent room
