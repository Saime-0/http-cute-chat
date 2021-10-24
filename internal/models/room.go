package models

// general room model
type Room struct {
	ID       int    `json:"id"`
	ChatID   int    `json:"chat_id"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
	Note     string `json:"note"`
}

// INTERNAL models
// RECEIVED models
type CreateRoom struct {
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
	Note     string `json:"note"`
}
type UpdateRoomData struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

// RESPONSE models
type RoomInfo struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
	Note     string `json:"note"`
}
type ListRoomInfo struct {
	Rooms []RoomInfo `json:"rooms"`
}
type RoomID struct {
	ID int `json:"id"`
}

// todo: change parent room and clear parent room
