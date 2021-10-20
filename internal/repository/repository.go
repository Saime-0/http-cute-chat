package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/models"
)

// revision
type Units interface {
	GetDomainByID(id int) (domain string, err error)
	GetIDByDomain(domain string) (id int, err error)
}

type Users interface {
	// user is not exists

	CreateUser(u *models.CreateUser) (id int, err error) // todo: проверка наличия дублирующей записи в бд
	GetUserData(user_id int) (user models.UserData, err error)
	GetUserIdByInput(input models.UserInput) (id int, err error)
	GetUserInfoByDomain(domain string) (user models.UserInfo, err error)
	GetUserInfoByID(id int) (user models.UserInfo, err error)
	GetListUsersByName(name string) (chats []models.UserInfo, err error)
	GetListChatsOwnedUser(user_id int) (chats []models.ChatInfo, err error)
	GetListChatsUser(user_id int) (chats []models.ChatInfo, err error)
	IsUserExistsByInput(input models.UserInput) (exists bool) // new
	// todo: get jwt: "id"
	GetUserSettings(user_id int) (settings *models.UserSettings, err error)
	UpdateUserData(user_id int, inp *models.UpdateUserData) error
	UpdateUserSettings(user_id int, inp *models.UpdateUserSettings) error

	// todo: jwt serv methods, create Auth or Sessions interface
	// ? todo: get by refresh token
	// todo: limit 5
	CreateNewUserRefreshSession(user_id int, s *models.RefreshSession) (err error)
	FindSessionByComparedToken(token string) (session_id int, user_id int, err error)
	UpdateRefreshSession(session_id int, s *models.RefreshSession) (err error)
}
type Chats interface {
	CreateChat(owner_id int, c *models.CreateChat) (id int, err error)
	GetChatInfoByDomain(domain string) (chat models.ChatInfo, err error)
	GetChatInfoByID(chat_id int) (chat models.ChatInfo, err error)
	// GetCountChatMembers(chat_id int) (count int, err error)
	GetListChatsByName(name string) (chats []models.ChatInfo, err error)
	GetListChatMembers(chat_id int) (members []models.UserInfo, err error)
	GetListChatRooms(chat_id int) (rooms []models.RoomInfo, err error)

	GetChatDataByID(id int) (chat models.ChatData, err error)
	UpdateChatData(chat_id int, inp *models.UpdateChatData) error

	UserIsChatOwner(user_id int, chat_id int) (res bool, err error)
	AddUserToChat(user_id int, chat_id int) error
}
type Rooms interface {
	CreateMessage(room_id int, m *models.CreateMessage) (message_id int, err error)
	GetListMessages(room_id int) (messages []models.MessageInfo, err error)
	// GetMessageInfo

	CreateRoom(r *models.Room) (id int, err error)
	UpdateRoomData(inp *models.UpdateUserData) error
}

// ? revision
type Dialogs interface {
	//
	CreateMessage(dialog_id int, m *models.CreateMessage) error
	//IsDialogBetweenUsersAvailable
	GetDialogIDBetweenUsers(user1_id int, user2_id int) (dialog_id int, err error)
	GetCompanionsIDs(user_id int) (id []int, err error)
	GetListMessages(dialog_id int) (messages []models.MessageInfo, err error)
}

// type Auth interface {
// 	SignIn(u models.CreateUser) error
// 	SignUp(u models.CreateUser) error
// 	RefreshToken(u models.RefreshToken) error
// }

type Repositories struct {
	Users Users
	Chats Chats
	Rooms Rooms
	// Dialogs Dialogs
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}
