package repository

import (
	"database/sql"
)

type Repositories struct {
	Auth     *AuthRepo
	Units    *UnitsRepo
	Users    *UsersRepo
	Chats    *ChatsRepo
	Rooms    *RoomsRepo
	Messages *MessagesRepo
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Auth:     NewAuthRepo(db),
		Units:    NewUnitsRepo(db),
		Users:    NewUsersRepo(db),
		Chats:    NewChatsRepo(db),
		Rooms:    NewRoomsRepo(db),
		Messages: NewMessagesRepo(db),
	}
}
