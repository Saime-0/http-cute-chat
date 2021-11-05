package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type Units interface {
	GetDomainByID(unit_id int) (domain string, err error)
	GetIDByDomain(unit_domain string) (id int, err error)
	UnitExistsByID(unit_id int) (exists bool)
	UnitExistsByDomain(unit_domain string) (exists bool)
}

type Users interface {
	CreateUser(user_model *models.CreateUser) (id int, err error) // todo: проверка наличия дублирующей записи в бд
	GetUserData(user_id int) (user models.UserData, err error)
	GetUserIdByInput(input_model *models.UserInput) (id int, err error)
	UserExistsByInput(input_model *models.UserInput) (exists bool)
	GetUserByDomain(user_domain string) (user models.UserInfo, err error)
	GetUserByID(user_id int) (user models.UserInfo, err error)
	GetUsersByNameFragment(fragment string, offset int) (users models.ListUserInfo, err error)

	GetUserSettings(user_id int) (settings models.UserSettings, err error)
	UpdateUserData(user_id int, user_model *models.UpdateUserData) (err error)
	UpdateUserSettings(user_id int, settings_model *models.UpdateUserSettings) error

	GetCountUserOwnedChats(user_id int) (count int, err error)
	UserExistsByID(user_id int) (exists bool)
	UserExistsByDomain(user_domain string) (exists bool)
}

type Auth interface {
	CreateNewUserRefreshSession(user_id int, session_model *models.RefreshSession) (sessions_count int, err error)
	DeleteOldestSession(user_id int) (err error)
	FindSessionByComparedToken(token string) (session_id int, user_id int, err error)
	UpdateRefreshSession(session_id int, session_model *models.RefreshSession) (err error)
}

type Chats interface {
	CreateChat(owner_id int, chat_model *models.CreateChat) (id int, err error)
	GetChatByDomain(chat_domain string) (chat models.ChatInfo, err error)
	GetChatByID(chat_id int) (chat models.ChatInfo, err error)
	GetCountChatMembers(chat_id int) (count int, err error)
	GetChatsByNameFragment(fragment string, offset int) (chats models.ListChatInfo, err error)
	GetChatMembers(chat_id int) (members models.ListUserInfo, err error)

	GetChatDataByID(chat_id int) (chat models.ChatData, err error)
	UpdateChatData(chat_id int, input_model *models.UpdateChatData) (err error)

	UserIsChatOwner(user_id int, chat_id int) bool
	UserIsChatMember(user_id int, chat_id int) bool
	AddUserToChat(user_id int, chat_id int) (err error)
	RemoveUserFromChat(user_id int, chat_id int) (err error)

	GetChatsOwnedUser(user_id int, offset int) (chats models.ListChatInfo, err error)
	GetChatsInvolvedUser(user_id int, offset int) (chats models.ListChatInfo, err error)
	GetCountRooms(chat_id int) (count int, err error)
	GetCountUserChats(user_id int) (count int, err error)
	ChatExistsByID(chat_id int) (exists bool)
	ChatExistsByDomain(chat_domain string) (exists bool)

	// invites
	GetCountLinks(chat_id int) (count int, err error)
	GetChatLinks(chat_id int) (links models.InviteLinks, err error)
	LinkExistsByCode(code string) (exists bool) // ! equal relevant
	FindInviteLinkByCode(code string) (link models.InviteLink, err error)
	DeleteInviteLinkByCode(code string) (err error)
	CreateInviteLink(link_model *models.CreateInviteLink) (link models.InviteLink, err error)
	InviteLinkIsRelevant(code string) (relevant bool)
	AddUserByCode(code string, user_id int) (chat_id int, err error)

	ChatIsPrivate(chat_id int) (private bool)

	BanUserInChat(user_id int, chat_id int) (err error)
	UnbanUserInChat(user_id int, chat_id int) (err error)
	UserIsBannedInChat(user_id int, chat_id int) (banned bool)
	GetChatBanlist(chat_id int) (users models.ListUserInfo, err error)

	//GetChatOwnerID(chat_id int)

	GetUserRoleData(user_id int, chat_id int) (role models.RoleData, err error)
	GetUserRoleInfo(user_id int, chat_id int) (role models.RoleInfo, err error)
	CreateRoleInChat(chat_id int, role_model *models.CreateRole) (role_id int, err error)
	GetChatRolesData(chat_id int) (roles models.ListRolesData, err error)
	GetChatRolesInfo(chat_id int) (roles models.ListRolesInfo, err error)
	GetCountChatRoles(chat_id int) (count int, err error)
	// ? GetCountUserRoles(user_id int, chat_id int) (count int, err error)
	GiveRole(user_id int, role_id int) (err error)
	RoleExistsByID(role_id int, chat_id int) (exists bool)
	TakeRole(user_id int, chat_id int) (err error)
	UpdateRoleData(role_id int, input_model *models.UpdateRole) (err error)
	DeleteRole(role_id int) (err error)
	// ? user have role?
}
type Rooms interface { //todo: get parent and child rooms
	RoomExistsByID(room_id int) (is_exists bool)
	CreateRoom(chat_id int, room_model *models.CreateRoom) (room_id int, err error)
	GetRoom(room_id int) (room models.RoomInfo, err error)
	UpdateRoomData(room_id int, input_model *models.UpdateRoomData) (err error)
	GetChatIDByRoomID(room_id int) (chat_id int, err error)
	GetChatRooms(chat_id int) (rooms models.ListRoomInfo, err error)
	RoomIsPrivate(room_id int) (private bool)
	// todo SetRoomParent(room_id int, parent_id int) (err error)
	// todo GetChildRooms(room_id int) (childs models.ListRoomInfo, err error)
}

type Dialogs interface {
	GetDialogIDBetweenUsers(user1_id int, user2_id int) (dialog_id int, err error)
	GetCompanions(user_id int) (users models.ListUserInfo, err error)
	DialogExistsBetweenUsers(user1_id int, user2_id int) (exits bool)
	CreateDialogBetweenUser(user1_id int, user2_id int) (dialog_id int, err error)
}

type Messages interface {
	CreateMessageInRoom(room_id int, message_model *models.CreateMessage) (message_id int, err error)
	CreateMessageInDialog(dialog_id int, message_model *models.CreateMessage) (message_id int, err error)
	GetMessagesFromRoom(room_id int, offset int) (messages models.MessagesList, err error)
	GetMessagesFromDialog(dialog_id int, offset int) (messages models.MessagesList, err error)
	GetMessageFromRoom(message_id int, room_id int) (message models.MessageInfo, err error)
	GetMessageFromDialog(message_id int, dialog_id int) (message models.MessageInfo, err error)
	MessageExistsByID(message_id int) (exists bool)
	MessageAvailableOnRoom(message_id int, room_id int) (exists bool)
	MessageAvailableOnDialog(message_id int, dialog_id int) (exists bool)
}

type Repositories struct {
	Auth     Auth
	Units    Units
	Users    Users
	Chats    Chats
	Rooms    Rooms
	Dialogs  Dialogs
	Messages Messages
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Auth:     NewAuthRepo(db),
		Units:    NewUnitsRepo(db),
		Users:    NewUsersRepo(db),
		Chats:    NewChatsRepo(db),
		Rooms:    NewRoomsRepo(db),
		Dialogs:  NewDialogsRepo(db),
		Messages: NewMessagesRepo(db),
	}
}
