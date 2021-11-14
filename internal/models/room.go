package models

// general room model
// type Room struct {
// 	ID       int    `json:"id"`
// 	ChatID   int    `json:"chat_id"`
// 	ParentID int    `json:"parent_id"`
// 	Name     string `json:"name"`
// 	Note     string `json:"note"`
// 	Private  bool   `json:"private"`
// }

// INTERNAL models
// RECEIVED models
type CreateRoom struct {
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
	Note     string `json:"note"`
	Private  bool   `json:"private"`
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
