package models

// general message model
type Message struct {
	ID      int    `json:"id"`
	ReplyTo int    `json:"reply_to"`
	Author  int    `json:"author"`
	Body    string `json:"body"`
	Type    int    `json:"type"`
}

// INTERNAL models

// RECEIVED models
type CreateMessage struct {
	ReplyTo int    `json:"reply_to"`
	Author  int    `json:"author"`
	Body    string `json:"body"`
	Type    int    `json:"type"`
}

// RESPONSE models
type MessageInfo struct {
	ID      int    `json:"id"`
	ReplyTo int    `json:"reply_to"`
	Author  int    `json:"author"`
	Body    string `json:"body"`
	Type    int    `json:"type"`
}
type MessagesList struct {
	Messages []MessageInfo `json:"messages"`
}
type MessageID struct {
	ID int `json:"id"`
}
