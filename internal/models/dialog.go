package models

// general room model
type Dialog struct {
	ID          int      `json:"id"`
	MessagesIDs []int    `json:"messages_ids"`
	Companion   UserInfo `json:"companion"`
}

// INTERNAL models
// RECEIVED models

// RESPONSE models
type DialogInfo struct {
	Companion UserInfo `json:"companion"`
}
type ListDialogs struct {
	Dialogs []DialogInfo `json:"dialogs"`
}

// ? todo: messageinfo by id (implementation "message.reply_to")
