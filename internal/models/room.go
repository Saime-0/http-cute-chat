package models

// general room model
type Room struct {
	ID         int    `json:"id"`
	ChatID     int    `json:"chat_id"`
	ParentRoom int    `json:"parent_room"`
	Name       string `json:"name"`
	Note       string `json:"note"`
}

// INTERNAL models
// RECEIVED models
type CreateRoom struct {
	ChatID     int    `json:"chat_id"`
	Name       string `json:"name"`
	ParentRoom int    `json:"parent_room"`
	Note       string `json:"note"`
}
type UpdateRoomData struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

// RESPONSE models
type RoomInfo struct {
	ID         int    `json:"id"`
	ParentRoom int    `json:"parent_room"`
	Name       string `json:"name"`
	Note       string `json:"note"`
}
type ListRoomInfo struct {
	Rooms []RoomInfo `json:"rooms"`
}
type RoomID struct {
	ID int `json:"id"`
}

// todo: change parent room and clear parent room
