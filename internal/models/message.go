package models

// general message model
type Message struct {
	ID         int    `json:"id"`
	ReplyTo    int    `json:"reply_to"`
	FromDomain int    `json:"From"`
	Body       string `json:"body"`
	Type       int    `json:"type"`
}

//  INTERNAL models

// RECEIVED models

// RESPONSE models
type MessageInfo struct {
	ReplyTo    int    `json:"reply_to"`
	FromDomain int    `json:"from_domain"`
	Body       string `json:"body"`
	Type       int    `json:"type"`
}
