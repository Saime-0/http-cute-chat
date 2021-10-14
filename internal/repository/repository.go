package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type Units interface {
	GetDomainByID(id int) (domain string, err error)
	GetIDByDomain(domain string) (id int, err error)
}

type Users interface {
	CreateUser(u models.User) error
	GetUserByDomain(domain string) (user models.User, err error)
	GetUserByID(id int) (user models.User, err error)
	GetUserSettings(id int) error
	UpdateUserInfo(inp models.UpdateUserInput) error
	UpdateUserSettings(inp models.UpdateUserInput) error
}
type Chats interface {
	CreateChat(c models.Chat) error
	GetChatByDomain(domain string) error
	GetChatByID(id int) error
	GetChatsByName(name string) error
	GetMembersRelatingChatID(chat_id int) error
	GetRoomsRelatingChatID(chat_id int) error
}
type Rooms interface {
	CreateRoom(chat_id int, r models.Room) (id int, err error)
	CreateMessage(m models.Message) error
	GetMessages(room_id int) error
}
type Dialogs interface {
	CreateMessage(m models.Message) error
	GetDialogIDBetweenUsers(user1_id int, user2_id int) error
	GetUserDialogsIDs(user_id int) error
	GetMessagesFromDialog(dialog_id int) error
}

type Repositories struct {
	Users Users
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}
