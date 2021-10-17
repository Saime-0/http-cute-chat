package models

// general message model
type Response struct {
	ID      int    `json:"id"`
	ReplyTo int    `json:"reply_to"`
	Author  int    `json:"author"`
	Body    string `json:"body"`
	Type    int    `json:"type"`
}

// INTERNAL models
// RECEIVED models

// RESPONSE models
