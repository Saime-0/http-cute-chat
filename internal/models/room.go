package models

// general room model
type Room struct {
	ID         int    `json:"id"`
	ChatID     int    `json:"chat_id"`
	ParentRoom int    `json:"parent_room"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
}

//  INTERNAL models

// RECEIVED models

// RESPONSE models
type RoomInfo struct {
	ParentRoom int    `json:"parent_room"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
}
