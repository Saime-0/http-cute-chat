package models

import "github.com/saime-0/http-cute-chat/internal/api/rules"

// general message model
// type Message struct {
// 	ID      int    `json:"id"`
// 	ReplyTo int    `json:"reply_to"`
// 	Author  int    `json:"author"`
// 	Body    string `json:"body"`
// 	UnitType    int    `json:"type"`
// }

// INTERNAL models

// RECEIVED models
type CreateMessage struct {
	ReplyTo int    `json:"reply_to"`
	Author  int    `json:"author"`
	Body    string `json:"body"`
	// UnitType    rules.MessageType `json:"type"`
}

// RESPONSE models
type MessageInfo struct {
	ID      int               `json:"id"`
	ReplyTo int               `json:"reply_to"`
	Author  int               `json:"author"`
	Body    string            `json:"body"`
	Type    rules.MessageType `json:"type"`
	Time    int               `json:"time"`
}
type MessagesList struct {
	Messages []MessageInfo `json:"messages"`
}
type MessageID struct {
	ID int `json:"id"`
}

/* // todo
msg_type:
	- system message
		- sender
		- event
		- body
	- user message
	- formatted message
msg_fields:
	- text
	- photo
	- file
	- vote
	- music
	- video
*/
