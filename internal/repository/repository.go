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
	CreateChat(u models.CreateUser) error
	GetUserDataByID(id int) (user models.UserData, err error)
	GetUserInfoByDomain(domain string) (user models.UserInfo, err error)
	GetUserInfoByID(id int) (user models.UserInfo, err error)

	GetUserSettings() (settings models.UserSettings, err error)
	UpdateUserData(inp models.UpdateUserData) error
	UpdateUserSettings(inp models.UpdateUserSettings) error

	// todo: jwt serv methods, create Auth or Sessions interface

}
type Chats interface {
	CreateChat(c models.Chat) error
	GetChatInfoByDomain(domain string) (chat models.ChatInfo, err error)
	GetChatInfoByID(id int) (chat models.ChatInfo, err error)
	GetCountChatMembers(chat_id int) (count int, err error)
	GetListChatsByName(name string) (chats []models.ChatInfo, err error)
	GetListChatMembers(chat_id int) (members []models.UserInfo, err error)
	GetListChatRooms(chat_id int) (rooms []models.RoomInfo, err error)

	GetChatDataByID(id int) (chat models.ChatData, err error)
	UpdateChatData(inp models.UpdateUserData) error
}
type Rooms interface {
	CreateMessage(room_id, m models.CreateMessage) error
	GetListMessages(room_id int) (messages []models.MessageInfo, err error)
	// GetMessageInfo

	CreateRoom(r models.Room) error
	UpdateRoomData(inp models.UpdateUserData) error
}
type Dialogs interface {
	//
	CreateMessage(dialog_id int, m models.CreateMessage) error
	GetDialogIDWithUser(user_id int) (id int, err error)
	GetUserDialogsIDs() error
	GetMessagesFromDialog(dialog_id int) error
}

// type Auth interface {
// 	SignIn(u models.CreateUser) error
// 	SignUp(u models.CreateUser) error
// 	RefreshToken(u models.RefreshToken) error
// }

type Repositories struct {
	Users Users
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}
