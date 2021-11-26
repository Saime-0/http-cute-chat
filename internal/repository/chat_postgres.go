package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"time"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type ChatsRepo struct {
	db *sql.DB
}

func NewChatsRepo(db *sql.DB) *ChatsRepo {
	return &ChatsRepo{
		db: db,
	}
}

func (r *ChatsRepo) CreateChat(ownerId int, chatModel *models.CreateChat) (id int, err error) {
	err = r.db.QueryRow(
		`WITH u AS (
			INSERT INTO units (domain, name) 
			VALUES ($1, $2) 
			RETURNING id
			) 
		INSERT INTO chats (id, owner_id, private) 
		SELECT u.id, $3, $4
		FROM u 
		RETURNING id`,
		chatModel.Domain,
		chatModel.Name,
		ownerId,
		chatModel.Private,
	).Scan(&id)

	return
}

func (r *ChatsRepo) GetChatByID(chatId int) (chat models.Chat, err error) {
	err = r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, chats.private
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id = $1`,
		chatId,
	).Scan(
		&chat.ID,
		&chat.Domain,
		&chat.Name,
		&chat.Private,
	)

	return
}
func (r *ChatsRepo) GetCountChatMembers(chatId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM chat_members  
		WHERE chat_id = $1`,
		chatId,
	).Scan(&count)

	return
}

func (r *ChatsRepo) GetChatsByNameFragment(fragment string, limit int, offset int) (chats []models.Chat, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, chats.private
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.name ILIKE $1 AND chats.private = FALSE
		LIMIT $2
		OFFSET $3`,
		"%"+fragment+"%",
		limit,
		offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.Chat{}
		if err = rows.Scan(&m.ID, &m.Domain, &m.Name, &m.Private); err != nil {
			return
		}
		if err != nil {
			return
		}
		chats = append(chats, m)
	}

	return

}
func (r *ChatsRepo) GetChatMembers(chatId int) (members models.ListChatMembers, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, member.role_id, member.char, member.joined_at
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.id
		WHERE member.chat_id = $1
		`,
		chatId,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.Member{}
		if err = rows.Scan(&m.User.ID, &m.User.Domain, &m.User.Name, &m.RoleID, &m.Char, &m.JoinedAt); err != nil {
			return
		}
		members.Members = append(members.Members, m)
	}

	return
}
func (r *ChatsRepo) GetChatMember(userId, chatId int) (member models.Member, err error) {
	err = r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, member.role_id, member.char, member.joined_at
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.id
		WHERE member.user_id = $1 AND member.chat_id = $2`,
		userId,
		chatId,
	).Scan(
		&member.User.ID,
		&member.User.Domain,
		&member.User.Name,
		&member.RoleID,
		&member.Char,
		&member.JoinedAt,
	)

	return
}

func (r *ChatsRepo) GetUserChar(userId int, chatId int) (char rules.CharType, err error) {
	err = r.db.QueryRow(
		`SELECT char
		FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
	).Scan(&char)

	return
}

func (r *ChatsRepo) UserIs(userId int, chatId int, char rules.CharType) (yes bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE user_id = $1 AND chat_id = $2 AND char = $3
    	)`,
		userId,
		chatId,
		char,
	).Scan(&yes)

	return
}

func (r *ChatsRepo) GetChars(chatId int, char rules.CharType) (admins models.ListChatMembers, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, member.role_id, member.char, member.joined_at
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.id
		WHERE member.chat_id = $1 AND member.char = $2`,
		chatId,
		char,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.Member{}
		if err = rows.Scan(&m.User.ID, &m.User.Domain, &m.User.Name, &m.RoleID, &m.Char, &m.JoinedAt); err != nil {
			return
		}
		admins.Members = append(admins.Members, m)
	}
	return
}

func (r *ChatsRepo) UpdateChatData(chatId int, inputModel *models.UpdateChatData) (err error) {
	if inputModel.Domain != "" {
		err = r.db.QueryRow(
			`UPDATE units
			SET domain = $2
			WHERE id = $1`,
			chatId,
			inputModel.Domain,
		).Err()
		if err != nil {
			return
		}
	}
	if inputModel.Name != "" {
		err = r.db.QueryRow(
			`UPDATE units
			SET name = $2
			WHERE id = $1`,
			chatId,
			inputModel.Name,
		).Err()
		if err != nil {
			return
		}
	}

	err = r.db.QueryRow(
		`UPDATE rooms
		SET private = $2
		WHERE id = $1`,
		chatId,
		inputModel.Private,
	).Err()

	return
}

func (r *ChatsRepo) UserIsChatOwner(userId int, chatId int) bool {
	isOwner := false
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM chats WHERE id = $1 AND owner_id = $2)`,
		chatId,
		userId,
	).Scan(&isOwner)
	if err != nil || !isOwner {
		return isOwner
	}
	return isOwner
}
func (r *ChatsRepo) UserIsChatMember(userId int, chatId int) bool {
	isMember := false
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM chat_members WHERE user_id = $1 AND chat_id = $2)`,
		userId,
		chatId,
	).Scan(&isMember)
	if err != nil || !isMember {
		return isMember
	}
	return isMember
}
func (r *ChatsRepo) AddUserToChat(userId int, chatId int) (err error) {
	err = r.db.QueryRow(
		`INSERT INTO chat_members (user_id, chat_id, joined_at)
		VALUES ($1, $2, $3)`,
		userId,
		chatId,
		time.Now().UTC().Unix(),
	).Err()
	if err != nil {
		return
	}
	return
}

// migrate from UsersRepo
func (r *ChatsRepo) GetChatsOwnedUser(userId int, offset int) (chats models.ListChatInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, chats.owner_id, units.domain,units.name
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id IN (
			SELECT chats.id
			FROM chats INNER JOIN chat_members 
			ON chats.owner_id = chat_members.user_id
			WHERE chats.owner_id = $1
			LIMIT 20
			OFFSET $2
			)`,
		userId,
		offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.ChatInfo{}
		if err = rows.Scan(&m.ID, &m.OwnerID, &m.Domain, &m.Name); err != nil {
			return
		}
		m.CountMembers, err = r.GetCountChatMembers(m.ID)
		if err != nil {
			return
		}
		chats.Chats = append(chats.Chats, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *ChatsRepo) GetChatsInvolvedUser(userId int, offset int) (chats models.ListChatInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, chats.owner_id, units.domain,units.name
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id IN (
			SELECT chat_id 
			FROM chat_members
			WHERE user_id = $1
			LIMIT 20
			OFFSET $2
			)`,
		userId,
		offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.ChatInfo{}
		if err = rows.Scan(&m.ID, &m.OwnerID, &m.Domain, &m.Name); err != nil {
			return
		}
		m.CountMembers, err = r.GetCountChatMembers(m.ID)
		if err != nil {
			return
		}
		chats.Chats = append(chats.Chats, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *ChatsRepo) GetCountUserChats(userId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id IN (
			SELECT chat_id 
			FROM chat_members
			WHERE user_id = $1
			)`,
		userId,
	).Scan(&count)
	return
}

func (r *ChatsRepo) GetCountRooms(chatId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM rooms
		WHERE chat_id = $1`,
		chatId,
	).Scan(&count)
	return
}

func (r *ChatsRepo) ChatExistsByID(chatId int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM chats
			WHERE id = $1
		)`,
		chatId,
	).Scan(&exists)

	return
}

func (r *ChatsRepo) ChatExistsByDomain(chatDomain string) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM units
			INNER JOIN chats
			ON chats.id = units.id
			WHERE units.domain = $1
		)`,
		chatDomain,
	).Scan(&exists)

	return
}

func (r *ChatsRepo) RemoveUserFromChat(userId int, chatId int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
	).Err()

	return
}

func (r *ChatsRepo) GetCountLinks(chatId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM invite_links  
		WHERE chat_id = $1`,
		chatId,
	).Scan(&count)

	return
}
func (r *ChatsRepo) GetChatLinks(chatId int) (links models.InviteLinks, err error) {
	rows, err := r.db.Query(
		`SELECT code, aliens, expires_at
		FROM invite_links
		WHERE chat_id = $1`,
		chatId,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.InviteLink{}
		if err = rows.Scan(&m.Code, &m.Aliens, &m.Exp); err != nil {
			return
		}
		links.Links = append(links.Links, m)
	}

	return
}
func (r *ChatsRepo) LinkExistsByCode(code string) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM invite_links
			WHERE code = $1
			)`,
		code,
	).Scan(&exists)

	return
}
func (r *ChatsRepo) FindInviteLinkByCode(code string) (link models.InviteLink, err error) {
	err = r.db.QueryRow(
		`SELECT code, aliens, expires_at
		FROM invite_links  
		WHERE code = $1`,
		code,
	).Scan(
		&link.Code,
		&link.Aliens,
		&link.Exp,
	)

	return
}
func (r *ChatsRepo) DeleteInviteLinkByCode(code string) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM invite_links
		WHERE code = $1`,
		code,
	).Err()

	return
}

func (r *ChatsRepo) CreateInviteLink(linkModel *models.CreateInviteLink) (link models.InviteLink, err error) {
	err = r.db.QueryRow(
		`INSERT INTO invite_links (chat_id, aliens, expires_at) 
		VALUES ($1, $2, $3)
		RETURNING code, aliens, expires_at`,
		linkModel.ChatID,
		linkModel.Aliens,
		linkModel.Exp,
	).Scan(
		&link.Code,
		&link.Aliens,
		&link.Exp,
	)

	return
}

func (r *ChatsRepo) InviteLinkIsRelevant(code string) (relevant bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM invite_links
			WHERE code = $1 AND aliens > 0 AND expires_at > $2
			)`,
		code,
		time.Now().UTC().Unix(),
	).Scan(&relevant)
	if !relevant {
		r.db.Exec(
			`DELETE FROM invite_links
			WHERE code = $1`,
			code,
		)
	}

	return
}

func (r *ChatsRepo) AddUserByCode(code string, userId int) (chatId int, err error) {

	err = r.db.QueryRow(
		`WITH l AS (
			UPDATE invite_links
			SET aliens = aliens - 1
			WHERE code = $1
			RETURNING chat_id
			)
		INSERT INTO chat_members (user_id, chat_id)
		VALUES ($2, l.chat_id)`,
		code,
		userId,
	).Scan(&chatId)

	return
}

func (r *ChatsRepo) ChatIsPrivate(chatId int) (private bool) {
	r.db.QueryRow(
		`SELECT private
		FROM chats
		WHERE id = $1`,
		chatId,
	).Scan(&private)

	return
}

func (r *ChatsRepo) AddToBanlist(userId int, chatId int) (err error) {
	err = r.db.QueryRow(
		`INSERT INTO chat_banlist (chat_id, user_id)
		VALUES ($1, $2)`,
		chatId,
		userId,
	).Err()

	return
}

func (r *ChatsRepo) RemoveFromBanlist(userId int, chatId int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM chat_banlist
		WHERE chat_id = $1 AND user_id = $2`,
		chatId,
		userId,
	).Err()

	return
}

func (r *ChatsRepo) UserIsBannedInChat(userId int, chatId int) (banned bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM chat_banlist
			WHERE chat_id = $1 AND user_id = $2
		)`,
		chatId,
		userId,
	).Scan(&banned)

	return
}

func (r *ChatsRepo) GetChatBanlist(chatId int) (users models.ListUserInfo, err error) {
	rows, err := r.db.Query(
		`SELECT id, domain, name
		FROM units
		WHERE id IN (
			SELECT user_id 
			FROM chat_banlist
			WHERE chat_id = $1
			)`,
		chatId,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.UserInfo{}
		if err = rows.Scan(&m.ID, &m.Domain, &m.Name); err != nil {
			return
		}
		users.Users = append(users.Users, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *ChatsRepo) UserRole(userId int, chatId int) (role models.Role, err error) {
	err = r.db.QueryRow(
		`SELECT id, name, color
		FROM roles 
		WHERE id = (
			SELECT role_id
			FROM chat_members
			WHERE user_id = $1 AND chat_id = $2
		)`,
		userId,
		chatId,
	).Scan(
		&role.ID,
		&role.Name,
		&role.Color,
	)

	return
}

func (r *ChatsRepo) CreateRoleInChat(chatId int, roleModel *models.CreateRole) (roleId int, err error) {
	err = r.db.QueryRow(
		`INSERT INTO roles
		(chat_id, role_name, color, visible, manage_rooms, room_id, manage_chat, manage_roles, manage_members)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, 0), $7, $8, $9)
		RETURNING id`,
		chatId,
		roleModel.RoleName,
		roleModel.Color,
		roleModel.Visible,
		roleModel.ManageRooms,
		roleModel.RoomID,
		roleModel.ManageChat,
		roleModel.ManageRoles,
		roleModel.ManageMembers,
	).Scan(&roleId)

	return
}

func (r *ChatsRepo) ChatRoles(chatId int) (roles []models.Role, err error) {
	rows, err := r.db.Query(
		`SELECT id, name, color
		FROM roles 
		WHERE chat_id = $1`,
		chatId,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.Role{}
		if err = rows.Scan(&m.ID, &m.Name, &m.Color); err != nil {
			return
		}
		roles = append(roles, m)
	}

	return
}

func (r *ChatsRepo) GetCountChatRoles(chatId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM roles 
		WHERE chat_id = $1`,
		chatId,
	).Scan(&count)

	return
}

func (r *ChatsRepo) GiveRole(userId int, roleId int) (err error) {
	err = r.db.QueryRow(
		`UPDATE chat_members
		SET role_id = $2
		WHERE user_id = $1`,
		userId,
		roleId,
	).Err()

	return
}

func (r *ChatsRepo) RoleExistsByID(roleId int, chatId int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM roles
			WHERE role_id = $1 AND chat_id = $2
		)`,
		roleId,
		chatId,
	).Scan(&exists)

	return
}

func (r *ChatsRepo) UpdateRoleData(roleId int, inputModel *models.UpdateRole) (err error) {
	err = r.db.QueryRow(
		`UPDATE roles
		SET role_name = $2, color = $3, visible = $4, manage_rooms = $5, room_id =  NULLIF($6, 0), manage_chat = $7, manage_roles = $8, manage_members = $9
		WHERE id = $1`,
		roleId,
		inputModel.RoleName,
		inputModel.Color,
		inputModel.Visible,
		inputModel.ManageRooms,
		inputModel.RoomID,
		inputModel.ManageChat,
		inputModel.ManageRoles,
		inputModel.ManageMembers,
	).Err()

	return
}

func (r *ChatsRepo) DeleteRole(roleId int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM roles
		WHERE id = $1`,
		roleId,
	).Err()

	return
}

func (r *ChatsRepo) TakeRole(userId int, chatId int) (err error) {
	err = r.db.QueryRow(
		`UPDATE chat_members
		SET role_id = NULL
		WHERE user_id = $1`,
		userId,
		chatId,
	).Err()

	return
}

func (r *ChatsRepo) GetMemberInfo(userId int, chatId int) (user models.MemberInfo, err error) {
	err = r.db.QueryRow(
		`SELECT role_id, joined_at
		FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
	).Scan(
		&user.RoleID,
		&user.JoinedAt,
	)

	return
}
