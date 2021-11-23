package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
)

type Units interface {
	GetDomainByID(unitId int) (domain string, err error)
	GetIDByDomain(unitDomain string) (id int, err error)
	UnitExistsByID(unitId int) (exists bool)
	UnitExistsByDomain(unitDomain string) (exists bool)
}

type Users interface {
	CreateUser(userModel *models.CreateUser) (id int, err error) // todo: проверка наличия дублирующей записи в бд
	GetUserData(userId int) (user models.UserData, err error)
	GetUserIdByInput(inputModel *models.UserInput) (id int, err error)
	UserExistsByInput(inputModel *models.UserInput) (exists bool)
	GetUserByDomain(userDomain string) (user models.UserInfo, err error)
	GetUserByID(userId int) (user models.UserInfo, err error)
	GetUsersByNameFragment(fragment string, offset int) (users models.ListUserInfo, err error)

	GetUserSettings(userId int) (settings models.UserSettings, err error)
	UpdateUserData(userId int, userModel *models.UpdateUserData) (err error)
	UpdateUserSettings(userId int, settingsModel *models.UpdateUserSettings) error

	GetCountUserOwnedChats(userId int) (count int, err error)
	UserExistsByID(userId int) (exists bool)
	UserExistsByDomain(userDomain string) (exists bool)
}

type Auth interface {
	CreateNewUserRefreshSession(userId int, sessionModel *models.RefreshSession) (sessionsCount int, err error)
	DeleteOldestSession(userId int) (err error)
	FindSessionByComparedToken(token string) (sessionId int, userId int, err error)
	UpdateRefreshSession(sessionId int, sessionModel *models.RefreshSession) (err error)
}

type Chats interface {
	CreateChat(ownerId int, chatModel *models.CreateChat) (id int, err error)
	GetChatByDomain(chatDomain string) (chat models.ChatInfo, err error)
	GetChatByID(chatId int) (chat models.ChatInfo, err error)
	GetCountChatMembers(chatId int) (count int, err error)
	GetChatsByNameFragment(fragment string, offset int) (chats models.ListChatInfo, err error)
	GetChatMembers(chatId int) (members models.ListUserInfo, err error)

	GetChatDataByID(chatId int) (chat models.ChatData, err error)
	UpdateChatData(chatId int, inputModel *models.UpdateChatData) (err error)

	UserIsChatOwner(userId int, chatId int) bool
	UserIsChatMember(userId int, chatId int) bool
	AddUserToChat(userId int, chatId int) (err error)
	RemoveUserFromChat(userId int, chatId int) (err error)

	GetChatsOwnedUser(userId int, offset int) (chats models.ListChatInfo, err error)
	GetChatsInvolvedUser(userId int, offset int) (chats models.ListChatInfo, err error)
	GetCountRooms(chatId int) (count int, err error)
	GetCountUserChats(userId int) (count int, err error)
	ChatExistsByID(chatId int) (exists bool)
	ChatExistsByDomain(chatDomain string) (exists bool)

	// invites
	GetCountLinks(chatId int) (count int, err error)
	GetChatLinks(chatId int) (links models.InviteLinks, err error)
	LinkExistsByCode(code string) (exists bool) // ! equal relevant
	FindInviteLinkByCode(code string) (link models.InviteLink, err error)
	DeleteInviteLinkByCode(code string) (err error)
	CreateInviteLink(linkModel *models.CreateInviteLink) (link models.InviteLink, err error)
	InviteLinkIsRelevant(code string) (relevant bool)
	AddUserByCode(code string, userId int) (chatId int, err error)

	ChatIsPrivate(chatId int) (private bool)

	BanUserInChat(userId int, chatId int) (err error)
	UnbanUserInChat(userId int, chatId int) (err error)
	UserIsBannedInChat(userId int, chatId int) (banned bool)
	GetChatBanlist(chatId int) (users models.ListUserInfo, err error)

	//GetChatOwnerID(chat_id int)

	GetUserRoleData(userId int, chatId int) (role models.RoleData, err error)
	GetUserRoleInfo(userId int, chatId int) (role models.RoleInfo, err error)
	CreateRoleInChat(chatId int, roleModel *models.CreateRole) (roleId int, err error)
	GetChatRolesData(chatId int) (roles models.ListRolesData, err error)
	GetChatRolesInfo(chatId int) (roles models.ListRolesInfo, err error)
	GetCountChatRoles(chatId int) (count int, err error)
	// ? GetCountUserRoles(user_id int, chat_id int) (count int, err error)
	GiveRole(userId int, roleId int) (err error)
	RoleExistsByID(roleId int, chatId int) (exists bool)
	TakeRole(userId int, chatId int) (err error)
	UpdateRoleData(roleId int, inputModel *models.UpdateRole) (err error)
	DeleteRole(roleId int) (err error)
	// ? user have role?

	GetMemberInfo(userId int, chatId int) (user models.MemberInfo, err error)
}
type Rooms interface { //todo: get parent and child rooms
	RoomExistsByID(roomId int) (isExists bool)
	CreateRoom(chatId int, roomModel *models.CreateRoom) (roomId int, err error)
	GetRoom(roomId int) (room models.RoomInfo, err error)
	UpdateRoomData(roomId int, inputModel *models.UpdateRoomData) (err error)
	GetChatIDByRoomID(roomId int) (chatId int, err error)
	GetChatRooms(chatId int) (rooms models.ListRoomInfo, err error)
	RoomIsPrivate(roomId int) (private bool)
	// todo SetRoomParent(room_id int, parent_id int) (err error)
	// todo GetChildRooms(room_id int) (childs models.ListRoomInfo, err error)

	RoomFormIsSet(roomId int) (isSet bool)
	GetRoomForm(roomId int) (form models.FormPattern, err error)
	UpdateRoomForm(roomId int, format string) (err error)
}

type Dialogs interface {
	GetDialogIDBetweenUsers(user1Id int, user2Id int) (dialogId int, err error)
	GetCompanions(userId int) (users models.ListUserInfo, err error)
	DialogExistsBetweenUsers(user1Id int, user2Id int) (exits bool)
	CreateDialogBetweenUser(user1Id int, user2Id int) (dialogId int, err error)
}

type Messages interface {
	CreateMessageInRoom(roomId int, msgType rules.MessageType, messageModel *models.CreateMessage) (messageId int, err error)
	CreateMessageInDialog(dialogId int, messageModel *models.CreateMessage) (messageId int, err error)
	GetMessagesFromRoom(roomId int, createdAfter int, offset int) (messages models.MessagesList, err error)
	GetMessagesFromDialog(dialogId int, offset int) (messages models.MessagesList, err error)
	GetMessageFromRoom(messageId int, roomId int) (message models.MessageInfo, err error)
	GetMessageFromDialog(messageId int, dialogId int) (message models.MessageInfo, err error)
	MessageExistsByID(messageId int) (exists bool)
	MessageAvailableOnRoom(messageId int, roomId int) (exists bool)
	MessageAvailableOnDialog(messageId int, dialogId int) (exists bool)
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
