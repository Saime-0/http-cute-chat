package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/models"
)

// revision
type Units interface {
	GetDomainByID(unit_id int) (domain string, err error)
	GetIDByDomain(unit_domain string) (id int, err error)
	IsUnitExistsByID(unit_id int) (exists bool)
	IsUnitExistsByDomain(unit_domain string) (exists bool)
}

type Users interface {
	// user is not exists

	CreateUser(user_model *models.CreateUser) (id int, err error) // todo: проверка наличия дублирующей записи в бд
	GetUserData(user_id int) (user models.UserData, err error)
	GetUserIdByInput(input_model models.UserInput) (id int, err error)
	GetUserInfoByDomain(domain string) (user models.UserInfo, err error)
	GetUserInfoByID(id int) (user models.UserInfo, err error)
	GetUsersByNameFragment(fragment string, offset int) (users models.ListUserInfo, err error)
	IsUserExistsByInput(input_model models.UserInput) bool // new
	// todo: get jwt: "id"
	GetUserSettings(user_id int) (settings models.UserSettings, err error)
	UpdateUserData(user_id int, user_model *models.UpdateUserData) error
	UpdateUserSettings(user_id int, settings_model *models.UpdateUserSettings) error

	// todo: jwt serv methods, create Auth or Sessions interface
	// ? todo: get by refresh token
	// todo: limit 5
	CreateNewUserRefreshSession(user_id int, session_model *models.RefreshSession) (sessions_count int, err error)
	DeleteOldestSession(user_id int) (err error)
	FindSessionByComparedToken(token string) (session_id int, user_id int, err error)
	UpdateRefreshSession(session_id int, session_model *models.RefreshSession) (err error)
	GetCountUserOwnedChats(user_id int) (count int, err error)
}
type Chats interface {
	CreateChat(owner_id int, chat_model *models.CreateChat) (id int, err error)
	GetChatInfoByDomain(domain string) (chat models.ChatInfo, err error)
	GetChatInfoByID(chat_id int) (chat models.ChatInfo, err error)
	GetCountChatMembers(chat_id int) (count int, err error)
	GetChatsByNameFragment(name string, offset int) (chats models.ListChatInfo, err error)
	GetChatMembers(chat_id int) (members models.ListUserInfo, err error)

	GetChatDataByID(chat_id int) (chat models.ChatData, err error)
	UpdateChatData(chat_id int, input_model *models.UpdateChatData) (err error)

	UserIsChatOwner(user_id int, chat_id int) bool
	UserIsChatMember(user_id int, chat_id int) bool
	AddUserToChat(user_id int, chat_id int) (err error)

	GetChatsOwnedUser(user_id int) (chats models.ListChatInfo, err error)
	GetChatsInvolvedUser(user_id int) (chats models.ListChatInfo, err error)
	GetCountRooms(chat_id int) (count int, err error)
	GetCountUserChats(user_id int) (count int, err error)
}
type Rooms interface { //todo: get parent and child rooms
	CreateMessage(room_id int, message_model *models.CreateMessage) (message_id int, err error)
	GetMessages(room_id int) (messages models.MessagesList, err error)
	IsRoomExistsByID(room_id int) (is_exists bool)

	CreateRoom(chat_id int, room_model *models.CreateRoom) (room_id int, err error)
	GetRoomInfo(room_id int) (room models.RoomInfo, err error)
	UpdateRoomData(room_id int, input_model *models.UpdateRoomData) (err error)

	GetChatIDByRoomID(room_id int) (chat_id int, err error)
	GetMessageInfo(message_id int, room_id int) (message models.MessageInfo, err error)

	// ! GetChildRooms(room_id int) (childs models.ListRoomInfo, err error)
	GetChatRooms(chat_id int) (rooms models.ListRoomInfo, err error)
}

// ? revision
type Dialogs interface {
	//
	CreateMessage(dialog_id int, message_model *models.CreateMessage) (message_id int, err error)
	//IsDialogBetweenUsersAvailable
	GetDialogIDBetweenUsers(user1_id int, user2_id int) (dialog_id int, err error)
	GetCompanions(user_id int) (users models.ListUserInfo, err error)
	GetMessages(dialog_id int) (messages models.MessagesList, err error)
	GetMessageInfo(message_id int, dialog_id int) (message models.MessageInfo, err error)
	DialogIsExistsBetweenUsers(user1_id int, user2_id int) (exits bool)
	CreateDialogBetweenUser(user1_id int, user2_id int) (dialog_id int, err error)
}

// type Auth interface {
// 	SignIn(u models.CreateUser) error
// 	SignUp(u models.CreateUser) error
// 	RefreshToken(u models.RefreshToken) error
// }

type Repositories struct {
	Units   Units
	Users   Users
	Chats   Chats
	Rooms   Rooms
	Dialogs Dialogs
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Units:   NewUnitsRepo(db),
		Users:   NewUsersRepo(db),
		Chats:   NewChatsRepo(db),
		Rooms:   NewRoomsRepo(db),
		Dialogs: NewDialogsRepo(db),
	}
}
